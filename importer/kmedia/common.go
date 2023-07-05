package kmedia

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	qm4 "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	CATALOGS_SOURCES_MAPPINGS_FILE = "importer/kmedia/data/Catalogs Sources Mappings - final.csv"
	CATALOGS_TAGS_MAPPINGS_FILE    = "importer/kmedia/data/catalogs_tags.csv"
)

var (
	mdb                     *sql.DB
	kmdb                    *sql.DB
	stats                   *ImportStatistics
	serverUrls              map[string]string
	catalogsSourcesMappings map[int]*models.Source
	catalogsTagsMappings    map[int]*models.Tag
)

func Init() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia migration")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	//defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmdb, err = sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmdb.Ping())
	//defer kmdb.Close()

	log.Info("Initializing static data from MDB")
	utils.Must(common.InitTypeRegistries(mdb))

	log.Info("Initializing static data from Kmedia")
	kmediaLessonCT, err = kmodels.ContentTypes(kmdb, qm.Where("name = ?", "Lesson")).One()
	utils.Must(err)
	serverUrls, err = initServers(kmdb)
	utils.Must(err)

	log.Info("Initializing catalogs sources mappings")
	catalogsSourcesMappings, err = initCatalogSourcesMappings()
	utils.Must(err)
	log.Infof("Got %d mappings", len(catalogsSourcesMappings))

	log.Info("Initializing catalogs tags mappings")
	catalogsTagsMappings, err = initCatalogTagsMappings()
	utils.Must(err)
	log.Infof("Got %d mappings", len(catalogsTagsMappings))

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
	utils.Must(kmdb.Close())
}

func importContainerWOCollectionNewCU(exec boil.Executor, container *kmodels.Container, cuType string) (*models.ContentUnit, error) {
	err := container.L.LoadFileAssets(kmdb, true, container)
	if err != nil {
		return nil, errors.Wrapf(err, "Load kmedia file assets %d", container.ID)
	}

	// Create import operation
	operation, err := api.CreateOperation(exec, common.OP_IMPORT_KMEDIA,
		api.Operation{WorkflowID: strconv.Itoa(container.ID)}, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Create operation %d", container.ID)
	}
	stats.OperationsCreated.Inc(1)

	// import container
	unit, err := importContainer(exec, container, nil, cuType, "", 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Import container %d", container.ID)
	}

	// import container files
	var file *models.File
	for _, fileAsset := range container.R.FileAssets {
		log.Infof("Processing file_asset %d", fileAsset.ID)
		stats.FileAssetsProcessed.Inc(1)

		// Create or update MDB file
		file, err = importFileAsset(exec, fileAsset, unit, operation)
		if err != nil {
			return nil, errors.Wrapf(err, "Import file_asset %d", fileAsset.ID)
		}
		if file != nil && file.Published {
			unit.Published = true
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Import container files %d", container.ID)
	}

	if unit.Published {
		_, err = unit.Update(exec, boil.Whitelist("published"))
		if err != nil {
			return nil, errors.Wrapf(err, "Update unit published column %d", container.ID)
		}
	}

	return unit, nil
}

func importContainerWCollection(exec boil.Executor, container *kmodels.Container, collection *models.Collection, cuType string) error {
	stats.ContainersProcessed.Inc(1)

	unit, err := models.ContentUnits(qm4.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One(mdb)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("New CU %d %s", container.ID, container.Name.String)
			return importContainerWCollectionNewCU(exec, container, collection, cuType)
		}
		return errors.Wrapf(err, "Lookup content unit kmid %d", container.ID)
	}

	log.Infof("CU exists [%d] container: %s %d", unit.ID, container.Name.String, container.ID)
	_, err = importContainer(exec, container, collection, cuType,
		strconv.Itoa(container.Position.Int), container.Position.Int)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	if cuType != common.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID].Name {
		log.Infof("Overriding CU Type to %s", cuType)
		unit.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[cuType].ID
		_, err = unit.Update(exec, boil.Whitelist("type_id"))
		if err != nil {
			return errors.Wrapf(err, "Update CU type %d", unit.ID)
		}
	}

	return nil
}

