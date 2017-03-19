package importer

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
)

var (
	kmediaLessonCT *kmodels.ContentType
	langID2Locale  map[string]string
	serverUrls     map[string]string
	stats          *ImportStatistics
)

type MediaType struct {
	Extension string
	Type      string
	SubType   string
	MimeType  string
}

// select asset_type, count(*) from file_assets group by asset_type order by count(*) desc;
var MEDIA_TYPES = map[string]MediaType{
	"mp4":  {Extension: "mp4", Type: "video", SubType: "", MimeType: "video/mp4"},
	"wmv":  {Extension: "wmv", Type: "video", SubType: "", MimeType: "video/x-ms-wmv"},
	"flv":  {Extension: "flv", Type: "video", SubType: "", MimeType: "video/x-flv"},
	"mov":  {Extension: "mov", Type: "video", SubType: "", MimeType: "video/quicktime"},
	"asf":  {Extension: "asf", Type: "video", SubType: "", MimeType: "video/x-ms-asf"},
	"mpg":  {Extension: "mpg", Type: "video", SubType: "", MimeType: "video/mpeg"},
	"avi":  {Extension: "avi", Type: "video", SubType: "", MimeType: "video/x-msvideo"},
	"mp3":  {Extension: "mp3", Type: "audio", SubType: "", MimeType: "audio/mpeg"},
	"wma":  {Extension: "wma", Type: "audio", SubType: "", MimeType: "audio/x-ms-wma"},
	"mid":  {Extension: "mid", Type: "audio", SubType: "", MimeType: "audio/midi"},
	"wav":  {Extension: "wav", Type: "audio", SubType: "", MimeType: "audio/x-wav"},
	"aac":  {Extension: "aac", Type: "audio", SubType: "", MimeType: "audio/aac"},
	"jpg":  {Extension: "jpg", Type: "image", SubType: "", MimeType: "image/jpeg"},
	"gif":  {Extension: "gif", Type: "image", SubType: "", MimeType: "image/gif"},
	"bmp":  {Extension: "bmp", Type: "image", SubType: "", MimeType: "image/bmp"},
	"tif":  {Extension: "tif", Type: "image", SubType: "", MimeType: "image/tiff"},
	"zip":  {Extension: "zip", Type: "image", SubType: "", MimeType: "application/zip"},
	"7z":   {Extension: "7z", Type: "image", SubType: "", MimeType: "application/x-7z-compressed"},
	"rar":  {Extension: "rar", Type: "image", SubType: "", MimeType: "application/x-rar-compressed"},
	"sfk":  {Extension: "sfk", Type: "image", SubType: "", MimeType: ""},
	"doc":  {Extension: "doc", Type: "text", SubType: "", MimeType: "application/msword"},
	"docx": {Extension: "docx", Type: "text", SubType: "", MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	"htm":  {Extension: "htm", Type: "text", SubType: "", MimeType: "text/html"},
	"html": {Extension: "htm", Type: "text", SubType: "", MimeType: "text/html"},
	"pdf":  {Extension: "pdf", Type: "text", SubType: "", MimeType: "application/pdf"},
	"epub": {Extension: "epub", Type: "text", SubType: "", MimeType: "application/epub+zip"},
	"rtf":  {Extension: "rtf", Type: "text", SubType: "", MimeType: "text/rtf"},
	"txt":  {Extension: "txt", Type: "text", SubType: "", MimeType: "text/plain"},
	"fb2":  {Extension: "fb2", Type: "text", SubType: "", MimeType: "text/xml"},
	"rb":  {Extension: "rb", Type: "text", SubType: "", MimeType: "application/x-rocketbook"},
	"xls":  {Extension: "xls", Type: "sheet", SubType: "", MimeType: "application/vnd.ms-excel"},
	"swf":  {Extension: "swf", Type: "banner", SubType: "", MimeType: "application/x-shockwave-flash"},
	"ppt":  {Extension: "ppt", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
	"pptx": {Extension: "pptx", Type: "presentation", SubType: "", MimeType: "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	"pps":  {Extension: "pps", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
}

type ImportStatistics struct {
	LessonsProcessed       int
	ValidLessons           int
	InvalidLessons         int
	ContainersProcessed    int
	ContainersVisited      int
	ContainersWithFiles    int
	ContainersWithoutFiles int
	FileAssetsProcessed    int
	FileAssetsMissingSHA1  int
	FileAssetsWInvalidMT   int
	FileAssetsMissingType  int

	OperationsCreated   int
	CollectionsCreated  int
	CollectionsUpdated  int
	ContentUnitsCreated int
	ContentUnitsUpdated int
	FilesCreated        int
	FilesUpdated        int

	TxCommitted  int
	TxRolledBack int
}

func (s *ImportStatistics) dump() {
	fmt.Println("Here comes import statistics:")

	fmt.Println("Kmedia:")
	fmt.Printf("ValidLessons            		%d\n", s.ValidLessons)
	fmt.Printf("InvalidLessons          		%d\n", s.InvalidLessons)
	fmt.Printf("LessonsProcessed        		%d\n", s.LessonsProcessed)
	fmt.Printf("ContainersVisited       		%d\n", s.ContainersVisited)
	fmt.Printf("ContainersWithFiles     		%d\n", s.ContainersWithFiles)
	fmt.Printf("ContainersWithoutFiles  		%d\n", s.ContainersWithoutFiles)
	fmt.Printf("ContainersProcessed     		%d\n", s.ContainersProcessed)
	fmt.Printf("FileAssetsProcessed     		%d\n", s.FileAssetsProcessed)
	fmt.Printf("FileAssetsMissingSHA1   		%d\n", s.FileAssetsMissingSHA1)
	fmt.Printf("FileAssetsWInvalidMT    		%d\n", s.FileAssetsWInvalidMT)
	fmt.Printf("FileAssetsMissingType   		%d\n", s.FileAssetsMissingType)

	fmt.Println("MDB:")
	fmt.Printf("OperationsCreated       		%d\n", s.OperationsCreated)
	fmt.Printf("CollectionsCreated      		%d\n", s.CollectionsCreated)
	fmt.Printf("CollectionsUpdated      		%d\n", s.CollectionsUpdated)
	fmt.Printf("ContentUnitsCreated     		%d\n", s.ContentUnitsCreated)
	fmt.Printf("ContentUnitsUpdated     		%d\n", s.ContentUnitsUpdated)
	fmt.Printf("FilesCreated            		%d\n", s.FilesCreated)
	fmt.Printf("FilesUpdated            		%d\n", s.FilesUpdated)

	fmt.Printf("TxCommitted             		%d\n", s.TxCommitted)
	fmt.Printf("TxRolledBack            		%d\n", s.TxRolledBack)
}

func ImportKmedia() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia migration")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmedia, err := sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmedia.Ping())
	defer kmedia.Close()

	log.Info("Initializing static data from MDB")
	utils.Must(api.CONTENT_TYPE_REGISTRY.Init())
	utils.Must(api.OPERATION_TYPE_REGISTRY.Init())

	log.Info("Initializing static data from Kmedia")
	kmediaLessonCT, err = kmodels.ContentTypes(kmedia, qm.Where("name = ?", "Lesson")).One()
	utils.Must(err)
	langID2Locale, err = initLangs(kmedia)
	utils.Must(err)
	serverUrls, err = initServers(kmedia)
	utils.Must(err)

	log.Info("Loading all virtual_lessons")
	vls, err := kmodels.VirtualLessons(kmedia).All()
	utils.Must(err)
	log.Infof("Got %d lessons", len(vls))

	// Process lessons
	var (
		tx         *sql.Tx
		containers []*kmodels.Container
		operation  *models.Operation
		collection *models.Collection
		unit       *models.ContentUnit
	)
	stats = new(ImportStatistics)
	for i, vl := range vls {
		log.Infof("%d Processing virtual_lesson %d", i, vl.ID)
		stats.LessonsProcessed++

		// Validate virtual lesson data
		containers, err = getValidContainers(kmedia, vl)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}
		if len(containers) == 0 {
			log.Warningf("Invalid lesson [%d]", vl.ID)
			stats.InvalidLessons++
			continue
		}
		stats.ValidLessons++

		// Begin mdb transaction
		tx, err = mdb.Begin()
		utils.Must(err)

		// Create import operation
		operation, err = api.CreateOperation(tx, api.OP_IMPORT_KMEDIA, api.Operation{WorkflowID: strconv.Itoa(vl.ID)}, nil)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack++
			log.Error(err)
			debug.PrintStack()
			continue
		}
		stats.OperationsCreated++

		// Create or update MDB collection
		collection, err = handleCollection(tx, vl)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack++
			log.Error(err)
			debug.PrintStack()
			continue
		}

		for _, container := range containers {
			log.Infof("Processing container %d", container.ID)
			stats.ContainersProcessed++

			// Create or update MDB content_unit
			unit, err = handleContentUnit(tx, kmedia, container, collection)
			if err != nil {
				break
			}

			for _, fileAsset := range container.R.FileAssets {
				log.Infof("Processing file_asset %d", fileAsset.ID)
				stats.FileAssetsProcessed++

				// Create or update MDB file
				_, err = handleFile(tx, fileAsset, unit, operation)
				if err != nil {
					log.Error(err)
					debug.PrintStack()
					break
				}
			}
			if err != nil {
				break
			}
		}

		if err == nil {
			utils.Must(tx.Commit())
			stats.TxCommitted++
		} else {
			utils.Must(tx.Rollback())
			stats.TxRolledBack++
			log.Error(err)
			debug.PrintStack()
		}
	}

	// TODO: clean mdb stale data that no longer exists in kmedia
	// This would require some good understanding of how data can go stale.
	// Files are removed from kmedia ?
	// Containers merged ? what happens to old container ? removed / flagged ?
	// Lessons change ?

	stats.dump()

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func getValidContainers(exec boil.Executor, vl *kmodels.VirtualLesson) ([]*kmodels.Container, error) {
	// Fetch containers with file assets
	containers, err := vl.Containers(exec,
		qm.Where("content_type_id = ?", kmediaLessonCT.ID),
		qm.Load("FileAssets")).
		All()
	if err != nil {
		return nil, err
	}
	stats.ContainersVisited += len(containers)

	// Filter out containers without file_assets
	validContainers := containers[:0]
	for _, x := range containers {
		if len(x.R.FileAssets) > 0 {
			validContainers = append(validContainers, x)
			stats.ContainersWithFiles++
		} else {
			log.Warningf("Empty container [%d]", x.ID)
			stats.ContainersWithoutFiles++
		}
	}

	return validContainers, nil
}

