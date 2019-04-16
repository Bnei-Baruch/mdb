package kabcoil

import (
	"net/url"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportTitles() {
	clock := Init()

	utils.Must(doImportTitles("importer/kabcoil/data/yoni_programs_titles_old.xlsx"))
	utils.Must(doImportTitles("importer/kabcoil/data/yoni_programs_titles_new.xlsx"))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

type UnitTitle struct {
	Link        string
	Name        string
	Description string
	rIdx        int
	cu          *models.ContentUnit
}

func doImportTitles(path string) error {
	log.Infof("Processing %s", path)

	xlFile, err := excelize.OpenFile(path)
	if err != nil {
		return errors.Wrapf(err, "xlsx.OpenFile: %s", path)
	}

	data := make(map[string][]*UnitTitle)

	wCU := 0
	woCU := 0

	for index, name := range xlFile.GetSheetMap() {
		rows, err := xlFile.GetRows(name)
		if err != nil {
			return errors.Wrapf(err, "xlFile.GetRows [%d] %s", index, name)
		}

		titles := make([]*UnitTitle, 0)

		for rIdx, row := range rows {
			if len(row) < 2 {
				continue
			}

			link := strings.TrimSpace(row[1])
			p, err := url.ParseRequestURI(link)
			if err != nil {
				//log.Infof("Sheet %s row %d invalid url: %s", sheet.Name, rIdx, link)
			} else if !strings.HasPrefix(p.Host, "files") {
				log.Infof("Sheet %s row %d bad host: %s", name, rIdx, p.Host)
			} else {
				name := ""
				if len(row) > 4 {
					name = strings.TrimSpace(row[4])
				}

				description := ""
				if len(row) > 5 {
					description = strings.TrimSpace(row[5])
				}

				if name == "" && description == "" {
					//log.Infof("Sheet %s row %d no values", sheet.Name, rIdx)
				} else {
					cu, err := linkToCU(link)
					if err != nil {
						log.Errorf("linkToCU: [%d] : %s", rIdx, err.Error())
						woCU++
						continue
					}

					wCU++
					titles = append(titles, &UnitTitle{
						Link:        link,
						Name:        name,
						Description: description,
						rIdx:        rIdx,
						cu:          cu,
					})
				}
			}
		}

		data[name] = titles
	}

	log.Infof("Data has %d entries (%d sheets)", len(data), xlFile.SheetCount)
	log.Infof("woCU: %d", woCU)
	log.Infof("wCU: %d", wCU)
	for k, v := range data {
		log.Infof("Sheet %s has %d valid entries", k, len(v))
		for i := range v {
			if err := updateCU(v[i]); err != nil {
				log.Errorf("updateCU: [row %d]: %s", i, err.Error())
			}
		}
	}

	return nil
}

func linkToCU(link string) (*models.ContentUnit, error) {
	s := strings.Split(link, "/")
	fname := s[len(s)-1]

	kmFile, err := kmodels.FileAssets(kmdb, qm.Where("name = ?", fname)).One()
	if err != nil {
		return nil, errors.Wrapf(err, "Find KM file %s", fname)
	}

	mFile, err := models.Files(mdb, qm.Where("(properties->>'kmedia_id')::int = ?", kmFile.ID)).One()
	if err != nil {
		return nil, errors.Wrapf(err, "Find MDB file %d", kmFile.ID)
	}

	if mFile.ContentUnitID.Valid {
		err = mFile.L.LoadContentUnit(mdb, true, mFile)
		if err != nil {
			return nil, errors.Wrapf(err, "mFile.L.LoadContentUnit %d", mFile.ContentUnitID.Int64)
		}
		return mFile.R.ContentUnit, nil
	}

	return nil, nil
}

func updateCU(ut *UnitTitle) error {
	i18n := models.ContentUnitI18n{
		ContentUnitID: ut.cu.ID,
		Language:      api.LANG_HEBREW,
		Name:          null.StringFrom(ut.Name),
		Description:   null.StringFrom(ut.Description),
	}

	err := i18n.Upsert(mdb, true,
		[]string{"content_unit_id", "language"}, []string{"name", "description"})
	if err != nil {
		return errors.Wrapf(err, "i18n.Upsert %d", ut.cu.ID)
	}

	return nil
}
