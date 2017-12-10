package roza

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type MFile struct {
	sha1      string
	ID        int64  `json:"id"`
	Name      string `json:"n"`
	CUID      int64  `json:"cuid,omitempty"`
	Published bool   `json:"p"`
}

type MiniKFile struct {
	ID    int    `json:"id"`
	Name  string `json:"n"`
	CnIDs []int  `json:"cnIDs,omitempty"`
}

type KFile struct {
	sha1   string
	Copies []*MiniKFile `json:"c,omitempty"`
}

type RFile struct {
	sha1   string
	Copies []*IdxFile `json:"c,omitempty"`
}

type MasterFile struct {
	sha1   string
	MDB    *MFile     `json:"m,omitempty"`
	KMedia *KFile     `json:"k,omitempty"`
	Roza   *RFile `json:"r,omitempty"`
}

func MatchFiles() {
	clock := Init()

	idx := new(RozaIndex)
	utils.Must(idx.Load(mdb))
	rFiles := idx.Sha1Map()

	kFiles, err := loadKFiles()
	utils.Must(err)

	cudMap, err := loadCUDs()
	utils.Must(err)

	mFiles, err := loadMFiles(cudMap)
	utils.Must(err)

	mm := masterMerge(mFiles, kFiles, rFiles)

	f, err := os.OpenFile("importer/roza/analysis/master.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	utils.Must(err)
	defer f.Close()
	utils.Must(json.NewEncoder(f).Encode(mm))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadMFiles(cudMap map[int64]int64) (map[string]*MFile, error) {
	files, err := models.Files(mdb, qm.Where("sha1 is not null")).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load MDB files")
	}

	m := make(map[string]*MFile, 600000)
	for i := range files {
		f := files[i]
		mFile := &MFile{
			sha1:      hex.EncodeToString(f.Sha1.Bytes),
			ID:        f.ID,
			Name:      f.Name,
			Published: f.Published,
		}

		if f.ContentUnitID.Valid {
			cuid, ok := cudMap[f.ContentUnitID.Int64]
			if ok {
				mFile.CUID = cuid
			} else {
				mFile.CUID = f.ContentUnitID.Int64
			}
		}
		m[mFile.sha1] = mFile
	}

	return m, nil
}

func loadKFiles() (map[string]*KFile, error) {
	rows, err := queries.Raw(kmdb, `SELECT
  fa.id,
  fa.sha1,
  fa.name,
  array_agg(DISTINCT cfa.container_id)
FROM file_assets fa INNER JOIN containers_file_assets cfa ON fa.id = cfa.file_asset_id AND fa.sha1 IS NOT NULL
GROUP BY fa.id`).Query()
	if err != nil {
		return nil, errors.Wrap(err, "Load KMedia files")
	}
	defer rows.Close()

	m := make(map[string]*KFile, 600000)
	for rows.Next() {
		var sha1 string
		var f MiniKFile
		var cnIDs pq.Int64Array
		if err := rows.Scan(&f.ID, &sha1, &f.Name, &cnIDs); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}

		if len(cnIDs) > 0 {
			f.CnIDs = make([]int, len(cnIDs))
			for i := range cnIDs {
				f.CnIDs[i] = int(cnIDs[i])
			}
		}

		k := sha1
		v, ok := m[k]
		if !ok {
			v = new(KFile)
			m[k] = v
		}
		v.Copies = append(v.Copies, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return m, nil
}

func masterMerge(mFiles map[string]*MFile, kFiles map[string]*KFile, rFiles map[string][]*IdxFile) map[string]*MasterFile {
	mm := make(map[string]*MasterFile, 600000)

	for k, v := range mFiles {
		vv, ok := mm[k]
		if !ok {
			vv = new(MasterFile)
			mm[k] = vv
		}
		vv.MDB = v
	}

	for k, v := range kFiles {
		vv, ok := mm[k]
		if !ok {
			vv = new(MasterFile)
			mm[k] = vv
		}
		vv.KMedia = v
	}

	for k, v := range rFiles {
		vv, ok := mm[k]
		if !ok {
			vv = new(MasterFile)
			mm[k] = vv
		}
		vv.Roza = &RFile{
			Copies: v,
		}
	}

	return mm
}
