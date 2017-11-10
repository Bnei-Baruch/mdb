package kmedia

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/hebcal"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
)

var knownKMContentTypes = map[int]string{
	1: api.CT_VIDEO_PROGRAM_CHAPTER,
	4: api.CT_LESSON_PART,
	5: api.CT_LECTURE,
}

type HolidayCollection struct {
	Holiday string
	Year    int
}

var HolidayCatalogs = map[int]HolidayCollection{
	8208: {Holiday: hebcal.H_CHANUKAH, Year: 2001},
	8210: {Holiday: hebcal.H_CHANUKAH, Year: 2002},
	8211: {Holiday: hebcal.H_CHANUKAH, Year: 2003},
	8212: {Holiday: hebcal.H_CHANUKAH, Year: 2004},
	8213: {Holiday: hebcal.H_CHANUKAH, Year: 2005},
	8214: {Holiday: hebcal.H_CHANUKAH, Year: 2006},
	8215: {Holiday: hebcal.H_CHANUKAH, Year: 2007},
	8216: {Holiday: hebcal.H_CHANUKAH, Year: 2008},
	8217: {Holiday: hebcal.H_CHANUKAH, Year: 2009},
	8218: {Holiday: hebcal.H_CHANUKAH, Year: 2010},
	8219: {Holiday: hebcal.H_CHANUKAH, Year: 2011},
	8220: {Holiday: hebcal.H_CHANUKAH, Year: 2012},
	8224: {Holiday: hebcal.H_CHANUKAH, Year: 2013},
	8225: {Holiday: hebcal.H_CHANUKAH, Year: 2014},
	8226: {Holiday: hebcal.H_CHANUKAH, Year: 2015},
	8227: {Holiday: hebcal.H_CHANUKAH, Year: 2016},

	8246: {Holiday: hebcal.H_PESACH, Year: 2001},
	8247: {Holiday: hebcal.H_PESACH, Year: 2002},
	8248: {Holiday: hebcal.H_PESACH, Year: 2003},
	8249: {Holiday: hebcal.H_PESACH, Year: 2004},
	8250: {Holiday: hebcal.H_PESACH, Year: 2005},
	8251: {Holiday: hebcal.H_PESACH, Year: 2006},
	8252: {Holiday: hebcal.H_PESACH, Year: 2007},
	8253: {Holiday: hebcal.H_PESACH, Year: 2008},
	8254: {Holiday: hebcal.H_PESACH, Year: 2009},
	8255: {Holiday: hebcal.H_PESACH, Year: 2010},
	8256: {Holiday: hebcal.H_PESACH, Year: 2011},
	8257: {Holiday: hebcal.H_PESACH, Year: 2012},
	8258: {Holiday: hebcal.H_PESACH, Year: 2013},
	8259: {Holiday: hebcal.H_PESACH, Year: 2014},
	8260: {Holiday: hebcal.H_PESACH, Year: 2015},
	8261: {Holiday: hebcal.H_PESACH, Year: 2016},
	8262: {Holiday: hebcal.H_PESACH, Year: 2017},

	8297: {Holiday: hebcal.H_PURIM, Year: 2001},
	8298: {Holiday: hebcal.H_PURIM, Year: 2002},
	8299: {Holiday: hebcal.H_PURIM, Year: 2003},
	8300: {Holiday: hebcal.H_PURIM, Year: 2004},
	8301: {Holiday: hebcal.H_PURIM, Year: 2014},
	8302: {Holiday: hebcal.H_PURIM, Year: 2015},
	8303: {Holiday: hebcal.H_PURIM, Year: 2016},
	8304: {Holiday: hebcal.H_PURIM, Year: 2017},

	8191: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2002},
	8190: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2003},
	8192: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2005},
	8193: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2006},
	8194: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2007},
	8195: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2008},
	8196: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2009},
	8197: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2010},
	8207: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2011},
	8198: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2012},
	8199: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2013},
	8200: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2014},
	8201: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2015},
	8202: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2016},
	8189: {Holiday: hebcal.H_ROSH_HASHANA, Year: 2017},

	8263: {Holiday: hebcal.H_SHAVUOT, Year: 2001},
	8264: {Holiday: hebcal.H_SHAVUOT, Year: 2002},
	8265: {Holiday: hebcal.H_SHAVUOT, Year: 2003},
	8272: {Holiday: hebcal.H_SHAVUOT, Year: 2004},
	8266: {Holiday: hebcal.H_SHAVUOT, Year: 2005},
	8267: {Holiday: hebcal.H_SHAVUOT, Year: 2006},
	8268: {Holiday: hebcal.H_SHAVUOT, Year: 2007},
	8269: {Holiday: hebcal.H_SHAVUOT, Year: 2008},
	8270: {Holiday: hebcal.H_SHAVUOT, Year: 2009},
	8271: {Holiday: hebcal.H_SHAVUOT, Year: 2010},
	8273: {Holiday: hebcal.H_SHAVUOT, Year: 2011},
	8274: {Holiday: hebcal.H_SHAVUOT, Year: 2012},
	8275: {Holiday: hebcal.H_SHAVUOT, Year: 2013},
	8276: {Holiday: hebcal.H_SHAVUOT, Year: 2014},
	8277: {Holiday: hebcal.H_SHAVUOT, Year: 2015},
	8278: {Holiday: hebcal.H_SHAVUOT, Year: 2016},
	8279: {Holiday: hebcal.H_SHAVUOT, Year: 2017},

	8230: {Holiday: hebcal.H_SUKKOT, Year: 2003},
	8242: {Holiday: hebcal.H_SUKKOT, Year: 2004},
	8232: {Holiday: hebcal.H_SUKKOT, Year: 2005},
	8233: {Holiday: hebcal.H_SUKKOT, Year: 2006},
	8234: {Holiday: hebcal.H_SUKKOT, Year: 2007},
	8241: {Holiday: hebcal.H_SUKKOT, Year: 2008},
	8235: {Holiday: hebcal.H_SUKKOT, Year: 2009},
	8236: {Holiday: hebcal.H_SUKKOT, Year: 2010},
	8243: {Holiday: hebcal.H_SUKKOT, Year: 2011},
	8244: {Holiday: hebcal.H_SUKKOT, Year: 2012},
	8237: {Holiday: hebcal.H_SUKKOT, Year: 2013},
	8238: {Holiday: hebcal.H_SUKKOT, Year: 2014},
	8239: {Holiday: hebcal.H_SUKKOT, Year: 2015},
	8240: {Holiday: hebcal.H_SUKKOT, Year: 2016},
	8228: {Holiday: hebcal.H_SUKKOT, Year: 2017},
}

