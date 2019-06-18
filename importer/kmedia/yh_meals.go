package kmedia

import (
	"database/sql"
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportFlatCatalogs() {
	clock := Init()

	stats = NewImportStatistics()
	utils.Must(importFlatCatalog(120, common.CT_FRIENDS_GATHERING))
	utils.Must(importFlatCatalog(4791, common.CT_MEAL))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importFlatCatalog(catalogID int, cuType string) error {
	catalog, err := kmodels.Catalogs(kmdb,
		qm.Where("id=?", catalogID),
		qm.Load("Containers")).One()
	if err != nil {
		return errors.Wrapf(err, "Load catalog %d", catalogID)
	}

	stats.CatalogsProcessed.Inc(1)

	for i := range catalog.R.Containers {
		tx, err := mdb.Begin()
		utils.Must(err)

		if err = importFlatContainer(tx, catalog.R.Containers[i], cuType); err != nil {
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

func importFlatContainer(exec boil.Executor, container *kmodels.Container, cuType string) error {
	stats.ContainersProcessed.Inc(1)

	unit, err := models.ContentUnits(mdb, qm.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("New CU %d %s", container.ID, container.Name.String)
			_, err := importContainerWOCollectionNewCU(exec, container, cuType)
			return err
		}
		return errors.Wrapf(err, "Lookup content unit kmid %d", container.ID)
	}

	log.Infof("CU exists [%d] container: %s %d", unit.ID, container.Name.String, container.ID)
	_, err = importContainer(exec, container, nil, cuType, "", 0)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	if cuType != common.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID].Name {
		log.Infof("Overriding CU Type to %s", cuType)
		unit.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[cuType].ID
		err = unit.Update(exec, "type_id")
		if err != nil {
			return errors.Wrapf(err, "Update CU type %d", unit.ID)
		}
	}

	return nil
}
