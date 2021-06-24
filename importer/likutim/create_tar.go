package likutim

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type CreateTar struct {
	duplicates []*Double
	mdb        *sql.DB
	inDir      string
}

func (c *CreateTar) Run() {
	err := c.duplicatesFromJSON()
	if err != nil {
		log.Fatalf("Error on read: %s", err)
	}

	c.openDB()
	defer c.mdb.Close()
	err = os.MkdirAll(viper.GetString("likutim.os-dir"), os.ModePerm)
	if err != nil {
		log.Errorf("Can't create directory: %s", err)
	}

	c.inDir = path.Join(viper.GetString("likutim.os-dir"), "likutimByCU")
	utils.Must(c.buildFolder())

	out := path.Join(viper.GetString("likutim.os-dir"), "fileByUIDLikutim.tar")
	cmd := exec.Command("tar", "-czvf", out, "-C", c.inDir, ".")
	r, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(r)
		return
	}
}

func (c *CreateTar) buildFolder() error {
	for _, d := range c.duplicates {
		f, err := models.Files(c.mdb,
			qm.Load("ContentUnit"),
			qm.Where("uid = ?", d.Save)).One()
		if err != nil {
			return err
		}
		url := "https://cdn.kabbalahmedia.info/" + f.UID
		pd := filepath.Join(c.inDir, f.R.ContentUnit.UID)
		if err := os.MkdirAll(pd, os.ModePerm); err != nil {
			return err
		}
		err = utils.DownloadUrl(url, filepath.Join(pd, f.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CreateTar) duplicatesFromJSON() error {
	p := path.Join(viper.GetString("likutim.os-dir"), "kitvei-makor-duplicates.json")
	j, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	var r []*Double
	err = json.Unmarshal(j, &r)
	if err != nil {
		return err
	}
	c.duplicates = r
	return nil
}

func (c *CreateTar) openDB() {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	c.mdb = mdb
}
