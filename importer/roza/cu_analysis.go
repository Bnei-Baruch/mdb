package roza

import (
	"encoding/hex"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

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

	utils.Must(compareRozaToUnits(idx, mdbFiles, kmFiles, cudMap))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadMDBFiles() (map[string]*models.File, error) {
	files, err := models.Files(mdb, qm.InnerJoin("roza_index r on files.sha1=r.sha1")).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load MDB files")
	}

	m := make(map[string]*models.File, 400000)
	for i := range files {
		f := files[i]
		m[hex.EncodeToString(f.Sha1.Bytes)] = f
	}

	return m, nil
}

func loadKMFiles() (map[string]*kmodels.FileAsset, error) {
	files, err := kmodels.FileAssets(kmdb, qm.Where("sha1 is not null")).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load KMedia files")
	}

	m := make(map[string]*kmodels.FileAsset, 600000)
	for i := range files {
		f := files[i]
		m[f.Sha1.String] = f
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
	kmFiles map[string]*kmodels.FileAsset,
	cudMap map[int64]int64) error {
	beavoda := idx.GetDir("/vfs/archive/Archive/____beavoda")
	if beavoda == nil {
		return errors.New("____beavoda not found")
	}

	s := []*IdxDirectory{beavoda}
	var x *IdxDirectory
	for len(s) > 0 {
		x, s = s[0], s[1:]
		compareIdxDir(x, mdbFiles, kmFiles, cudMap)

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

	return nil
}

func compareIdxDir(
	d *IdxDirectory,
	mdbFiles map[string]*models.File,
	kmFiles map[string]*kmodels.FileAsset,
	cudMap map[int64]int64) {
	if len(d.Files) == 0 {
		log.Infof("%s", d.Name)
		return
	}

	cuIDs := make(map[int64]int)
	notInMDB := make([]*IdxFile, 0)
	for i := range d.Files {
		f := d.Files[i]
		if mdbF, ok := mdbFiles[f.SHA1]; ok {
			var cuID int64
			if cuID, ok = cudMap[mdbF.ContentUnitID.Int64]; !ok {
				cuID = mdbF.ContentUnitID.Int64
			}
			cuIDs[cuID]++
		} else {
			notInMDB = append(notInMDB, f)
		}
	}

	fcIn := len(d.Files) - len(notInMDB)
	fcInPercent := float64(fcIn) / float64(len(d.Files))

	// Skip logging if all files in MDB and all of them are in the same CU (ID not 0)
	if _, ok := cuIDs[0]; !ok && len(cuIDs) == 1 && len(d.Files) == fcIn {
		return
	}

	log.Infof("%s %.2f %d/%d files in mdb %v", d.Name, fcInPercent, fcIn, len(d.Files), cuIDs)
	for i := range notInMDB {
		f := notInMDB[i]

		if kmF, ok := kmFiles[f.SHA1]; ok {
			log.Infof("\t\t%s\t%d", f.Name, kmF.ID)
		} else {
			log.Infof("\t\t%s", f.Name)
		}
	}
}
