package kmedia

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/pkg/errors"
	"sort"
)

const CATALOGS_MAPPINGS_FILE = "importer/kmedia/data/Catalogs Sources Mappings - final.csv"

var (
	mdb                     *sql.DB
	kmdb                    *sql.DB
	kmediaLessonCT          *kmodels.ContentType
	langID2Locale           map[string]string
	serverUrls              map[string]string
	catalogsSourcesMappings map[int]*models.Source
	stats                   *ImportStatistics
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
	"rb":   {Extension: "rb", Type: "text", SubType: "", MimeType: "application/x-rocketbook"},
	"xls":  {Extension: "xls", Type: "sheet", SubType: "", MimeType: "application/vnd.ms-excel"},
	"swf":  {Extension: "swf", Type: "banner", SubType: "", MimeType: "application/x-shockwave-flash"},
	"ppt":  {Extension: "ppt", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
	"pptx": {Extension: "pptx", Type: "presentation", SubType: "", MimeType: "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	"pps":  {Extension: "pps", Type: "presentation", SubType: "", MimeType: "application/vnd.ms-powerpoint"},
}

type AtomicInt32 struct {
	value int32
}

func (a *AtomicInt32) Inc(delta int32) {
	atomic.AddInt32(&a.value, delta)
}

func (a *AtomicInt32) Get() int32 {
	return atomic.LoadInt32(&a.value)
}

type AtomicHistogram struct {
	value map[string]*AtomicInt32
}

func NewAtomicHistogram() *AtomicHistogram {
	return &AtomicHistogram{value: make(map[string]*AtomicInt32)}
}

func (h *AtomicHistogram) Inc(key string, delta int32) {
	v, ok := h.value[key]
	if !ok {
		v = new(AtomicInt32)
		h.value[key] = v
	}
	v.Inc(delta)
}

func (h *AtomicHistogram) Get() map[string]int32 {
	r := make(map[string]int32, len(h.value))
	for k, v := range h.value {
		r[k] = v.Get()
	}
	return r
}

type ImportStatistics struct {
	LessonsProcessed       AtomicInt32
	ValidLessons           AtomicInt32
	InvalidLessons         AtomicInt32
	ContainersProcessed    AtomicInt32
	ContainersVisited      AtomicInt32
	ContainersWithFiles    AtomicInt32
	ContainersWithoutFiles AtomicInt32
	FileAssetsProcessed    AtomicInt32
	FileAssetsMissingSHA1  AtomicInt32
	FileAssetsWInvalidMT   AtomicInt32
	FileAssetsMissingType  AtomicInt32

	OperationsCreated   AtomicInt32
	CollectionsCreated  AtomicInt32
	CollectionsUpdated  AtomicInt32
	ContentUnitsCreated AtomicInt32
	ContentUnitsUpdated AtomicInt32
	FilesCreated        AtomicInt32
	FilesUpdated        AtomicInt32

	TxCommitted  AtomicInt32
	TxRolledBack AtomicInt32

	UnkownCatalogs AtomicHistogram
}

func NewImportStatistics() *ImportStatistics {
	return &ImportStatistics{UnkownCatalogs: *NewAtomicHistogram()}
}

func (s *ImportStatistics) dump() {
	fmt.Println("Here comes import statistics:")

	fmt.Println("Kmedia:")
	fmt.Printf("ValidLessons            		%d\n", s.ValidLessons.Get())
	fmt.Printf("InvalidLessons          		%d\n", s.InvalidLessons.Get())
	fmt.Printf("LessonsProcessed        		%d\n", s.LessonsProcessed.Get())
	fmt.Printf("ContainersVisited       		%d\n", s.ContainersVisited.Get())
	fmt.Printf("ContainersWithFiles     		%d\n", s.ContainersWithFiles.Get())
	fmt.Printf("ContainersWithoutFiles  		%d\n", s.ContainersWithoutFiles.Get())
	fmt.Printf("ContainersProcessed     		%d\n", s.ContainersProcessed.Get())
	fmt.Printf("FileAssetsProcessed     		%d\n", s.FileAssetsProcessed.Get())
	fmt.Printf("FileAssetsMissingSHA1   		%d\n", s.FileAssetsMissingSHA1.Get())
	fmt.Printf("FileAssetsWInvalidMT    		%d\n", s.FileAssetsWInvalidMT.Get())
	fmt.Printf("FileAssetsMissingType   		%d\n", s.FileAssetsMissingType.Get())

	fmt.Println("MDB:")
	fmt.Printf("OperationsCreated       		%d\n", s.OperationsCreated.Get())
	fmt.Printf("CollectionsCreated      		%d\n", s.CollectionsCreated.Get())
	fmt.Printf("CollectionsUpdated      		%d\n", s.CollectionsUpdated.Get())
	fmt.Printf("ContentUnitsCreated     		%d\n", s.ContentUnitsCreated.Get())
	fmt.Printf("ContentUnitsUpdated     		%d\n", s.ContentUnitsUpdated.Get())
	fmt.Printf("FilesCreated            		%d\n", s.FilesCreated.Get())
	fmt.Printf("FilesUpdated            		%d\n", s.FilesUpdated.Get())

	fmt.Printf("TxCommitted             		%d\n", s.TxCommitted.Get())
	fmt.Printf("TxRolledBack            		%d\n", s.TxRolledBack.Get())

	uc := s.UnkownCatalogs.Get()
	fmt.Printf("UnkownCatalogs            		%d\n", len(uc))
	keys := make([]string, len(uc))
	i := 0
	for k := range uc {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, uc[k])
	}
}

func ImportKmedia() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia migration")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmdb, err = sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmdb.Ping())
	defer kmdb.Close()

	log.Info("Initializing static data from MDB")
	utils.Must(api.CONTENT_TYPE_REGISTRY.Init())
	utils.Must(api.OPERATION_TYPE_REGISTRY.Init())

	log.Info("Initializing static data from Kmedia")
	kmediaLessonCT, err = kmodels.ContentTypes(kmdb, qm.Where("name = ?", "Lesson")).One()
	utils.Must(err)
	langID2Locale, err = initLangs(kmdb)
	utils.Must(err)
	serverUrls, err = initServers(kmdb)
	utils.Must(err)

	log.Info("Initializing catalogs sources mappings")
	catalogsSourcesMappings, err = initCatalogMappings()
	utils.Must(err)
	log.Infof("Got %d mappings", len(catalogsSourcesMappings))

	log.Info("Loading all virtual_lessons")
	vls, err := kmodels.VirtualLessons(kmdb).All()
	utils.Must(err)
	log.Infof("Got %d lessons", len(vls))

	// Process lessons
	stats = NewImportStatistics()

	log.Info("Setting up workers")
	jobs := make(chan *kmodels.VirtualLesson, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		workersWG.Add(1)
		go worker(jobs, &workersWG)
	}

	log.Info("Queing work")
	for _, vl := range vls {
		jobs <- vl
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	// TODO: clean mdb stale data that no longer exists in kmedia
	// This would require some good understanding of how data can go stale.
	// Files are removed from kmedia ?
	// Containers merged ? what happens to old container ? removed / flagged ?
	// Lessons change ?

	stats.dump()

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func worker(jobs <-chan *kmodels.VirtualLesson, wg *sync.WaitGroup) {
	for vl := range jobs {
		log.Infof("Processing virtual_lesson %d", vl.ID)
		stats.LessonsProcessed.Inc(1)

		// Validate virtual lesson data
		containers, err := getValidContainers(kmdb, vl)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			continue
		}
		if len(containers) == 0 {
			log.Warnf("Invalid lesson [%d]", vl.ID)
			stats.InvalidLessons.Inc(1)
			continue
		}
		stats.ValidLessons.Inc(1)

		// Begin mdb transaction
		tx, err := mdb.Begin()
		utils.Must(err)

		// Create import operation
		operation, err := api.CreateOperation(tx, api.OP_IMPORT_KMEDIA,
			api.Operation{WorkflowID: strconv.Itoa(vl.ID)}, nil)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		}
		stats.OperationsCreated.Inc(1)

		// Handle MDB collection
		collection, err := handleCollection(tx, vl)
		if err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		}

		var unit *models.ContentUnit
		for _, container := range containers {
			log.Infof("Processing container %d", container.ID)
			stats.ContainersProcessed.Inc(1)

			// Create or update MDB content_unit
			unit, err = handleContentUnit(tx, container, collection)
			if err != nil {
				log.Error(err)
				debug.PrintStack()
				break
			}

			for _, fileAsset := range container.R.FileAssets {
				log.Infof("Processing file_asset %d", fileAsset.ID)
				stats.FileAssetsProcessed.Inc(1)

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
			stats.TxCommitted.Inc(1)
		} else {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
		}
	}

	wg.Done()
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
	stats.ContainersVisited.Inc(int32(len(containers)))

	// Filter out containers without file_assets
	validContainers := containers[:0]
	for _, x := range containers {
		if len(x.R.FileAssets) > 0 {
			validContainers = append(validContainers, x)
			stats.ContainersWithFiles.Inc(1)
		} else {
			log.Warningf("Empty container [%d]", x.ID)
			stats.ContainersWithoutFiles.Inc(1)
		}
	}

	return validContainers, nil
}

