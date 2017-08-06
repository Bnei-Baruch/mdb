package kmedia

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"sort"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const CONGRESSES_FILE = "importer/kmedia/data/Conventions - congresses.csv"
const NEW_UNITS_FILE = "importer/kmedia/data/congress_units_new.csv"
const EXISTING_UNITS_FILE = "importer/kmedia/data/congress_units_exists.csv"

type EventPart struct {
	KMediaID  int
	Name      string
	Position  int
	Container *kmodels.Container
}

type Congress struct {
	KMediaID    int
	Country     string
	City        string
	FullAddress string
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Year        int
	Month       int
	Events      []EventPart
	Catalog     *kmodels.Catalog
}

var CountryCatalogs = map[int]string{
	7900: "Ukraine",
	8091: "Kazakhstan",
	4556: "Canada",
	4537: "United States",
	4536: "Russia",
	4543: "Estonia",
	4553: "Austria",
	4680: "Colombia",
	4563: "Switzerland",
	4544: "Spain",
	4549: "Argentina",
	4787: "Sweden",
	4788: "Bulgaria",
	2323: "Mexico",
	4545: "Italy",
	4552: "Chile",
	4667: "Brazil",
	4551: "England",
	4547: "Germany",
	4550: "Turkey",
	7872: "Romania",
	7910: "Czech",
	4658: "Lithuania",
	4741: "France",
	4710: "Georgia",
	4554: "Poland",
	//4555: "Israel",
	//8030: "America_2017",
}

func ImportCongresses() {
	clock := Init()

	log.Infof("Loading congresses file")
	congresses, err := initCongresses()
	utils.Must(err)

	stats = NewImportStatistics()

	//utils.Must(dumpCongresses())
	//utils.Must(loadEventParts(congresses))
	//utils.Must(analyzeExisting(congresses))
	utils.Must(importNewCongresses(congresses))
	utils.Must(importNewUnits())
	utils.Must(importExistingUnits())

	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func dumpCongresses() error {
	all := make([]*Congress, 0)
	for cID, country := range CountryCatalogs {
		log.Infof("Processing %s", country)

		catalogs, err := kmodels.Catalogs(kmdb, qm.Where("parent_id = ?", cID)).All()
		if err != nil {
			return errors.Wrap(err, "Load catalogs from db")
		}
		log.Infof("Got %d catalogs", len(catalogs))

		for _, catalog := range catalogs {
			c := &Congress{
				KMediaID: catalog.ID,
				Country:  country,
				Name:     catalog.Name,
			}
			all = append(all, c)
		}
	}

	log.Infof("Found %d total congresses", len(all))
	for _, congress := range all {
		fmt.Printf("%d\t%s\t%s\n", congress.KMediaID, congress.Country, congress.Name)
	}

	return nil
}

func initCongresses() (map[int]*Congress, error) {
	// Read mappings file
	records, err := utils.ReadCSV(CONGRESSES_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read congresses")
	}
	log.Infof("Congresses file has %d rows", len(records))

	// Create mappings
	mappings := make(map[int]*Congress, len(records)-1)
	for i, r := range records[1:] {
		kmid, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, errors.Wrapf(err, "Bad kmedia_id, row [%d]", i)
		}

		var year int
		if r[6] != "" {
			year, err = strconv.Atoi(r[6])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad year, row [%d]", i)
			}
		}

		var month int
		if r[7] != "" {
			month, err = strconv.Atoi(r[7])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad month, row [%d]", i)
			}
		}

		var start time.Time
		if r[4] != "" {
			start, err = time.Parse("2006-01-02", r[4])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad start_date, row [%d]", i)
			}
		}

		var end time.Time
		if r[5] != "" {
			end, err = time.Parse("2006-01-02", r[5])
			if err != nil {
				return nil, errors.Wrapf(err, "Bad end_date, row [%d]", i)
			}
		}

		mappings[kmid] = &Congress{
			KMediaID:    kmid,
			Country:     r[1],
			City:        r[2],
			FullAddress: r[3],
			StartDate:   start,
			EndDate:     end,
			Year:        year,
			Month:       month,
		}
	}

	return mappings, nil
}