func importContainerWCollectionNewCU(exec boil.Executor, container *kmodels.Container, collection *models.Collection, cuType string) error {
	err := container.L.LoadFileAssets(kmdb, true, container)
	if err != nil {
		return errors.Wrapf(err, "Load kmedia file assets %d", container.ID)
	}

	// Create import operation
	operation, err := api.CreateOperation(exec, common.OP_IMPORT_KMEDIA,
		api.Operation{WorkflowID: strconv.Itoa(container.ID)}, nil)
	if err != nil {
		return errors.Wrapf(err, "Create operation %d", container.ID)
	}
	stats.OperationsCreated.Inc(1)

	// import container
	unit, err := importContainer(exec, container, collection, cuType,
		strconv.Itoa(container.Position.Int), container.Position.Int)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	// import container files
	var file *models.File
	for _, fileAsset := range container.R.FileAssets {
		log.Infof("Processing file_asset %d", fileAsset.ID)
		stats.FileAssetsProcessed.Inc(1)

		// Create or update MDB file
		file, err = importFileAsset(exec, fileAsset, unit, operation)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			break
		}
		if file != nil && file.Published {
			unit.Published = true
		}
	}
	if err != nil {
		return errors.Wrapf(err, "Import container files %d", container.ID)
	}

	if unit.Published {
		_, err = unit.Update(exec, boil.Whitelist("published"))
		if err != nil {
			return errors.Wrapf(err, "Update unit published column %d", container.ID)
		}
	}

	return nil
}

