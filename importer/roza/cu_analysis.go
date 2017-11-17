package roza

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var skipRE = regexp.MustCompile("\\.(eps|gif|tif|tiff|png|jpg|pdf|xls|htm|html|xml|zip|mso|db|wmz)$")

const (
	EMPTY                 = "empty"
	ROZA_ONLY             = "rozaOnly"
	ROZA_KM_ONLY          = "rozaKmOnly"
	ROZA_MDB_ONLY         = "rozaMdbOnly"
	PERFECT_STRIKE        = "perfectStrike"
	ALL_IN_NO_UNIT        = "allInNoUnit"
	ALL_IN_MISSING_UNIT   = "allInMissingUnit"
	ALL_IN_TOO_MANY_UNITS = "allInTooManyUnits"
	MIXED                 = "mixed"
)

type MatchAnalysis map[string][]*IdxDirDiff

type IdxDirDiff struct {
	status string
	Path   string         `json:"path"`
	Files  []*IdxFileDiff `json:"files,omitempty"`
}

type IdxFileDiff struct {
	Name    string `json:"name"`
	SHA1    string `json:"-"`
	MdbID   int64  `json:"mdbID,omitempty"`
	MdbCUID int64  `json:"mdbCUID,omitempty"`
	KmID    int    `json:"kmID,omitempty"`
	KmCnIDs []int  `json:"kmCnIDs,omitempty"`
}

type MiniKMFile struct {
	ID    int
	Sha1  string
	CnIDs []int
}

func MatchUnits() {
	clock := Init()

	idx := new(RozaIndex)
	utils.Must(idx.Load(mdb))

	mdbFiles, err := loadMDBFiles()
	utils.Must(err)

	kmFiles, err := loadKMFiles()
	utils.Must(err)

	cudMap, err := loadCUDs()
	utils.Must(err)

	ma, err := compareRozaToUnits(idx, mdbFiles, kmFiles, cudMap)
	utils.Must(err)
	for k, v := range ma {
		log.Infof("MA[%s] = %d", k, len(v))
	}

	f, err := os.OpenFile("importer/roza/analysis/all.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	utils.Must(err)
	defer f.Close()
	utils.Must(json.NewEncoder(f).Encode(ma))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadMDBFiles() (map[string]*models.File, error) {
	files, err := models.Files(mdb, qm.InnerJoin("roza_index r on files.sha1=r.sha1")).All()
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

func loadKMFiles() (map[string]*MiniKMFile, error) {
	rows, err := queries.Raw(kmdb, `SELECT
  fa.id,
  fa.sha1,
  array_agg(DISTINCT cfa.container_id)
FROM file_assets fa INNER JOIN containers_file_assets cfa ON fa.id = cfa.file_asset_id AND fa.sha1 IS NOT NULL
GROUP BY fa.id`).Query()
	if err != nil {
		return nil, errors.Wrap(err, "Load KMedia files")
	}
	defer rows.Close()

	m := make(map[string]*MiniKMFile, 600000)
	for rows.Next() {
		var f MiniKMFile
		var cnIDs pq.Int64Array
		if err := rows.Scan(&f.ID, &f.Sha1, &cnIDs); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}

		if len(cnIDs) > 0 {
			f.CnIDs = make([]int, len(cnIDs))
			for i := range cnIDs {
				f.CnIDs[i] = int(cnIDs[i])
			}
		}
		m[f.Sha1] = &f
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return m, nil
}

func loadCUDs() (map[int64]int64, error) {
	cuds, err := models.ContentUnitDerivations(mdb).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load MDB CUDs")
	}

	m := make(map[int64]int64, 500)
	for i := range cuds {
		x := cuds[i]
		m[x.DerivedID] = x.SourceID
	}

	return m, nil
}

func compareRozaToUnits(
	idx *RozaIndex,
	mdbFiles map[string]*models.File,
	kmFiles map[string]*MiniKMFile,
	cudMap map[int64]int64) (MatchAnalysis, error) {
	beavoda := idx.GetDir("/vfs/archive/Archive/____beavoda")
	if beavoda == nil {
		return nil, errors.New("____beavoda not found")
	}

	ma := make(MatchAnalysis)

	s := []*IdxDirectory{beavoda}
	var x *IdxDirectory
	for len(s) > 0 {
		x, s = s[0], s[1:]
		if len(x.Files) > 0 {
			diff := compareIdxDir(x, mdbFiles, kmFiles, cudMap)
			k := diff.status
			v, ok := ma[k]
			if !ok {
				v = make([]*IdxDirDiff, 0)
			}
			ma[k] = append(v, diff)
		}

		if x.Children.Empty() {
			continue
		}

		values := make([]*IdxDirectory, x.Children.Size())
		it := x.Children.Iterator()
		for i := 0; it.Next(); i++ {
			values[i] = it.Value().(*IdxDirectory)
		}
		s = append(values, s...)
	}

	return ma, nil
}

func compareIdxDir(
	d *IdxDirectory,
	mdbFiles map[string]*models.File,
	kmFiles map[string]*MiniKMFile,
	cudMap map[int64]int64) *IdxDirDiff {

	diff := &IdxDirDiff{Path: d.path(), Files: make([]*IdxFileDiff, 0)}

	// create file diffs
	for i := range d.Files {
		f := d.Files[i]

		if skipRE.MatchString(f.Name) {
			continue
		}

		fDiff := &IdxFileDiff{
			Name: f.Name,
			SHA1: f.SHA1,
		}

		if mdbF, ok := mdbFiles[f.SHA1]; ok {
			fDiff.MdbID = mdbF.ID
			if mdbF.ContentUnitID.Valid {
				if cuID, ok := cudMap[mdbF.ContentUnitID.Int64]; ok {
					fDiff.MdbCUID = cuID
				} else {
					fDiff.MdbCUID = mdbF.ContentUnitID.Int64
				}
			}
		}

		if kmF, ok := kmFiles[f.SHA1]; ok {
			fDiff.KmID = kmF.ID
			fDiff.KmCnIDs = kmF.CnIDs
		}

		diff.Files = append(diff.Files, fDiff)
	}

	if len(diff.Files) == 0 {
		diff.status = EMPTY
		return diff
	}

	// figure out total status
	cuIDs := make(map[int64]int)
	inMDB := 0
	inKM := 0
	for i := range diff.Files {
		fDiff := diff.Files[i]
		if fDiff.MdbID != 0 {
			inMDB++
			if fDiff.MdbCUID != 0 {
				cuIDs[fDiff.MdbCUID]++
			}
		}
		if fDiff.KmID != 0 {
			inKM++
		}
	}

	if inMDB == 0 && inKM == 0 {
		diff.status = ROZA_ONLY
	} else if inMDB == 0 && inKM == len(diff.Files) {
		diff.status = ROZA_KM_ONLY
	} else if inKM == 0 && inMDB == len(diff.Files) {
		diff.status = ROZA_MDB_ONLY
	} else if inMDB == inKM && inMDB == len(diff.Files) {
		if len(cuIDs) == 0 {
			diff.status = ALL_IN_NO_UNIT
		} else if len(cuIDs) == 1 {
			var inCU int
			for _, v := range cuIDs {
				inCU = v
			}
			if inCU == len(diff.Files) {
				diff.status = PERFECT_STRIKE
			} else {
				diff.status = ALL_IN_MISSING_UNIT
			}
		} else {
			diff.status = ALL_IN_TOO_MANY_UNITS
		}
	} else {
		diff.status = MIXED
	}

	return diff
}
