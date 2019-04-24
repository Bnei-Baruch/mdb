package dgima

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
)

const LABELS_FILE = "importer/dgima/data/capture_label_id.csv"

func Import() {
	clock, _ := Init()

	utils.Must(doImport())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doImport() error {
	records, err := utils.ReadCSV(LABELS_FILE)
	if err != nil {
		return errors.Wrap(err, "Read labels")
	}
	log.Infof("Labels file has %d rows", len(records))

	for i := range records {
		f, _, err := api.FindFileBySHA1(mdb, records[i][0])
		if err != nil {
			log.Warnf("FindFileBySHA1 err %s [%s]", err.Error(), records[i][0])
			continue
		}

		err = api.UpdateFileProperties(mdb, f, map[string]interface{}{"label_id": records[i][1]})
		if err != nil {
			log.Errorf("UpdateFileProperties err %s [%s] => %s", err.Error(), records[i][0], records[i][1])
		}
	}

	return nil
}

