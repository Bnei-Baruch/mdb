package likutim

import (
	"database/sql"
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"io/ioutil"
	"path"
	"strings"
)

type CreateUnits struct {
	duplicates []Double
	mdb        *sql.DB
}

func (c *CreateUnits) Run() {
	//err := c.duplicatesFromJSON()
	//if err != nil {
	//log.Errorf("Error on read", err)
	compare := new(Compare)
	compare.Run()
	c.duplicates = compare.result
	//}

	c.openDB()
	defer c.mdb.Close()

	c.fetchUnits()
}

func (c *CreateUnits) fetchUnits() {
	for _, d := range c.duplicates {
		fbase, err := models.Files(c.mdb,
			qm.Where("uid = ?", d.Save),
		).One()
		if err != nil {
			log.Errorf("cant find base file by uid: ", d.Save, err)
			continue
		}
		cukm, err := models.ContentUnits(c.mdb,
			qm.Where("id = ?", fbase.ContentUnitID.Int64),
			qm.Load("Files"),
		).One()
		if err != nil {
			log.Errorf("cant find base unit by file: %v ", fbase, err)
			continue
		}

		for _, uid := range d.Doubles {
			tx, err := c.mdb.Begin()
			if err != nil {
				log.Errorf("problem open transaction", err)
				tx.Rollback()
				continue
			}
			f, err := models.Files(c.mdb, qm.Where("uid = ?", uid)).One()
			if err != nil {
				log.Errorf("cant find file by uid: ", uid, err)
				tx.Rollback()
				continue
			}

			u, err := models.ContentUnits(c.mdb, qm.Where("id = ?", f.ContentUnitID.Int64)).One()
			if err != nil {
				log.Errorf("cant find unit by file: %v ", f, err)
				tx.Rollback()
				continue
			}

			cuo, err := FindOrigin(c.mdb, u.ID)
			if err != nil {
				log.Errorf("cant find origin unit by unit: %v ", u, err)
				tx.Rollback()
				continue
			}

			cu, err := c.createCU(cuo)
			if err != nil {
				log.Errorf("cant create unit %v to unit: %v ", cu, cuo, err)
				tx.Rollback()
				continue
			}
			c.moveFiles(cukm, &cu)

			d := &models.ContentUnitDerivation{
				SourceID:  cuo.ID,
				DerivedID: cu.ID,
				Name:      "test",
			}
			err = cuo.AddSourceContentUnitDerivations(c.mdb, true, d)
			if err != nil {
				log.Errorf("cant derive unit %v to unit: %v ", u, cuo, err)
				tx.Rollback()
				continue
			}
			err = tx.Commit()
			if err != nil {
				log.Errorf("problem commit transaction", err)
			}

		}
	}
}

func (c *CreateUnits) moveFiles(cukm, newu *models.ContentUnit) {
	for _, f := range cukm.R.Files {
		if !strings.Contains(f.Name, ".doc") {
			continue
		}
		f.ContentUnitID = null.Int64{Int64: newu.ID, Valid: true}
		err := f.Update(c.mdb, "content_unit_id")
		if err != nil {
			log.Errorf("Cant insert files from kitvei makor unit: %s  to new CU: %s", cukm.UID, newu.UID)
		}
	}
}

func (c *CreateUnits) createCU(cuo *models.ContentUnit) (models.ContentUnit, error) {
	cu := models.ContentUnit{
		UID:       utils.GenerateUID(8),
		TypeID:    common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID,
		Secure:    0,
		Published: true,
	}
	err := cu.Insert(c.mdb)
	if err != nil {
		log.Errorf("Cant add tags for CU id %d", cu.ID, err)
		return cu, err
	}

	//take data from origin for new unit
	err = cu.AddTags(c.mdb, false, cuo.R.Tags...)
	if err != nil {
		log.Errorf("Cant add tags for CU id %d", cu.ID, err)
		return cu, err
	}

	var i18n models.ContentUnitI18nSlice
	for _, i := range cuo.R.ContentUnitI18ns {
		n := &models.ContentUnitI18n{
			ContentUnitID: cu.ID,
			Language:      i.Language,
			Name:          i.Name,
		}
		i18n = append(i18n, n)
	}
	err = cu.AddContentUnitI18ns(c.mdb, false, i18n...)
	if err != nil {
		log.Errorf("Cant add i18n for CU id %d", cu.ID, err)
		return cu, err
	}
	return cu, nil
}

func (c *CreateUnits) duplicatesFromJSON() error {
	p := path.Join(viper.GetString("source-import.results-dir"), "kitvei-makor-duplicates.json")
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

func (c *CreateUnits) openDB() {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	c.mdb = mdb
}
