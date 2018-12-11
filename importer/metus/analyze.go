package metus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/mdb/utils"
	"sort"
	"strings"
)

func Analyze() {
	clock := Init()

	utils.Must(doAnalyze())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doAnalyze() error {
	fh := new(FieldsHelper)
	if err := fh.Load(); err != nil {
		return errors.Wrap(err, "FieldsHelper.Load")
	}
	fh.dump()

	//return dbToDisk()
	return loadDisk()
}

func loadDisk() error {
	oh := new(ObjectsHelper)
	if err := oh.LoadFromDisk("importer/metus/data"); err != nil {
		return errors.Wrap(err, "ObjectsHelper.LoadFromDisk")
	}
	log.Infof("Loaded %d objects from disk. roots: %d", len(oh.byID), len(oh.roots))

	fileObjects := make([]*Object, 0)
	err := oh.WalkDFS(func(o *Object) error {
		if o.Type == 4 {
			fileObjects = append(fileObjects, o)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "ObjectHelper.Walk")
	}

	log.Infof("%d file objects (type 4)", len(fileObjects))
	sort.Slice(fileObjects, func(i, j int) bool {
		return fileObjects[i].getPhysicalFilename() < fileObjects[j].getPhysicalFilename()
	})

	inFilteredBins := 0
	for i := range fileObjects {
		f := fileObjects[i]
		idPath := oh.getIDPath(f.ID)
		for j := range idPath {
			if idPath[j] == "1011917" || // Archive Special
				idPath[j] == "1026289" { // Ligdol Bins
				f.IsInFilteredBin = true
				inFilteredBins++
				break
			}
		}
	}

	pi := new(PhysicalIndex)
	if err := pi.Load(); err != nil {
		return errors.Wrap(err, "PhysicalIndex.Load")
	}
	pi.Match(fileObjects)

	fm := new(FilenameMatcher)
	if err := fm.Load(); err != nil {
		return errors.Wrap(err, "FilenameMatcher.Load")
	}

	match := make([]*Object, 0)
	matchNoCU := make([]*Object, 0)
	noMatch := make([]*Object, 0)
	byCU := make(map[int64][]*Object)
	for i := range fileObjects {
		f := fileObjects[i]
		name := f.getPhysicalFilename()
		//log.Infof("%s\t%d",name, f.ID)

		cuID := fm.Match(name)
		f.matchedCU = cuID
		if cuID > 0 {
			match = append(match, f)
			byCU[cuID] = append(byCU[cuID], f)
		} else if cuID == 0 {
			matchNoCU = append(matchNoCU, f)
		} else {
			noMatch = append(noMatch, f)
		}
	}

	log.Infof("\n\n\n\n")
	log.Infof("match %d, len(byCU) %d, matchNoCU %d, noMatch %d", len(match), len(byCU), len(matchNoCU), len(noMatch))
	//log.Infof("\n\nmatch")
	//for i := range match {
	//	log.Infof("%s\t%d", match[i].getPhysicalFilename(), match[i].ID)
	//}
	//log.Infof("\n\nmatchNoCU")
	//for i := range matchNoCU {
	//	log.Infof("%s\t%d", matchNoCU[i].getPhysicalFilename(), matchNoCU[i].ID)
	//}
	//log.Infof("\n\nnoMatch")
	//for i := range noMatch {
	//	log.Infof("%s\t%d", noMatch[i].getPhysicalFilename(), noMatch[i].ID)
	//}

	//byCUKeys := make([]int64, 0)
	//for k, _ := range byCU {
	//	byCUKeys = append(byCUKeys, k)
	//}
	//sort.Slice(byCUKeys, func(i, j int) bool {
	//	return len(byCU[byCUKeys[i]]) < len(byCU[byCUKeys[j]])
	//})
	//for i := range byCUKeys {
	//	v := byCU[byCUKeys[i]]
	//	log.Infof("CU %d\t%d", byCUKeys[i], len(v))
	//	for j := range v {
	//		log.Infof("\t%s\t%d", v[j].getPhysicalFilename(), v[j].ID)
	//	}
	//}

	var phys, dup, missing int
	for i := range fileObjects {
		f := fileObjects[i]
		if f.FileRecord != nil {
			if f.IsDuplicate {
				dup++
			} else {
				phys++
			}
		} else {
			missing++
		}
	}
	log.Infof("phys %d\tdup %d\t missing %d", phys, dup, missing)

	var matchPhys, matchNoPhys, noMatchPhys, noMatchNoPhys int
	for i := range match {
		f := match[i]
		if f.FileRecord != nil {
			matchPhys++
		} else {
			matchNoPhys++
		}
	}
	for i := range noMatch {
		f := noMatch[i]
		if f.FileRecord != nil {
			noMatchPhys++
		} else {
			noMatchNoPhys++
		}
	}
	log.Infof("matchPhys %d\tmatchNoPhys %d\t noMatchPhys %d\t noMatchNoPhys %d", matchPhys, matchNoPhys, noMatchPhys, noMatchNoPhys)

	log.Infof("inFilteredBins %d", inFilteredBins)
	ifbWithCU := make([]*Object, 0)
	ifbWithoutCU := make([]*Object, 0)
	for i := range fileObjects {
		f := fileObjects[i]
		if !f.IsInFilteredBin {
			continue
		}
		if f.matchedCU > 0 {
			ifbWithCU = append(ifbWithCU, f)
			//log.Infof("%s\t=>\t%s", strings.Join(oh.getNamePath(f.ID), "/"), f.getPhysicalFilename())
		} else {
			ifbWithoutCU = append(ifbWithoutCU, f)
		}
	}

	log.Infof("ifbWithCU %d, ifbWithoutCU %d", len(ifbWithCU), len(ifbWithoutCU))
	for i := range ifbWithCU {
		f := ifbWithCU[i]
		log.Infof("%s\t=>\t%s", strings.Join(oh.getNamePath(f.ID), "/"), f.getPhysicalFilename())
	}
	for i := range ifbWithoutCU {
		f := ifbWithoutCU[i]
		log.Infof("%s\t=>\t%s", strings.Join(oh.getNamePath(f.ID), "/"), f.getPhysicalFilename())
	}

	return nil
}

func dbToDisk() error {
	fh := new(FieldsHelper)
	if err := fh.Load(); err != nil {
		return errors.Wrap(err, "FieldsHelper.Load")
	}

	// load OBJECTS from DB
	oh := new(ObjectsHelper)
	if err := oh.LoadFromDB(); err != nil {
		return errors.Wrap(err, "ObjectsHelper.LoadFromDB")
	}

	fh.nonMissingMetadata = 0
	for _, o := range oh.byID {
		o.MetadataJson = fh.getMetadataAsJson(o)
	}
	log.Infof("fh.nonMissingMetadata %d", fh.nonMissingMetadata)

	if err := oh.WalkBFS(func(o *Object) error {
		if err := saveToFile(o, oh); err != nil {
			return errors.Wrapf(err, "saveToFile [%d]", o.ID)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "ObjectsHelper.Walk")
	}

	return nil
}

func saveToFile(o *Object, oh *ObjectsHelper) error {
	dirPath := oh.getIDPath(o.ID)
	dir := fmt.Sprintf("importer/metus/data/%s", filepath.Join(dirPath...))
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return errors.Wrapf(err, "os.MkdirAll %s", dir)
	}

	// create output file
	idStr := strconv.Itoa(o.ID)
	path := fmt.Sprintf("%s/%s.json", dir, idStr)
	output, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "os.Create %s", path)
	}
	defer output.Close()

	err = json.NewEncoder(output).Encode(o)
	if err != nil {
		return errors.Wrap(err, "json.Encode")
	}

	return nil
}
