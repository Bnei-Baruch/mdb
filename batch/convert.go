package batch

import (
	"database/sql"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

/*

-- batch convert

DROP TABLE IF EXISTS batch_convert;
CREATE TABLE batch_convert (
  file_id       BIGINT REFERENCES files                          NOT NULL,
  operation_id  BIGINT REFERENCES operations                     NULL,
  request_at    TIMESTAMP WITH TIME ZONE                         NULL,
  request_error TEXT                         					 NULL
);

*/

var MT_MP4 string
var MT_WMV string
var MT_FLV string

func PrepareFilesForConvert() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting prepare files for convert")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	log.Info("Initializing static data from MDB")
	utils.Must(api.InitTypeRegistries(mdb))

	MT_MP4 = api.MEDIA_TYPE_REGISTRY.ByExtension["mp4"].MimeType
	MT_WMV = api.MEDIA_TYPE_REGISTRY.ByExtension["wmv"].MimeType
	MT_FLV = api.MEDIA_TYPE_REGISTRY.ByExtension["flv"].MimeType

	log.Info("Loading video files")
	files, err := models.Files(mdb, qm.Where("type=?", "video")).All()
	utils.Must(err)
	log.Infof("Got %d video files", len(files))

	utils.Must(populateDBQueue(files))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func populateDBQueue(files []*models.File) error {
	cuMap := make(map[int64]map[string][]*models.File)
	noCU := make([]*models.File, 0)
	withChilds := make(map[int64]bool)
	var ok bool
	for i := range files {
		f := files[i]

		// we only care about wmv, flv and mp4
		switch f.MimeType.String {
		case MT_MP4, MT_WMV, MT_FLV:
		default:
			continue
		}

		if f.ParentID.Valid {
			withChilds[f.ParentID.Int64] = true
		}

		if f.ContentUnitID.Valid {
			var byMime map[string][]*models.File
			if byMime, ok = cuMap[f.ContentUnitID.Int64]; !ok {
				byMime = make(map[string][]*models.File)
				cuMap[f.ContentUnitID.Int64] = byMime
			}
			byMime[f.MimeType.String] = append(byMime[f.MimeType.String], f)
		} else {
			noCU = append(noCU, f)
		}
	}

	log.Infof("noCU: %d", len(noCU))
	log.Infof("cuMap: %d", len(cuMap))
	log.Infof("withChilds: %d", len(withChilds))

	cmp4 := 0
	cwmv := 0
	cflv := 0
	for _, v := range cuMap {
		cmp4 += len(v[MT_MP4])
		cwmv += len(v[MT_WMV])
		cflv += len(v[MT_FLV])
	}

	log.Infof("cmp4: %d", cmp4)
	log.Infof("cwmv: %d", cwmv)
	log.Infof("cflv: %d", cflv)

	var finalFiles = make([]*models.File, 0)
	skippedFlv := 0
	skippedWmv := 0
	for _, v := range cuMap {
		mp4Files := v[MT_MP4]
		mp4s := make(map[string]*models.File, len(mp4Files))
		for i := range mp4Files {
			s := strings.Split(mp4Files[i].Name, ".")
			mp4s[strings.Join(s[0:len(s)-1], ".")] = mp4Files[i]
		}

		wmvFiles := v[MT_WMV]
		wmvs := make(map[string]*models.File, len(wmvFiles))
		for i := range wmvFiles {
			s := strings.Split(wmvFiles[i].Name, ".")
			wmvs[strings.Join(s[0:len(s)-1], ".")] = wmvFiles[i]
		}

		flvFiles := v[MT_FLV]
		flvs := make(map[string]*models.File, len(flvFiles))
		for i := range flvFiles {
			s := strings.Split(flvFiles[i].Name, ".")
			flvs[strings.Join(s[0:len(s)-1], ".")] = flvFiles[i]
		}

		// flvs first
		for k, v := range flvs {
			if _, ok := withChilds[v.ID]; ok {
				//log.Infof("Skip %s it has children", v.Name)
				skippedFlv += 1
			} else if _, ok := mp4s[k]; ok {
				//log.Infof("Skip %s for %s", v.Name, mp4.Name)
				skippedFlv += 1
			} else {
				//log.Infof("Adding file [flv]: %s [%d]", v.Name, v.ID)
				finalFiles = append(finalFiles, v)
			}
		}

		// wmvs second
		for k, v := range wmvs {
			if _, ok := withChilds[v.ID]; ok {
				//log.Infof("Skip %s it has children", v.Name)
				skippedWmv += 1
			} else if _, ok := mp4s[k]; ok {
				//log.Infof("Skip %s for %s", v.Name, mp4.Name)
				skippedWmv += 1
			} else if _, ok := flvs[k]; ok {
				//log.Infof("Skip %s for %s", v.Name, flv.Name)
				skippedWmv += 1
			} else {
				//log.Infof("Adding file [wmv]: %s [%d]", v.Name, v.ID)
				finalFiles = append(finalFiles, v)
			}
		}
	}

	log.Infof("len(finalFiles): %d", len(finalFiles))
	log.Infof("skippedFlv: %d", skippedFlv)
	log.Infof("skippedWmv: %d", skippedWmv)

	return fillQueue(finalFiles)
}

func fillQueue(files []*models.File) error {
	tx, err := mdb.Begin()
	if err != nil {
		return errors.Wrap(err, "Begin transaction")
	}

	stmt, err := tx.Prepare(pq.CopyIn("batch_convert", "file_id"))
	if err != nil {
		return errors.Wrap(err, "Prepare statement")
	}

	for i := range files {
		f := files[i]
		_, err = stmt.Exec(f.ID)
		if err != nil {
			return errors.Wrapf(err, "Exec [%d]", f.ID)
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "Final exec")
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "Close statement")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "Commit transaction")
	}

	return nil
}

