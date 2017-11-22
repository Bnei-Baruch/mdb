package roza

import (
	"fmt"
	"os"
	"time"
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries"

	"github.com/Bnei-Baruch/mdb/utils"
)

type UploadFile struct {
	IdxFile  *IdxFile
	KmCnID   int
	KmCnName string
}

func PrepareUpoad() {
	clock := Init()

	idx := new(RozaIndex)
	utils.Must(idx.Load(mdb))

	kmFiles, err := loadKMFiles()
	utils.Must(err)

	kmContainers, err := loadKMContainers()
	utils.Must(err)

	err = prepareUpload(idx, kmFiles, kmContainers)
	utils.Must(err)

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadKMContainers() (map[int]string, error) {
	rows, err := queries.Raw(kmdb, "select id, name from containers").Query()
	if err != nil {
		return nil, errors.Wrap(err, "Load Kmedia containers")
	}
	defer rows.Close()

	cns := make(map[int]string, 60000)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		cns[id] = name
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return cns, nil
}

func prepareUpload(
	idx *RozaIndex,
	kmFiles map[string]*MiniKMFile,
	kmContainers map[int]string) error {

	beavoda := idx.GetDir("/vfs/archive/Archive/____beavoda")
	if beavoda == nil {
		return errors.New("____beavoda not found")
	}

	forUpload := make(map[string][]*UploadFile, 50000)

	s := []*IdxDirectory{beavoda}
	var x *IdxDirectory
	for len(s) > 0 {
		x, s = s[0], s[1:]

		if len(x.Files) > 0 {
			cnIDs := hashset.New()
			missing := make([]*UploadFile, 0)
			for i := range x.Files {
				f := x.Files[i]
				if skipRE.MatchString(f.Name) {
					continue
				}

				if kmF, ok := kmFiles[f.SHA1]; ok {
					cnIDs.Add(utils.ConvertArgsInt(kmF.CnIDs...)...)
				} else {
					missing = append(missing, &UploadFile{IdxFile: f})
				}
			}

			if len(missing) > 0 && cnIDs.Size() == 1 {
				for _, cnID := range cnIDs.Values() {
					if cnName, ok := kmContainers[cnID.(int)]; ok {
						for i := range missing {
							y := &UploadFile{IdxFile: missing[i].IdxFile}
							y.KmCnID = cnID.(int)
							y.KmCnName = cnName

							k := y.IdxFile.SHA1
							v, ok := forUpload[k]
							if !ok {
								v = make([]*UploadFile, 0)
							}
							forUpload[k] = append(v, y)
						}
					} else {
						log.Infof("Unknown KmCnID %d: %s", cnID.(int), x.path())
					}
				}
			}
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

	// dedup each file by cnID
	finalFiles := make([]*UploadFile, 0)
	for _, v := range forUpload {
		cnIDs := hashset.New()
		for i := range v {
			f := v[i]
			if !cnIDs.Contains(f.KmCnID) {
				finalFiles = append(finalFiles, v[0])
				cnIDs.Add(f.KmCnID)
			}
		}
	}

	sort.Slice(finalFiles, func(i, j int) bool {
		return finalFiles[i].KmCnID < finalFiles[j].KmCnID
	})

	var totalBytes int64
	out, err := os.OpenFile("importer/roza/analysis/upload.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "os.OpenFile")
	}
	defer out.Close()

	for i := range finalFiles {
		f := finalFiles[i]
		_, err := fmt.Fprintf(out, "%d,\"%s\",\"%s/%s\"\n", f.KmCnID, f.KmCnName, f.IdxFile.Directory.path(), f.IdxFile.Name)
		if err != nil {
			return errors.Wrapf(err, "fmt.Fprintf %s/%s\t%s", f.IdxFile.Directory.path(), f.IdxFile.Name, f.KmCnName)
		}
		totalBytes += f.IdxFile.Size
	}
	log.Infof("%d files for upload [%d bytes]", len(finalFiles), totalBytes)

	return nil
}
