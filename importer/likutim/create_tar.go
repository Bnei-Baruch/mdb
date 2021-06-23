package likutim

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
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
	"path"
	"path/filepath"
	"strings"
)

type CreateTar struct {
	duplicates []Double
	mdb        *sql.DB
	baseDir    string
}

func (c *CreateTar) Run() {
	err := c.duplicatesFromJSON()
	if err != nil {
		log.Errorf("Error on read: %s", err)
	}

	c.openDB()
	defer c.mdb.Close()
	err = os.MkdirAll(viper.GetString("likutim.os-dir"), os.ModePerm)
	if err != nil {
		log.Errorf("Can't create directory: %s", err)
	}

	c.baseDir = path.Join(viper.GetString("likutim.os-dir"), "likutimByCU.tar")
	utils.Must(c.buildFolder())
	utils.Must(c.buildTar())
}

func (c CreateTar) buildFolder() error {
	for _, d := range c.duplicates {
		f, err := models.Files(c.mdb,
			qm.Load("ContentUnit"),
			qm.Where("uid = ?", d.Save)).One()
		if err != nil {
			return err
		}
		err = c.fetchFile(f)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c CreateTar) fetchFile(f *models.File) error {
	dirName := f.R.ContentUnit.UID
	p := filepath.Join(c.baseDir, dirName)
	err := os.MkdirAll(p, 0755)
	if err != nil {
		return err
	}

	resp, err := http.Get("https://files.kabbalahmedia.info/files/" + f.Name)
	if err != nil {
		log.Errorf("Not find file on link: %s, error: %s", f.Name, err)
		return err
	}
	file, err := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile(filepath.Join(p, f.Name), file, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (c CreateTar) buildTar() error {
	out := path.Join(viper.GetString("likutim.os-dir"), "fileByUIDLikutim.tar")

	var buf bytes.Buffer
	err := compress(c.baseDir, &buf)
	if err != nil {
		log.Errorf("Error on compress to tar. Error: %s", err)
		return err
	}

	f, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	if err != nil {
		log.Errorf("Error on create tar file. Error: %s", err)
		return err
	}
	if _, err := io.Copy(f, &buf); err != nil {
		log.Errorf("Error on insert data to the tar file. Error: %s", err)
		return err
	}
	return nil
}

func compress(src string, buf io.Writer) error {

	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			p := strings.Split(file, "/")
			header.Name = path.Join("/", p[len(p)-2], p[len(p)-1])
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CreateTar) duplicatesFromJSON() error {
	p := path.Join(viper.GetString("likutim.os-dir"), "kitvei-makor-duplicates.json")
	j, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	var r []Double
	err = json.Unmarshal(j, &r)
	if err != nil {
		return err
	}
	c.duplicates = r
	return nil
}

func (c *CreateTar) openDB() {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	c.mdb = mdb
}