func loadEventParts(congresses map[int]*Congress) error {
	log.Infof("Loading event parts")

	for kmid, congress := range congresses {
		if congress.StartDate.IsZero() || congress.EndDate.IsZero() {
			log.Infof("Skipping congress: %s [%d]", congress.Country, kmid)
			continue
		}

		catalog, err := kmodels.FindCatalog(kmdb, kmid)
		if err != nil {
			return errors.Wrapf(err, "Load congress catalog from kmdb [%d]", kmid)
		}
		congress.Catalog = catalog

		err = catalog.L.LoadContainers(kmdb, true, catalog)
		if err != nil {
			return errors.Wrapf(err, "Load containers from kmdb [%d]", kmid)
		}

		log.Infof("[%d] %s-%d has %d containers",
			congress.KMediaID, congress.Country, congress.Year, len(catalog.R.Containers))

		for _, container := range catalog.R.Containers {
			ep := EventPart{
				KMediaID:  container.ID,
				Name:      container.Name.String,
				Position:  container.Position.Int,
				Container: container,
			}
			congress.Events = append(congress.Events, ep)
		}

		sort.Slice(congress.Events, func(i, j int) bool {
			return congress.Events[i].Position < congress.Events[j].Position
		})

		//for _, ep := range congress.Events {
		//	log.Infof("%d. %s [%d]", ep.Position, ep.Name, ep.KMediaID)
		//}
	}

	return nil
}

func analyzeExisting(congresses map[int]*Congress) error {
	fExists, err := ioutil.TempFile("/tmp", "congress_units_exists")
	if err != nil {
		return errors.Wrap(err, "Create temp file: fExists")
	}
	defer fExists.Close()
	log.Infof("fExists file: %s", fExists.Name())

	fNew, err := ioutil.TempFile("/tmp", "congress_units_new")
	if err != nil {
		return errors.Wrap(err, "Create temp file: fNew")
	}
	defer fNew.Close()
	log.Infof("fNew file: %s", fNew.Name())

	for kmid, congress := range congresses {
		if congress.StartDate.IsZero() || congress.EndDate.IsZero() {
			log.Infof("Skipping congress: %s [%d]", congress.Country, kmid)
			continue
		}

		//collection, err := models.Collections(mdb,
		//	qm.Where("(properties->>'kmedia_id')::int = ?", kmid),
		//	qm.Load("CollectionsContentUnits", "CollectionsContentUnits.ContentUnit"),
		//).One()
		//if err != nil {
		//	if err == sql.ErrNoRows {
		//		// create new congress
		//		log.Infof("Create new collection for %s [%d]", congress.Catalog.Name, kmid)
		//		collection, err = api.CreateCollection(mdb, api.CT_CONGRESS, map[string]interface{}{
		//			"kmedia_id":  kmid,
		//			"active":     false,
		//			"country":    congress.Country,
		//			"city":       congress.City,
		//			"start_date": congress.StartDate,
		//			"end_date":   congress.EndDate,
		//		})
		//		if err != nil {
		//			return errors.Wrapf(err, "Create collection, [kmid %d]", kmid)
		//		}
		//
		//		// TODO: i18n ??
		//	} else {
		//		return errors.Wrapf(err, "Load collection from mdb [kmid %d]", kmid)
		//	}
		//}
		//
		//log.Infof("MDB Collection [%d] kmedia_id %d", collection.ID, kmid)

		for _, event := range congress.Events {
			cu, err := models.ContentUnits(mdb,
				qm.Where("(properties->>'kmedia_id')::int = ?", event.KMediaID),
				qm.Load("ContentUnitI18ns", "CollectionsContentUnits"),
			).One()
			if err != nil {
				if err == sql.ErrNoRows {
					cu = nil
				} else {
					return errors.Wrapf(err, "Lookup content unit kmid %d", event.KMediaID)
				}
			}

			if cu == nil {
				// new
				_, err = fmt.Fprintf(fNew, "%d,%s,%d,%s,%d\n",
					congress.KMediaID, congress.Catalog.Name, event.KMediaID, event.Name, event.Container.ContentTypeID.Int)
				if err != nil {
					return errors.Wrap(err, "write to new units file")
				}
			} else {
				// exists
				//var name string
				//for _, i18n := range cu.R.ContentUnitI18ns {
				//	if i18n.Language == api.LANG_HEBREW {
				//		name = i18n.Name.String
				//		break
				//	}
				//}

				var previousCollection int64
				if len(cu.R.CollectionsContentUnits) > 0 {
					previousCollection = cu.R.CollectionsContentUnits[0].CollectionID
				}

				_, err = fmt.Fprintf(fExists, "%d,%s,%d,%s,%d,%d,%d,%d\n",
					congress.KMediaID, congress.Catalog.Name, event.KMediaID, event.Name, event.Container.ContentTypeID.Int,
					cu.ID, cu.TypeID, previousCollection)
				if err != nil {
					return errors.Wrap(err, "write to new units file")
				}
			}
		}

	}

	return nil
}

