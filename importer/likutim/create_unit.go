package likutim

import (
	"database/sql"
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type CreateUnits struct {
	duplicates []*Double
	mdb        *sql.DB
}

func (c *CreateUnits) Run() {
	err := c.duplicatesFromJSON()
	if err != nil {
		log.Errorf("Error on read. Error: %s", err)
		compare := new(Compare)
		compare.Run()
		c.duplicates = compare.result
	}

	c.openDB()
	defer c.mdb.Close()

	err = os.MkdirAll(viper.GetString("likutim.results-dir"), os.ModePerm)
	if err != nil {
		log.Errorf("Can't create directory: %s", err)
	}
	c.fetchUnits()
}

func (c *CreateUnits) fetchUnits() {
	for _, d := range c.duplicates {
		log.Debugf("\n\nStart create new unit type LIKUTIM, file: %s all files: %v", d.Save, d.Doubles)
		fbase, err := models.Files(c.mdb,
			qm.Where("uid = ?", d.Save),
		).One()
		if err != nil {
			log.Errorf("Can't find base file by uid: %s. Error: %s", d.Save, err)
			continue
		}
		cukn, err := models.ContentUnits(c.mdb,
			qm.Where("id = ?", fbase.ContentUnitID.Int64),
			qm.Load("Files"),
		).One()
		if err != nil {
			log.Errorf("Can't find base unit by file: %v. Error: %s", fbase, err)
			continue
		}

		for _, uid := range d.Doubles {
			tx, err := c.mdb.Begin()
			if err != nil {
				log.Errorf("problem open transaction. Error: %s", err)
				continue
			}
			f, err := models.Files(tx, qm.Where("uid = ?",
				qm.Load("ContentUnit", "ContentUnit.DerivedContentUnitDerivations", "ContentUnit.DerivedContentUnitDerivations.Source"),
				uid)).One()
			if err != nil {
				log.Errorf("Can't find file by uid: %s. Error: %v", uid, err)
				tx.Rollback()
				continue
			}

			if f.R.ContentUnit == nil {
				log.Errorf("Can't find unit by file: %v. Error: %s", f, err)
				tx.Rollback()
				continue
			}

			if len(f.R.ContentUnit.R.DerivedContentUnitDerivations) == 0 {
				log.Errorf("Can't find origin unit by unit: %v. Error: %s", f.R.ContentUnit, err)
				tx.Rollback()
				continue
			}

			cuo, err := models.ContentUnits(tx,
				qm.Where("id = ?", f.R.ContentUnit.R.DerivedContentUnitDerivations[0].SourceID),
				qm.Load("ContentUnitI18ns", "Tags"),
			).One()
			if err != nil {
				log.Errorf("Can't find origin unit by unit: %v. Error: %s", f.R.ContentUnit, err)
				tx.Rollback()
				continue
			}

			cu, err := c.createCU(cuo)
			if err != nil {
				log.Errorf("Can't create unit %v to unit: %v. Error: %s", cu, cuo, err)
				tx.Rollback()
				continue
			}
			if err := c.moveFiles(cukn, &cu); err != nil {
				log.Error(err)
				tx.Rollback()
				continue
			}

			d := &models.ContentUnitDerivation{
				SourceID:  cuo.ID,
				DerivedID: cu.ID,
			}
			err = cuo.AddSourceContentUnitDerivations(tx, true, d)
			if err != nil {
				log.Errorf("Can't derive unit %v to unit: %v. Error: %s", cu, cuo, err)
				tx.Rollback()
				continue
			}
			err = tx.Commit()
			if err != nil {
				log.Errorf("problem commit transaction. Error: %s", err)
			}
			log.Debugf("End create new unit type LIKUTIM id: %d", cu.ID)
		}
	}
}

func (c *CreateUnits) moveFiles(cukm, newu *models.ContentUnit) error {
	for _, f := range cukm.R.Files {
		if !strings.Contains(f.Name, ".doc") {
			continue
		}
		f.ContentUnitID = null.Int64{Int64: newu.ID, Valid: true}
		if err := f.Update(c.mdb, "content_unit_id"); err != nil {
			return errors.Wrapf(err, "Can't insert files from kitvei makor unit: %s  to new CU: %s", cukm.UID, newu.UID)
		}
	}
	return nil
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
		log.Errorf("Can't add tags for CU id %d. Error: %s", cu.ID, err)
		return cu, err
	}

	//take data from origin for new unit
	err = cu.AddTags(c.mdb, false, cuo.R.Tags...)
	if err != nil {
		log.Errorf("Can't add tags for CU id %d. Error: %s", cu.ID, err)
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
	err = cu.AddContentUnitI18ns(c.mdb, true, i18n...)
	if err != nil {
		log.Errorf("Can't add i18n for CU id %d. Error: %s", cu.ID, err)
		return cu, err
	}
	return cu, nil
}

func (c *CreateUnits) duplicatesFromJSON() error {
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

func (c *CreateUnits) openDB() {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	c.mdb = mdb
}