func importContainer(exec boil.Executor,
	container *kmodels.Container,
	collection *models.Collection,
	contentType string,
	ccuName string,
	ccuPosition int,
) (*models.ContentUnit, error) {

	// Get or create content unit by kmedia_id
	unit, err := models.ContentUnits(qm4.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One(exec)
	if err == nil {
		stats.ContentUnitsUpdated.Inc(1)
		if contentType != "" && contentType != common.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID].Name {
			log.Warnf("Different CU type %d %s != %s", unit.ID, common.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID].Name, contentType)
		}
	} else {
		if err == sql.ErrNoRows {
			// Create new content unit
			unit, err = api.CreateContentUnit(exec, contentType, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "Insert unit, container_id [%d]", container.ID)
			}
			stats.ContentUnitsCreated.Inc(1)
		} else {
			return nil, errors.Wrapf(err, "Lookup unit, container_id [%d]", container.ID)
		}
	}

	// Secure
	unit.Secure = mapSecure(container.Secure)
	_, err = unit.Update(exec, boil.Whitelist("secure"))
	if err != nil {
		return nil, errors.Wrapf(err, "Update secure, unit [%d]", unit.ID)
	}

	// Properties
	props := make(map[string]interface{})
	if unit.Properties.Valid {
		unit.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = container.ID
	if container.LangID.Valid {
		props["original_language"] = common.StdLang(container.LangID.String)
	}
	if container.Filmdate.Valid {
		props["film_date"] = container.Filmdate.Time.Format("2006-01-02")
	}
	if container.PlaytimeSecs.Valid {
		props["duration"] = container.PlaytimeSecs.Int
	}
	err = api.UpdateContentUnitProperties(exec, unit, props)
	if err != nil {
		return nil, errors.Wrapf(err, "Update properties, unit [%d]", unit.ID)
	}

	// TODO: what to do with censor workflow information ?

	// I18n
	descriptions, err := container.ContainerDescriptions(kmdb).All()
	if err != nil {
		return nil, errors.Wrapf(err, "Lookup container descriptions, container_id [%d]", container.ID)
	}

	hasI18n := false
	for _, d := range descriptions {
		if (d.ContainerDesc.Valid && d.ContainerDesc.String != "") ||
			(d.Descr.Valid && d.Descr.String != "") {
			hasI18n = true
			cui18n := models.ContentUnitI18n{
				ContentUnitID: unit.ID,
				Language:      common.LANG_MAP[d.LangID.String],
				Name:          null.NewString(d.ContainerDesc.String, d.ContainerDesc.Valid),
				Description:   null.NewString(d.Descr.String, d.Descr.Valid),
			}
			err = cui18n.Upsert(exec,
				true,
				[]string{"content_unit_id", "language"},
				boil.Whitelist("name", "description"),
				boil.Infer())
			if err != nil {
				return nil, errors.Wrapf(err, "Upsert unit i18n, unit [%d]", unit.ID)
			}
		}
	}

	// no i18n - use container name
	if !hasI18n && container.Name.Valid {
		for _, lang := range []string{common.LANG_ENGLISH, common.LANG_HEBREW, common.LANG_RUSSIAN, common.LANG_SPANISH} {
			cui18n := models.ContentUnitI18n{
				ContentUnitID: unit.ID,
				Language:      lang,
				Name:          null.NewString(container.Name.String, container.Name.Valid),
			}
			err = cui18n.Upsert(exec,
				true,
				[]string{"content_unit_id", "language"},
				boil.Whitelist("name", "description"),
				boil.Infer())
			if err != nil {
				return nil, errors.Wrapf(err, "Upsert unit i18n, unit [%d]", unit.ID)
			}
		}
	}

	if collection != nil {
		err := createOrUpdateCCU(exec, unit, models.CollectionsContentUnit{
			CollectionID:  collection.ID,
			ContentUnitID: unit.ID,
			Name:          ccuName,
			Position:      ccuPosition,
		})
		if err != nil {
			return nil, err
		}
	}

	// Associate sources & tags
	// we combine existing mappings with catalogs mappings
	err = container.L.LoadCatalogs(kmdb, true, container)
	if err != nil {
		return nil, errors.Wrapf(err, "Load catalogs, container [%d]", container.ID)
	}

	err = unit.L.LoadSources(exec, true, unit, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Load CU sources %d", unit.ID)
	}

	err = unit.L.LoadTags(exec, true, unit, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Load CU tags %d", unit.ID)
	}

	srcMap := make(map[int64]*models.Source)
	tagMap := make(map[int64]*models.Tag)
	for _, x := range container.R.Catalogs {
		if s, ok := catalogsSourcesMappings[x.ID]; ok {
			srcMap[s.ID] = s
		} else if t, ok := catalogsTagsMappings[x.ID]; ok {
			tagMap[t.ID] = t
		} else {
			stats.UnkownCatalogs.Inc(fmt.Sprintf("%s [%d]", x.Name, x.ID), 1)
		}
	}
	for _, s := range unit.R.Sources {
		srcMap[s.ID] = s
	}
	for _, t := range unit.R.Tags {
		tagMap[t.ID] = t
	}

	// Set sources
	src := make([]*models.Source, len(srcMap))
	i := 0
	for _, v := range srcMap {
		src[i] = v
		i++
	}
	err = unit.SetSources(exec, false, src...)
	if err != nil {
		return nil, errors.Wrapf(err, "Set sources, unit [%d]", unit.ID)
	}

	// Set tags
	tags := make([]*models.Tag, len(tagMap))
	i = 0
	for _, v := range tagMap {
		tags[i] = v
		i++
	}
	err = unit.SetTags(exec, false, tags...)
	if err != nil {
		return nil, errors.Wrapf(err, "Set tags, unit [%d]", unit.ID)
	}

	// person
	if container.LecturerID.Valid {
		person := mapPerson(container.LecturerID.Int)
		if person != nil {
			cup := models.ContentUnitsPerson{
				PersonID: person.ID,
				RoleID:   1, // lecturer
			}
			err := unit.AddContentUnitsPersons(exec, true, &cup)
			if err != nil {
				return nil, errors.Wrapf(err, "Add person, unit [%d]", unit.ID)
			}
		}
	}

	return unit, nil
}

func createOrUpdateCCU(exec boil.Executor, unit *models.ContentUnit, ccu models.CollectionsContentUnit) error {
	x, err := models.FindCollectionsContentUnit(exec, ccu.CollectionID, ccu.ContentUnitID)
	if err != nil {
		if err == sql.ErrNoRows {
			// create
			err = unit.AddCollectionsContentUnits(exec, true, &models.CollectionsContentUnit{
				CollectionID:  ccu.CollectionID,
				ContentUnitID: ccu.ContentUnitID,
				Name:          ccu.Name,
				Position:      ccu.Position,
			})
			if err != nil {
				return errors.Wrapf(err, "Create CCU [c,cu]=[%d,%d]", ccu.CollectionID, ccu.ContentUnitID)
			}
		} else {
			return errors.Wrapf(err, "Find CCU [c,cu]=[%d,%d]", ccu.CollectionID, ccu.ContentUnitID)
		}
	} else {
		// update if name or position changed
		if ccu.Name != x.Name || ccu.Position != x.Position {
			x.Name = ccu.Name
			x.Position = ccu.Position
			_, err = x.Update(exec, boil.Whitelist("name", "position"))
			if err != nil {
				return errors.Wrapf(err, "Update CCU [c,cu]=[%d,%d]", ccu.CollectionID, ccu.ContentUnitID)
			}
		}
	}

	return nil
}