var HolidayTags = map[string]string{
	hebcal.H_ROSH_HASHANA: "PkEfPB9i",
	hebcal.H_SUKKOT:       "Q2ZFsb9a",
	hebcal.H_CHANUKAH:     "rxNl0zXg",
	hebcal.H_PURIM:        "ZjqGWdYE",
	hebcal.H_PESACH:       "RWqjxgkj",
	hebcal.H_SHAVUOT:      "MyLcuAgH",
}

func ImportHolidays() {
	clock := Init()

	stats = NewImportStatistics()

	csMap, err := loadAndImportMissingHolidayCollections()
	utils.Must(err)

	utils.Must(importHolidaysContainers(csMap))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadAndImportMissingHolidayCollections() (map[int]*models.Collection, error) {
	hcal := new(hebcal.Hebcal)
	if err := hcal.Load(); err != nil {
		return nil, errors.Wrap(err, "Load hebcal")
	}

	cs, err := models.Collections(mdb,
		qm.Where("type_id = ?", api.CONTENT_TYPE_REGISTRY.ByName[api.CT_HOLIDAY].ID)).
		All()
	if err != nil {
		return nil, errors.Wrap(err, "Load collections")
	}

	csMap := make(map[int]*models.Collection)
	for i := range cs {
		c := cs[i]
		if !c.Properties.Valid {
			continue
		}

		var props map[string]interface{}
		if err := json.Unmarshal(c.Properties.JSON, &props); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal collection properties %d", c.ID)
		}

		if kmid, ok := props["kmedia_id"]; ok {
			csMap[int(kmid.(float64))] = c
		}
	}

	for k, v := range HolidayCatalogs {
		if _, ok := csMap[k]; ok {
			continue
		}

		start, end := hcal.GetPeriod(v.Holiday, v.Year)
		if start == "" || end == "" {
			return nil, errors.Errorf("Holiday missing period %s %d", v.Holiday, v.Year)
		}

		props := map[string]interface{}{
			"kmedia_id":   k,
			"holiday_tag": HolidayTags[v.Holiday],
			"start_date":  start,
			"end_date":    end,
		}
		log.Infof("Create collection %v", props)
		c, err := api.CreateCollection(mdb, api.CT_HOLIDAY, props)
		if err != nil {
			return nil, errors.Wrapf(err, "Create collection %v", props)
		}
		stats.CollectionsCreated.Inc(1)
		csMap[k] = c

		// I18n
		descriptions, err := kmodels.CatalogDescriptions(kmdb, qm.Where("catalog_id = ?", k)).All()
		if err != nil {
			return nil, errors.Wrapf(err, "Lookup catalog descriptions, [kmid %d]", k)
		}
		for _, d := range descriptions {
			if d.Name.Valid && d.Name.String != "" {
				ci18n := models.CollectionI18n{
					CollectionID: c.ID,
					Language:     api.LANG_MAP[d.LangID.String],
					Name:         d.Name,
				}
				err = ci18n.Upsert(mdb,
					true,
					[]string{"collection_id", "language"},
					[]string{"name"})
				if err != nil {
					return nil, errors.Wrapf(err, "Upsert collection i18n, collection [%d]", c.ID)
				}
			}
		}
	}

	return csMap, nil
}

