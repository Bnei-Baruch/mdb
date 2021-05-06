package cusource

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
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
	result map[string]diff
}

type diff struct {
	sha1            string
	sUid            string
	isOnDb          bool
	isOnFolder      bool
	isOnFileStorage bool
}

func (c *ComparatorDbVsFolder) Run() {
	c.result = make(map[string]diff)

	dir := c.untar()
	mdb := c.openDB()
	c.uploadData(mdb, dir)
	utils.Must(os.RemoveAll(dir))
	utils.Must(mdb.Close())

}

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
		//qm.Limit(10),
	).All()
	utils.Must(err)
	for _, u := range units {
		c.unitDataFromDB(*u)
	}

	filesOS, err := ioutil.ReadDir(dir)
	utils.Must(err)

	for _, sDir := range filesOS {
		c.unitDataFromFolder(sDir, dir)
	}
	c.printResults()
	fmt.Print("\n******** End process ********\n")
}

func (c *ComparatorDbVsFolder) unitDataFromDB(unit models.ContentUnit) {
	sUid := unit.UID
	for _, f := range unit.R.Files {
		shaDb := hex.EncodeToString(f.Sha1.Bytes)
		kdb := sUid + "_" + shaDb
		if _, ok := c.result[kdb]; !ok {
			c.result[kdb] = diff{
				sha1:   shaDb,
				sUid:   sUid,
				isOnDb: true,
			}
		} else {
			t := c.result[kdb]
			t.isOnDb = true
			c.result[kdb] = t
		}

		//check if have file at Shirai
		var p map[string]string
		err := f.Properties.Unmarshal(&p)
		if err != nil {
			log.Debugf("CU id - %s, file name - %s have no Properties. error: %s", sUid, f.Name, err.Error())
			continue
		}

		url := "https://files.kabbalahmedia.info/files/" + f.Name
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("Not find file on link: %s, error: %s", url, err.Error())
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Debugf("CU id - %s, file name - %s is not on FileStorage. error: %s", sUid, f.Name, err)
			continue
		}
		utils.Must(resp.Body.Close())
		b1 := sha1.Sum(b)
		var shab = make(chan []byte, 1)
		shab <- b1[:]
		shaFS := hex.EncodeToString(<-shab)
		kfs := sUid + "_" + shaFS
		if _, ok := c.result[kfs]; !ok {
			c.result[kfs] = diff{
				sha1:            shaDb,
				sUid:            sUid,
				isOnFileStorage: true,
			}
		} else {
			t := c.result[kfs]
			t.isOnFileStorage = true
			c.result[kfs] = t
		}
	}
}

func (c *ComparatorDbVsFolder) unitDataFromFolder(dir os.FileInfo, path string) {
	if !dir.IsDir() {
		log.Errorf("Not folder %s", dir.Name())
		return
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
		b1 := sha1.Sum(b)
		var shab = make(chan []byte, 1)
		shab <- b1[:]
		sha := hex.EncodeToString(<-shab)
		utils.Must(f.Close())

		key := sUid + "_" + sha

		if _, ok := c.result[key]; !ok {
			c.result[key] = diff{
				sha1:       sha,
				sUid:       sUid,
				isOnFolder: true,
			}
		} else {
			t := c.result[key]
			t.isOnFolder = true
			c.result[key] = t
		}
	}
}

func (c *ComparatorDbVsFolder) printResults() {
	lines := []string{"sha1, sUid, isOnDb, isOnFileStorage, isOnFolder"}
	for _, d := range c.result {
		l := fmt.Sprintf("\n%s, %s, %t, %t, %t", d.sha1, d.sUid, d.isOnDb, d.isOnFileStorage, d.isOnFolder)
		lines = append(lines, l)
	}
	b := []byte(strings.Join(lines, ","))
	err := ioutil.WriteFile(viper.GetString("source-import.output"), b, 0644)
	utils.Must(err)
}
