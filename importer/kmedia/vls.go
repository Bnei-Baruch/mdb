package kmedia

import (
	"encoding/json"
	"github.com/Bnei-Baruch/mdb/common"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var VLCatalogs = []int{
	8154, //virtual_lessons
	4654, //virtual_lessons/program_lc
	27,   //virtual_lessons/program_lc_rus
	//2054, //virtual_lessons/series_2002-2004
	277,  //virtual_lessons/series_2002-2004/VirtualLessonRus/HakdamaLeSeferHaZohar
	276,  //virtual_lessons/series_2002-2004/VirtualLessonRus/HakdamaLeTaas
	1171, //virtual_lessons/series_2002-2004/VirtualLessonRus/MahutHochmatKabbala
	281,  //virtual_lessons/series_2002-2004/VirtualLessonRus/MatanTora
	278,  //virtual_lessons/series_2002-2004/VirtualLessonRus/MavoLeSeferHaZohar
	284,  //virtual_lessons/series_2002-2004/VirtualLessonRus/PanimMeirotVeMasbirot
	292,  //virtual_lessons/series_2002-2004/VirtualLessonRus/PiHaham
	279,  //virtual_lessons/series_2002-2004/VirtualLessonRus/PriHaham
	291,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Pticha
	285,  //virtual_lessons/series_2002-2004/VirtualLessonRus/PtichaKolelet
	283,  //virtual_lessons/series_2002-2004/VirtualLessonRus/PtichaLePirushHaSulam
	286,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Q&A
	280,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Shamati
	288,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Shonot
	282,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Taas
	1376, //virtual_lessons/series_2002-2004/VirtualLessonRus/The Freedom
	275,  //virtual_lessons/series_2002-2004/VirtualLessonRus/Vvedenie
	287,  //virtual_lessons/series_2002-2004/VirtualLessonRus/ZitutimKabbalahLe Matchil
	2053, //virtual_lessons/series_2005-2007
	3559, //virtual_lessons/series_2007-2008
	3892, //virtual_lessons/series_2009
	4164, //virtual_lessons/series_2010
	8017, //virtual_lessons/vl_chn_fundamental-course
	8018, //virtual_lessons/vl_chn_intermediate-course
	8019, //virtual_lessons/vl_chn_preface-to-the-wisdom-of-kabbalah
	725,  //virtual_lessons/vl_EuroKabSeries
	213,  //virtual_lessons/vl_heb_kurs-esod
	3589, //virtual_lessons/vl_heb_virtual-group-israel
	3920, //virtual_lessons/vl_rus
	4760, //virtual_lessons/vl_rus_ahana-le-zohar
	4351, //virtual_lessons/vl_rus_osnovi-kabbali
	4036, //virtual_lessons/vl_rus_otkrivaem-zohar
	//8068, //virtual_lessons/vl_webinar
	//8067, //virtual_lessons/vl_zman-kabbalah

}

func ImportVLs() {
	clock := Init()

	stats = NewImportStatistics()

	csMap, err := loadAndImportMissingVLCollections()
	utils.Must(err)

	utils.Must(importVLsContainers(csMap))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadAndImportMissingVLCollections() (map[int]*models.Collection, error) {

	cs, err := models.Collections(mdb,
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_VIRTUAL_LESSONS].ID)).
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

	for _, kmID := range VLCatalogs {
		if _, ok := csMap[kmID]; ok {
			continue
		}

		props := map[string]interface{}{
			"kmedia_id": kmID,
			"active": false,
		}
		log.Infof("Create collection %v", props)
		c, err := api.CreateCollection(mdb, common.CT_VIRTUAL_LESSONS, props)
		if err != nil {
			return nil, errors.Wrapf(err, "Create collection %v", props)
		}
		stats.CollectionsCreated.Inc(1)
		csMap[kmID] = c

		// I18n
		descriptions, err := kmodels.CatalogDescriptions(kmdb, qm.Where("catalog_id = ?", kmID)).All()
		if err != nil {
			return nil, errors.Wrapf(err, "Lookup catalog descriptions, [kmid %d]", kmID)
		}
		for _, d := range descriptions {
			if d.Name.Valid && d.Name.String != "" {
				ci18n := models.CollectionI18n{
					CollectionID: c.ID,
					Language:     common.LANG_MAP[d.LangID.String],
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

func importVLsContainers(csMap map[int]*models.Collection) error {
	cnMap, cuMap, err := loadContainersInCatalogsAndCUs(8154)
	if err != nil {
		return errors.Wrap(err, "Load containers")
	}

	csToPublish := make(map[int64]*models.Collection, len(csMap))

	for cnID, cn := range cnMap {
		tx, err := mdb.Begin()
		utils.Must(err)

		// import container
		cu, ok := cuMap[cnID]
		if ok {
			// update - for tags
			cu, err = importContainer(tx, cn, nil, common.CT_VIRTUAL_LESSON, "", 0)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Import existing container %d", cnID)
			}

			if common.CT_VIRTUAL_LESSON != common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name {
				cu.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[common.CT_VIRTUAL_LESSON].ID
				err = cu.Update(tx, "type_id")
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "Update CU type %d", cu.ID)
				}
			}
		} else {
			if _, ok := knownKMContentTypes[cn.ContentTypeID.Int]; ok {
				// create - CU doesn't exist
				cu, err = importContainerWOCollectionNewCU(tx, cn, common.CT_VIRTUAL_LESSON)
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
			c, ok := csMap[catalog.ID]
			if !ok {
				continue
			}

			if cu.Published && !c.Published{
				csToPublish[c.ID] = c
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

	// publish collections
	for _, v := range csToPublish {
		v.Published = true
		if err :=v.Update(mdb,"published"); err != nil {
			return errors.Wrapf(err, "Publish collection %d", v.ID)
		}
	}

	return nil
}
