package likutim

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type PrintWithDoc struct{}

type printData struct {
	cuUid    string
	oName    string
	filmDate time.Time
	topic    string
}

func (c *PrintWithDoc) Run() {
	mdb := c.openDB()
	defer mdb.Close()

	err := os.MkdirAll(viper.GetString("likutim.results-dir"), os.ModePerm)
	if err != nil {
		log.Errorf("Can't create directory: %s", err)
	}
	dir, err := ioutil.TempDir(viper.GetString("likutim.results-dir"), "temp_kitvei_makor")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cus, err := models.ContentUnits(mdb,
		qm.Select("distinct on (\"content_units\".id) \"content_units\".*"),
		//qm.InnerJoin("files f ON f.content_unit_id = \"content_units\".id"),
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_KITEI_MAKOR].ID),
		qm.Load("Files", "Tags", "Tags.TagI18ns", "DerivedContentUnitDerivations", "DerivedContentUnitDerivations.Source", "DerivedContentUnitDerivations.Source.ContentUnitI18ns"),
	).All()
	if err != nil {
		log.Errorf("can't load units. Error: %s", err)
	}

	forPrint := make([]printData, 0)

	for _, cu := range cus {
		if len(cu.R.DerivedContentUnitDerivations) == 0 {
			log.Errorf("Can't find origin unit by unit: %v. Error: %s", cu, err)
			continue
		}

		if err != nil {
			log.Errorf("Error on looking for origins for unit Error: %s", err)
		}
		forPrint = append(forPrint, c.prepareForPrint(cu, cu.R.DerivedContentUnitDerivations[0].R.Source))
	}
	c.printToCSV(forPrint)
}

func (c *PrintWithDoc) prepareForPrint(cu, cuo *models.ContentUnit) printData {
	var cuProps map[string]interface{}
	err := json.Unmarshal(cu.Properties.JSON, &cuProps)
	if err != nil {
		log.Errorf("json.Unmarshal cu properties %d, Error: %s", cu.ID, err)
	}
	var film time.Time
	if cuProps != nil {
		film, err = time.Parse("2006-01-02", cuProps["film_date"].(string))
		if err != nil {
			log.Errorf("time.Parse cu %d film_date %s", cu.ID, cuProps["film_date"])
		}
	}

	name := null.String{}
	if cuo != nil {
		for _, o := range cuo.R.ContentUnitI18ns {
			if o.Language == "he" {
				name = o.Name
			}
		}
	}

	label := null.String{}
	for _, t := range cu.R.Tags {
		for _, ttr := range t.R.TagI18ns {
			if ttr.Language == "he" {
				label = ttr.Label
			}
		}
	}

	return printData{
		cuUid:    cu.UID,
		oName:    name.String,
		filmDate: film,
		topic:    label.String,
	}
}

func (c *PrintWithDoc) printToCSV(data []printData) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].filmDate.Before(data[j].filmDate)
	})

	lines := []string{"Unit UID", "Original Name", "film date", "topic name"}
	for _, d := range data {
		l := fmt.Sprintf("\n%s, %s, %s, %s", d.cuUid, strings.ReplaceAll(d.oName, ",", "->"), d.filmDate.String(), d.topic)
		lines = append(lines, l)
	}
	b := []byte(strings.Join(lines, ","))
	p := path.Join(viper.GetString("likutim.os-dir"), "move-kitvei-makor.csv")
	err := ioutil.WriteFile(p, b, 0644)
	utils.Must(err)
}

func (c *PrintWithDoc) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}
