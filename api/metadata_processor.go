package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ProcessCITMetadata(exec boil.Executor, metadata CITMetadata, original, proxy, source *models.File) ([]events.Event, error) {
	return doProcess(exec, metadata, original, proxy, source, nil)
}

// Do all stuff for processing metadata coming from Content Identification Tool.
// 	1. Update properties for original and proxy (film_date, capture_date)
//	2. Update language of original
// 	3. Create content_unit (content_type, dates)
// 	4. Describe content unit (i18ns)
//	5. Add files to new unit
// 	6. Add ancestor files to unit
//  7. Add peer ancestor (related captures)
// 	8. Associate unit with sources, tags, and persons
// 	9. Get or create collection
// 	10. Update collection (content_type, dates, number) if full lesson or new lesson
// 	11. Associate collection and unit
// 	12. Associate unit and derived units
// 	13. Set default permissions ?!
func doProcess(exec boil.Executor, metadata CITMetadata, original, proxy, source *models.File, cu *models.ContentUnit) ([]events.Event, error) {
	isUpdate := cu != nil
	log.Infof("Processing CITMetadata, isUpdate: %t", isUpdate)

	// Update properties for original and proxy (film_date, capture_date)
	filmDate := metadata.CaptureDate
	//if metadata.WeekDate != nil {
	//	filmDate = *metadata.WeekDate
	//}
	if metadata.FilmDate != nil {
		filmDate = *metadata.FilmDate
	}

	evnts := make([]events.Event, 0)

	props := map[string]interface{}{
		"capture_date":      metadata.CaptureDate,
		"film_date":         filmDate,
		"original_language": common.StdLang(metadata.Language),
	}
	log.Infof("Updating files properties: %v", props)
	err := UpdateFileProperties(exec, original, props)
	if err != nil {
		return nil, err
	}
	evnts = append(evnts, events.FileUpdateEvent(original))
	if proxy != nil {
		err = UpdateFileProperties(exec, proxy, props)
		if err != nil {
			return nil, err
		}
		evnts = append(evnts, events.FileUpdateEvent(proxy))
	}

	// Update language of original.
	// TODO: What about proxy !?
	if metadata.HasTranslation {
		original.Language = null.StringFrom(common.LANG_MULTI)
	} else {
		l := common.StdLang(metadata.Language)
		if l == common.LANG_UNKNOWN {
			log.Warnf("Unknown language in metadata %s", metadata.Language)
		}
		original.Language = null.StringFrom(l)
	}
	log.Infof("Updating original.Language to %s", original.Language.String)
	_, err = original.Update(exec, boil.Whitelist("language"))
	if err != nil {
		return nil, errors.Wrap(err, "Save original to DB")
	}

	// Create content_unit (content_type, dates)
	isDerived := metadata.ArtifactType.Valid && metadata.ArtifactType.String != "main"
	ct := metadata.ContentType
	if isDerived {
		// User input is verified below
		ct = strings.ToUpper(metadata.ArtifactType.String)
	}

	var originalProps map[string]interface{}
	err = original.Properties.Unmarshal(&originalProps)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal original properties")
	}
	if duration, ok := originalProps["duration"]; ok {
		props["duration"] = int(duration.(float64))
	} else {
		log.Warnf("Original is missing duration property [%d]", original.ID)
	}

	if metadata.LabelID.Valid {
		props["label_id"] = metadata.LabelID.Int
	}
	if metadata.Number.Valid {
		props["number"] = metadata.Number.Int
	}
	if metadata.Part.Valid {
		props["part"] = metadata.Part.Int
	}

	if isUpdate {
		// content_type
		if ctVal, ok := common.CONTENT_TYPE_REGISTRY.ByName[ct]; !ok {
			return nil, errors.Errorf("Unknown content type %s", ct)
		} else if ctVal.ID != cu.TypeID {
			// update unit's content type
			cu.TypeID = ctVal.ID
			_, err = cu.Update(exec, boil.Whitelist("type_id"))
			if err != nil {
				return nil, errors.Wrap(err, "Update unit type in DB")
			}
		}

		// props
		propsBytes, err := json.Marshal(props)
		if err != nil {
			return nil, errors.Wrap(err, "json Marshal")
		}
		cu.Properties = null.JSONFrom(propsBytes)
		_, err = cu.Update(exec, boil.Whitelist("properties"))
		if err != nil {
			return nil, errors.Wrap(err, "Update unit properties in DB")
		}
	} else {
		log.Infof("Creating content unit of type %s", ct)
		cu, err = CreateContentUnit(exec, ct, props)
		if err != nil {
			return nil, errors.Wrap(err, "Create content unit")
		}
		evnts = append(evnts, events.ContentUnitCreateEvent(cu))

		log.Infof("Describing content unit [%d]", cu.ID)
		err = DescribeContentUnit(exec, cu, metadata)
		if err != nil {
			log.Errorf("Error describing content unit: %s", err.Error())
		}
	}

	// we lookup Original's capture_stop operation as it holds required information below
	var captureStopProps map[string]interface{}
	captureStop, err := FindUpChainOperation(exec, original.ID, common.OP_CAPTURE_STOP)
	if err != nil {
		if ex, ok := err.(UpChainOperationNotFound); ok {
			log.Warnf("capture_stop operation not found for original: %s", ex.Error())
		}
	} else {
		if captureStop.Properties.Valid {
			err = json.Unmarshal(captureStop.Properties.JSON, &captureStopProps)
			if err != nil {
				return nil, errors.Wrap(err, "json Unmarshal")
			}
		}
	}

	// Add files to new unit
	log.Info("Adding files to unit")
	err = cu.AddFiles(exec, false, original)
	if err != nil {
		return nil, errors.Wrap(err, "Add original to unit")
	}
	if proxy != nil {
		err = cu.AddFiles(exec, false, proxy)
		if err != nil {
			return nil, errors.Wrap(err, "Add proxy to unit")
		}
	}
	if source != nil {
		err = cu.AddFiles(exec, false, source)
		if err != nil {
			return nil, errors.Wrap(err, "Add source to unit")
		}
	}

	// Add ancestor files to unit (not for derived units)
	if !isDerived && !isUpdate {
		log.Info("Main unit, adding ancestors...")
		ancestors, err := FindFileAncestors(exec, original.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Find original's ancestors")
		}

		if proxy != nil {
			err = proxy.L.LoadParent(exec, true, proxy, nil)
			if err != nil {
				return nil, errors.Wrap(err, "Load proxy's parent")
			}
			if proxy.R.Parent != nil {
				ancestors = append(ancestors, proxy.R.Parent)
			}
		}

		err = cu.AddFiles(exec, false, ancestors...)
		if err != nil {
			return nil, errors.Wrap(err, "Add ancestors to unit")
		}
		log.Infof("Added %d ancestors", len(ancestors))
		for i := range ancestors {
			x := ancestors[i]
			evnts = append(evnts, events.FileUpdateEvent(x))
			log.Infof("%s [%d]", x.Name, x.ID)
		}

		// Add peer ancestor (related captures)
		if workflowID, ok := captureStopProps["workflow_id"]; ok {
			// find other captures
			relatedCaptures, err := FindOperationsByWorkflowID(exec, workflowID, common.OP_CAPTURE_STOP)
			if err != nil {
				return nil, errors.Wrap(err, "find related captures")
			}

			// find other captures files and add them all to this unit
			var relatedCapturesFiles []*models.File
			for _, capture := range relatedCaptures {
				if capture.ID == captureStop.ID {
					continue
				}
				if err := capture.L.LoadFiles(exec, true, capture, nil); err != nil {
					return nil, errors.Wrapf(err, "load related capture files %d", capture.ID)
				}
				captureFile := capture.R.Files[0]
				relatedCapturesFiles = append(relatedCapturesFiles, captureFile)
				if files, err := FindFileDescendants(exec, captureFile.ID); err != nil {
					return nil, errors.Wrapf(err, "load descendants of related capture file %d", captureFile.ID)
				} else {
					relatedCapturesFiles = append(relatedCapturesFiles, files...)
				}
			}

			err = cu.AddFiles(exec, false, relatedCapturesFiles...)
			if err != nil {
				return nil, errors.Wrap(err, "Add related capture files to unit")
			}
			log.Infof("Added %d related capture files", len(relatedCapturesFiles))
			for _, f := range relatedCapturesFiles {
				evnts = append(evnts, events.FileUpdateEvent(f))
				log.Infof("%s [%d]", f.Name, f.ID)
			}
		} else {
			log.Info("capture_stop not found or its missing workflow_id. Skipping related captures associations")
		}
	}

	// Associate unit with sources, tags, and persons
	if len(metadata.Sources) > 0 {
		log.Infof("Associating %d sources", len(metadata.Sources))
		sources, err := models.Sources(
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Sources)...)).
			All(exec)
		if err != nil {
			return nil, errors.Wrap(err, "Lookup sources in DB")
		}

		// are we missing some source ?
		if len(sources) != len(metadata.Sources) {
			missing := make([]string, 0)
			for _, x := range metadata.Sources {
				found := false
				for _, y := range sources {
					if x == y.UID {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, x)
				}
			}
			log.Warnf("Unknown sources: %s", missing)
		}

		err = cu.SetSources(exec, false, sources...)
		if err != nil {
			return nil, errors.Wrap(err, "Associate sources")
		}

		seriesEvnts, err := associateLessonsSeriesSources(exec, cu, metadata.Sources)
		if err != nil {
			return nil, errors.Wrap(err, "Associate Lessons series collection by sources")
		}

		evnts = append(evnts, seriesEvnts...)
	}

	if len(metadata.Tags) > 0 {
		log.Infof("Associating %d tags", len(metadata.Tags))
		tags, err := models.Tags(
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Tags)...)).
			All(exec)
		if err != nil {
			return nil, errors.Wrap(err, "Lookup tags  in DB")
		}

		// are we missing some tag ?
		if len(tags) != len(metadata.Tags) {
			missing := make([]string, 0)
			for _, x := range metadata.Tags {
				found := false
				for _, y := range tags {
					if x == y.UID {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, x)
				}
			}
			log.Warnf("Unknown sources: %s", missing)
		}
		err = cu.SetTags(exec, false, tags...)
		if err != nil {
			return nil, errors.Wrap(err, "Associate tags")
		}
	}
	if len(metadata.Likutim) > 0 {
		log.Infof("Associating %d likutim", len(metadata.Likutim))
		likutim, err := models.ContentUnits(
			qm.Select("distinct on (\"content_units\".id) \"content_units\".*"),
			qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID),
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Likutim)...)).
			All(exec)
		if err != nil {
			return nil, errors.Wrap(err, "Lookup tags  in DB")
		}

		// are we missing some unit ?
		if len(likutim) != len(metadata.Likutim) {
			missing := make([]string, 0)
			for _, x := range metadata.Likutim {
				found := false
				for _, y := range likutim {
					if x == y.UID {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, x)
				}
			}
			log.Warnf("Unknown likutim: %s", missing)
		}

		derivations := make([]*models.ContentUnitDerivation, len(likutim))
		for i, l := range likutim {
			cud := &models.ContentUnitDerivation{
				SourceID:  cu.ID,
				DerivedID: l.ID,
			}
			derivations[i] = cud
		}
		err = cu.AddSourceContentUnitDerivations(exec, true, derivations...)
		if err != nil {
			return nil, errors.Wrap(err, "Associate likutim")
		}
		for _, l := range likutim {
			evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(l))
		}
		evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(cu))
		seriesEvnts, err := associateLessonsSeriesLikutim(exec, cu, metadata.Likutim)
		if err != nil {
			return nil, errors.Wrap(err, "Associate Lessons series collection by likutim")
		}

		evnts = append(evnts, seriesEvnts...)

	}

	// Handle persons ...
	if strings.ToLower(metadata.Lecturer) == common.P_RAV {
		log.Info("Associating unit to rav")
		cup := &models.ContentUnitsPerson{
			ContentUnitID: cu.ID,
			PersonID:      common.PERSON_REGISTRY.ByPattern[common.P_RAV].ID,
			RoleID:        common.CONTENT_ROLE_TYPE_REGISTRY.ByName[common.CR_LECTURER].ID,
		}

		// upsert make sure we either have such relation or insert a new one
		err = cup.Upsert(exec, false, nil, boil.Infer(), boil.Infer())
		if err != nil {
			return nil, errors.Wrap(err, "Associate persons")
		}
	} else if isUpdate && strings.ToLower(metadata.Lecturer) == "norav" {
		// in update mode, if norav so we remove relation to rav (if any)
		cup := &models.ContentUnitsPerson{
			ContentUnitID: cu.ID,
			PersonID:      common.PERSON_REGISTRY.ByPattern[common.P_RAV].ID,
		}
		_, err = cup.Delete(exec)
		if err != nil {
			return nil, errors.Wrap(err, "Delete Rav association")
		}
	} else {
		log.Infof("Unknown lecturer %s, skipping person association.", metadata.Lecturer)
	}

	// Get or create collection
	if metadata.CollectionUID.Valid {
		log.Infof("Specific collection %s", metadata.CollectionUID.String)

		// find collection
		c, err := models.Collections(qm.Where("uid = ?", metadata.CollectionUID.String)).One(exec)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("No such collection %s", metadata.CollectionUID.String)
			} else {
				return nil, errors.Wrap(err, "Lookup collection in DB")
			}
		}

		// Associate unit to collection
		if c != nil &&
			(!metadata.ArtifactType.Valid || metadata.ArtifactType.String == "main") {
			err := associateUnitToCollection(exec, cu, c, metadata)
			if err != nil {
				return nil, errors.Wrap(err, "associate content_unit to collection")
			}
			evnts = append(evnts, events.CollectionContentUnitsChangeEvent(c))
		}
	}

	// Update mode ends here
	if isUpdate {
		return evnts, nil
	}

	if ct == common.CT_LESSON_PART ||
		ct == common.CT_FULL_LESSON ||
		ct == common.CT_KTAIM_NIVCHARIM {
		log.Info("Lesson reconciliation")

		// Reconcile or create new
		// Reconciliation is done by looking up the operation chain of original to capture_stop.
		// There we have a property of saying the capture_id of the full lesson capture.
		if captureID, ok := captureStopProps["collection_uid"]; ok {
			log.Infof("Reconcile by capture_id %s", captureID)
			var cct string
			if metadata.WeekDate == nil {
				cct = common.CT_DAILY_LESSON
			} else {
				cct = common.CT_SPECIAL_LESSON
			}

			// Keep this property on the collection for other parts to find it
			props["capture_id"] = captureID
			if metadata.Number.Valid {
				props["number"] = metadata.Number.Int
			}
			delete(props, "duration")
			delete(props, "part")

			// get or create collection
			c, err := FindCollectionByCaptureID(exec, captureID)
			if err != nil {
				if _, ok := err.(CollectionNotFound); !ok {
					return nil, err
				}

				// Create new collection
				log.Info("Creating new collection")
				c, err = CreateCollection(exec, cct, props)
				if err != nil {
					return nil, err
				}
				evnts = append(evnts, events.CollectionCreateEvent(c))
			} else if ct == common.CT_FULL_LESSON {
				// Update collection properties to those of full lesson
				log.Info("Full lesson, overriding collection properties")
				if c.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[cct].ID {
					log.Infof("Full lesson, content_type changed to %s", cct)
					c.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[cct].ID
					_, err = c.Update(exec, boil.Whitelist("type_id"))
					if err != nil {
						return nil, errors.Wrap(err, "Update collection type in DB")
					}
				}

				err = UpdateCollectionProperties(exec, c, props)
				if err != nil {
					return nil, err
				}
				evnts = append(evnts, events.CollectionUpdateEvent(c))
			}

			// Associate unit to collection
			if c != nil &&
				(!metadata.ArtifactType.Valid ||
					metadata.ArtifactType.String == "main" ||
					metadata.ArtifactType.String == "KTAIM_NIVCHARIM") {
				err := associateUnitToCollection(exec, cu, c, metadata)
				if err != nil {
					return nil, errors.Wrap(err, "associate content_unit to collection")
				}
				evnts = append(evnts, events.CollectionContentUnitsChangeEvent(c))
			}
		} else {
			log.Warn("capture_stop not found or its missing collection_uid. Skipping lesson reconciliation")
		}
	}

	// Associate unit and derived units
	// We take into account that a derived content unit arrives before it's source content unit.
	// Such cases are possible due to the studio operator actions sequence.
	err = original.L.LoadParent(exec, true, original, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Load original's parent")
	}

	if original.R.Parent == nil {
		log.Warn("We don't have original's parent file. Skipping derived units association.")
	} else {
		log.Info("Processing derived units associations")
		if !metadata.ArtifactType.Valid ||
			metadata.ArtifactType.String == "main" {
			// main content unit
			log.Info("We're the main content unit")

			log.Info("Looking up pending derived units")
			derivedCUs, err := mainToDerived(exec, metadata, original)
			if err != nil {
				return nil, err
			}

			log.Infof("%d derived units pending our association", len(derivedCUs))
			for k, v := range derivedCUs {
				log.Infof("DerivedID: %d, Name: %s", k, v)
				cud := &models.ContentUnitDerivation{
					DerivedID: k,
					Name:      v,
				}
				err = cu.AddSourceContentUnitDerivations(exec, true, cud)
				if err != nil {
					return nil, errors.Wrap(err, "Save derived unit association in DB")
				}

				_, err = queries.Raw(`UPDATE content_units SET properties = properties - 'artifact_type' WHERE id = $1`, k).
					Exec(exec)
				if err != nil {
					return nil, errors.Wrap(err, "Delete derived unit artifact_type property from DB")
				}
			}

			if len(derivedCUs) > 0 {
				evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(cu))
			}

		} else {
			// derived content unit
			log.Info("We're the derived content unit")

			mainCUID, err := derivedToMain(exec, metadata, cu, original)
			if err != nil {
				return nil, err
			}

			if mainCUID == 0 {
				// save artifact type for later use (when main unit appears)
				log.Info("Main content unit not found, saving artifact_type property")
				err = UpdateContentUnitProperties(exec, cu, map[string]interface{}{
					"artifact_type": metadata.ArtifactType.String,
				})
				if err != nil {
					return nil, err
				}
			} else {
				// main content unit already exists
				log.Infof("Main content unit exists %d", mainCUID)
				cud := &models.ContentUnitDerivation{
					SourceID: mainCUID,
					Name:     metadata.ArtifactType.String,
				}
				err = cu.AddDerivedContentUnitDerivations(exec, true, cud)
				if err != nil {
					return nil, errors.Wrap(err, "Save source unit in DB")
				}
				evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(cu))
			}
		}
	}

	// set default permissions ?!

	return evnts, nil
}

