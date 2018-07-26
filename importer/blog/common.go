package blog

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/models"
)

var (
	mdb *sql.DB
	currentBlog *models.Blog
	allBlogs map[int64]*models.Blog
)

func Init() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting blog import")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))


	blogs, err := models.Blogs(mdb).All()
	utils.Must(err)
	allBlogs = make(map[int64]*models.Blog, len(blogs))
	for i := range blogs {
		allBlogs[blogs[i].ID] = blogs[i]
	}

	currentBlog = allBlogs[1]
	log.Infof("current blog is %s", currentBlog.Name)

	return clock
}

func Shutdown() {
	utils.Must(mdb.Close())
}

func traverse(walkFunc filepath.WalkFunc) error {
	baseDir := fmt.Sprintf("importer/blog/data/%s", currentBlog.Name)

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("prevent panic by handling failure accessing a path %q: %v\n", baseDir, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		return walkFunc(path, info, err)
	})

	if err != nil {
		return errors.Wrap(err, "filepath.Walk")
	}

	return nil
}
