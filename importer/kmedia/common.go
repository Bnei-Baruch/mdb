package kmedia

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
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
	utils.Must(api.InitTypeRegistries(mdb))

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

func importContainer(exec boil.Executor,
	container *kmodels.Container,
	collection *models.Collection,
	contentType string,
	ccuName string,
	ccuPosition int,
) (*models.ContentUnit, error) {

	// Get or create content unit by kmedia_id
	unit, err := models.ContentUnits(exec, qm.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One()
	if err == nil {
		stats.ContentUnitsUpdated.Inc(1)
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
	err = unit.Update(exec, "secure")
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
		props["original_language"] = api.StdLang(container.LangID.String)
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
	for _, d := range descriptions {
		if (d.ContainerDesc.Valid && d.ContainerDesc.String != "") ||
			(d.Descr.Valid && d.Descr.String != "") {
			cui18n := models.ContentUnitI18n{
				ContentUnitID: unit.ID,
				Language:      api.LANG_MAP[d.LangID.String],
				Name:          d.ContainerDesc,
				Description:   d.Descr,
			}
			err = cui18n.Upsert(exec,
				true,
				[]string{"content_unit_id", "language"},
				[]string{"name", "description"})
			if err != nil {
				return nil, errors.Wrapf(err, "Upsert unit i18n, unit [%d]", unit.ID)
			}
		}
	}

	// Associate content_unit with collection , name = position
	err = unit.L.LoadCollectionsContentUnits(exec, true, unit)
	if err != nil {
		return nil, errors.Wrapf(err, "Fetch unit collections, unit [%d]", unit.ID)
	}

	// lookup existing association
	var ccu *models.CollectionsContentUnit
	if unit.R.CollectionsContentUnits != nil {
		for _, x := range unit.R.CollectionsContentUnits {
			if x.CollectionID == collection.ID {
				ccu = x
				break
			}
		}
	}

	if ccu == nil {
		// Create
		err = unit.AddCollectionsContentUnits(exec, true, &models.CollectionsContentUnit{
			CollectionID:  collection.ID,
			ContentUnitID: unit.ID,
			Name:          ccuName,
			Position:      ccuPosition,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "Add unit collections, unit [%d]", unit.ID)
		}
	} else {
		// update if name changed
		if ccu.Name != ccuName || ccu.Position != ccuPosition {
			ccu.Name = ccuName
			ccu.Position = ccuPosition
			err = ccu.Update(exec, "name", "position")
			if err != nil {
				return nil, errors.Wrapf(err,
					"Update unit collection association, unit [%d], collection [%d]",
					unit.ID, collection.ID)
			}
		}
	}

	// Associate sources & tags
	err = container.L.LoadCatalogs(kmdb, true, container)
	if err != nil {
		return nil, errors.Wrapf(err, "Load catalogs, container [%d]", container.ID)
	}

	// Dedup list of matches
	unqSrc := make(map[*models.Source]bool, 0)
	unqTags := make(map[*models.Tag]bool, 0)
	for _, x := range container.R.Catalogs {
		if s, ok := catalogsSourcesMappings[x.ID]; ok {
			unqSrc[s] = true
		} else if t, ok := catalogsTagsMappings[x.ID]; ok {
			unqTags[t] = true
		} else {
			stats.UnkownCatalogs.Inc(fmt.Sprintf("%s [%d]", x.Name, x.ID), 1)
		}
	}

	// Set sources
	src := make([]*models.Source, len(unqSrc))
	i := 0
	for k := range unqSrc {
		src[i] = k
		i++
	}
	err = unit.SetSources(exec, false, src...)
	if err != nil {
		return nil, errors.Wrapf(err, "Set sources, unit [%d]", unit.ID)
	}

	// Set tags
	tags := make([]*models.Tag, len(unqTags))
	i = 0
	for k := range unqTags {
		tags[i] = k
		i++
	}
	err = unit.SetTags(exec, false, tags...)
	if err != nil {
		return nil, errors.Wrapf(err, "Set tags, unit [%d]", unit.ID)
	}

	return unit, nil
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
		if mt, ok := api.MEDIA_TYPE_REGISTRY.ByExtension[strings.ToLower(fileAsset.AssetType.String)]; ok {
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
		l := api.LANG_MAP[fileAsset.LangID.String]
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

	err = file.Update(exec)
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
	_, err = queries.Raw(exec,
		`INSERT INTO files_operations (file_id, operation_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		file.ID, operation.ID).Exec()

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
	rows, err := queries.Raw(mdb, `WITH RECURSIVE rec_sources AS (
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
FROM rec_sources;`).Query()
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
	tags, err := models.Tags(mdb).All()
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
		return api.SEC_PUBLIC
	} else if kmVal < 4 {
		return api.SEC_SENSITIVE
	}
	return api.SEC_PRIVATE
}