func associateUnitToCollection(exec boil.Executor, cu *models.ContentUnit, c *models.Collection, metadata CITMetadata) error {
	log.Infof("Associating unit and collection [c-cu]=[%d-%d]", c.ID, cu.ID)

	ccu := &models.CollectionsContentUnit{
		CollectionID:  c.ID,
		ContentUnitID: cu.ID,
	}

	switch common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name {
	case common.CT_FULL_LESSON:
		if c.TypeID == common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID ||
			c.TypeID == common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID {
			ccu.Name = "full"
		} else if metadata.Number.Valid {
			ccu.Name = strconv.Itoa(metadata.Number.Int)
		}
		break
	case common.CT_LESSON_PART:
		if metadata.Part.Valid {
			ccu.Name = strconv.Itoa(metadata.Part.Int)
		}
		break
	case common.CT_VIDEO_PROGRAM_CHAPTER:
		if metadata.Episode.Valid {
			ccu.Name = metadata.Episode.String
		}
		break
	default:
		if metadata.Number.Valid {
			ccu.Name = strconv.Itoa(metadata.Number.Int)
		}

		// first 3 event part types are lesson, YH and meal, we skip them.
		if metadata.PartType.Valid && metadata.PartType.Int > 2 {
			idx := metadata.PartType.Int - 3
			if idx < len(common.MISC_EVENT_PART_TYPES) {
				ccu.Name = common.MISC_EVENT_PART_TYPES[idx] + ccu.Name
			} else {
				log.Warnf("Unknown event part type: %d", metadata.PartType.Int)
			}
		}
		break
	}
	if metadata.ArtifactType.Valid &&
		metadata.ArtifactType.String != "main" {
		ccu.Name = fmt.Sprintf("%s_%s", metadata.ArtifactType.String, ccu.Name)
	}

	// Make this new unit the last one in this collection
	var err error
	ccu.Position, err = GetNextPositionInCollection(exec, c.ID)
	if err != nil {
		return errors.Wrap(err, "Get last position in collection")
	}

	log.Infof("Association name: %s", ccu.Name)
	err = ccu.Upsert(exec, true,
		[]string{"collection_id", "content_unit_id"},
		boil.Whitelist("name", "position"),
		boil.Infer())
	if err != nil {
		return errors.Wrap(err, "Save collection and content unit association in DB")
	}

	return nil
}

