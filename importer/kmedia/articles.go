package kmedia

import (
	"encoding/json"
	"runtime/debug"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var ArticlesCatalogs = []int{
	8182, //	itonut/itonut_facebook
	7916, //	itonut/itonut_haaretz
	7915, //	itonut/itonut_jpost
	8183, //	itonut/itonut_newsmax
	8101, //	itonut/itonut_other
	7902, //	itonut/itonut_ynet
}

var publishersPatterns = []string{
	//"15min-lt",
	"algemeiner",
	"basicincome",
	"blogactiv",
	//"blogactive",
	//"blogavtive",
	"breakingisraelnews",
	"dreuz",
	"haaretz",
	//"huffipost",
	"huffpost",
	"iarchon",
	"internationaljpost",
	"iton-jpost",
	"jewishbusinessnews",
	"jewishvoiceny",
	"jpost",
	"medium",
	"newsmax",
	//"passover",
	"russiancanadianinfo",
	//"shavuon",
	"timesofisrael",
	"unitedwithisrael",
	"win",
	"ynet",
}

var publishersPatternsMappings = map[string]string{
	"blogactive": "blogactiv",
	"blogavtive": "blogactiv",
	"huffipost":  "huffpost",
}

func ImportArticles() {
	clock := Init()

	stats = NewImportStatistics()

	csMap, err := loadAndImportMissingArticlesCollections()
	utils.Must(err)

	utils.Must(importArticlesContainers(csMap))

	publishersMap, err := createMissingPublishers()
	utils.Must(err)

	utils.Must(splitArticlesToPublications(publishersMap))

	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadAndImportMissingArticlesCollections() (map[int]*models.Collection, error) {

	cs, err := models.Collections(mdb,
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_ARTICLES].ID)).
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

	for _, kmID := range ArticlesCatalogs {
		if _, ok := csMap[kmID]; ok {
			continue
		}

		props := map[string]interface{}{
			"kmedia_id": kmID,
		}
		log.Infof("Create collection %v", props)
		c, err := api.CreateCollection(mdb, common.CT_ARTICLES, props)
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

func importArticlesContainers(csMap map[int]*models.Collection) error {
	cnMap, cuMap, err := loadContainersInCatalogsAndCUs(7957)
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
			// skip existing containers
			utils.Must(tx.Rollback())
			continue
		} else {
			// create - CU doesn't exist
			cu, err = importContainerWOCollectionNewCU(tx, cn, common.CT_ARTICLE)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Import new container %d", cnID)
			}
		}

		// create or update CCUs
		for i := range cn.R.Catalogs {
			catalog := cn.R.Catalogs[i]
			c, ok := csMap[catalog.ID]
			if !ok {
				continue
			}

			if cu.Published && !c.Published {
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
		if err := v.Update(mdb, "published"); err != nil {
			return errors.Wrapf(err, "Publish collection %d", v.ID)
		}
	}

	return nil
}

func createMissingPublishers() (map[string]*models.Publisher, error) {
	pMap := make(map[string]*models.Publisher)
	publishers, err := models.Publishers(mdb).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load publishers")
	}
	for i := range publishers {
		p := publishers[i]
		pMap[p.Pattern.String] = p
	}

	tx, err := mdb.Begin()
	utils.Must(err)

	for i := range publishersPatterns {
		pattern := publishersPatterns[i]
		if _, ok := pMap[pattern]; !ok {
			uid, err := api.GetFreeUID(tx, new(api.PublisherUIDChecker))
			if err != nil {
				utils.Must(tx.Rollback())
				return nil, errors.Wrapf(err, "generate uid %s", pattern)
			}

			p := models.Publisher{
				UID:     uid,
				Pattern: null.StringFrom(pattern),
			}

			err = p.Insert(tx)
			if err != nil {
				utils.Must(tx.Rollback())
				return nil, errors.Wrapf(err, "save publisher %s", pattern)
			}

			for _, lang := range [...]string{common.LANG_HEBREW, common.LANG_RUSSIAN, common.LANG_ENGLISH, common.LANG_SPANISH} {
				pI18n := models.PublisherI18n{
					PublisherID: p.ID,
					Language:    lang,
					Name:        null.StringFrom(pattern),
				}
				err = pI18n.Insert(tx)
				if err != nil {
					utils.Must(tx.Rollback())
					return nil, errors.Wrapf(err, "save publisher  I18n %s %s", pattern, lang)
				}
			}

			pMap[pattern] = &p
		}
	}

	utils.Must(tx.Commit())

	return pMap, nil
}

func splitArticlesToPublications(publishersMap map[string]*models.Publisher) error {
	cus, err := models.ContentUnits(mdb,
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_ARTICLE].ID),
		qm.Load("Files")).
		All()
	if err != nil {
		return errors.Wrap(err, "fetch content units")
	}

	log.Infof("%d article units", len(cus))

	for i := range cus {
		cu := cus[i]
		log.Infof("Article cuID %d", cu.ID)

		for j := range cu.R.Files {
			f := cu.R.Files[j]
			if f.Type != "image" {
				continue
			}

			n := strings.TrimSuffix(f.Name, ".zip")
			ns := strings.Split(n, "_")
			pattern := ns[len(ns)-1]
			log.Infof("%s\tpattern:%s %s", f.Name, pattern, f.Language.String)

			if x, ok := publishersPatternsMappings[pattern]; ok {
				pattern = x
			}

			publisher, ok := publishersMap[pattern]
			if !ok {
				log.Infof("Unknown publisher %s", pattern)
				continue
			}

			tx, err := mdb.Begin()
			utils.Must(err)

			pCU, err := api.CreateContentUnit(tx, common.CT_PUBLICATION, map[string]interface{}{
				"original_language": common.StdLang(f.Language.String),
			})
			if err != nil {
				utils.Must(tx.Rollback())
				log.Error(err)
				debug.PrintStack()
				continue
			}

			pCU.Published = f.Published
			err = pCU.Update(tx, "published")
			if err != nil {
				utils.Must(tx.Rollback())
				log.Error(err)
				debug.PrintStack()
				continue
			}

			err = pCU.AddPublishers(tx, false, publisher)
			if err != nil {
				utils.Must(tx.Rollback())
				log.Error(err)
				debug.PrintStack()
				continue
			}

			cud := &models.ContentUnitDerivation{
				SourceID: cu.ID,
				Name:     common.CT_PUBLICATION,
			}
			err = pCU.AddDerivedContentUnitDerivations(tx, true, cud)
			if err != nil {
				utils.Must(tx.Rollback())
				log.Error(err)
				debug.PrintStack()
				continue
			}

			f.ContentUnitID = null.Int64From(pCU.ID)
			err = f.Update(tx, "content_unit_id")
			if err != nil {
				utils.Must(tx.Rollback())
				log.Error(err)
				debug.PrintStack()
				continue
			}

			utils.Must(tx.Commit())
		}
	}

	return nil
}
