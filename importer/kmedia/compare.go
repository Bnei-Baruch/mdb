package kmedia

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func Compare() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting Kmedia compare")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Setting up connection to Kmedia")
	kmdb, err = sql.Open("postgres", viper.GetString("kmedia.url"))
	utils.Must(err)
	utils.Must(kmdb.Ping())
	defer kmdb.Close()

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	log.Info("Loading kmedia catalogs hierarchy")
	utils.Must(loadCatalogHierarchy())

	//log.Info("Compare collections")
	//utils.Must(compareCollections())

	//log.Info("Compare units")
	//utils.Must(compareUnits())

	log.Info("Missing containers")
	utils.Must(missingContainers())

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func compareCollections() error {
	log.Info("Loading all collections with kmedia_id")
	collections, err := models.Collections(mdb,
		//qm.Where("properties -> 'kmedia_id' is not null"),
		qm.Where("type_id=4 and properties -> 'kmedia_id' is not null"),
		qm.Load("CollectionsContentUnits", "CollectionsContentUnits.ContentUnit")).
		All()
	if err != nil {
		return errors.Wrap(err, "Load collections from mdb")
	}

	log.Infof("Got %d collections", len(collections))

	for i := range collections {
		if err := compareCollection(collections[i]); err != nil {
			return errors.Wrapf(err, "Compare collection %d [%d]", i, collections[i].ID)
		}
	}

	return nil
}

func compareUnits() error {
	log.Info("Loading all content units with kmedia_id")
	units, err := models.ContentUnits(mdb,
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	if err != nil {
		return errors.Wrap(err, "Load content_units from mdb")
	}

	log.Infof("Got %d units", len(units))

	for i := range units {
		cu := units[i]
		var props map[string]interface{}
		err := json.Unmarshal(cu.Properties.JSON, &props)
		if err != nil {
			return errors.Wrapf(err, "json.Unmarshal unit properties %d", cu.ID)
		}

		cn, err := kmodels.FindContainer(kmdb, int(props["kmedia_id"].(float64)))
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("Container doesn't exists %d [cu %d]", props["kmedia_id"], cu.ID)
				continue
			} else {
				return errors.Wrapf(err, "Load container %d", props["kmedia_id"])
			}
		}

		if err := compareUnit(cu, cn); err != nil {
			return errors.Wrapf(err, "Compare unit %d [%d]", i, units[i].ID)
		}
	}

	return nil
}