func mainToDerived(exec boil.Executor, metadata CITMetadata, original *models.File) (map[int64]string, error) {
	part := -888 // something we never use for part
	if metadata.Part.Valid {
		part = metadata.Part.Int
	}

	// We lookup original's siblings for derived content units that arrived before us.
	// We then associate them with us and remove their "unprocessed" mark.
	// Meaning, the presence of "artifact_type" property
	rows, err := queries.Raw(
		`SELECT
  cu.id,
  cu.properties ->> 'artifact_type'
FROM content_units cu
  INNER JOIN files f ON f.content_unit_id = cu.id AND f.parent_id = $1
WHERE cu.properties ? 'artifact_type' AND (cu.properties ->> 'part') :: INT = $2`,
		original.ParentID.Int64, part).
		Query(exec)
	if err != nil {
		return nil, errors.Wrap(err, "Load derived content units")
	}

	// put results in a map first since we can't process them while iterating.
	// see this bug:  https://github.com/lib/pq/issues/81
	derivedCUs := make(map[int64]string)
	for rows.Next() {
		var cuid int64
		var artifactType string
		err = rows.Scan(&cuid, &artifactType)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}
		derivedCUs[cuid] = artifactType
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Iter rows")
	}
	err = rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Close rows")
	}

	return derivedCUs, nil
}

