package kmedia

import (
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
	qm4 "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportLectures() {
	clock := Init()

	stats = NewImportStatistics()
	utils.Must(importLectures(709))
	utils.Must(importLectures(4574))
	utils.Must(importLectures(4508))
	utils.Must(importLectures(2186))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importLectures(catalogID int) error {
	catalog, err := kmodels.Catalogs(kmdb,
		qm.Where("id=?", catalogID),
		qm.Load("Containers")).One()
	if err != nil {
		return errors.Wrapf(err, "Load catalog %d", catalogID)
	}

	collection, err := models.Collections(
		qm4.Where("(properties->>'kmedia_id')::int = ?", catalogID)).
		One(mdb)
	if err != nil {
		return errors.Wrapf(err, "Lookup collection in mdb [kmid %d]", catalogID)
	}

	stats.CatalogsProcessed.Inc(1)

	for i := range catalog.R.Containers {
		tx, err := mdb.Begin()
		utils.Must(err)

		if err = importContainerWCollection(tx, catalog.R.Containers[i], collection, common.CT_LECTURE); err != nil {
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
