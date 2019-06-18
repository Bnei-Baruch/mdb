package kmedia

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func Compare() {
	clock := Init()

	stats = NewImportStatistics()

	log.Info("Loading kmedia catalogs hierarchy")
	utils.Must(loadCatalogHierarchy())

	//log.Info("Compare collections")
	//utils.Must(compareCollections())

	//log.Info("Compare units")
	//utils.Must(compareUnits())

	log.Info("Missing containers")
	//_, err := missingContainers()
	missing, err := missingContainers()
	utils.Must(err)

	declamations := map[string][]*kmodels.Container{
		"7": missing["7"],
	}

	dumpMissingContainers(declamations)
	//dumpMissingContainers(missing)

	utils.Must(importMissingContainers(declamations))
	//utils.Must(importMissingContainers(missing))

	stats.dump()

	Shutdown()

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

type Grouper interface {
	Match(*kmodels.Container) bool
	String() string
}

type ContentTypeGrouper struct {
	TypeID int
}

func NewContentTypeGrouper(typeID int) *ContentTypeGrouper {
	g := new(ContentTypeGrouper)
	g.TypeID = typeID
	return g
}

func (g *ContentTypeGrouper) String() string {
	return fmt.Sprintf("%d", g.TypeID)
}

func (g *ContentTypeGrouper) Match(cn *kmodels.Container) bool {
	return g.TypeID == cn.ContentTypeID.Int && (g.TypeID != 7 || cn.ID != 24502)
}

type RegexGrouper struct {
	re   *regexp.Regexp
	name string
}

func NewRegexGrouper(name, re string) *RegexGrouper {
	g := new(RegexGrouper)
	g.name = name
	g.re = regexp.MustCompile(re)
	return g
}

func (g *RegexGrouper) String() string {
	return g.name
}

func (g *RegexGrouper) Match(cn *kmodels.Container) bool {
	return g.re.MatchString(cn.Name.String)
}

type PredicateGrouper struct {
	fn   func(cn *kmodels.Container) bool
	name string
}

func NewPredicateGrouper(name string, fn func(cn *kmodels.Container) bool) *PredicateGrouper {
	g := new(PredicateGrouper)
	g.name = name
	g.fn = fn
	return g
}

func (g *PredicateGrouper) String() string {
	return g.name
}

func (g *PredicateGrouper) Match(cn *kmodels.Container) bool {
	return g.fn(cn)
}

var groupers = []Grouper{
	NewContentTypeGrouper(2),  // clip
	NewContentTypeGrouper(3),  // song
	NewContentTypeGrouper(5),  // lecture
	NewContentTypeGrouper(6),  // book
	NewContentTypeGrouper(7),  // declamation
	NewContentTypeGrouper(10), // text

	NewRegexGrouper("rabash", "(?i)(^rb_)|(_rb_)"),
	NewRegexGrouper("YH", "(?i)yeshivat-haverim"),
	NewRegexGrouper("maamar", "(?i)Maamar_zohoraim|mzohoraim"),
	NewRegexGrouper("hodaot", "(?i)hodaot"),
	NewRegexGrouper("lesson summary", "(?i)lesson-summary"),
	NewRegexGrouper("ulpan", "(?i)UlpanIvrit"),
	NewRegexGrouper("websites", "(?i)website"),

	//NewPredicateGrouper("misc - public", func(cn *kmodels.Container) bool {
	//	return cn.Secure == 0
	//}),
	//NewPredicateGrouper("misc - private", func(cn *kmodels.Container) bool {
	//	return cn.Secure != 0
	//}),
}

func missingContainers() (map[string][]*kmodels.Container, error) {
	cuMap := make(map[int]bool)

	log.Info("Loading all content units with kmedia_id")
	units, err := models.ContentUnits(mdb,
		qm.Where("properties -> 'kmedia_id' is not null")).
		All()
	if err != nil {
		return nil, errors.Wrap(err, "Load content_units from mdb")
	}
	log.Infof("Got %d units", len(units))

	for i := range units {
		cu := units[i]
		var props map[string]interface{}
		err := json.Unmarshal(cu.Properties.JSON, &props)
		if err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal unit properties %d", cu.ID)
		}
		cuMap[int(props["kmedia_id"].(float64))] = true
	}

	containers, err := kmodels.Containers(kmdb).All()
	if err != nil {
		return nil, errors.Wrap(err, "Load containers")
	}

	missing := make(map[string][]*kmodels.Container, 0)
	for i := range containers {
		cn := containers[i]
		if _, ok := cuMap[cn.ID]; !ok {
			err := cn.L.LoadFileAssets(kmdb, true, cn)
			if err != nil {
				return nil, errors.Wrapf(err, "Load File Assets cnID %d", cn.ID)
			}

			if len(cn.R.FileAssets) == 0 {
				continue
			}

			var key string
			for j := range groupers {
				if groupers[j].Match(cn) {
					key = groupers[j].String()
					break
				}
			}
			if key == "" {
				key = "misc"
			}

			if v, ok := missing[key]; ok {
				missing[key] = append(v, cn)
			} else {
				missing[key] = []*kmodels.Container{cn}
			}
		}
	}

	for _, v := range missing {
		sort.Slice(v, func(i, j int) bool {
			a := v[i]
			b := v[j]

			if a.ContentTypeID.Int != b.ContentTypeID.Int {
				return a.ContentTypeID.Int < b.ContentTypeID.Int
			}

			if !a.Filmdate.Time.Equal(b.Filmdate.Time) {
				return a.Filmdate.Time.Before(b.Filmdate.Time)
			}

			return a.Name.String < b.Name.String
		})
	}

	f, err := ioutil.TempFile("/tmp", "kmedia_compare")
	if err != nil {
		return nil, errors.Wrap(err, "Create temp file")
	}
	defer f.Close()
	log.Infof("Report file: %s", f.Name())

	names := make([]string, len(groupers)+1)

	for i := range groupers {
		names[i] = groupers[i].String()
	}
	names[len(names)-1] = "misc"

	for i := range names {
		v := missing[names[i]]
		fmt.Fprintf(f, "%s\t\t\t%d\n", names[i], len(v))
	}

	for i := range names {
		v := missing[names[i]]
		fmt.Fprintf(f, "\n\n#################\n%s\t%d\n#################\n\n", names[i], len(v))
		for i := range v {
			cn := v[i]
			_, err := fmt.Fprintf(f, "%d\t%d\t%s\t%s\t%d\t%d\n",
				cn.ID,
				cn.ContentTypeID.Int,
				cn.Filmdate.Time.Format("2006-01-02"),
				cn.Name.String,
				cn.VirtualLessonID.Int,
				cn.Secure,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "Write tsv row %d", cn.ID)
			}
		}
	}

	return missing, nil
}