func derivedToMain(exec boil.Executor, metadata CITMetadata, cu *models.ContentUnit, original *models.File) (int64, error) {
	part := -888 // something we never use for part
	if metadata.Part.Valid {
		part = metadata.Part.Int
	}

	mainCT, ok := common.CONTENT_TYPE_REGISTRY.ByName[metadata.ContentType]
	if !ok {
		return 0, errors.Errorf("Unknown content type %s", metadata.ContentType)
	}

	var cuID int64
	err := queries.Raw(`
SELECT cu.id
FROM content_units cu
  INNER JOIN files f ON f.content_unit_id = cu.id AND f.parent_id = $1 AND cu.id != $2 AND cu.type_id = $3
WHERE (cu.properties ->> 'part') :: INT = $4`,
		original.ParentID.Int64, cu.ID, mainCT.ID, part).QueryRow(exec).Scan(&cuID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, errors.Wrap(err, "Query main CU ID")
		}
	}

	return cuID, nil
}

var MIN_CU_NUMBER_FOR_NEW_LESSON_SERIES = 3
var DAYS_CHECK_FOR_LESSONS_SERIES = 30

func associateLessonsSeriesSources(exec boil.Executor, cu *models.ContentUnit, sUids []string) ([]events.Event, error) {
	evnts := make([]events.Event, 0)
	sByLeaf, err := MapParentByLeaf(exec, sUids)
	if err != nil {
		return nil, NewInternalError(err)
	}
	q := fmt.Sprintf(`
  SELECT DISTINCT ON(s.id) s.uid, array_agg(DISTINCT cu.id)
  FROM content_units cu
  INNER JOIN content_units_sources cus ON cu.id = cus.content_unit_id
  INNER JOIN sources s ON s.id = cus.source_id
  WHERE cu.type_id = $1 
  AND coalesce((cu.properties->>'film_date')::date, cu.created_at) > (CURRENT_DATE - '%d day'::interval)
  AND cu.published = TRUE AND cu.secure = 0
  AND s.uid IN (%s)
  GROUP BY  s.id
`, DAYS_CHECK_FOR_LESSONS_SERIES, fmt.Sprintf("'%s'", strings.Join(sUids, "','")))

	rows, err := queries.Raw(q, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID).Query(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}
	defer rows.Close()

	cuByS := make(map[string][]int64)
	var sId string
	var cuIdsByS pq.Int64Array
	var cuIds []int64
	for rows.Next() {
		err = rows.Scan(&sId, &cuIdsByS)
		if err != nil {
			return nil, NewInternalError(err)
		}
		cuByS[sByLeaf[sId]] = cuIdsByS
		cuIds = append(cuIds, cuIdsByS...)
	}

	cus, err := models.ContentUnits(
		models.ContentUnitWhere.ID.IN(cuIds),
		qm.Load("Sources"),
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
	).All(exec)

	if err != nil {
		return nil, NewInternalError(err)
	}

	var cByS = make(map[string]*models.Collection)
	var startDateByS = make(map[string]string)
	//group lesson series by cus sources
	for _, _cu := range append(cus, cu) {
		for _, s := range _cu.R.Sources {
			sUid := sByLeaf[s.UID]
			if _, ok := startDateByS[s.UID]; !ok && cuByS[sUid] != nil && cuByS[sUid][0] == _cu.ID {
				var _props map[string]interface{}
				if err := _cu.Properties.Unmarshal(&_props); err != nil {
					continue
				}
				if fd, ok := _props["film_date"]; ok {
					startDateByS[s.UID] = fd.(string)
				}
			}
		}
		for _, ccu := range _cu.R.CollectionsContentUnits {
			if !ccu.R.Collection.Properties.Valid || ccu.R.Collection.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID {
				continue
			}
			var props map[string]interface{}
			if err := ccu.R.Collection.Properties.Unmarshal(&props); err != nil {
				continue
			}
			_sUid := props["source"].(string)
			if _, ok := cByS[_sUid]; ok {
				continue
			}
			if _sUid != "" {
				cByS[_sUid] = ccu.R.Collection
			}
		}
	}

	var cuProps map[string]interface{}
	if err := cu.Properties.Unmarshal(&cuProps); err != nil {
		return nil, err
	}
	needAddByC := make(map[int64][]*models.CollectionsContentUnit)

	//find and if need create lessons series per new cu's sources
	for _, s := range cu.R.Sources {
		sUid := sByLeaf[s.UID]
		if _, ok := cByS[s.UID]; len(cuByS[sUid]) < MIN_CU_NUMBER_FOR_NEW_LESSON_SERIES && !ok {
			continue
		}
		if _, ok := cByS[s.UID]; !ok {
			props := map[string]interface{}{
				"source":     s.UID,
				"end_date":   cuProps["film_date"],
				"start_date": startDateByS[s.UID],
			}
			c, err := CreateCollection(exec, common.CT_LESSONS_SERIES, props)
			if err != nil {
				return nil, NewInternalError(err)
			}
			c.Published = true
			_, err = c.Update(exec, boil.Whitelist("published"))
			if err != nil {
				return nil, NewInternalError(err)
			}
			cByS[s.UID] = c
			//add prev cus to new lessons series
			for i, id := range cuByS[sUid] {
				ccu := &models.CollectionsContentUnit{
					ContentUnitID: id,
					Position:      i,
				}
				needAddByC[cByS[s.UID].ID] = append(needAddByC[cByS[s.UID].ID], ccu)
			}
			evnts = append(evnts, events.CollectionCreateEvent(c))
		}
		ccu := &models.CollectionsContentUnit{
			ContentUnitID: cu.ID,
			Position:      len(cuByS[sUid]) + 1,
		}
		needAddByC[cByS[s.UID].ID] = append(needAddByC[cByS[s.UID].ID], ccu)
	}

	for _, c := range cByS {
		if err := c.AddCollectionsContentUnits(exec, true, needAddByC[c.ID]...); err != nil {
			return nil, err
		}
		if err := UpdateCollectionProperties(exec, c, map[string]interface{}{"end_date": cuProps["film_date"]}); err != nil {
			return nil, err
		}
		if _, err := c.Update(exec, boil.Infer()); err != nil {
			return nil, err
		}
		c, _ := models.Collections(
			models.CollectionWhere.ID.EQ(c.ID),
			qm.Load("CollectionsContentUnits"),
			qm.Load("CollectionsContentUnits.ContentUnit"),
		).One(exec)
		evnts = append(evnts, events.CollectionContentUnitsChangeEvent(c))
	}
	return evnts, nil
}