func missingContainers() error {
	cuMap := make(map[int]bool)

	log.Info("Loading all content units with kmedia_id")
	units, err := models.ContentUnits(mdb,
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	if err != nil {
		return errors.Wrap(err, "Load content_units from mdb")
	}
	log.Infof("Got %d units", len(units))

	for i := range units {
		cu := units[i]
		var props map[string]interface{}
		err := json.Unmarshal(cu.Properties.JSON, &props)
		if err != nil {
			return errors.Wrapf(err, "json.Unmarshal unit properties %d", cu.ID)
		}
		cuMap[int(props["kmedia_id"].(float64))] = true
	}

	containers, err := kmodels.Containers(kmdb).All()
	if err != nil {
		return errors.Wrap(err, "Load containers")
	}

	f, err := ioutil.TempFile("/tmp", "kmedia_compare")
	if err != nil {
		return errors.Wrap(err, "Create temp file")
	}
	defer f.Close()

	log.Infof("Report file: %s", f.Name())

	for i := range containers {
		cn := containers[i]
		if _, ok := cuMap[cn.ID]; !ok {
			_, err := fmt.Fprintf(f, "%d\t%d\t%s\t%s\t%d\n", cn.ID, cn.ContentTypeID.Int, cn.Filmdate.Time.Format("2006-01-02"), cn.Name.String, cn.VirtualLessonID.Int)
			if err != nil {
				return errors.Wrapf(err, "Write tsv row %d", cn.ID)
			}
		}
	}

	return nil
}

func compareCollection(c *models.Collection) error {
	var props map[string]interface{}
	if err := json.Unmarshal(c.Properties.JSON, &props); err != nil {
		return errors.Wrap(err, "json.Unmarshal properties")
	}

	kmid := int(props["kmedia_id"].(float64))

	isLesson := api.CONTENT_TYPE_REGISTRY.ByName[api.CT_DAILY_LESSON].ID == c.TypeID ||
		api.CONTENT_TYPE_REGISTRY.ByName[api.CT_SPECIAL_LESSON].ID == c.TypeID

	// lessons are compared to virtual lesson all others are compared to catalog
	var containers []*kmodels.Container
	if isLesson {
		log.Infof("Compare collection %d [%s] to virtual_lesson %d", c.ID, api.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name, kmid)

		vl, err := kmodels.VirtualLessons(kmdb,
			qm.Where("id=?", kmid),
			qm.Load("Containers")).
			One()
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("Kmedia virtual_lesson %d doesn't exist %d", kmid, c.ID)
				return nil
			}
			return errors.Wrapf(err, "Load kmedia catalog %d", kmid)
		}
		log.Infof("Virtual Lesson %d %s", kmid, vl.FilmDate.Time.Format("2006-01-02"))
		containers = vl.R.Containers
	} else {
		log.Infof("Compare collection %d [%s] to catalog %d", c.ID, api.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name, kmid)

		catalog, err := kmodels.Catalogs(kmdb,
			qm.Where("id=?", kmid),
			qm.Load("Containers")).
			One()
		if err != nil {
			if err == sql.ErrNoRows {
				log.Warnf("Kmedia catalog %d doesn't exist %d", kmid, c.ID)
				return nil
			}
			return errors.Wrapf(err, "Load kmedia catalog %d", kmid)
		}

		log.Infof("Catalog %d %s", kmid, catalog.Name)
		containers = catalog.R.Containers
	}

	// compare list of CU <-> Containers

	// kmedia container by id
	cnMap := make(map[int]*kmodels.Container, len(containers))
	for i := range containers {
		cn := containers[i]
		cnMap[cn.ID] = cn
	}

	matching := make(map[*models.ContentUnit]*kmodels.Container)
	for i := range c.R.CollectionsContentUnits {
		cu := c.R.CollectionsContentUnits[i].R.ContentUnit
		var kmCnID int
		if cu.Properties.Valid {
			var cuProps map[string]interface{}
			if err := json.Unmarshal(cu.Properties.JSON, &cuProps); err != nil {
				return errors.Wrapf(err, "json.Unmarshal CU properties %d", cu.ID)
			}
			if v, ok := cuProps["kmedia_id"]; ok {
				kmCnID = int(v.(float64))
			}
		}

		if kmCnID > 0 {
			if cn, ok := cnMap[kmCnID]; ok {
				matching[cu] = cn
				delete(cnMap, kmCnID)
			} else {
				// kmedia container not in catalog, does it exists at all ?
				cn, err := kmodels.Containers(kmdb,
					qm.Where("id=?", kmCnID),
					qm.Load("Catalogs")).
					One()
				if err != nil {
					if err == sql.ErrNoRows {
						log.Warnf("Container %d doesn't exists. kmid %d. [c,cu] = [%d,%d]", kmCnID, kmid, c.ID, cu.ID)
					} else {
						return errors.Wrapf(err, "Check container exists %d", kmCnID)
					}
				} else {
					// if this container belongs to a child of this catalog, it's ok for us.
					isValid := false
					if !isLesson {
						if cn.R.Catalogs == nil {
							log.Warnf("Container %d exists but doesn't belong to any catalog. [c,cu] = [%d,%d]", kmCnID, c.ID, cu.ID)
						} else {
							for j := range cn.R.Catalogs {
								if isAncestor(kmid, cn.R.Catalogs[j].ID) {
									isValid = true
									break
								}
							}
						}
					}

					if !isValid {
						log.Warnf("Container %d exists but not in kmid %d. [c,cu] = [%d,%d]", kmCnID, kmid, c.ID, cu.ID)
					}
				}
			}
		} else {
			// CU of full lesson (backup) is not expected to be in kmedia DB
			if api.CONTENT_TYPE_REGISTRY.ByName[api.CT_FULL_LESSON].ID == cu.TypeID &&
				(api.CONTENT_TYPE_REGISTRY.ByName[api.CT_DAILY_LESSON].ID == c.TypeID ||
					api.CONTENT_TYPE_REGISTRY.ByName[api.CT_SPECIAL_LESSON].ID == c.TypeID) {
				continue
			} else {
				log.Infof("CU exists only in MDB %d [%s]", cu.ID, api.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)
			}
		}
	}

	// remaining kmedia containers not in mdb collection. Are they in MDB at all ?
	if len(cnMap) > 0 {
		for k, v := range cnMap {
			cu, err := models.ContentUnits(mdb,
				qm.Where("(properties->>'kmedia_id')::int = ?", k)).
				One()
			if err != nil {
				if err == sql.ErrNoRows {
					log.Warnf("CU doesn't exists. [kmid, container_id] = [%d,%d] %s", kmid, k, v.Name.String)
				} else {
					return errors.Wrapf(err, "check CU exists by kmedia_id %d", k)
				}
			} else {
				log.Warnf("CU %d exists but not in collection %d, [kmid, container_id]= [%d,%d] %s",
					cu.ID, c.ID, kmid, k, v.Name.String)
			}
		}
	}

	if len(matching) > 0 {
		log.Infof("%d matching", len(matching))
		for k, v := range matching {
			if err := compareUnit(k, v); err != nil {
				return errors.Wrapf(err, "compare cu %d to cn %d", k.ID, v.ID)
			}
		}
	}

	return nil
}

