package batch

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"encoding/json"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"strconv"
)

func EventsSubcollections() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting to organize events subcollections")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	utils.Must(doEventsSubcollections([]int64{10999, 11091, 11273}))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doEventsSubcollections(cIDs []int64) error {
	for _, cID := range cIDs {

		c, err := models.CollectionsG(
			qm.Where("id=?", cID),
			qm.Load("CollectionsContentUnits",
				"CollectionsContentUnits.ContentUnit",
				"CollectionsContentUnits.ContentUnit.ContentUnitI18ns"),
		).One()
		if err != nil {
			return errors.Wrapf(err, "Load collection %d", cID)
		}

		log.Infof("Collection %d [%d unit]", cID, len(c.R.CollectionsContentUnits))

		var cProps map[string]interface{}
		err = json.Unmarshal(c.Properties.JSON, &cProps)
		if err != nil {
			return errors.Wrapf(err, "json.Unmarshal collection properties %d", cID)
		}

		start, err := time.Parse("2006-01-02", cProps["start_date"].(string))
		if err != nil {
			return errors.Wrapf(err, "time.Parse start_date %s", cProps["start_date"])
		}
		end, err := time.Parse("2006-01-02", cProps["end_date"].(string))
		if err != nil {
			return errors.Wrapf(err, "time.Parse end_date %s", cProps["end_date"])
		}

		cuByCaptureID := make(map[string][]*models.ContentUnit)
		for _, ccu := range c.R.CollectionsContentUnits {
			cu := ccu.R.ContentUnit

			// filter all but LESSON_PART
			if cu.TypeID != 11 {
				continue
			}

			// filter related to congress units
			var cuProps map[string]interface{}
			err := json.Unmarshal(cu.Properties.JSON, &cuProps)
			if err != nil {
				return errors.Wrapf(err, "json.Unmarshal cu properties %d", cu.ID)
			}

			film, err := time.Parse("2006-01-02", cuProps["film_date"].(string))
			if err != nil {
				return errors.Wrapf(err, "time.Parse cu %d film_date %s", cu.ID, cProps["film_date"])
			}

			if film.Before(start) || film.After(end) {
				log.Infof("Skipping cu %d %s not in [%s-%s]", cu.ID,
					film.Format("2006-01-02"), start.Format("2006-01-02"), end.Format("2006-01-02"))
				continue
			}

			err = cu.L.LoadFiles(boil.GetDB(), true, cu)
			if err != nil {
				return errors.Wrapf(err, "Load CU files %d", cu.ID)
			}

			for _, f := range cu.R.Files {
				if f.Type == "video" {
					op, err := api.FindUpChainOperation(boil.GetDB(), f.ID, api.OP_CAPTURE_STOP)
					if err != nil {
						return errors.Wrapf(err, "find upchain op for file %d", f.ID)
					}

					if op.Properties.Valid {
						var oProps map[string]interface{}
						err = json.Unmarshal(op.Properties.JSON, &oProps)
						if err != nil {
							return errors.Wrapf(err, "json Unmarshal op properties %d", op.ID)
						}
						captureID, ok := oProps["collection_uid"]
						if ok {
							k := captureID.(string)
							v, ok := cuByCaptureID[k]
							if !ok {
								v = make([]*models.ContentUnit, 0)
							}
							cuByCaptureID[k] = append(v, cu)
						} else {
							log.Warnf("op has no collection_uid property")
						}
					}
					break
				}
			}
		}

		log.Infof("len(cuByCaptureID) %d", len(cuByCaptureID))
		for k, v := range cuByCaptureID {
			// see if we already have this collection
			cc, err := api.FindCollectionByCaptureID(boil.GetDB(), k)
			if err != nil {
				if _, ok := err.(api.CollectionNotFound); !ok {
					return errors.Wrapf(err, "FindCollectionByCaptureID %s", k)
				}
			}

			if cc != nil {
				log.Infof("capture_id %s collection exist \t%d\t%d", k, cc.ID, cc.TypeID)
				continue
			}

			cu := v[0]
			var props map[string]interface{}
			if err := json.Unmarshal(cu.Properties.JSON, &props); err != nil {
				return errors.Wrapf(err, "json.Unmarshal cu props %d", cu.ID)
			}
			delete(props, "duration")
			delete(props, "kmedia_id")

			captureDate, err := time.Parse("2006-01-02", props["capture_date"].(string))
			if err != nil {
				return errors.Wrapf(err, "time.Parse cu %d capture_date %s", cu.ID, props["capture_date"])
			}

			cct := api.CT_DAILY_LESSON
			if captureDate.Weekday() == time.Saturday {
				cct = api.CT_SPECIAL_LESSON
			}

			tx, err := boil.GetDB().(*sql.DB).Begin()
			utils.Must(err)

			log.Infof("Creating collection %s %v", cct, props)
			c, err = api.CreateCollection(tx, cct, props)
			if err != nil {
				utils.Must(tx.Rollback())
				return err
			}
			log.Infof("Created collection %d", c.ID)

			c.Published = true
			if err := c.Update(tx, "published"); err != nil {
				return errors.Wrapf(err, "update collection published cID %d", c.ID)
			}

			for _, cu := range v {
				var name string
				for _, i18n := range cu.R.ContentUnitI18ns {
					if i18n.Language == api.LANG_HEBREW {
						name = i18n.Name.String
						break
					}
				}

				if err := cu.L.LoadCollectionsContentUnits(tx, true, cu); err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "Load CCU's for cu %d", cu.ID)
				}

				ccu := cu.R.CollectionsContentUnits[0]
				ccuName := ccu.Name
				position, err := strconv.Atoi(ccuName)
				if err != nil {
					log.Errorf("strconv.Atoi(ccuName) cu id %d", cu.ID)
				}

				log.Infof("%d\t%d\t%s\t%s\t%d", cu.ID, cu.TypeID, name, ccuName, position)

				nccu := &models.CollectionsContentUnit{
					CollectionID:  c.ID,
					ContentUnitID: cu.ID,
					Name:          ccuName,
					Position:      position,
				}
				err = c.AddCollectionsContentUnits(tx, true, nccu)
				if err != nil {
					return errors.Wrapf(err, "Save ccu in DB cID %d cuID %d", c.ID, cu.ID)
				}
			}

			utils.Must(tx.Commit())
		}
	}

	return nil
}
