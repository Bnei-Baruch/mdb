package cusource

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ComparatorDbVsFolder struct {
}

type diffWithKeys struct {
	sha1 string
	sUid string
}

type diff struct {
	sha1            string
	sUid            string
	isOnDb          bool
	isOnFolder      bool
	isOnFileStorage bool
}

func (c *ComparatorDbVsFolder) Run() {
	dir := c.untar()
	mdb := c.openDB()
	c.uploadData(mdb, dir)
	utils.Must(os.RemoveAll(dir))
	utils.Must(mdb.Close())

}

//prepare upload

func (c *ComparatorDbVsFolder) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

func (c *ComparatorDbVsFolder) untar() string {
	path := viper.GetString("source-import.source-dir")

	r, err := os.Open(path)
	utils.Must(err)
	gzr, err := gzip.NewReader(r)
	utils.Must(err)
	defer utils.Must(gzr.Close())

	dir, _ := ioutil.TempDir(filepath.Dir(path), "un_tared")

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Must(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(filepath.Join(dir, header.Name), 0755)
			utils.Must(err)
		case tar.TypeReg:
			f, err := os.Create(filepath.Join(dir, header.Name))
			utils.Must(err)
			_, err = io.Copy(f, tr)
			utils.Must(err)
			err = f.Close()
			utils.Must(err)
		default:
			continue
		}
	}
	return dir
}

func (c *ComparatorDbVsFolder) uploadData(mdb *sql.DB, dir string) {
	units, err := models.ContentUnits(mdb,
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SOURCE].ID),
		qm.Load("Files"),
	).All()
	utils.Must(err)

	chDB := make(chan []diffWithKeys)
	chFStorage := make(chan []diffWithKeys)
	chFolder := make(chan []diffWithKeys)

	for _, u := range units {
		go func(unit models.ContentUnit) {
			db, fs := c.unitDataFromDB(unit)
			chDB <- db
			chFStorage <- fs
		}(*u)
	}

	filesOS, err := ioutil.ReadDir(dir)
	utils.Must(err)

	for _, sDir := range filesOS {
		go func() {
			d := c.unitDataFromFolder(sDir, dir)
			chFolder <- d
		}()
	}
	result := make(map[string]diff)
	for _, r := range <-chDB {
		key := r.sUid + "_" + r.sha1
		if v, ok := result[key]; ok {
			v.isOnDb = true
			result[key] = v
		} else {
			d := diff{sha1: r.sha1, sUid: r.sUid, isOnDb: true}
			result[key] = d
		}
	}

	for _, r := range <-chFStorage {
		key := r.sUid + "_" + r.sha1
		if v, ok := result[key]; ok {
			v.isOnFileStorage = true
			result[key] = v
		} else {
			d := diff{sha1: r.sha1, sUid: r.sUid, isOnFileStorage: true}
			result[key] = d
		}
	}

	for _, r := range <-chFolder {
		key := r.sUid + "_" + r.sha1
		if v, ok := result[key]; ok {
			v.isOnFolder = true
			result[key] = v
		} else {
			d := diff{sha1: r.sha1, sUid: r.sUid, isOnFolder: true}
			result[key] = d
		}
	}

	fmt.Print("\n\n******** Pint results ********\n\n")
	for _, d := range result {
		fmt.Printf("\nsha1: %s, sUid: %s, isOnDb: %t, isOnFolder: %t, isOnFileStorage: %t \n", d.sha1, d.sUid, d.isOnDb, d.isOnFolder, d.isOnFileStorage)
	}
	fmt.Print("\n\n******** End results ********\n\n")
}

func (c *ComparatorDbVsFolder) unitDataFromDB(unit models.ContentUnit) ([]diffWithKeys, []diffWithKeys) {
	sUid := unit.UID
	isOnDB := make([]diffWithKeys, 0)
	isOnFStorage := make([]diffWithKeys, 0)
	for _, f := range unit.R.Files {
		isOnDB = append(isOnDB, diffWithKeys{sha1: fmt.Sprintf("%x", f.Sha1.Bytes), sUid: sUid})

		var p map[string]string
		err := f.Properties.Unmarshal(&p)
		if err != nil {
			fmt.Printf("\nCU id - %s, file name - %s have no Properties. error: %s", sUid, f.Name, err.Error())
			continue
		}

		var url string
		if u, ok := p["url"]; !ok {
			fmt.Printf("\nCU id - %s, file name - %s have no url on Properties. error: %s", sUid, f.Name, err)
			continue
		} else {
			url = u
		}
		resp, err := http.Get(url)
		utils.Must(err)

		var b []byte
		_, err = resp.Body.Read(b)
		if err != nil {
			fmt.Printf("\nCU id - %s, file name - %s is not on FileStorage. error: %s", sUid, f.Name, err)
			continue
		}
		sha := sha1.Sum(b)
		isOnFStorage = append(isOnFStorage, diffWithKeys{sUid: sUid, sha1: fmt.Sprintf("%x", sha)})
		utils.Must(resp.Body.Close())
	}
	return isOnDB, isOnFStorage
}

func (c *ComparatorDbVsFolder) unitDataFromFolder(dir os.FileInfo, path string) []diffWithKeys {

	isOnFolder := make([]diffWithKeys, 0)
	if !dir.IsDir() {
		fmt.Printf("\nNot forder %s", dir.Name())
		return isOnFolder
	}
	sUid := dir.Name()
	dPath := filepath.Join(path, dir.Name())
	files, err := ioutil.ReadDir(dPath)
	utils.Must(err)

	for _, file := range files {
		spl := strings.Split(file.Name(), "/")
		name := spl[len(spl)-1]
		if isDoc := strings.Contains(name, ".docx"); !isDoc {
			continue
		}

		f, err := os.Open(filepath.Join(dPath, file.Name()))
		utils.Must(err)
		var b []byte
		_, err = f.Read(b)
		utils.Must(err)
		sha := sha1.Sum(b)
		utils.Must(f.Close())
		isOnFolder = append(isOnFolder, diffWithKeys{sUid: sUid, sha1: fmt.Sprintf("%x", sha)})

	}
	return isOnFolder
}
