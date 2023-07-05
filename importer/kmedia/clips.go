package kmedia

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	qm4 "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var ClipsCatalogs = []int{
	3673, //	tv-clip
	4447, //	tv-clip/
	4661, //	tv-clip/channel-kabbalah_clip_RUS
	4660, //	tv-clip/channel-kabbalah_clips_ENG
	4659, //	tv-clip/channel-kabbalah_clips_HEB
	3959, //	tv-clip/channel-kabbalah_clips_HEB/Channel Kabbalah/Clips/60 sec al mashber olami
	3953, //	tv-clip/channel-kabbalah_clips_HEB/Channel Kabbalah/Clips/Begova einaim
	4663, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_aliha-al-dereh-emet
	4703, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_campus
	4665, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays
	4678, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_Haatzmaut
	4672, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_Hanuka
	4675, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_Pesach
	4674, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_Purim
	4669, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_roshhashana
	4676, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_Shavuot
	4671, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_sukkot
	4677, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_TishaaBeAv
	4673, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_TuBiShvat
	4670, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_Holidays/channel-kabbalah_clips_heb_holidays_YomKipur
	4662, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_no-program
	4664, //	tv-clip/channel-kabbalah_clips_HEB/channel-kabbalah_clips_HEB_rabot-mahshavot-balev-ish
	3954, //	tv-clip/channel-kabbalah_clips_HEB/Channel Kabbalah/Clips/Klipim muzikalim
	3951, //	tv-clip/channel-kabbalah_clips_HEB/Channel Kabbalah/Clips/Mtza hevdelim
	3952, //	tv-clip/channel-kabbalah_clips_HEB/clip_60-sec-al-kabalah
	4318, //	tv-clip/channel-kabbalah_clips_HEB/clip_teleconference_heb
	3935, //	tv-clip/Channel Kabbalah_plus/rus/Theater
	3868, //	tv-clip/channel-kabbalah/theater-bb
	7907, //	tv-clip/clip_30-sec
	7911, //	tv-clip/clip_ahshav-ani
	6935, //	tv-clip/clip_am-israel
	4754, //	tv-clip/clip_bait-shelanu
	6955, //	tv-clip/clip_chelovek-matritzi
	4756, //	tv-clip/clip_congress
	4770, //	tv-clip/clip_crossroads
	4738, //	tv-clip/clip_daily-lesson_2
	4758, //	tv-clip/clip_general
	7868, //	tv-clip/clip_globalniy-krizis
	6953, //	tv-clip/clip_goryachaya-tema
	4757, //	tv-clip/clip_gosudarstvo-narod
	8039, //	tv-clip/clip_gurman
	7914, //	tv-clip/clip_hadashot
	4751, //	tv-clip/clip_haim-hadashim
	//7891, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_ani-vehevra
	//7893, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_ani-ve-teva
	//7897, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_aruts_haim
	//7876, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_arvut
	//7883, //	tv-clip/clip_haim-hadashim/clip_haim_hadashim_bitahon
	//7879, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_briut
	//7896, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_gisha-haim
	//7886, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_hagim-ve-moadim
	//7877, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_herum
	//7885, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_hevra-israelit
	//7894, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_hohmat-hakabbalah
	//7888, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_horut-mishpaha
	//7895, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_kabbalah-ve-mistika
	//7881, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_kariera
	//7890, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_kehila
	//7878, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_kesef
	//7880, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_leida
	//7884, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_megamot-olamiyot
	//7892, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_mimshal
	//7889, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_osher
	//7887, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_tarbut-yehudit
	//7882, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_tikshoret
	//7875, //	tv-clip/clip_haim-hadashim/clip_haim-hadashim_zugiut
	8147, //	tv-clip/clip_haverim-shelanu
	4759, //	tv-clip/clip_integral-mir
	4762, //	tv-clip/clip_integral-social
	3906, //	tv-clip/clip_kabbalah-za-90min
	4775, //	tv-clip/clip_kadrovie-sekreti
	6962, //	tv-clip/clip_korotkie-istorii
	4744, //	tv-clip/clip_kurs-arvut
	7901, //	tv-clip/clip_lc
	6966, //	tv-clip/clip_lifnei-shina
	8126, //	tv-clip/clip_limud-atzmi
	7904, //	tv-clip/clip_masa-olami
	4774, //	tv-clip/clip_medizina-budushego
	8149, //	tv-clip/clip_mnenie
	7909, //	tv-clip/clip_mudrost-dnia
	4812, //	tv-clip/clip_mudrost-tolpi
	6951, //	tv-clip/clip_nifgashim-im-kabbalah
	7899, //	tv-clip/clip_novosti
	4823, //	tv-clip/clip_osnovi-kabbali
	4745, //	tv-clip/clip_pitaron-behibur
	8188, //	tv-clip/clip_poslednee-pokolenie
	3914, //	tv-clip/clip-rus
	3947, //	tv-clip/clips-eng
	4267, //	tv-clip/clip_sheela
	6939, //	tv-clip/clip_sila-knigi-zohar
	8305, //	tv-clip/clip_sipur-shelanu
	4765, //	tv-clip/clip_skvoz-vremia
	4517, //	tv-clip/clip_sprosi-kabbalista
	4746, //	tv-clip/clip_taini-vechnoy-knigi
	4439, //	tv-clip/clip_teleconference_eng
	4007, //	tv-clip/clip_umenia-zazvonil-telefon
	6950, //	tv-clip/clip_vavilon-vchera-isegodnia
	4444, //	tv-clip/program_contundente
	4445, //	tv-clip/program_senderos
	3909, //	tv-clip/promo
	4067, //	tv-clip/spa_clip
}

