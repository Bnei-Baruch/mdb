package batch

import (
	"database/sql"
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/mdb/utils"
)

type PostRegexpReplacer struct {
	DB        *sql.DB
	RegStr    string
	NewStr    string
	Limit     int
	TableName string
	ColName   string
}

type entity struct {
	content string
	id      int64
}

func (a *PostRegexpReplacer) Init() error {

	if a.Limit == 0 {
		a.Limit = 100
	}

	if a.TableName == "" || a.ColName == "" || a.NewStr == "" || a.RegStr == "" {
		return fmt.Errorf("missing one of the fields: a.tableName = %s a.colName = %s a.NewStr = %s a.RegStr = %s", a.TableName, a.ColName, a.NewStr, a.RegStr)
	}

	mdb, err := sql.Open("postgres", viper.GetString("DB.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	a.DB = mdb
	return nil
}

func (a *PostRegexpReplacer) Run() {
	defer a.shutdown()
	log.Infof("Start replace %s to %s for regexp %s", a.RegStr, a.NewStr, a.RegStr)
	a.Run()
	log.Infof("End replace")
}

func (a *PostRegexpReplacer) Do() {
	var total int
	utils.Must(queries.Raw(a.DB, fmt.Sprintf("Select count(id) FROM %s", a.TableName)).QueryRow().Scan(&total))

	log.Debugf("Total blogs for replace %d", total)
	iterations := total / a.Limit

	for i := 0; i <= iterations; i++ {
		log.Infof("\nStart transaction")
		tx, err := a.DB.Begin()
		utils.Must(err)
		err = a.updateDB(tx, i)
		if err != nil {
			log.Errorf("Exception %v on iteration %d", err, i)
			utils.Must(tx.Rollback())
			continue
		}
		utils.Must(tx.Commit())
		log.Infof("Transaction committed")
	}
}

func (a *PostRegexpReplacer) updateDB(tx boil.Transactor, iteration int) error {
	log.Infof("Start replace on iteration %d", iteration)
	rows, err := queries.Raw(tx, fmt.Sprintf(`Select %s, id FROM %s ORDER BY id`, a.ColName, a.TableName)).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	update := make([]*entity, 0)
	for rows.Next() {
		e := &entity{}
		if err := rows.Scan(&e.content, &e.id); err != nil {
			return err
		}
		var ok bool
		if e.content, ok = a.replace(e.content); ok {
			update = append(update, e)
		}
	}

	for _, e := range update {
		q := fmt.Sprintf(`UPDATE %s SET %s = '%s' WHERE id = %d`, a.TableName, a.ColName, e.content, e.id)
		_, err = queries.Raw(tx, q).Exec()
		if err != nil {
			log.Errorf("Error on Update post with id %d error %v", e.id, err)
			return err
		}
		log.Infof("Successfully replace for post id %d", e.id)
	}
	log.Infof("End replace on iteration %d", iteration)
	return nil
}

func (a *PostRegexpReplacer) replace(cont string) (string, bool) {
	re := regexp.MustCompile(a.RegStr)
	if !re.MatchString(cont) {
		return "", false
	}
	cont = re.ReplaceAllString(cont, a.NewStr)
	log.Infof("New content %s", cont)
	return cont, true
}

func (a *PostRegexpReplacer) shutdown() {
	if err := a.DB.Close(); err != nil {
		log.Errorf("DB.close %v", err)
	}
}