func importFileAsset(exec boil.Executor, fileAsset *kmodels.FileAsset, unit *models.ContentUnit,
	operation *models.Operation) (*models.File, error) {

	// Get or Create MDB file by SHA1
	var hash string
	if fileAsset.Sha1.Valid {
		hash = fileAsset.Sha1.String

		// This sha1 is for empty files, i.e. physical size = 0
		// We skip them for now.
		// Hoping someone will find these files and get their real sha1...
		if hash == "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
			return nil, nil
		}
	} else {
		hash = fmt.Sprintf("%x", sha1.Sum([]byte(strconv.Itoa(fileAsset.ID))))
		stats.FileAssetsMissingSHA1.Inc(1)
	}

	//file, _, err := api.FindFileBySHA1(exec, hash)
	file, hashB, err := api.FindFileBySHA1(exec, hash)
	if err == nil {
		stats.FilesUpdated.Inc(1)
	} else {
		if _, ok := err.(api.FileNotFound); ok {
			shouldCreate := true

			// For unknown file assets with valid sha1 do second lookup before we create a new file.
			// This time with the fake sha1, if exists we replace fake hash with valid hash.
			// Note: this paragraph should not be executed on first import.
			if fileAsset.Sha1.Valid {
				file, _, err = api.FindFileBySHA1(exec,
					fmt.Sprintf("%x", sha1.Sum([]byte(strconv.Itoa(fileAsset.ID)))))
				if err == nil {
					file.Sha1 = null.BytesFrom(hashB)
					shouldCreate = false
					stats.FilesUpdated.Inc(1)
				} else {
					if _, ok := err.(api.FileNotFound); !ok {
						return nil, errors.Wrapf(err, "Second file lookup, file_asset [%d]", fileAsset.ID)
					}
				}
			}

			if shouldCreate {
				// Create new file
				f := api.File{
					FileName:  fileAsset.Name.String,
					Sha1:      hash,
					Size:      int64(fileAsset.Size.Int),
					CreatedAt: &api.Timestamp{Time: fileAsset.Date.Time},
				}
				file, err = api.CreateFile(exec, nil, f, nil)
				if err != nil {
					return nil, errors.Wrapf(err, "Create file")
				}
				stats.FilesCreated.Inc(1)
			}
		} else {
			return nil, errors.Wrapf(err, "Lookup file %s", hash)
		}
	}

	// Media types
	if fileAsset.AssetType.Valid {
		if mt, ok := common.MEDIA_TYPE_REGISTRY.ByExtension[strings.ToLower(fileAsset.AssetType.String)]; ok {
			file.Type = mt.Type
			file.SubType = mt.SubType
			file.MimeType = null.NewString(mt.MimeType, mt.MimeType != "")
		} else {
			stats.FileAssetsWInvalidMT.Inc(1)
		}
	} else {
		stats.FileAssetsMissingType.Inc(1)
	}

	// Language
	if fileAsset.LangID.Valid {
		l := common.LANG_MAP[fileAsset.LangID.String]
		file.Language = null.NewString(l, l != "")
	}

	// Secure
	if fileAsset.Secure.Valid {
		file.Secure = mapSecure(fileAsset.Secure.Int)
	}
	file.Published = file.Secure == 0

	// Properties
	props := make(map[string]interface{})
	if file.Properties.Valid {
		file.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = fileAsset.ID
	props["url"] = serverUrls[fileAsset.ServernameID.String] + "/" + file.Name
	if fileAsset.PlaytimeSecs.Valid {
		props["duration"] = fileAsset.PlaytimeSecs.Int
	}
	p, _ := json.Marshal(props)
	file.Properties = null.JSONFrom(p)

	_, err = file.Update(exec, boil.Infer())
	if err != nil {
		return nil, errors.Wrapf(err, "Update file [%d]", file.ID)
	}

	// i18n
	// We don't take anything from file_asset_descriptions as it`s mostly junk

	// Associate files with content unit
	if file.ContentUnitID.Valid && file.ContentUnitID.Int64 != unit.ID {
		log.Warnf("Changing file's unit association from %d to %d", file.ContentUnitID.Int64, unit.ID)
	}
	err = unit.AddFiles(exec, false, file)
	if err != nil {
		return nil, errors.Wrapf(err, "Associate file [%d] to unit [%d]", file.ID, unit.ID)
	}

	// Associate files with operation

	// We use a raw query here to do nothing on conflicts
	// These conflicts happen when different file_assets in the same lesson have identical SHA1
	_, err = queries.Raw(
		`INSERT INTO files_operations (file_id, operation_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		file.ID, operation.ID).Exec(exec)

	return file, err
}

func initServers(exec boil.Executor) (map[string]string, error) {
	servers, err := kmodels.Servers(exec).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load kmedia servers")
	}

	serverUrls := make(map[string]string)
	for _, s := range servers {
		serverUrls[s.Servername] = s.Httpurl.String
	}
	return serverUrls, nil
}

func initCatalogSourcesMappings() (map[int]*models.Source, error) {
	// read mappings file
	records, err := utils.ReadCSV(CATALOGS_SOURCES_MAPPINGS_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read catalogs sources mappings")
	}
	log.Infof("Catalogs Sources Mappings has %d rows", len(records))

	// read MDB sources
	rows, err := queries.Raw(`WITH RECURSIVE rec_sources AS (
  SELECT
    s.id,
    concat(a.code, '/', s.name) path
  FROM sources s INNER JOIN authors_sources x ON s.id = x.source_id
    INNER JOIN authors a ON x.author_id = a.id
  WHERE s.parent_id IS NULL
  UNION
  SELECT
    s.id,
    concat(rs.path, '/', s.name)
  FROM sources s INNER JOIN rec_sources rs ON s.parent_id = rs.id
)
SELECT *
FROM rec_sources;`).Query(mdb)
	if err != nil {
		return nil, errors.Wrap(err, "Read MDB sources")
	}

	defer rows.Close()

	tmp := make(map[string]*models.Source)
	for rows.Next() {
		var id int64
		var path string
		err := rows.Scan(&id, &path)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}
		tmp[path] = &models.Source{ID: id}
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Iterating MDB sources")
	}
	log.Infof("%d MDB Sources", len(tmp))

	mappings := make(map[int]*models.Source)
	for i, r := range records[1:] {
		catalogID, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, errors.Wrapf(err, "Bad catalog_id, row [%d]", i)
		}

		sourcePath := strings.TrimSpace(r[2])
		s, ok := tmp[sourcePath]
		if !ok {
			log.Warnf("Unknown source, path=%s", sourcePath)
		}
		mappings[catalogID] = s
	}

	return mappings, nil
}

func initCatalogTagsMappings() (map[int]*models.Tag, error) {
	// Read mappings file
	records, err := utils.ReadCSV(CATALOGS_TAGS_MAPPINGS_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read catalogs tags mappings")
	}
	log.Infof("Catalogs Tags Mappings has %d rows", len(records))

	// Read all tags from MDB
	tags, err := models.Tags().All(mdb)
	if err != nil {
		return nil, errors.Wrap(err, "Fetch tags from MDB")
	}
	tmp := make(map[int64]*models.Tag, len(tags))
	for _, t := range tags {
		tmp[t.ID] = t
	}

	// Create mappings
	mappings := make(map[int]*models.Tag, len(records)-1)
	for i, r := range records[1:] {
		catalogID, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, errors.Wrapf(err, "Bad catalog_id, row [%d]", i)
		}
		tagID, err := strconv.Atoi(r[1])
		if err != nil {
			return nil, errors.Wrapf(err, "Bad tag_id, row [%d]", i)
		}

		mappings[catalogID] = tmp[int64(tagID)]
	}

	return mappings, nil
}

func mapSecure(kmVal int) int16 {
	if kmVal == 0 {
		return common.SEC_PUBLIC
	} else if kmVal < 4 {
		return common.SEC_SENSITIVE
	}
	return common.SEC_PRIVATE
}

func mapPerson(kmID int) *models.Person {
	switch kmID {
	case 1:
		return common.PERSON_REGISTRY.ByPattern["rav"]
	case 8:
		return common.PERSON_REGISTRY.ByPattern["rb"]
	default:
		return nil
	}
}

func loadContainersInCatalogsAndCUs(catalogIDs ...int) (map[int]*kmodels.Container, map[int]*models.ContentUnit, error) {
	q := `WITH RECURSIVE rec_catalogs AS (
  SELECT c.id
  FROM catalogs c
  WHERE id IN (%s)
  UNION
  SELECT c.id
  FROM catalogs c INNER JOIN rec_catalogs rc ON c.parent_id = rc.id
)
SELECT
  DISTINCT cc.container_id
FROM rec_catalogs rc INNER JOIN catalogs_containers cc ON rc.id = cc.catalog_id`

	catIDs := make([]string, len(catalogIDs))
	for i := range catalogIDs {
		catIDs[i] = strconv.Itoa(catalogIDs[i])
	}

	rows, err := queries.Raw(fmt.Sprintf(q, strings.Join(catIDs, ","))).Query(kmdb)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load containers")
	}
	defer rows.Close()

	cnIDs := make([]int, 0)
	for rows.Next() {
		var x int
		if err := rows.Scan(&x); err != nil {
			return nil, nil, errors.Wrap(err, "rows.Scan")
		} else {
			cnIDs = append(cnIDs, x)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.Wrap(err, "rows.Err")
	}
	log.Infof("len(cnIDS) = %d", len(cnIDs))

	pageSize := 2500
	page := 0
	cnMap := make(map[int]*kmodels.Container, len(cnIDs))
	cuMap := make(map[int]*models.ContentUnit, len(cnIDs))
	for page*pageSize < len(cnIDs) {
		s := page * pageSize
		e := utils.Min(len(cnIDs), s+pageSize)
		ids := utils.ConvertArgsInt(cnIDs[s:e]...)

		cns, err := kmodels.Containers(kmdb,
			qm.WhereIn("id in ?", ids...),
			qm.Load("Catalogs")).
			All()
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Load containers page %d", page)
		}
		for i := range cns {
			cn := cns[i]
			cnMap[cn.ID] = cn
		}

		cus, err := models.ContentUnits(
			qm4.WhereIn("(properties->>'kmedia_id')::int in ?", ids...)).
			All(mdb)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Load content units page %d", page)
		}
		for i := range cus {
			cu := cus[i]
			var props map[string]interface{}
			if err := json.Unmarshal(cu.Properties.JSON, &props); err != nil {
				return nil, nil, errors.Wrapf(err, "json.Unmarshal CU properties %d", cu.ID)
			}
			cuMap[int(props["kmedia_id"].(float64))] = cu
		}

		page++
	}
	log.Infof("len(cnMap) = %d", len(cnMap))
	log.Infof("len(cuMap) = %d", len(cuMap))

	return cnMap, cuMap, nil
}

func loadContainersByTypeAndCUs(typeID int) (map[int]*kmodels.Container, map[int]*models.ContentUnit, error) {
	containers, err := kmodels.Containers(kmdb,
		qm.Where("content_type_id = ?", typeID),
		qm.Load("Catalogs")).
		All()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load containers")
	}

	cnMap := make(map[int]*kmodels.Container)
	for i := range containers {
		cn := containers[i]
		cnMap[cn.ID] = cn
	}

	cnIDs := make([]interface{}, len(cnMap))
	i := 0
	for k := range cnMap {
		cnIDs[i] = k
		i++
	}

	pageSize := 2500
	page := 0
	cuMap := make(map[int]*models.ContentUnit)
	for page*pageSize < len(cnIDs) {
		s := page * pageSize
		e := utils.Min(len(cnIDs), s+pageSize)

		cus, err := models.ContentUnits(
			qm4.WhereIn("(properties->>'kmedia_id')::int in ?", cnIDs[s:e]...)).
			All(mdb)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Load content units page %d", page)
		}
		for i := range cus {
			cu := cus[i]
			var props map[string]interface{}
			if err := json.Unmarshal(cu.Properties.JSON, &props); err != nil {
				return nil, nil, errors.Wrapf(err, "json.Unmarshal CU properties %d", cu.ID)
			}
			cuMap[int(props["kmedia_id"].(float64))] = cu
		}

		page++
	}
	log.Infof("len(cnMap) = %d", len(cnMap))
	log.Infof("len(cuMap) = %d", len(cuMap))

	return cnMap, cuMap, nil
}
