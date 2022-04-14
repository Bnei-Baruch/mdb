package cleanup

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const CLIP_COLLECTION_UID = "MISCCLIP"

func Analyze() {
	clock, _ := Init()

	utils.Must(doAnalyze())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func Import() {
	clock, _ := Init()

	utils.Must(doImport())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

type CUAnalysis struct {
	*common.UnitWName
	clipFiles   []*models.File
	noClipFiles []*models.File
}

func doAnalyze() error {
	cuMap, err := loadUnits()
	if err != nil {
		return errors.WithMessage(err, "load CUs")
	}

	alerts, err := findClipsOutsideClip(cuMap)
	if err != nil {
		return errors.WithMessage(err, "find alerts")
	}

	err = dumpExcel(alerts)
	if err != nil {
		return errors.WithMessage(err, "dump excel")
	}

	return nil
}

func doImport() error {
	cuMap, err := loadUnits()
	if err != nil {
		return errors.WithMessage(err, "load CUs")
	}

	alerts, err := findClipsOutsideClip(cuMap)
	if err != nil {
		return errors.WithMessage(err, "find alerts")
	}

	c, err := api.FindCollectionByUID(mdb, CLIP_COLLECTION_UID)
	if err != nil {
		return errors.WithMessage(err, "find clips collection")
	}

	for k, v := range alerts {
		if err := doImportAlert(v, c); err != nil {
			return errors.WithMessagef(err, "import alert cu [%d]", k)
		}
	}

	return err
}

func loadUnits() (map[int64]*CUAnalysis, error) {
	cuCount, err := models.ContentUnits().Count(mdb)
	if err != nil {
		return nil, errors.Wrapf(err, "Load content units count")
	}
	log.Infof("MDB has %d content units", cuCount)

	pageSize := 2500
	page := 0
	cuMap := make(map[int64]*CUAnalysis, cuCount)
	for page*pageSize < int(cuCount) {
		log.Infof("Loading page #%d", page)
		s := page * pageSize

		cus, err := models.ContentUnits(
			qm.Offset(s),
			qm.Limit(pageSize),
			qm.Load("Files"),
			qm.Load("ContentUnitI18ns")).
			All(mdb)
		if err != nil {
			return nil, errors.Wrapf(err, "Load content units page %d", page)
		}
		for i := range cus {
			cuMap[cus[i].ID] = &CUAnalysis{UnitWName: &common.UnitWName{ContentUnit: cus[i]}}
		}
		page++
	}

	return cuMap, nil
}

func findClipsOutsideClip(cuMap map[int64]*CUAnalysis) (map[int64]*CUAnalysis, error) {
	clipCT := common.CONTENT_TYPE_REGISTRY.ByName[common.CT_CLIP].ID
	clipRE := regexp.MustCompile("clip")
	lessonOrProgramRE := regexp.MustCompile("_(lesson|program)_")

	alerts := make(map[int64]*CUAnalysis)
	for _, cu := range cuMap {
		if cu.TypeID == clipCT {
			continue
		}

		for i := range cu.R.Files {
			f := cu.R.Files[i]
			if clipRE.MatchString(f.Name) {
				cu.clipFiles = append(cu.clipFiles, f)
			} else if lessonOrProgramRE.MatchString(f.Name) {
				cu.noClipFiles = append(cu.noClipFiles, f)
			}
		}
		if len(cu.clipFiles) > 0 {
			alerts[cu.ID] = cu
		}
	}

	return alerts, nil
}

func dumpExcel(alerts map[int64]*CUAnalysis) error {
	log.Infof("%d units has unexpected clips in them", len(alerts))

	out := excelize.NewFile()
	out.SetCellStr("Sheet1", "A1", "ID")
	out.SetCellStr("Sheet1", "B1", "Name")
	out.SetCellStr("Sheet1", "C1", "Type")
	out.SetCellStr("Sheet1", "D1", "Full Files")
	out.SetCellStr("Sheet1", "E1", "Clip Files")

	row := 1
	for _, cu := range alerts {
		if len(cu.noClipFiles) == 0 {
			continue
		}

		row++

		url := fmt.Sprintf("http://app.mdb.bbdomain.org/admin/content_units/%d", cu.ID)
		out.SetCellStr("Sheet1", fmt.Sprintf("A%d", row), fmt.Sprintf("%d", cu.ID))
		out.SetCellHyperLink("Sheet1", fmt.Sprintf("A%d", row), url, "External")
		out.SetCellStr("Sheet1", fmt.Sprintf("B%d", row), cu.Name())
		out.SetCellStr("Sheet1", fmt.Sprintf("C%d", row), common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)
		out.SetCellInt("Sheet1", fmt.Sprintf("D%d", row), len(cu.noClipFiles))
		out.SetCellInt("Sheet1", fmt.Sprintf("E%d", row), len(cu.clipFiles))
	}

	err := out.SaveAs("organize_clips.xlsx")
	if err != nil {
		return errors.Wrap(err, "out.SaveAs")
	}

	return nil
}

func doImportAlert(cu *CUAnalysis, clipsCollection *models.Collection) error {
	// group clip files by name stripped from language and extension
	fileGroups := make(map[string][]*models.File)
	for i := range cu.clipFiles {
		k := importer.NormalizedFileName(cu.clipFiles[i].Name)
		fileGroups[k] = append(fileGroups[k], cu.clipFiles[i])
	}
	log.Infof("CU [%d] has %d clips file groups", cu.ID, len(fileGroups))

	// for each group, create a new, sensitive, derived unit
	// associate the files with that unit
	// associate the unit with the clips collection
	for k, v := range fileGroups {
		log.Infof("\t%s\t%d", k, len(v))

		filmDates := make(map[string]int)
		langs := make(map[string]int)
		for i := range v {
			line := importer.ParseLine(v[i].Name)

			if line.Language == "" {
				log.Warnf("Empty language, [%d] %s", v[i].ID, v[i].Name)
			}
			langs[line.Language]++

			if line.FilmDate == "" {
				log.Warnf("Empty film_date, [%d] %s", v[i].ID, v[i].Name)
			}
			filmDates[line.FilmDate]++
		}

		props := make(map[string]interface{})

		if len(filmDates) != 1 {
			log.Warnf("Unexpected number of film_date %d", len(filmDates))
			props["film_date"] = "1970-01-01"
		} else {
			for k := range filmDates {
				props["film_date"] = k
				break
			}
		}

		if len(langs) == 0 {
			log.Warn("No languages")
		} else if len(langs) == 1 {
			for k := range langs {
				props["original_language"] = common.StdLang(k)
				break
			}
		} else {
			props["original_language"] = common.LANG_UNKNOWN
		}

		tx, err := mdb.Begin()
		utils.Must(err)

		clipCU, err := api.CreateContentUnit(tx, common.CT_CLIP, props)
		if err != nil {
			utils.Must(tx.Rollback())
			log.Error(err)
			debug.PrintStack()
			continue
		}
		log.Infof("New CU ID is %d", clipCU.ID)

		clipCU.Secure = common.SEC_SENSITIVE
		_, err = clipCU.Update(tx, boil.Whitelist("secure"))
		if err != nil {
			utils.Must(tx.Rollback())
			log.Error(err)
			debug.PrintStack()
			continue
		}

		i18ns := make([]*models.ContentUnitI18n, 0)
		//for k, v := range names {
		for _, lang := range [...]string{common.LANG_HEBREW, common.LANG_ENGLISH, common.LANG_RUSSIAN, common.LANG_SPANISH} {
			i18n := &models.ContentUnitI18n{
				ContentUnitID: clipCU.ID,
				Language:      lang,
				Name:          null.StringFrom(k),
			}
			i18ns = append(i18ns, i18n)
		}
		err = clipCU.AddContentUnitI18ns(tx, true, i18ns...)
		if err != nil {
			return errors.Wrap(err, "Save to DB")
		}

		err = clipCU.AddFiles(tx, false, v...)
		if err != nil {
			utils.Must(tx.Rollback())
			log.Error(err)
			debug.PrintStack()
			continue
		}

		cud := &models.ContentUnitDerivation{
			SourceID: cu.ID,
			Name:     k,
		}
		err = clipCU.AddDerivedContentUnitDerivations(tx, true, cud)
		if err != nil {
			utils.Must(tx.Rollback())
			log.Error(err)
			debug.PrintStack()
			continue
		}

		ccu := &models.CollectionsContentUnit{
			CollectionID:  clipsCollection.ID,
			ContentUnitID: clipCU.ID,
		}

		err = clipsCollection.AddCollectionsContentUnits(tx, true, ccu)
		if err != nil {
			utils.Must(tx.Rollback())
			log.Error(err)
			debug.PrintStack()
			continue
		}

		utils.Must(tx.Commit())
	}

	return nil
}

//func doChunkSize() error {
//	fCount, err := models.Files(mdb).Count()
//	if err != nil {
//		return errors.Wrapf(err, "Load files count")
//	}
//	log.Infof("MDB has %d files", fCount)
//
//	pageSize := 5000
//	page := 0
//	fMap := make(map[int64]*models.File, fCount)
//	for page*pageSize < int(fCount) {
//		log.Infof("Loading page #%d", page)
//		s := page * pageSize
//
//		files, err := models.Files(mdb,
//			qm.Select("id", "name", "size"),
//			qm.Offset(s),
//			qm.Limit(pageSize)).
//			All()
//		if err != nil {
//			return errors.Wrapf(err, "Load files page %d", page)
//		}
//		for i := range files {
//			fMap[files[i].ID] = files[i]
//		}
//		page++
//	}
//
//	return nil
//}
//
//type ReedSolomnParams struct {
//	Data   int
//	Parity int
//}
//
//var RS9x3 = &ReedSolomnParams{Data: 9, Parity:3}
//var RS4x2 = &ReedSolomnParams{Data: 4, Parity:2}
//var (
//	c64k = 64 * 1024
//	c128k = 2 * c64k
//	c256k = 2 * c128k
//	c512k = 2 * c256k
//	c1m = 2 * c512k
//	c2m = 2 * c1m
//	c4m = 2 * c2m
//	c8m = 2 * c4m
//	c16m = 2 * c8m
//	c32m = 2 * c16m
//	c64m = 2 * c32m
//	chunkSizes = [...]int{c64k, c128k, c256k, c512k, c1m, c2m, c4m, c8m,c16m, c32m, c64m}
//	)
//
//type FileChunkStats struct {
//	ChunkSize int
//	Data int
//	Parity int
//	Empty int
//}
//
//func BestChunking(fSize int64, rsParams *ReedSolomnParams) FileChunkStats {
//	minChunkSize := int(math.Floor(float64(fSize) / float64(rsParams.Data)))
//
//	bestChunk :=-1
//	for i := range chunkSizes {
//		if chunkSizes[i] < minChunkSize {
//			bestChunk = i
//		}
//	}
//
//	if bestChunk > 0 {
//		//TODO: files larger then 9 * max chunk
//		stats := FileChunkStats{ChunkSize: chunkSizes[bestChunk]}
//		stats.Data = int(math.Ceil(float64(fSize) / float64(stats.ChunkSize)))
//		modData := stats.Data % rsParams.Data
//		stats.Data += modData
//		if modData > 0 {
//			stats.Empty = rsParams.Data - modData
//		}
//		stats.Parity = rsParams.Parity * ((stats.Data + stats.Empty) % rsParams.Data)
//		return stats
//	}
//
//	return FileChunkStats{}
//}
