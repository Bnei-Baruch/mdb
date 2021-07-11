package batch

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func PostsHttpToHttps() {
	domain := "youtube"
	log.Infof("Start replace http to https for domain %s", domain)
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()

	oldStr := "http://www." + domain
	newStr := "https://www." + domain
	utils.Must(logCounts(mdb, newStr, oldStr))

	q := fmt.Sprintf("UPDATE \"blog_posts\" SET content = REPLACE(content, '%s', '%s')", oldStr, newStr)
	_, err = mdb.Exec(q)
	utils.Must(err)

	log.Infof("Replace event was ended")
	utils.Must(logCounts(mdb, newStr, oldStr))

	log.Infof("End replace http to https for domain %s", domain)
}

func logCounts(mdb boil.Executor, newStr, oldStr string) error {

	forReplace, err := models.BlogPosts(mdb,
		qm.Where(fmt.Sprintf("content ~ '.*%s*.'", oldStr)),
	).Count()
	if err != nil {
		return err
	}
	log.Debugf("Need to replace %d posts", forReplace)

	notReplace, err := models.BlogPosts(mdb,
		qm.Where(fmt.Sprintf("content ~ '.*%s*.'", newStr)),
	).Count()

	if err != nil {
		return err
	}
	log.Debugf("Not need to replace: %d posts. Total: %d", notReplace, notReplace+forReplace)
	return nil
}