var TES_PARTS_UIDS = []string{"9xNFLSSp", "XlukqLH8", "AerA1hNN", "1kDKQxJb", "o5lXptLo", "eNwJXy4s", "ahipVtPu", "Pscnn3pP", "Lfu7W3CD", "n03vXCJl", "UGcGGSpP", "NpLQT0LX", "AUArdCkH", "tit6XNAo", "FaKUG7ru", "mW6eON0z"}
var ZOAR_UID = "AwGBQX2L"
var ZOAR_PART_ONE_UID = "cSyh3vQM"

func MapParentByLeaf(exec boil.Executor, uids []string) (map[string]string, error) {

	q := fmt.Sprintf(`
WITH RECURSIVE recurcive_s(id, uid, parent_id, start_uid) AS(
	SELECT id, uid, parent_id, uid
		FROM sources where uid IN (%s)
	UNION
	SELECT s.id, s.uid, s.parent_id, rs.start_uid
		FROM recurcive_s rs, sources s where rs.parent_id = s.id
)
SELECT start_uid, uid FROM recurcive_s WHERE uid IN (%s)
`,
		fmt.Sprintf("'%s'", strings.Join(uids, "','")),
		fmt.Sprintf("'%s'", strings.Join(append(TES_PARTS_UIDS, ZOAR_UID, ZOAR_PART_ONE_UID), "','")),
	)
	rows, err := queries.Raw(q).Query(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}
	defer rows.Close()

	var uid string
	var nUid string
	newByOldUid := make(map[string]string)
	checkIsNewOnce := make(map[string]bool)
	for rows.Next() {
		err = rows.Scan(&uid, &nUid)
		if err != nil {
			return nil, NewInternalError(err)
		}
		if _, ok := checkIsNewOnce[nUid]; ok {
			continue
		}
		checkIsNewOnce[nUid] = true
		if _, ok := newByOldUid[uid]; !ok || nUid == ZOAR_PART_ONE_UID {
			newByOldUid[uid] = nUid
		}
	}

	for _, uid := range uids {
		if _, ok := newByOldUid[uid]; !ok {
			newByOldUid[uid] = uid
		}
	}
	return newByOldUid, nil
}

