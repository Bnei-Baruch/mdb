package batch

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func OrganizeKiteiMakor() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting to organize kitei makor")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	log.Info("Loading kitei-makor files")
	ktFiles, err := models.Files(mdb,
		qm.Where("name ~ ?", "kitei-makor"),
		qm.Load("ContentUnit")).
		All()
	utils.Must(err)
	log.Infof("Got %d files", len(ktFiles))

	utils.Must(doOrganizeKiteiMakor(ktFiles))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doOrganizeKiteiMakor(ktFiles []*models.File) error {
	cuByID := make(map[int64]*models.ContentUnit)
	cuMap := make(map[int64][]*models.File)
	noCU := make([]*models.File, 0)
	for i := range ktFiles {
		f := ktFiles[i]

		// skip sketches with buggy file name
		if f.Type == "image" {
			continue
		}

		if f.ContentUnitID.Valid {
			cuID := f.ContentUnitID.Int64
			cuByID[cuID] = f.R.ContentUnit
			v := cuMap[cuID]
			if v == nil {
				v = make([]*models.File, 0)
			}
			v = append(v, f)
			cuMap[cuID] = v
		} else {
			noCU = append(noCU, f)
		}
	}

	log.Infof("Here comes %d noCU", len(noCU))
	for i := range noCU {
		f := noCU[i]
		log.Infof("%d %s", f.ID, f.Name)
	}

	log.Infof("len(cuMap) %d", len(cuMap))
	for k, v := range cuMap {
		cu := cuByID[k]
		log.Infof("CU [%d] type_id %d has %d kitei-makor files", k, cu.TypeID, len(v))
		//for i := range v {
		//	f := v[i]
		//	log.Infof("%d %s", f.ID, f.Name)
		//}

		// we need to keep a state where main CU --derives--> KITEI_MAKOR CU
		// so follow link from either end.
		// 1. if KT CU doesn't exist then create it and link it
		// 2. move all 'kitei-makor' files into KT CU
		if api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name == api.CT_KITEI_MAKOR {
			err := cu.L.LoadDerivedContentUnitDerivations(mdb, true, cu)
			if err != nil {
				return errors.Wrapf(err, "Load Source CUs for %d", cu.ID)
			}

			if len(cu.R.DerivedContentUnitDerivations) == 0 {
				log.Infof("KT CU has no source: %d", cu.ID)
			} else if len(cu.R.DerivedContentUnitDerivations) > 1 {
				log.Warnf("KT CU has too many source CU %d", cu.ID, len(cu.R.DerivedContentUnitDerivations))
			} else {
				cud := cu.R.DerivedContentUnitDerivations[0]
				if vv, ok := cuMap[cud.SourceID]; ok {
					log.Infof("KT CU %d has Source CU %d with %d KT files", cu.ID, cud.SourceID, len(vv))
				} else {
					log.Infof("KT CU %d, Source CU %d has no KT files", cu.ID, cud.SourceID)
				}
			}
		} else {
			err := cu.L.LoadSourceContentUnitDerivations(mdb, true, cu)
			if err != nil {
				return errors.Wrapf(err, "Load Derived CUs for %d", cu.ID)
			}

			if len(cu.R.SourceContentUnitDerivations) == 0 {
				log.Infof("Main CU has no derived CUs: %d", cu.ID)

				tx, err := mdb.Begin()
				utils.Must(err)

				// create and associate KT CU and move KT files there
				ktCU, err := createKTCU(tx, cu, v)
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "mainCU %d", cu.ID)
				}

				// move files from main cu to kt cu
				err = moveKTFiles(tx, v, ktCU.ID)
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "mainCU %d", cu.ID)
				}

				utils.Must(tx.Commit())

			} else if len(cu.R.SourceContentUnitDerivations) > 1 {
				log.Warnf("Main CU %d has too many derived CUs %d", cu.ID, len(cu.R.SourceContentUnitDerivations))
			} else {
				cud := cu.R.SourceContentUnitDerivations[0]
				if vv, ok := cuMap[cud.DerivedID]; ok {
					log.Infof("Main CU %d has KT CU %d with %d files", cu.ID, cud.DerivedID, len(vv))
				} else {
					log.Infof("Main CU %d, KT CU %d has no files", cu.ID, cud.DerivedID)
				}

				tx, err := mdb.Begin()
				utils.Must(err)

				// move files from main cu to kt cu
				err = moveKTFiles(tx, v, cud.DerivedID)
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "mainCU %d", cu.ID)
				}

				utils.Must(tx.Commit())
			}
		}
	}

	return nil
}

func createKTCU(exec boil.Executor, mainCU *models.ContentUnit, ktFiles []*models.File) (*models.ContentUnit, error) {
	log.Infof("Create KT CU for %d", mainCU.ID)

	var props map[string]interface{}
	if !mainCU.Properties.Valid {
		return nil, errors.Errorf("mainCU invalid props %d", mainCU.ID)
	}

	err := json.Unmarshal(mainCU.Properties.JSON, &props)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal mainCU properties %d", mainCU.ID)
	}
	delete(props, "kmedia_id")

	// calculate duration if possible
	var sum float64
	var count int
	for i := range ktFiles {
		f := ktFiles[i]
		if f.Properties.Valid {
			var fProps map[string]interface{}
			err := json.Unmarshal(f.Properties.JSON, &fProps)
			if err != nil {
				return nil, errors.Wrapf(err, "json.Unmarshal file properties %d", f.ID)
			}

			if d, ok := fProps["duration"]; ok {
				if d.(float64) > 0 {
					sum += d.(float64)
					count += 1
				}
			}
		}
	}
	if sum > 0 {
		props["duration"] = int64(sum / float64(count))
		log.Infof("mainCU %d duration: (%f/%d)=%d", mainCU.ID, sum, count, props["duration"])
	} else {
		delete(props, "duration")
		log.Infof("mainCU %d duration invalid", mainCU.ID)
	}

	ktCU, err := api.CreateContentUnit(exec, api.CT_KITEI_MAKOR, props)
	if err != nil {
		return nil, errors.Wrap(err, "Create KT CU")
	}

	ktCU.Published = true
	err = ktCU.Update(exec, "published")
	if err != nil {
		return nil, errors.Wrapf(err, "Update KT CU published %d", ktCU)
	}

	cud := &models.ContentUnitDerivation{
		SourceID: mainCU.ID,
		Name:     api.CT_KITEI_MAKOR,
	}
	err = ktCU.AddDerivedContentUnitDerivations(exec, true, cud)
	if err != nil {
		return nil, errors.Wrap(err, "Save CUD in DB")
	}

	return ktCU, nil
}

func moveKTFiles(exec boil.Executor, ktFiles []*models.File, cuID int64) error {
	log.Infof("Moving %d files to %d", len(ktFiles), cuID)

	ids := make([]string, len(ktFiles))
	for i := range ktFiles {
		ids[i] = strconv.Itoa(int(ktFiles[i].ID))
	}

	_, err := queries.Raw(exec,
		fmt.Sprintf("UPDATE files SET content_unit_id=%d WHERE id IN (%s)", cuID, strings.Join(ids, ","))).
		Exec()
	if err != nil {
		return errors.Wrap(err, "move KT files")
	}

	return nil
}