func dumpMissingContainers(missing map[string][]*kmodels.Container) {
	log.Info("Here comes missing containers")
	for k, v := range missing {
		log.Infof("\n\ntype %s [%d]:", k, len(v))
		for i := range v {
			log.Infof("%d\t%s\t%s", v[i].ID, v[i].Filmdate.Time, v[i].Name.String)
		}
	}
}

func importMissingContainers(missing map[string][]*kmodels.Container) error {

	// clips
	cns := missing["2"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_CLIP)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import clip %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// songs
	cns = missing["3"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_SONG)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import song %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// lectures
	cns = missing["5"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_LECTURE)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import lecture %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// books
	cns = missing["6"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_BOOK)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import book %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// declamation
	cns = missing["7"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_BLOG_POST)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import declamation %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// texts
	collection, err := models.FindCollection(mdb, 11825)
	utils.Must(err)
	cns = missing["10"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_UNKNOWN)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import text %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// rabash
	cns = missing["rabash"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		collection, err = api.CreateCollection(tx, common.CT_DAILY_LESSON, map[string]interface{}{
			"film_date":         "1970-01-01",
			"original_language": common.LANG_HEBREW,
			"kmedia_rabash":     true,
		})
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Create rabash lesson %d", cns[i].ID)
		}

		err = importContainerWCollection(tx, cns[i], collection, common.CT_LESSON_PART)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import rabash %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// YH
	cns = missing["YH"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		stats.ContainersProcessed.Inc(1)
		_, err = importContainerWOCollectionNewCU(tx, cns[i], common.CT_FRIENDS_GATHERING)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import YH %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// maamar
	collection, err = models.FindCollection(mdb, 11821)
	utils.Must(err)
	cns = missing["maamar"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_CLIP)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import maamar %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// hodaot
	collection, err = models.FindCollection(mdb, 11822)
	utils.Must(err)
	cns = missing["hodaot"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_CLIP)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import hodaot %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// lesson summary
	collection, err = models.FindCollection(mdb, 11823)
	utils.Must(err)
	cns = missing["lesson summary"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_CLIP)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import lesson summary %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// ulpan ivrit
	collection, err = models.FindCollection(mdb, 11824)
	utils.Must(err)
	cns = missing["ulpan"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_VIDEO_PROGRAM_CHAPTER)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import ulpan ivrit %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// websites
	collection, err = models.FindCollection(mdb, 11820)
	utils.Must(err)
	cns = missing["websites"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(tx, cns[i], collection, common.CT_UNKNOWN)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import website %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	// misc
	collection, err = models.FindCollection(mdb, 11826)
	utils.Must(err)
	cns = missing["misc"]
	for i := range cns {
		tx, err := mdb.Begin()
		utils.Must(err)

		err = importContainerWCollection(mdb, cns[i], collection, common.CT_UNKNOWN)
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Import misc %d", cns[i].ID)
		}
		utils.Must(tx.Commit())
	}

	return nil
}

func compareCollection(c *models.Collection) error {
	var props map[string]interface{}
	if err := json.Unmarshal(c.Properties.JSON, &props); err != nil {
		return errors.Wrap(err, "json.Unmarshal properties")
	}

	kmid := int(props["kmedia_id"].(float64))

	isLesson := common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID == c.TypeID ||
		common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID == c.TypeID

	// lessons are compared to virtual lesson all others are compared to catalog
	var containers []*kmodels.Container
	if isLesson {
		log.Infof("Compare collection %d [%s] to virtual_lesson %d", c.ID, common.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name, kmid)

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
		log.Infof("Compare collection %d [%s] to catalog %d", c.ID, common.CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name, kmid)

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
			if common.CONTENT_TYPE_REGISTRY.ByName[common.CT_FULL_LESSON].ID == cu.TypeID &&
				(common.CONTENT_TYPE_REGISTRY.ByName[common.CT_DAILY_LESSON].ID == c.TypeID ||
					common.CONTENT_TYPE_REGISTRY.ByName[common.CT_SPECIAL_LESSON].ID == c.TypeID) {
				continue
			} else {
				log.Infof("CU exists only in MDB %d [%s]", cu.ID, common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name)
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
				return errors.Wrapf(err, "check File exists by sha1 %s", k)
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