func associateLessonsSeriesLikutim(exec boil.Executor, cu *models.ContentUnit, lUids []string) ([]events.Event, error) {
	evnts := make([]events.Event, 0)

	q := fmt.Sprintf(`
  SELECT DISTINCT ON(lcu.id) lcu.id, array_agg(DISTINCT cu.id)
  FROM content_units cu
  INNER JOIN content_unit_derivations dcu ON cu.id = dcu.source_id
  INNER JOIN content_units lcu ON lcu.id = dcu.derived_id
  WHERE cu.type_id = $1
  AND coalesce((cu.properties->>'film_date')::date, cu.created_at) > (CURRENT_DATE - '%d day'::interval)
  AND cu.published = TRUE AND cu.secure = 0
  AND lcu.uid IN (%s)
  GROUP BY  lcu.id
`, DAYS_CHECK_FOR_LESSONS_SERIES, fmt.Sprintf("'%s'", strings.Join(lUids, "','")))
	rows, err := queries.Raw(q, common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSON_PART].ID).Query(exec)

	if err != nil {
		return nil, NewInternalError(err)
	}
	defer rows.Close()

	cuByL := make(map[int64][]int64)
	var lId int64
	var cuIdsByL pq.Int64Array
	var cuIds []int64
	for rows.Next() {
		err = rows.Scan(&lId, &cuIdsByL)
		if err != nil {
			return nil, NewInternalError(err)
		}
		cuByL[lId] = cuIdsByL
		cuIds = append(cuIds, cuIdsByL...)
	}

	cus, err := models.ContentUnits(
		models.ContentUnitWhere.ID.IN(append(cuIds, cu.ID)),
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.Collection"),
		qm.Load("SourceContentUnitDerivations"),
		qm.Load("SourceContentUnitDerivations.Derived"),
	).All(exec)
	if err != nil {
		return nil, NewInternalError(err)
	}
	var cByL = make(map[string]*models.Collection)
	var startDateByL = make(map[string]string)
	//group lesson series by cus sources
	for _, _cu := range cus {
		for _, l := range _cu.R.SourceContentUnitDerivations {
			if l.R.Derived == nil {
				continue
			}
			if _, ok := startDateByL[l.R.Derived.UID]; !ok && cuByL[l.R.Derived.ID] != nil && cuByL[l.R.Derived.ID][0] == _cu.ID {
				var _props map[string]interface{}
				if err := _cu.Properties.Unmarshal(&_props); err != nil {
					continue
				}

				if fd, ok := _props["film_date"]; ok {
					startDateByL[l.R.Derived.UID] = fd.(string)
				}
			}
		}
		for _, ccu := range _cu.R.CollectionsContentUnits {
			if !ccu.R.Collection.Properties.Valid || ccu.R.Collection.TypeID != common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LESSONS_SERIES].ID {
				continue
			}
			var props map[string]interface{}
			if err := ccu.R.Collection.Properties.Unmarshal(&props); err != nil {
				continue
			}
			_sUid := props["source"].(string)
			if _, ok := cByL[_sUid]; ok {
				continue
			}
			if _sUid != "" {
				cByL[_sUid] = ccu.R.Collection
			}
		}
	}

	var cuProps map[string]interface{}
	if err := cu.Properties.Unmarshal(&cuProps); err != nil {
		return nil, err
	}
	needAddByC := make(map[int64][]*models.CollectionsContentUnit)

	//find and if need create lessons series per new cu's sources
	for _, dcu := range cu.R.SourceContentUnitDerivations {
		if err := dcu.L.LoadDerived(exec, true, dcu, nil); err != nil {
			return nil, err
		}
		if _, ok := cByL[dcu.R.Derived.UID]; len(cuByL[dcu.DerivedID]) < MIN_CU_NUMBER_FOR_NEW_LESSON_SERIES && !ok {
			continue
		}
		if _, ok := cByL[dcu.R.Derived.UID]; !ok {
			props := map[string]interface{}{
				"source":     dcu.R.Derived.UID,
				"end_date":   cuProps["film_date"],
				"start_date": startDateByL[dcu.R.Derived.UID],
			}
			c, err := CreateCollection(exec, common.CT_LESSONS_SERIES, props)
			if err != nil {
				return nil, NewInternalError(err)
			}
			c.Published = true
			_, err = c.Update(exec, boil.Whitelist("published"))
			if err != nil {
				return nil, NewInternalError(err)
			}
			cByL[dcu.R.Derived.UID] = c
			//add prev cus to new lessons series
			for i, id := range cuByL[dcu.R.Derived.ID] {
				ccu := &models.CollectionsContentUnit{
					ContentUnitID: id,
					Position:      i,
				}
				needAddByC[cByL[dcu.R.Derived.UID].ID] = append(needAddByC[cByL[dcu.R.Derived.UID].ID], ccu)
			}
			evnts = append(evnts, events.CollectionCreateEvent(c))
		}
		ccu := &models.CollectionsContentUnit{
			ContentUnitID: cu.ID,
			Position:      len(cuByL[dcu.R.Derived.ID]) + 1,
		}
		needAddByC[cByL[dcu.R.Derived.UID].ID] = append(needAddByC[cByL[dcu.R.Derived.UID].ID], ccu)
	}

	for _, c := range cByL {
		if err := c.AddCollectionsContentUnits(exec, true, needAddByC[c.ID]...); err != nil {
			return nil, err
		}
		if err := UpdateCollectionProperties(exec, c, map[string]interface{}{"end_date": cuProps["film_date"]}); err != nil {
			return nil, err
		}
		if _, err := c.Update(exec, boil.Infer()); err != nil {
			return nil, err
		}
		evnts = append(evnts, events.CollectionContentUnitsChangeEvent(c))
	}
	return evnts, nil
}

