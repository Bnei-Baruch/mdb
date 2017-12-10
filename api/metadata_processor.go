package api

import (
	"database/sql"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"encoding/json"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/volatiletech/sqlboiler/queries"
	"strings"
)

// Do all stuff for processing metadata coming from Content Identification Tool.
// 	1. Update properties for original and proxy (film_date, capture_date)
//	2. Update language of original
// 	3. Create content_unit (content_type, dates)
// 	4. Describe content unit (i18ns)
//	5. Add files to new unit
// 	6. Add ancestor files to unit
// 	7. Associate unit with sources, tags, and persons
// 	8. Get or create collection
// 	9. Update collection (content_type, dates, number) if full lesson or new lesson
// 	10. Associate collection and unit
// 	11. Associate unit and derived units
// 	12. Set default permissions ?!
func ProcessCITMetadata(exec boil.Executor, metadata CITMetadata, original, proxy *models.File) ([]events.Event, error) {
	log.Info("Processing CITMetadata")

	// Update properties for original and proxy (film_date, capture_date)
	filmDate := metadata.CaptureDate
	if metadata.WeekDate != nil {
		filmDate = *metadata.WeekDate
	}
	if metadata.FilmDate != nil {
		filmDate = *metadata.FilmDate
	}

	props := map[string]interface{}{
		"capture_date":      metadata.CaptureDate,
		"film_date":         filmDate,
		"original_language": StdLang(metadata.Language),
	}
	log.Infof("Updating files properties: %v", props)
	err := UpdateFileProperties(exec, original, props)
	if err != nil {
		return nil, err
	}
	err = UpdateFileProperties(exec, proxy, props)
	if err != nil {
		return nil, err
	}

	evnts := make([]events.Event, 2)
	evnts[0] = events.FileUpdateEvent(original)
	evnts[1] = events.FileUpdateEvent(proxy)

	// Update language of original.
	// TODO: What about proxy !?
	if metadata.HasTranslation {
		original.Language = null.StringFrom(LANG_MULTI)
	} else {
		l := StdLang(metadata.Language)
		if l == LANG_UNKNOWN {
			log.Warnf("Unknown language in metadata %s", metadata.Language)
		}
		original.Language = null.StringFrom(l)
	}
	log.Infof("Updating original.Language to %s", original.Language.String)
	err = original.Update(exec, "language")
	if err != nil {
		return nil, errors.Wrap(err, "Save original to DB")
	}

	// Create content_unit (content_type, dates)
	isDerived := metadata.ArtifactType.Valid && metadata.ArtifactType.String != "main"
	ct := metadata.ContentType
	if isDerived {
		// TODO: verify user input. artifact_type should be either invalid, "main" or a known content_type
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
		props["label_id"] = metadata.LabelID.String
	}

	log.Infof("Creating content unit of type %s", ct)
	cu, err := CreateContentUnit(exec, ct, props)
	if err != nil {
		return nil, errors.Wrap(err, "Create content unit")
	}
	evnts = append(evnts, events.ContentUnitCreateEvent(cu))

	log.Infof("Describing content unit [%d]", cu.ID)
	err = DescribeContentUnit(exec, cu, metadata)
	if err != nil {
		log.Errorf("Error describing content unit: %s", err.Error())
	}

	// Add files to new unit
	log.Info("Adding files to unit")
	err = cu.AddFiles(exec, false, original, proxy)
	if err != nil {
		return nil, errors.Wrap(err, "Add files to unit")
	}

	// Add ancestor files to unit (not for derived units)
	if !isDerived {
		log.Info("Main unit, adding ancestors...")
		ancestors, err := FindFileAncestors(exec, original.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Find original's ancestors")
		}

		err = proxy.L.LoadParent(exec, true, proxy)
		if err != nil {
			return nil, errors.Wrap(err, "Load proxy's parent")
		}
		if proxy.R.Parent != nil {
			ancestors = append(ancestors, proxy.R.Parent)
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
	}

	// Associate unit with sources, tags, and persons
	if len(metadata.Sources) > 0 {
		log.Infof("Associating %d sources", len(metadata.Sources))
		sources, err := models.Sources(exec,
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Sources)...)).
			All()
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

		err = cu.AddSources(exec, false, sources...)
		if err != nil {
			return nil, errors.Wrap(err, "Associate sources")
		}
	}

	if len(metadata.Tags) > 0 {
		log.Infof("Associating %d tags", len(metadata.Tags))
		tags, err := models.Tags(exec,
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Tags)...)).
			All()
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
		err = cu.AddTags(exec, false, tags...)
		if err != nil {
			return nil, errors.Wrap(err, "Associate tags")
		}
	}

	// Handle persons ...
	if strings.ToLower(metadata.Lecturer) == P_RAV {
		log.Info("Associating unit to rav")
		cup := &models.ContentUnitsPerson{
			PersonID: PERSON_REGISTRY.ByPattern[P_RAV].ID,
			RoleID:   CONTENT_ROLE_TYPE_REGISTRY.ByName[CR_LECTURER].ID,
		}
		err = cu.AddContentUnitsPersons(exec, true, cup)
		if err != nil {
			return nil, errors.Wrap(err, "Associate persons")
		}
	} else {
		log.Infof("Unknown lecturer %s, skipping person association.", metadata.Lecturer)
	}

	// Get or create collection
	//var c *models.Collection
	if metadata.CollectionUID.Valid {
		log.Infof("Specific collection %s", metadata.CollectionUID.String)

		// find collection
		c, err := models.Collections(exec, qm.Where("uid = ?", metadata.CollectionUID.String)).One()
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

	if ct == CT_LESSON_PART || ct == CT_FULL_LESSON {
		log.Info("Lesson reconciliation")

		// Reconcile or create new
		// Reconciliation is done by looking up the operation chain of original to capture_stop.
		// There we have a property of saying the capture_id of the full lesson capture.
		captureStop, err := FindUpChainOperation(exec, original.ID, OP_CAPTURE_STOP)
		if err != nil {
			if ex, ok := err.(UpChainOperationNotFound); ok {
				log.Warnf("capture_stop operation not found for original: %s", ex.Error())
			} else {
				return nil, err
			}
		} else if captureStop.Properties.Valid {
			var oProps map[string]interface{}
			err = json.Unmarshal(captureStop.Properties.JSON, &oProps)
			if err != nil {
				return nil, errors.Wrap(err, "json Unmarshal")
			}

			captureID, ok := oProps["collection_uid"]
			if ok {
				log.Infof("Reconcile by capture_id %s", captureID)
				var cct string
				if metadata.WeekDate == nil {
					cct = CT_DAILY_LESSON
				} else {
					cct = CT_SPECIAL_LESSON
				}

				// Keep this property on the collection for other parts to find it
				props["capture_id"] = captureID
				if metadata.Number.Valid {
					props["number"] = metadata.Number.Int
				}
				delete(props, "duration")

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
				} else if ct == CT_FULL_LESSON {
					// Update collection properties to those of full lesson
					log.Info("Full lesson, overriding collection properties")
					if c.TypeID != CONTENT_TYPE_REGISTRY.ByName[cct].ID {
						log.Infof("Full lesson, content_type changed to %s", cct)
						c.TypeID = CONTENT_TYPE_REGISTRY.ByName[cct].ID
						err = c.Update(exec, "type_id")
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
					(!metadata.ArtifactType.Valid || metadata.ArtifactType.String == "main") {
					err := associateUnitToCollection(exec, cu, c, metadata)
					if err != nil {
						return nil, errors.Wrap(err, "associate content_unit to collection")
					}
					evnts = append(evnts, events.CollectionContentUnitsChangeEvent(c))
				}
			} else {
				log.Warnf("No collection_uid in capture_stop [%d] properties", captureStop.ID)
			}
		} else {
			log.Warnf("Invalid properties in capture_stop [%d]", captureStop.ID)
		}
	}

	// Associate unit and derived units
	// We take into account that a derived content unit arrives before it's source content unit.
	// Such cases are possible due to the studio operator actions sequence.
	err = original.L.LoadParent(exec, true, original)
	if err != nil {
		return nil, errors.Wrap(err, "Load original's parent")
	}

	if original.R.Parent == nil {
		log.Warn("We don't have original's parent file. Skipping derived units association.")
	} else {
		log.Info("Processing derived units associations")
		mainCUID := original.R.Parent.ContentUnitID
		if !metadata.ArtifactType.Valid ||
			metadata.ArtifactType.String == "main" {
			// main content unit
			log.Info("We're the main content unit")

			// We lookup original's siblings for derived content units that arrived before us.
			// We then associate them with us and remove their "unprocessed" mark.
			// Meaning, the presence of "artifact_type" property
			rows, err := queries.Raw(exec,
				`SELECT
				  cu.id,
				  cu.properties ->> 'artifact_type'
				FROM files f
				  INNER JOIN content_units cu ON f.content_unit_id = cu.id
				    AND cu.id != $1
				    AND cu.properties ? 'artifact_type'
				WHERE f.parent_id = $2`,
				original.ContentUnitID.Int64, original.ParentID.Int64).
				Query()
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

				_, err = queries.Raw(exec,
					`UPDATE content_units SET properties = properties - 'artifact_type' WHERE id = $1`,
					k).Exec()
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

			if mainCUID.Valid {
				// main content unit already exists
				log.Infof("Main content unit exists %d", mainCUID.Int64)
				cud := &models.ContentUnitDerivation{
					SourceID: mainCUID.Int64,
					Name:     metadata.ArtifactType.String,
				}
				err = cu.AddDerivedContentUnitDerivations(exec, true, cud)
				if err != nil {
					return nil, errors.Wrap(err, "Save source unit in DB")
				}
				evnts = append(evnts, events.ContentUnitDerivativesChangeEvent(cu))
			} else {
				// save artifact type for later use (when main unit appears)
				log.Info("Main content unit not found, saving artifact_type property")
				err = UpdateContentUnitProperties(exec, cu, map[string]interface{}{
					"artifact_type": metadata.ArtifactType.String,
				})
				if err != nil {
					return nil, err
				}
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

	switch CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name {
	case CT_FULL_LESSON:
		if c.TypeID == CONTENT_TYPE_REGISTRY.ByName[CT_DAILY_LESSON].ID ||
			c.TypeID == CONTENT_TYPE_REGISTRY.ByName[CT_SPECIAL_LESSON].ID {
			ccu.Name = "full"
		} else if metadata.Number.Valid {
			ccu.Name = strconv.Itoa(metadata.Number.Int)
		}
		break
	case CT_LESSON_PART:
		if metadata.Part.Valid {
			ccu.Name = strconv.Itoa(metadata.Part.Int)
		}
		break
	case CT_VIDEO_PROGRAM_CHAPTER:
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
			if idx < len(MISC_EVENT_PART_TYPES) {
				ccu.Name = MISC_EVENT_PART_TYPES[idx] + ccu.Name
			} else {
				log.Warnf("Unknown event part type: %d", metadata.PartType.Int)
			}
		}
		break
	}

	// Make this new unit the last one in this collection
	var err error
	ccu.Position, err = GetNextPositionInCollection(exec, c.ID)
	if err != nil {
		return errors.Wrap(err, "Get last position in collection")
	}

	log.Infof("Association name: %s", ccu.Name)
	err = c.AddCollectionsContentUnits(exec, true, ccu)
	if err != nil {
		return errors.Wrap(err, "Save collection and content unit association in DB")
	}

	return nil
}
