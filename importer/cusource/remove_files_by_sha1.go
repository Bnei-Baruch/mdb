package cusource

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io"
	"os"
	"strconv"
	"strings"
)

func RemoveFilesBySHA1() {
	log.SetLevel(log.InfoLevel)
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	defer mdb.Close()
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))

	run(mdb)
}

func run(mdb *sql.DB) {
	names := getNames()
	log.Infof("Take file names from tar: %v\n", names)
	// actual removal
	files, err := models.Files(mdb,
		qm.WhereIn("name in ?", utils.ConvertArgsString(names)...),
		qm.Load("Operations"),
	).All()
	utils.Must(err)
	log.Infof("Found %d files for remove", len(files))
	removeCount := len(files)
	for _, f := range files {
		oIds := make([]string, 0)
		for _, o := range f.R.Operations {
			oIds = append(oIds, strconv.FormatInt(o.ID, 10))
		}
		err = deleteFileOnTransaction(mdb, f.ID, oIds)
		if err != nil {
			log.Infof("File uid %s was not removed", f.UID)
			removeCount--
			continue
		}
	}
	log.Infof("All procceses ended. Removed %d files", removeCount)
}
func deleteFileOnTransaction(mdb *sql.DB, fId int64, oIds []string) error {
	log.Info("Open transaction")
	tx, err := mdb.Begin()
	defer func() {
		if err == nil {
			log.Info("Send commit of transaction ")
			err = tx.Commit()
			if err != nil {
				log.Infof("Commit error %s", err)
			}
		} else {
			log.Info("Send Rollback of transaction ")
			err = tx.Rollback()
			if err != nil {
				log.Infof("Commit error %s", err)
			}
		}
	}()

	if err != nil {
		log.Infof("Problem on create transaction: %s ", err)
		return err
	}

	log.Infof("Start remove files_operations with files ids:%v ", fId)
	qfo := fmt.Sprintf("DELETE FROM files_operations fo where  fo.file_id = %d", fId)
	rfo, err := queries.Raw(mdb, qfo).Query()
	if err != nil {
		log.Infof("Problem on delete files_operations: %s ", err)
		return err
	}
	defer rfo.Close()

	log.Infof("Start remove operations with ids: %v ", oIds)

	if len(oIds) > 0 {
		qo := fmt.Sprintf("DELETE FROM operations o where  o.id IN (%s)", strings.Join(oIds, ","))
		ro, err := queries.Raw(mdb, qo).Query()
		if err != nil {
			log.Infof("Problem on delete operations: %s ", err)
			return err
		}
		defer ro.Close()
	}

	qfs := fmt.Sprintf("DELETE FROM files_storages fs where  fs.file_id = %d", fId)
	rfs, err := queries.Raw(mdb, qfs).Query()
	if err != nil {
		log.Infof("Problem on delete files_storages: %s ", err)
		return err
	}
	defer rfs.Close()

	qf := fmt.Sprintf("DELETE FROM files f where  f.id = %d", fId)
	rf, err := queries.Raw(mdb, qf).Query()
	if err != nil {
		log.Infof("Problem on delete files: %s ", err)
		return err
	}
	defer rf.Close()

	log.Info("Delete successful and committed transaction ")
	return nil
}

func getNames() []string {
	path := viper.GetString("source-import.source-dir")
	r, err := os.Open(path)
	utils.Must(err)
	gzr, err := gzip.NewReader(r)
	utils.Must(err)
	defer utils.Must(gzr.Close())

	tr := tar.NewReader(gzr)
	result := make([]string, 0)

	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Must(err)
		}

		if isDoc := strings.Contains(header.Name, ".doc"); header.Typeflag == tar.TypeReg && isDoc {
			spl := strings.Split(header.Name, "/")
			result = append(result, spl[len(spl)-1])
		}
	}
	return result
}
