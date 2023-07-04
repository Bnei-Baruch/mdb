package mergelanguages

import (
	"database/sql"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type MergeLanguages struct {
	langFrom string
	langTo   string
	tx       *sql.Tx
}

func (m *MergeLanguages) Init(from, to string) {
	m.langFrom = from
	m.langTo = to
}
func (m *MergeLanguages) Run() {
	mdb := m.openDB()
	defer mdb.Close()
	var err error
	m.tx, err = mdb.Begin()
	utils.Must(err)
	// recover from panics in transaction
	defer func() {
		if p := recover(); p != nil {
			if ex := m.tx.Rollback(); ex != nil {
				log.Error("Couldn't roll back transaction")
			}
			panic(p) // re-throw panic after Rollback
		} else {
			utils.Must(m.tx.Commit())
		}
	}()
	err, evntsCU := m.changeCUI18ns()
	utils.Must(err)
	evnts := make([]events.Event, 0)
	evnts = append(evnts, evntsCU...)

	err, evntsF := m.changeFiles()
	utils.Must(err)
	evnts = append(evnts, evntsF...)

	emitter, err := events.InitEmitter()
	utils.Must(err)
	emitter.Emit(evnts...)
}

func (m *MergeLanguages) changeCUI18ns() (error, []events.Event) {
	evnts := make([]events.Event, 0)
	cuI18ns, err := models.ContentUnitI18ns(
		qm.InnerJoin(`files f ON f.content_unit_id = "content_unit_i18n".content_unit_id`),
		models.ContentUnitI18nWhere.Language.EQ(m.langFrom),
		qm.Load(models.ContentUnitI18nRels.ContentUnit),
	).All(m.tx)
	if err != nil {
		log.Fatalf("Can't load units: %s", err)
		return err, nil
	}
	_, err = cuI18ns.UpdateAll(m.tx, models.M{"language": m.langTo})
	utils.Must(err)
	for _, i18n := range cuI18ns {
		evnts = append(evnts, events.ContentUnitUpdateEvent(i18n.R.ContentUnit))
	}
	return nil, evnts
}

func (m *MergeLanguages) changeFiles() (error, []events.Event) {
	evnts := make([]events.Event, 0)
	files, err := models.Files(
		models.FileWhere.Language.EQ(null.StringFrom(m.langFrom)),
		qm.Load(models.ContentUnitI18nRels.ContentUnit),
	).All(m.tx)
	if err != nil {
		log.Fatalf("Can't load files: %s", err)
		return err, nil
	}
	_, err = files.UpdateAll(m.tx, models.M{"language": m.langTo})
	utils.Must(err)
	for _, f := range files {
		evnts = append(evnts, events.FileUpdateEvent(f))
	}
	return nil, evnts
}

func (m *MergeLanguages) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}