func compareUnit(cu *models.ContentUnit, cn *kmodels.Container) error {
	if err := cu.L.LoadFiles(mdb, true, cu); err != nil {
		return errors.Wrapf(err, "Load CU files %d", cu.ID)
	}

	if err := cn.L.LoadFileAssets(kmdb, true, cn); err != nil {
		return errors.Wrapf(err, "Load container files %d", cn.ID)
	}

	cnFiles := make(map[string]*kmodels.FileAsset, len(cn.R.FileAssets))
	for i := range cn.R.FileAssets {
		f := cn.R.FileAssets[i]
		if f.Sha1.Valid {
			cnFiles[f.Sha1.String] = f
		}
	}

	for i := range cu.R.Files {
		f := cu.R.Files[i]
		if !f.Published || !f.Sha1.Valid {
			continue
		}

		checksum := hex.EncodeToString(f.Sha1.Bytes)

		var fKMID int
		if f.Properties.Valid {
			var props map[string]interface{}
			if err := json.Unmarshal(f.Properties.JSON, &props); err != nil {
				return errors.Wrapf(err, "json.Unmarshal file properties %d", f.ID)
			}

			if v, ok := props["kmedia_id"]; ok {
				fKMID = int(v.(float64))
			}
		}

		if cnFile, ok := cnFiles[checksum]; ok {
			delete(cnFiles, checksum)

			if fKMID == 0 {
				log.Warnf("File %d is not mapped to kmedia file asset %d. %s", f.ID, cnFile.ID, f.Name)
			} else if fKMID != cnFile.ID {
				log.Warnf("wrong kmedia_id %d for file %d %s", fKMID, f.ID, f.Name)
			}
		} else {
			if fKMID == 0 {
				continue
			}

			isMissingSHA1 := false
			for i := range cn.R.FileAssets {
				f := cn.R.FileAssets[i]
				if f.ID == fKMID && !f.Sha1.Valid {
					isMissingSHA1 = true
				}
			}

			if isMissingSHA1 {
				continue
			}

			// File is not in kmedia container, is it in kmedia at all ?
			exists, err := kmodels.FileAssets(kmdb, qm.Where("id=?", fKMID)).Exists()
			if err != nil {
				return errors.Wrapf(err, "Check file_asset exists %d", fKMID)
			}

			if exists {
				log.Warnf("file_asset %d exists but not in container %d. [cu,f] = [%d,%d] %s", fKMID, cn.ID, cu.ID, f.ID, f.Name)
			} else {
				log.Warnf("file_asset doesn't exists %d. container %d. [cu,f] = [%d,%d] %s", fKMID, cn.ID, cu.ID, f.ID, f.Name)
			}
		}
	}

	// remaining kmedia file assets not in CU. Are they in MDB at all ?
	if len(cnFiles) > 0 {
		for k, v := range cnFiles {

			s, err := hex.DecodeString(k)
			if err != nil {
				return errors.Wrapf(err, "hex.DecodeString %s", k)
			}

			exists, err := models.Files(mdb, qm.Where("sha1 = ?", s)).Exists()
			if err != nil {
				return errors.Wrapf(err, "check File exists by sha1 %d", k)
			}

			if exists {
				if !strings.Contains(v.Name.String, "kitei-makor") {
					log.Warnf("File exists but not in CU %d, %s %s", cu.ID, k, v.Name.String)
				}
			} else {
				log.Warnf("File doesn't exists %s %s", k, v.Name.String)
			}
		}
	}

	return nil
}

var catalogsH map[int]int

func loadCatalogHierarchy() error {

	catalogs, err := kmodels.Catalogs(kmdb).All()
	if err != nil {
		return errors.Wrap(err, "Load all catalogs")
	}

	catalogsH = make(map[int]int, len(catalogs))

	for i := range catalogs {
		c := catalogs[i]
		if c.ParentID.Valid {
			catalogsH[c.ID] = c.ParentID.Int
		}
	}

	return nil
}

func isAncestor(parent int, child int) bool {
	//log.Infof("isAncestor: %d %d", parent, child)
	x := child
	for x > 0 && x != parent {
		//log.Infof("isAncestor in loop x: %d", x)
		if p, ok := catalogsH[x]; ok {
			//log.Infof("isAncestor in loop p: %d", p)
			x = p
		} else {
			x = 0
		}
	}

	//log.Infof("isAncestor: %d %d %t", parent, child, x == parent)
	return x == parent
}