func importHolidaysContainers(csMAp map[int]*models.Collection) error {
	cnMap, cuMap, err := loadContainersInCatalogsAndCUs(12)
	if err != nil {
		return errors.Wrap(err, "Load containers")
	}

	for cnID, cn := range cnMap {
		tx, err := mdb.Begin()
		utils.Must(err)

		// import container
		cu, ok := cuMap[cnID]
		if ok {
			// update - for tags
			cu, err = importContainer(tx, cn, nil, api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name, "", 0)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Import existing container %d", cnID)
			}
		} else {
			if ct, ok := knownKMContentTypes[cn.ContentTypeID.Int]; ok {
				// create - CU doesn't exist
				cu, err = importContainerWOCollectionNewCU(tx, cn, ct)
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "Import new container %d", cnID)
				}
			} else {
				utils.Must(tx.Rollback())
				continue
			}
		}

		// create or update CCUs
		for i := range cn.R.Catalogs {
			catalog := cn.R.Catalogs[i]
			c, ok := csMAp[catalog.ID]
			if !ok {
				continue
			}

			log.Infof("Associating %d %s to %d %s: [cu,c]=[%d,%d]", cn.ID, cn.Name.String, catalog.ID, catalog.Name, cu.ID, c.ID)

			if tx == nil {
				tx, err = mdb.Begin()
				utils.Must(err)
			}

			err = createOrUpdateCCU(tx, cu, models.CollectionsContentUnit{
				CollectionID:  c.ID,
				ContentUnitID: cu.ID,
			})
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Create or update CCU %d", cnID)
			}
		}

		if tx != nil {
			utils.Must(tx.Commit())
		}
	}

	return nil
}