var collapsedClipsCatalogs = map[int]int{
	7891: 4751,
	7893: 4751,
	7897: 4751,
	7876: 4751,
	7883: 4751,
	7879: 4751,
	7896: 4751,
	7886: 4751,
	7877: 4751,
	7885: 4751,
	7894: 4751,
	7888: 4751,
	7895: 4751,
	7881: 4751,
	7890: 4751,
	7878: 4751,
	7880: 4751,
	7884: 4751,
	7892: 4751,
	7889: 4751,
	7887: 4751,
	7882: 4751,
	7875: 4751,
}

func ImportClips() {
	clock := Init()

	stats = NewImportStatistics()

	csMap, err := loadAndImportMissingClipsCollections()
	utils.Must(err)

	utils.Must(importClipsContainers(csMap))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadAndImportMissingClipsCollections() (map[int]*models.Collection, error) {

	cs, err := models.Collections(
		qm4.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_CLIPS].ID)).
		All(mdb)
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

	for _, kmID := range ClipsCatalogs {
		if _, ok := csMap[kmID]; ok {
			continue
		}

		props := map[string]interface{}{
			"kmedia_id": kmID,
		}
		log.Infof("Create collection %v", props)
		c, err := api.CreateCollection(mdb, common.CT_CLIPS, props)
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
					Name:         null.NewString(d.Name.String, d.Name.Valid),
				}
				err = ci18n.Upsert(mdb,
					true,
					[]string{"collection_id", "language"},
					boil.Whitelist("name"),
					boil.Infer())
				if err != nil {
					return nil, errors.Wrapf(err, "Upsert collection i18n, collection [%d]", c.ID)
				}
			}
		}
	}

	return csMap, nil
}

func importClipsContainers(csMap map[int]*models.Collection) error {
	cnMap, cuMap, err := loadContainersByTypeAndCUs(2)
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
			cu, err = importContainer(tx, cn, nil, common.CT_CLIP, "", 0)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Import existing container %d", cnID)
			}

			if common.CT_CLIP != common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name {
				cu.TypeID = common.CONTENT_TYPE_REGISTRY.ByName[common.CT_CLIP].ID
				_, err = cu.Update(tx, boil.Whitelist("type_id"))
				if err != nil {
					utils.Must(tx.Rollback())
					return errors.Wrapf(err, "Update CU type %d", cu.ID)
				}
			}
		} else {
			// create - CU doesn't exist
			cu, err = importContainerWOCollectionNewCU(tx, cn, common.CT_CLIP)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Import new container %d", cnID)
			}
		}

		// create or update CCUs
		for i := range cn.R.Catalogs {
			catalog := cn.R.Catalogs[i]

			catalogID := catalog.ID
			collapsedID, ok := collapsedClipsCatalogs[catalog.ID]
			if ok {
				catalogID = collapsedID
			}

			c, ok := csMap[catalogID]
			if !ok {
				continue
			}

			if collapsedID != 0 {
				log.Infof("Associating %d %s to collapsed catalog %d: [cu,c]=[%d,%d]", cn.ID, cn.Name.String, collapsedID, cu.ID, c.ID)
			} else {
				log.Infof("Associating %d %s to %d %s: [cu,c]=[%d,%d]", cn.ID, cn.Name.String, catalog.ID, catalog.Name, cu.ID, c.ID)
			}

			if cu.Published && !c.Published {
				csToPublish[c.ID] = c
			}

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
		if _, err := v.Update(mdb, boil.Whitelist("published")); err != nil {
			return errors.Wrapf(err, "Publish collection %d", v.ID)
		}
	}

	return nil
}