func QueueWork() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting queue work")

	log.Info("Setting up connection to MDB")
	mdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	utils.Must(doQueueWork())

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doQueueWork() error {
	pageSize := 100
	query := `select f.id, encode(f.sha1,'hex') from batch_convert bc inner join files f on bc.file_id = f.id and bc.request_at is null and request_error is null limit $1`
	rows, err := queries.Raw(mdb, query, pageSize).Query()
	if err != nil {
		return errors.Wrap(err, "Load page")
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var sha1 string
		if err := rows.Scan(&id, &sha1); err != nil {
			return errors.Wrap(err, "rows.Scan")
		}

		log.Infof("queueing %d %s", id, sha1)
		if err := queueFile(id, sha1); err != nil {
			return err
		}
	}
	if err = rows.Err(); err != nil {
		return errors.Wrap(err, "rows.Err")
	}

	return nil
}

type TranscodeRequest struct {
	Sha1   string `json:"sha1"`
	Format string `json:"format"`
}

func queueFile(id int64, sha1 string) error {
	url := "http://files.kabbalahmedia.info/api/v1/transcode"

	resp, err := utils.HttpPostJson(url, TranscodeRequest{Sha1: sha1, Format: "mp4"})
	if err != nil {
		return errors.Wrap(err, "call conversion service")
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "Read response body")
		}
		log.Warnf("HTTP Error [%d]: %s", resp.StatusCode, string(b))

		_, err = queries.Raw(mdb, "update batch_convert set request_at=utc_now(), request_error=$1 where file_id=$2", string(b), id).Exec()
		if err != nil {
			return errors.Wrapf(err, "Update db queue request_error [%d]: %s", id, string(b))
		}
	} else {
		_, err := queries.Raw(mdb, "update batch_convert set request_at=utc_now() where file_id=$1", id).Exec()
		if err != nil {
			return errors.Wrapf(err, "Update db queue [%d]", id)
		}
	}

	return nil
}
