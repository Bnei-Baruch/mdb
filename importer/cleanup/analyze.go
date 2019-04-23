package cleanup

import (
	"fmt"
	"regexp"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func Analyze() {
	clock, _ := Init()

	utils.Must(doAnalyze())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

type CUAnalysis struct {
	*models.ContentUnit
	clipFiles   []*models.File
	noClipFiles []*models.File
}

func doAnalyze() error {
	cuCount, err := models.ContentUnits(mdb).Count()
	if err != nil {
		return errors.Wrapf(err, "Load content units count")
	}
	log.Infof("MDB has %d content units", cuCount)

	pageSize := 2500
	page := 0
	cuMap := make(map[int64]*CUAnalysis, cuCount)
	for page*pageSize < int(cuCount) {
		log.Infof("Loading page #%d", page)
		s := page * pageSize

		cus, err := models.ContentUnits(mdb,
			qm.Offset(s),
			qm.Limit(pageSize),
			qm.Load("Files")).
			All()
		if err != nil {
			return errors.Wrapf(err, "Load content units page %d", page)
		}
		for i := range cus {
			cuMap[cus[i].ID] = &CUAnalysis{ContentUnit: cus[i]}
		}
		page++
	}

	return findClipsOutsideClip(cuMap)
}

func findClipsOutsideClip(cuMap map[int64]*CUAnalysis) error {
	clipCT := api.CONTENT_TYPE_REGISTRY.ByName[api.CT_CLIP].ID
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

	log.Infof("%d units has unexpected clips in them", len(alerts))
	out := excelize.NewFile()
	out.SetCellStr("Sheet1", "A1", "Content Unit")
	out.SetCellStr("Sheet1", "B1", "Type")
	out.SetCellStr("Sheet1", "C1", "Full Files")
	out.SetCellStr("Sheet1", "D1", "Clip Files")

	row := 1
	for _, cu := range alerts {
		if len(cu.noClipFiles) == 0 {
			continue
		}

		row++

		url := fmt.Sprintf("http://app.mdb.bbdomain.org/admin/content_units/%d", cu.ID)
		out.SetCellStr("Sheet1", fmt.Sprintf("A%d", row), fmt.Sprintf("%d", cu.ID))
		out.SetCellHyperLink("Sheet1", fmt.Sprintf("A%d", row), url, "External")
		out.SetCellStr("Sheet1", fmt.Sprintf("B%d", row), api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)
		out.SetCellInt("Sheet1", fmt.Sprintf("C%d", row), len(cu.noClipFiles))
		out.SetCellInt("Sheet1", fmt.Sprintf("D%d", row), len(cu.clipFiles))


		//fileStrings := make([])
		//xlRow.AddCell().SetString(strings.Join("\n"))

		//log.Infof("%d\t%s\t%d\t%d",
		//	cu.ID, api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name, len(cu.noClipFiles), len(cu.clipFiles))
		//for i := range cu.clipFiles {
		//	f := cu.clipFiles[i]
		//	log.Infof("\t%d\t%s", f.ID, f.Name)
		//}
		//
		//f := cu.noClipFiles[0]
		//log.Infof("\t\t%d\t%s", f.ID, f.Name)
	}

	err := out.SaveAs("organize_clips.xlsx")
	if err != nil {
		return errors.Wrap(err, "out.SaveAs")
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
