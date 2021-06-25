package likutim

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
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
	c.updateUnits()
}

func (c *CreateUnits) updateUnits() {
	for _, d := range c.duplicates {
		log.Debugf("\n\nStart create new unit type LIKUTIM, file: %s all files: %v", d.Save, d.Doubles)

		tx, err := c.mdb.Begin()
		if err != nil {
			log.Error(err)
			continue
		}
		err, cu := c.createUnitOnTransaction(tx, d.Save)
		if err != nil {
			log.Error(err)
			tx.Rollback()
			continue
		}
		err = tx.Commit()
		// Derivation unit type LIKITIM to origins of doubles
		for _, uid := range d.Doubles {
			if uid == d.Save {
				continue
			}
			tx, err := c.mdb.Begin()
			if err != nil {
				log.Error(err)
				tx.Rollback()
				continue
			}
			err = moveDataOnTransaction(tx, uid, cu)
			err = tx.Commit()
			if err != nil {
				log.Errorf("problem commit transaction. Error: %s", err)
			}
			log.Debugf("End create new unit type LIKUTIM id: %d", cu.ID)
		}
	}
}

func (c *CreateUnits) createUnitOnTransaction(tx *sql.Tx, uid string) (error, *models.ContentUnit) {
	//fetch files and get origin unit
	f, err := models.Files(tx,
		qm.Where("uid = ?", uid),
		qm.Load("ContentUnit", "ContentUnit.Files", "ContentUnit.DerivedContentUnitDerivations", "ContentUnit.DerivedContentUnitDerivations.Source"),
	).One()
	if err != nil {
		return errors.Wrapf(err, "Can't find base file by uid: %s.", uid), nil
	}
	if f.R.ContentUnit == nil {
		return errors.Wrapf(err, "Can't find base unit by file: %v.", f), nil
	}

	if len(f.R.ContentUnit.R.DerivedContentUnitDerivations) == 0 {
		return errors.Wrapf(err, "Can't find origin unit by unit: %v.", f.R.ContentUnit), nil
	}

	cuo, err := models.ContentUnits(tx,
		qm.Where("id = ?", f.R.ContentUnit.R.DerivedContentUnitDerivations[0].SourceID),
		qm.Load("ContentUnitI18ns", "Tags"),
	).One()
	if err != nil {
		return errors.Wrapf(err, "Can't find origin unit by unit: %v.", f.R.ContentUnit), nil
	}
	//create unit type LIKUTIM
	cu, err := c.createCU(tx, cuo)
	if err != nil {
		return errors.Wrapf(err, "Can't create unit %v to unit: %v.", cu, cuo), nil
	}
	if err := c.moveFiles(tx, f.R.ContentUnit, &cu); err != nil {
		return err, nil
	}
	//derivation to origin
	cud := &models.ContentUnitDerivation{
		SourceID:  cuo.ID,
		DerivedID: cu.ID,
	}
	err = cuo.AddSourceContentUnitDerivations(tx, true, cud)
	if err != nil {
		return errors.Wrapf(err, "Can't derive unit %v to unit: %v.", cu, cuo), nil
	}
	return nil, &cu
}

func moveDataOnTransaction(tx *sql.Tx, uid string, cu *models.ContentUnit) error {

	f, err := models.Files(tx, qm.Where("uid = ?", uid),
		qm.Load("ContentUnit", "ContentUnit.DerivedContentUnitDerivations", "ContentUnit.DerivedContentUnitDerivations.Source")).One()
	if err != nil {
		return errors.Wrapf(err, "Can't find file by uid: %s.", uid)
	}

	if f.R.ContentUnit == nil {
		return errors.Wrapf(err, "Can't find unit by file: %v.", f)
	}

	if len(f.R.ContentUnit.R.DerivedContentUnitDerivations) == 0 {
		return errors.Wrapf(err, "Can't find origin unit by unit: %v.", f.R.ContentUnit)
	}

	d := &models.ContentUnitDerivation{
		SourceID:  f.R.ContentUnit.R.DerivedContentUnitDerivations[0].SourceID,
		DerivedID: cu.ID,
	}
	err = cu.AddDerivedContentUnitDerivations(tx, true, d)
	if err != nil {
		return errors.Wrapf(err, "Can't derive unit %v to unit id: %d.", cu, f.R.ContentUnit.R.DerivedContentUnitDerivations[0].SourceID)
	}
	return nil
}

func (c *CreateUnits) moveFiles(tx *sql.Tx, cukm, newu *models.ContentUnit) error {
	for _, f := range cukm.R.Files {
		if !strings.Contains(f.Name, ".doc") {
			continue
		}
		f.ContentUnitID = null.Int64{Int64: newu.ID, Valid: true}
		if err := f.Update(tx, "content_unit_id"); err != nil {
			return errors.Wrapf(err, "Can't insert files from kitvei makor unit: %s  to new CU: %s", cukm.UID, newu.UID)
		}
	}
	return nil
}

func (c *CreateUnits) createCU(tx *sql.Tx, cuo *models.ContentUnit) (models.ContentUnit, error) {
	cu := models.ContentUnit{
		UID:       utils.GenerateUID(8),
		TypeID:    common.CONTENT_TYPE_REGISTRY.ByName[common.CT_LIKUTIM].ID,
		Secure:    0,
		Published: true,
	}
	err := cu.Insert(tx)
	if err != nil {
		log.Errorf("Can't add tags for CU id %d. Error: %s", cu.ID, err)
		return cu, err
	}
	log.Debugf("Unit was inserted CU %v", cu)

	//take data from origin for new unit
	err = cu.AddTags(tx, false, cuo.R.Tags...)
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
	err = cu.AddContentUnitI18ns(tx, true, i18n...)
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
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	c.mdb = mdb
}
