package roza

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type UnitAnalysis map[string][]*UnitDiff

type UnitDiff struct {
	CUID   int64 `json:"cuID"`
	status string
	Files  []*MDBFileDiff `json:"files"`
}

type MDBFileDiff struct {
	MdbName  string     `json:"name"`
	MdbID    int64      `json:"mdbID"`
	MdbCUID  int64      `json:"mdbCUID,omitempty"`
	KmID     int        `json:"kmID,omitempty"`
	KmCnIDs  []int      `json:"kmCnIDs,omitempty"`
	IdxFiles []*IdxFile `json:"directories,omitempty"`
}

func MatchDirectories() {
	clock := Init()

	idx := new(RozaIndex)
	utils.Must(idx.Load(mdb))

	mdbFiles, err := loadMDBPublishedFiles()
	utils.Must(err)

	kmFiles, err := loadKMFiles()
	utils.Must(err)

	cudMap, err := loadCUDs()
	utils.Must(err)

	ua, err := compareUnitsToRoza(idx, mdbFiles, kmFiles, cudMap)
	utils.Must(err)
	for k, v := range ua {
		log.Infof("UA[%s] = %d", k, len(v))
	}

	f, err := os.OpenFile("importer/roza/analysis/mdb.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	utils.Must(err)
	defer f.Close()
	utils.Must(json.NewEncoder(f).Encode(ua))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadMDBPublishedFiles() (map[string]*models.File, error) {
	files, err := models.Files(mdb, qm.Where("sha1 is not null and content_unit_id is not null and published is true")).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load MDB files")
	}

	m := make(map[string]*models.File, 600000)
	for i := range files {
		f := files[i]
		m[hex.EncodeToString(f.Sha1.Bytes)] = f
	}

	return m, nil
}

func compareUnitsToRoza(
	idx *RozaIndex,
	mdbFiles map[string]*models.File,
	kmFiles map[string]*MiniKMFile,
	cudMap map[int64]int64) (UnitAnalysis, error) {

	sMap := idx.Sha1Map()

	// mdb file diff per unit
	unitDiffs := make(map[int64]*UnitDiff, 40000)
	for i := range mdbFiles {
		f := mdbFiles[i]
		cuID := f.ContentUnitID.Int64
		if v, ok := cudMap[cuID]; ok {
			cuID = v
		}

		cuDiff, ok := unitDiffs[cuID]
		if !ok {
			cuDiff = &UnitDiff{
				CUID:  cuID,
				Files: make([]*MDBFileDiff, 0),
			}
			unitDiffs[cuID] = cuDiff
		}

		fDiff := &MDBFileDiff{
			MdbName:  f.Name,
			MdbID:    f.ID,
			MdbCUID:  cuID,
			IdxFiles: sMap[i],
		}

		// kmedia
		if kmF, ok := kmFiles[i]; ok {
			fDiff.KmID = kmF.ID
			fDiff.KmCnIDs = kmF.CnIDs
		}

		cuDiff.Files = append(cuDiff.Files, fDiff)
	}

	ua := make(UnitAnalysis)

	for i := range unitDiffs {
		diff := unitDiffs[i]
		inRoza := 0
		dirMap := make(map[string]int, 0)
		for j := range diff.Files {
			fDiff := diff.Files[j]
			if len(fDiff.IdxFiles) > 0 {
				inRoza++
				for k := range fDiff.IdxFiles {
					dirMap[fDiff.IdxFiles[k].Directory.path()]++
				}
			}
		}

		if inRoza == 0 {
			diff.status = MDB_ONLY
		} else if inRoza == len(diff.Files) {
			if len(dirMap) == 1 {
				diff.status = PERFECT_STRIKE
			} else {
				diff.status = ALL_IN_TOO_MANY_FOLDERS
			}
		} else {
			diff.status = MIXED
		}

		k := diff.status
		v, ok := ua[k]
		if !ok {
			v = make([]*UnitDiff, 0)
		}
		ua[k] = append(v, diff)
	}

	return ua, nil
}