/* send-fix

Sometimes, after a unit was created in a send operation,
we need to fix it.

Either the metadata that was given is wrong or we might need a different trim.
For such cases a new button in the trim admin "fix" is made.

The workflow simply shows the same CIT screen and unit selection.
We should:

1. re-process the metadata in "update" mode
2. figure out files for removal
3. mark those as removed
4. update the unit's published status

*/
func ProcessCITMetadataUpdate(exec boil.Executor, metadata CITMetadata, original, proxy, source *models.File) ([]events.Event, error) {
	unit, err := models.ContentUnits(qm.Where("uid = ?", metadata.UnitToFixUID.String)).One(exec)
	if err != nil {
		return nil, errors.Wrapf(err, "lookup unit UID %s", metadata.UnitToFixUID.String)
	}

	evnts, err := doProcess(exec, metadata, original, proxy, source, unit)
	if err != nil {
		return nil, errors.Wrap(err, "doProcess")
	}

	// We remove only files generated in convert (carbon)
	// and previous trimmed not in our path.
	// Other, manually inserted, files are not touched and are left to admin
	// to figure out what to do with them.

	// Figure out merged set of file IDs
	// which are either ancestor of original or proxy of source
	// These should be excluded from removal.
	mutualAncestors := hashset.New()

	oPath, err := FindFileAncestors(exec, original.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "lookup original ancestors %d", original.ID)
	}
	for i := range oPath {
		mutualAncestors.Add(oPath[i].ID)
	}

	if proxy != nil {
		pPath, err := FindFileAncestors(exec, proxy.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "lookup proxy ancestors %d", original.ID)
		}
		for i := range pPath {
			mutualAncestors.Add(pPath[i].ID)
		}
	}

	if source != nil {
		sPath, err := FindFileAncestors(exec, source.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "lookup source ancestors %d", original.ID)
		}
		for i := range sPath {
			mutualAncestors.Add(sPath[i].ID)
		}
	}

	ancestorsIDs := mutualAncestors.Values()
	// These are the fix. Not the problem. Don't remove them
	ancestorsIDs = append(ancestorsIDs, original.ID)
	if proxy != nil {
		ancestorsIDs = append(ancestorsIDs, original.ID, proxy.ID)
	}

	log.Infof("ancestorsIDs: %v", ancestorsIDs)

	// fetch file IDs to remove
	var fIDs pq.Int64Array
	q := `SELECT array_agg(distinct f.id)
FROM files f
  INNER JOIN files_operations fo ON f.id = fo.file_id
  INNER JOIN operations o ON fo.operation_id = o.id AND o.type_id = ANY($1)
WHERE f.content_unit_id = $2 AND NOT f.id = ANY($3) 
`
	err = queries.Raw(q, pq.Array([]int64{
		common.OPERATION_TYPE_REGISTRY.ByName[common.OP_TRIM].ID,
		common.OPERATION_TYPE_REGISTRY.ByName[common.OP_CONVERT].ID,
	}), unit.ID, pq.Array(ancestorsIDs)).QueryRow(exec).Scan(&fIDs)
	if err != nil {
		return nil, errors.Wrap(err, "fetch file IDs to remove")
	}

	log.Infof("%d files to remove: %v", len(fIDs), fIDs)
	wasPublished := false
	if len(fIDs) > 0 {
		// actual removal
		_, err = models.Files(
			qm.WhereIn("id in ?", utils.ConvertArgsInt64(fIDs)...)).
			UpdateAll(exec, models.M{
				"removed_at": null.TimeFrom(time.Now().UTC()),
			})
		if err != nil {
			return nil, errors.Wrap(err, "Update files to remove")
		}

		// file removed events
		removedFiles, err := models.Files(
			qm.Select("id", "uid", "published"),
			qm.WhereIn("id in ?", utils.ConvertArgsInt64(fIDs)...)).
			All(exec)
		if err != nil {
			return nil, errors.Wrap(err, "Refresh files to remove")
		}

		for i := range removedFiles {
			evnts = append(evnts, events.FileRemoveEvent(removedFiles[i]))
			wasPublished = wasPublished || removedFiles[i].Published
		}

	}

	// unit published status change
	impact, err := FileLeftUnitImpact(exec, wasPublished, unit.ID)
	if err != nil {
		return nil, errors.Wrap(err, "File left impact")
	}
	evnts = append(evnts, impact.Events()...)

	return evnts, nil
}