func handleCollection(exec boil.Executor, vl *kmodels.VirtualLesson) (*models.Collection, error) {
	collection, err := models.Collections(exec, qm.Where("(properties->>'kmedia_id')::int = ?", vl.ID)).One()
	if err == nil {
		stats.CollectionsUpdated++
	} else {
		if err == sql.ErrNoRows {
			// Create new collection
			collection = &models.Collection{
				UID:    utils.GenerateUID(8),
				TypeID: api.CONTENT_TYPE_REGISTRY.ByName[api.CT_DAILY_LESSON].ID,
			}
			err = collection.Insert(exec)
			if err != nil {
				return nil, err
			}
			stats.CollectionsCreated++
		} else {
			return nil, err
		}
	}

	if vl.FilmDate.Time.Weekday() == 6 {
		collection.TypeID = api.CONTENT_TYPE_REGISTRY.ByName[api.CT_SATURDAY_LESSON].ID
	} else {
		collection.TypeID = api.CONTENT_TYPE_REGISTRY.ByName[api.CT_DAILY_LESSON].ID
	}

	var props = make(map[string]interface{})
	if collection.Properties.Valid {
		collection.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = vl.ID
	props["film_date"] = vl.FilmDate.Time.Format("2006-01-02")
	p, _ := json.Marshal(props)
	collection.Properties = null.JSONFrom(p)

	// TODO: what to do with name and description ?

	return collection, collection.Update(exec)
}

func handleContentUnit(exec boil.Executor, kmedia boil.Executor, container *kmodels.Container,
	collection *models.Collection) (*models.ContentUnit, error) {

	// Get or create content unit by kmedia_id
	unit, err := models.ContentUnits(exec, qm.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One()
	if err == nil {
		stats.ContentUnitsUpdated++
	} else {
		if err == sql.ErrNoRows {
			// Create new content unit
			unit = &models.ContentUnit{
				UID:    utils.GenerateUID(8),
				TypeID: api.CONTENT_TYPE_REGISTRY.ByName[api.CT_LESSON_PART].ID,
			}
			err = unit.Insert(exec)
			if err != nil {
				return nil, err
			}
			stats.ContentUnitsCreated++
		} else {
			return nil, err
		}
	}

	// Properties
	props := make(map[string]interface{})
	if unit.Properties.Valid {
		unit.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = container.ID
	props["secure"] = container.Secure
	if container.LangID.Valid {
		props["original_language"] = langID2Locale[container.LangID.String]
	}
	if container.Filmdate.Valid {
		props["film_date"] = container.Filmdate.Time.Format("2006-01-02")
	}
	if container.PlaytimeSecs.Valid {
		props["duration"] = container.PlaytimeSecs.Int
	}
	p, _ := json.Marshal(props)
	unit.Properties = null.JSONFrom(p)

	// TODO: what to do with censor workflow information ?

	err = unit.Update(exec, "properties")
	if err != nil {
		return nil, err
	}

	// I18n
	descriptions, err := container.ContainerDescriptions(kmedia).All()
	if err != nil {
		return nil, err
	}
	for _, d := range descriptions {
		cui18n := models.ContentUnitI18n{
			ContentUnitID: unit.ID,
			Language:      langID2Locale[d.LangID.String],
			Name:          d.ContainerDesc,
			Description:   d.Descr,
		}
		err = cui18n.Upsert(exec,
			true,
			[]string{"content_unit_id", "language"},
			[]string{"name", "description"})
		if err != nil {
			return nil, err
		}
	}

	// Associate content_unit with collection , name = position
	err = unit.L.LoadCollectionsContentUnits(exec, true, unit)
	if err != nil {
		return nil, err
	}

	position := strconv.Itoa(container.Position.Int)
	if unit.R.CollectionsContentUnits == nil {
		// Create
		err = unit.AddCollectionsContentUnits(exec, true, &models.CollectionsContentUnit{
			CollectionID:  collection.ID,
			ContentUnitID: unit.ID,
			Name:          position,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Update
		for _, x := range unit.R.CollectionsContentUnits {
			if x.CollectionID == collection.ID && x.Name != position {
				x.Name = position
				err = x.Update(exec, "name")
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return unit, nil
}

func handleFile(exec boil.Executor, fileAsset *kmodels.FileAsset, unit *models.ContentUnit,
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
		stats.FileAssetsMissingSHA1++
	}

	//file, _, err := api.FindFileBySHA1(exec, hash)
	file, hashB, err := api.FindFileBySHA1(exec, hash)
	if err == nil {
		stats.FilesUpdated++
	} else {
		if _, ok := err.(api.FileNotFound); ok {
			shouldCreate := true

			//For unknown file assets with valid sha1 do second lookup before we create a new file.
			//This time with the fake sha1, if exists we replace fake hash with valid hash.
			//Note: this paragraph should not be executed on first import.
			if fileAsset.Sha1.Valid {
				file, _, err = api.FindFileBySHA1(exec,
					fmt.Sprintf("%x", sha1.Sum([]byte(strconv.Itoa(fileAsset.ID)))))
				if err == nil {
					file.Sha1 = null.BytesFrom(hashB)
					shouldCreate = false
					stats.FilesUpdated++
				} else {
					if _, ok := err.(api.FileNotFound); !ok {
						return nil, err
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
					return nil, err
				}
				stats.FilesCreated++
			}
		} else {
			return nil, err
		}
	}

	// Media types
	if fileAsset.AssetType.Valid {
		if mt, ok := MEDIA_TYPES[strings.ToLower(fileAsset.AssetType.String)]; ok {
			file.Type = mt.Type
			file.SubType = mt.SubType
			file.MimeType = null.NewString(mt.MimeType, mt.MimeType != "")
		} else {
			stats.FileAssetsWInvalidMT++
		}
	} else {
		stats.FileAssetsMissingType++
	}

	// Language
	if fileAsset.LangID.Valid {
		l := langID2Locale[fileAsset.LangID.String]
		file.Language = null.NewString(l, l != "")
	}

	// Properties
	props := make(map[string]interface{})
	if file.Properties.Valid {
		file.Properties.Unmarshal(&props)
	}
	props["kmedia_id"] = fileAsset.ID
	props["secure"] = fileAsset.Secure
	props["url"] = serverUrls[fileAsset.ServernameID.String] + "/" + file.Name
	if fileAsset.PlaytimeSecs.Valid {
		props["duration"] = fileAsset.PlaytimeSecs.Int
	}
	p, _ := json.Marshal(props)
	file.Properties = null.JSONFrom(p)

	err = file.Update(exec)
	if err != nil {
		return nil, err
	}

	// i18n
	// We don't take anything from file_asset_descriptions as it`s mostly junk

	// Associate files with content unit
	err = unit.AddFiles(exec, false, file)
	if err != nil {
		return nil, err
	}

	// Associate files with operation

	// We use a raw query here to do nothing on conflicts
	// These conflicts happen when different file_assets in the same lesson have identical SHA1
	_, err = queries.Raw(exec,
		`INSERT INTO files_operations (file_id, operation_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		file.ID, operation.ID).Exec()

	return file, err
}

func initLangs(exec boil.Executor) (map[string]string, error) {
	languages, err := kmodels.Languages(exec).All()
	if err != nil {
		return nil, err
	}

	langID2Locale := make(map[string]string)
	for _, l := range languages {
		langID2Locale[l.Code3.String] = l.Locale.String
	}
	return langID2Locale, nil
}

func initServers(exec boil.Executor) (map[string]string, error) {
	servers, err := kmodels.Servers(exec).All()
	if err != nil {
		return nil, err
	}

	serverUrls := make(map[string]string)
	for _, s := range servers {
		serverUrls[s.Servername] = s.Httpurl.String
	}
	return serverUrls, nil
}