func importNewCongresses(congresses map[int]*Congress) error {
	tx, err := mdb.Begin()
	if err != nil {
		return errors.Wrap(err, "Start transaction")
	}

	for kmid, congress := range congresses {
		if congress.StartDate.IsZero() || congress.EndDate.IsZero() {
			log.Infof("Skipping congress: %s [%d]", congress.Country, kmid)
			continue
		}
		stats.CatalogsProcessed.Inc(1)

		collection, err := models.Collections(mdb,
			qm.Where("(properties->>'kmedia_id')::int = ?", kmid),
			qm.Load("CollectionsContentUnits", "CollectionsContentUnits.ContentUnit"),
		).One()
		if err != nil {
			if err == sql.ErrNoRows {
				// create new congress
				log.Infof("Create new collection for %s %s [%d]", congress.Country, congress.StartDate, kmid)
				collection, err = api.CreateCollection(mdb, api.CT_CONGRESS, map[string]interface{}{
					"kmedia_id":  kmid,
					"active":     false,
					"country":    congress.Country,
					"city":       congress.City,
					"start_date": congress.StartDate,
					"end_date":   congress.EndDate,
				})
				if err != nil {
					return errors.Wrapf(err, "Create collection, [kmid %d]", kmid)
				}
				stats.CollectionsCreated.Inc(1)

				// I18n
				descriptions, err := kmodels.CatalogDescriptions(kmdb, qm.Where("catalog_id = ?", kmid)).All()
				if err != nil {
					return errors.Wrapf(err, "Lookup catalog descriptions, [kmid %d]", kmid)
				}
				for _, d := range descriptions {
					if d.Name.Valid && d.Name.String != "" {
						ci18n := models.CollectionI18n{
							CollectionID: collection.ID,
							Language:     api.LANG_MAP[d.LangID.String],
							Name:         d.Name,
						}
						err = ci18n.Upsert(mdb,
							true,
							[]string{"collection_id", "language"},
							[]string{"name"})
						if err != nil {
							return errors.Wrapf(err, "Upsert collection i18n, collection [%d]", collection.ID)
						}
					}
				}
			} else {
				return errors.Wrapf(err, "Load collection from mdb [kmid %d]", kmid)
			}
		}

		log.Infof("MDB Collection [%d] kmedia_id %d", collection.ID, kmid)

	}

	if err == nil {
		utils.Must(tx.Commit())
		stats.TxCommitted.Inc(1)
	} else {
		utils.Must(tx.Rollback())
		stats.TxRolledBack.Inc(1)
		return err
	}

	return nil
}

func importNewUnits() error {
	// Read mappings file
	records, err := utils.ReadCSV(NEW_UNITS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read NEW_UNITS_FILE")
	}
	log.Infof("NEW_UNITS_FILE file has %d rows", len(records))

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	for _, x := range records[1:] {
		tx, err := mdb.Begin()
		if err != nil {
			return errors.Wrap(err, "Start transaction")
		}

		if err = doNewUnit(tx, h, x); err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		} else {
			utils.Must(tx.Commit())
			stats.TxCommitted.Inc(1)
		}
	}

	return nil
}

func doNewUnit(exec boil.Executor, h map[string]int, x []string) error {
	// input validation
	ct := x[h["content_type"]]
	if ct == "" {
		log.Infof("Empty content_type, skipping %s", x)
		return nil
	}

	containerID := x[h["container_id"]]
	exists, err := models.ContentUnits(exec,
		qm.Where("(properties->>'kmedia_id')::int = ?", containerID)).Exists()
	if err != nil {
		return errors.Wrapf(err, "Check unit exists in mdb [kmid %d]", containerID)
	}
	if exists {
		log.Infof("Unit exists [kmid %d]", containerID)
		return nil
	}

	catalogID := x[h["catalog_id"]]
	collection, err := models.Collections(exec,
		qm.Where("(properties->>'kmedia_id')::int = ?", catalogID),
	).One()
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrapf(err, "Collection not found [kmid %d]", catalogID)
		} else {
			return errors.Wrapf(err, "Load collection from mdb [kmid %d]", catalogID)
		}
	}

	cID, err := strconv.Atoi(containerID)
	if err != nil {
		return errors.Wrapf(err, "bad container_id %s", containerID)
	}
	container, err := kmodels.Containers(kmdb,
		qm.Where("id = ?", cID),
		qm.Load("FileAssets")).One()
	if err != nil {
		return errors.Wrapf(err, "Lookup container %d", cID)
	}

	ccuName := strconv.Itoa(container.Position.Int)
	if ct == api.CT_EVENT_PART {
		ccuName = x[h["event_part_type"]] + ccuName
	}

	// Create import operation
	operation, err := api.CreateOperation(exec, api.OP_IMPORT_KMEDIA,
		api.Operation{WorkflowID: strconv.Itoa(cID)}, nil)
	if err != nil {
		return errors.Wrapf(err, "Create operation %d", cID)
	}
	stats.OperationsCreated.Inc(1)

	// import container
	log.Infof("Processing container %d", cID)
	stats.ContainersProcessed.Inc(1)
	unit, err := importContainer(exec, container, collection, ct, ccuName)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", cID)
	}

	// import container files
	var file *models.File
	for _, fileAsset := range container.R.FileAssets {
		log.Infof("Processing file_asset %d", fileAsset.ID)
		stats.FileAssetsProcessed.Inc(1)

		// Create or update MDB file
		file, err = importFileAsset(exec, fileAsset, unit, operation)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			break
		}
		if file.Published {
			unit.Published = true
		}
	}
	if err != nil {
		return errors.Wrapf(err, "Import container files %d", cID)
	}

	if unit.Published {
		collection.Published = true
		err = unit.Update(exec, "published")
		if err != nil {
			return errors.Wrapf(err, "Update unit published column %d", cID)
		}
	}

	return nil
}

