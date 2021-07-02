package cusource

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func RemoveFilesByFileName() {
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
		oIds := make([]int64, 0)
		for _, o := range f.R.Operations {
			oIds = append(oIds, o.ID)
		}
		err = deleteFileOnTransaction(mdb, f.ID, oIds)
		if err != nil {
			log.Debugf("File uid %s was not removed", f.UID)
			removeCount--
			continue
		}
	}
	log.Infof("All procceses ended. Removed %d files", removeCount)
}
func deleteFileOnTransaction(mdb *sql.DB, fId int64, oIds []int64) error {
	log.Info("Open transaction")
	tx, err := mdb.Begin()
	if err != nil {
		log.Errorf("Problem on create transaction: %s ", err)
		return err
	}

	log.Infof("Start remove files_operations with files ids:%v ", fId)
	qfo := fmt.Sprintf("DELETE FROM files_operations fo where  fo.file_id = %d", fId)
	_, err = queries.Raw(mdb, qfo).Exec()
	if err != nil {
		log.Errorf("Problem on delete files_operations: %s ", err)
		errR := tx.Rollback()
		if errR != nil {
			log.Errorf("Rollback error %s", errR)
			return errR
		}
		return err
	}

	log.Infof("Start remove operations with ids: %v ", oIds)

	if len(oIds) > 0 {
		err = models.Operations(mdb, qm.WhereIn("id IN ?", utils.ConvertArgsInt64(oIds)...)).DeleteAll()
		if err != nil {
			log.Errorf("Problem on delete operations: %s ", err)
			errR := tx.Rollback()
			if errR != nil {
				log.Errorf("Rollback error %s", errR)
				return errR
			}
			return err
		}
	}

	qfs := fmt.Sprintf("DELETE FROM files_storages fs where  fs.file_id = %d", fId)
	_, err = queries.Raw(mdb, qfs).Exec()
	if err != nil {
		log.Errorf("Problem on delete files_storages: %s ", err)
		errR := tx.Rollback()
		if errR != nil {
			log.Errorf("Rollback error %s", errR)
			return errR
		}
		return err
	}

	err = models.Files(mdb, qm.WhereIn("id = ?", fId)).DeleteAll()
	if err != nil {
		log.Errorf("Problem on delete files: %s ", err)
		errR := tx.Rollback()
		if errR != nil {
			log.Errorf("Rollback error %s", errR)
			return errR
		}
		return err
	}

	log.Info("Delete successful and committed transaction ")
	err = tx.Commit()
	if err != nil {
		log.Errorf("Commit error %s", err)
		errR := tx.Rollback()
		if errR != nil {
			log.Errorf("Rollback error %s", errR)
			return errR
		}
		return err
	}
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
