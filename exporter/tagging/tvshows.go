package tagging

import (
	"fmt"
	"sort"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
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
	cs, err := models.Collections(
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_VIDEO_PROGRAM].ID),
		qm.Load("CollectionsContentUnits"),
		qm.Load("CollectionsContentUnits.ContentUnit"),
		qm.Load("CollectionsContentUnits.ContentUnit.ContentUnitI18ns"),
		qm.Load("CollectionI18ns"),
	).All(mdb)
	if err != nil {
		return errors.Wrap(err, "Load collections")
	}

	log.Infof("%d TV Shows in MDB", len(cs))
	collections := make([]*common.CollectionWName, len(cs))
	for i := range cs {
		collections[i] = &common.CollectionWName{Collection: cs[i]}
	}

	sort.Slice(collections, func(i, j int) bool {
		return collections[i].Name() < collections[j].Name()
	})

	out := excelize.NewFile()

	for i := range collections {
		c := collections[i]

		name := CleanSheetName(fmt.Sprintf("%s (%s)", c.Name(), c.UID))
		log.Infof("%s", name)
		out.NewSheet(name)

		// sort units by position
		sort.Slice(c.R.CollectionsContentUnits, func(i, j int) bool {
			return c.R.CollectionsContentUnits[i].Position < c.R.CollectionsContentUnits[j].Position
		})

		for j := range c.R.CollectionsContentUnits {
			cu := common.UnitWName{ContentUnit: c.R.CollectionsContentUnits[j].R.ContentUnit}

			url := fmt.Sprintf("https://archive.kbb1.com/programs/cu/%s", cu.UID)
			out.SetCellHyperLink(name, fmt.Sprintf("A%d", j+1), url, "External")
			out.SetCellStr(name, fmt.Sprintf("B%d", j+1), cu.Name())
			out.SetCellStr(name, fmt.Sprintf("C%d", j+1), cu.Description())
		}
	}

	err = out.SaveAs("MDB_export_tagging_tvshows.xlsx")
	if err != nil {
		return errors.Wrap(err, "out.SaveAs")
	}

	return nil
}