func importExistingUnits() error {
	// Read mappings file
	records, err := utils.ReadCSV(EXISTING_UNITS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read EXISTING_UNITS_FILE")
	}
	log.Infof("EXISTING_UNITS_FILE file has %d rows", len(records))

	h, err := utils.ParseCSVHeader(records[0])
	if err != nil {
		return errors.Wrap(err, "Bad header")
	}

	for _, x := range records[1:] {
		tx, err := mdb.Begin()
		if err != nil {
			return errors.Wrap(err, "Start transaction")
		}

		if err = doExistingUnit(tx, h, x); err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		} else {
			utils.Must(tx.Commit())
			stats.TxCommitted.Inc(1)
		}
	}

	return nil
}

func doExistingUnit(exec boil.Executor, h map[string]int, x []string) error {
	unitID := x[h["unit.id"]]
	cuID, err := strconv.Atoi(unitID)
	if err != nil {
		return errors.Wrapf(err, "bad unit.id %s", cuID)
	}
	unit, err := models.FindContentUnit(exec, int64(cuID))
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrapf(err, "Unit doesn't exist in mdb [%d]", unitID)
		} else {
			return errors.Wrapf(err, "Load unit from mdb [%d]", unitID)
		}
	}

	catalogID := x[h["catalog_id"]]
	collection, err := models.Collections(exec,
		qm.Where("(properties->>'kmedia_id')::int = ?", catalogID),
	).One()
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrapf(err, "Collection not found [kmid %d]", catalogID)
		} else {
			return errors.Wrapf(err, "Load collection from mdb [kmid %d]", catalogID)
		}
	}

	containerID := x[h["container_id"]]
	cID, err := strconv.Atoi(containerID)
	if err != nil {
		return errors.Wrapf(err, "bad container_id %s", containerID)
	}
	container, err := kmodels.Containers(kmdb, qm.Where("id = ?", cID)).One()
	if err != nil {
		return errors.Wrapf(err, "Lookup container %d", cID)
	}

	ccuName := strconv.Itoa(container.Position.Int)

	// import container
	log.Infof("Processing container %d", cID)
	stats.ContainersProcessed.Inc(1)

	log.Infof("Associating unit %d [kmid %d] with collection %d [kmid %s] ccuName: %s",
		unit.ID, container.ID, collection.ID, catalogID, ccuName)

	ccu, err := models.CollectionsContentUnits(exec,
		qm.Where("collection_id = ? and content_unit_id = ?", collection.ID, unit.ID)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			// create new association
			err = unit.AddCollectionsContentUnits(exec, true, &models.CollectionsContentUnit{
				CollectionID:  collection.ID,
				ContentUnitID: unit.ID,
				Name:          ccuName,
			})
			if err != nil {
				return errors.Wrapf(err, "Associating unit %d [kmid %d] with collection %d [kmid %s]",
					unit.ID, container.ID, collection.ID, catalogID)
			}
			return nil
		} else {
			return errors.Wrapf(err, "Load CCU from mdb")
		}
	}

	// update existing association name
	ccu.Name = ccuName
	err = ccu.Update(exec, "name")
	if err != nil {
		return errors.Wrapf(err, "Update ccu name")
	}

	return nil
}
