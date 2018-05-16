package tagging

import (
	"fmt"
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ExportTVShows() {
	clock := Init()

	utils.Must(doExportTVShows())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doExportTVShows() error {
	cs, err := models.Collections(mdb,
		qm.Where("type_id = ?", api.CONTENT_TYPE_REGISTRY.ByName[api.CT_VIDEO_PROGRAM].ID),
		qm.Load("CollectionsContentUnits",
			"CollectionsContentUnits.ContentUnit",
			"CollectionsContentUnits.ContentUnit.ContentUnitI18ns",
			"CollectionI18ns"),
	).All()
	if err != nil {
		return errors.Wrap(err, "Load collections")
	}

	log.Infof("%d TV Shows in MDB", len(cs))
	collections := make([]*CollectionWName, len(cs))
	for i := range cs {
		collections[i] = &CollectionWName{Collection: cs[i]}
	}

	sort.Slice(collections, func(i, j int) bool {
		return collections[i].Name() < collections[j].Name()
	})

	out := xlsx.NewFile()

	for i := range collections {
		c := collections[i]

		sheet := &xlsx.Sheet{}

		// sort units by position
		sort.Slice(c.R.CollectionsContentUnits, func(i, j int) bool {
			return c.R.CollectionsContentUnits[i].Position < c.R.CollectionsContentUnits[j].Position
		})

		for j := range c.R.CollectionsContentUnits {
			row := sheet.AddRow()
			cell := row.AddCell()
			cu := UnitWName{ContentUnit: c.R.CollectionsContentUnits[j].R.ContentUnit}
			cell.Value = cu.Name()

			cell = row.AddCell()
			url := fmt.Sprintf("https://archive.kbb1.com/programs/cu/%s", cu.UID)
			cell.SetStringFormula(fmt.Sprintf("HYPERLINK(\"%s\")", url))
		}

		name := CleanSheetName(fmt.Sprintf("%s (%s)", c.Name(), c.UID))
		log.Infof("%s", name)
		sheet, err = out.AppendSheet(*sheet, name)
		if err != nil {
			return errors.Wrapf(err, "out.AddSheet %d", i)
		}
	}

	err = out.Save("MDB_export_tagging_tvshows_rus.xlsx")
	if err != nil {
		return errors.Wrap(err, "out.Save")
	}

	return nil
}