func handleCollection(exec boil.Executor, vl *kmodels.VirtualLesson) (*models.Collection, error) {
	collection, err := models.Collections(exec, qm.Where("(properties->>'kmedia_id')::int = ?", vl.ID)).One()
	if err == nil {
		stats.CollectionsUpdated.Inc(1)
	} else {
		if err == sql.ErrNoRows {
			// Create new collection
			collection = &models.Collection{
				UID:    utils.GenerateUID(8),
				TypeID: api.CONTENT_TYPE_REGISTRY.ByName[api.CT_DAILY_LESSON].ID,
			}
			err = collection.Insert(exec)
			if err != nil {
				return nil, errors.Wrapf(err, "Insert collection, virtual lesson [%d]", vl.ID)
			}
			stats.CollectionsCreated.Inc(1)
		} else {
			return nil, errors.Wrapf(err, "Lookup collection, virtual lesson [%d]", vl.ID)
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

func handleContentUnit(exec boil.Executor, container *kmodels.Container,
	collection *models.Collection) (*models.ContentUnit, error) {

	// Get or create content unit by kmedia_id
	unit, err := models.ContentUnits(exec, qm.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One()
	if err == nil {
		stats.ContentUnitsUpdated.Inc(1)
	} else {
		if err == sql.ErrNoRows {
			// Create new content unit
			unit = &models.ContentUnit{
				UID:    utils.GenerateUID(8),
				TypeID: api.CONTENT_TYPE_REGISTRY.ByName[api.CT_LESSON_PART].ID,
			}
			err = unit.Insert(exec)
			if err != nil {
				return nil, errors.Wrapf(err, "Insert unit, container_id [%d]", container.ID)
			}
			stats.ContentUnitsCreated.Inc(1)
		} else {
			return nil, errors.Wrapf(err, "Lookup unit, container_id [%d]", container.ID)
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
		return nil, errors.Wrapf(err, "Update properties, unit [%d]", unit.ID)
	}

	// I18n
	descriptions, err := container.ContainerDescriptions(kmdb).All()
	if err != nil {
		return nil, errors.Wrapf(err, "Lookup container descriptions, container_id [%d]", container.ID)
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
			return nil, errors.Wrapf(err, "Upsert unit i18n, unit [%d]", unit.ID)
		}
	}

	// Associate content_unit with collection , name = position
	err = unit.L.LoadCollectionsContentUnits(exec, true, unit)
	if err != nil {
		return nil, errors.Wrapf(err, "Fetch unit collections, unit [%d]", unit.ID)
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
			return nil, errors.Wrapf(err, "Add unit collections, unit [%d]", unit.ID)
		}
	} else {
		// Update
		for _, x := range unit.R.CollectionsContentUnits {
			if x.CollectionID == collection.ID && x.Name != position {
				x.Name = position
				err = x.Update(exec, "name")
				if err != nil {
					return nil, errors.Wrapf(err,
						"Update unit collection association, unit [%d], collection [%d]",
						unit.ID, collection.ID)
				}
			}
		}
	}

	// Associate sources
	err = container.L.LoadCatalogs(kmdb, true, container)
	if err != nil {
		return nil, errors.Wrapf(err, "Load catalogs, container [%d]", container.ID)
	}

	// Use a map as a set for uniqueness
	unqSrc := make(map[*models.Source]bool, 0)
	for _, x := range container.R.Catalogs {
		s, ok := catalogsSourcesMappings[x.ID]
		if ok {
			unqSrc[s] = true
		} else {
			stats.UnkownCatalogs.Inc(fmt.Sprintf("%s [%d]", x.Name, x.ID), 1)
		}
	}

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
		if mt, ok := MEDIA_TYPES[strings.ToLower(fileAsset.AssetType.String)]; ok {
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
		return nil, errors.Wrapf(err, "Update file [%d]", file.ID)
	}

	// i18n
	// We don't take anything from file_asset_descriptions as it`s mostly junk

	// Associate files with content unit
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

func initCatalogMappings() (map[int]*models.Source, error) {
	// read catalogs mappings
	records, err := utils.ReadCSV(CATALOGS_MAPPINGS_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read catalogs mappings")
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
