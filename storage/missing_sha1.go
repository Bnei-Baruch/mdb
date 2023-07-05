package storage

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/Bnei-Baruch/mdb/utils"
)

func MissingSha1Analysis() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting import of storage status")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	fileMap, err := loadFilesByKMID(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Load files from MDB"))
	}

	utils.Must(analyzeSha1CSV(fileMap, mdb))
	utils.Must(analyzeMissingSha1CSV(fileMap, mdb))

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadFilesByKMID(db *sql.DB) (map[int64]string, error) {
	fileMap := make(map[int64]string, 500000)

	rows, err := queries.Raw(`SELECT (properties->>'kmedia_id')::int, encode(sha1, 'hex') FROM files WHERE sha1 IS NOT NULL AND properties ? 'kmedia_id'`).
		Query(db)
	if err != nil {
		return nil, errors.Wrap(err, "Load files")
	}

	for rows.Next() {
		var kmID int64
		var sha1 string
		err = rows.Scan(&kmID, &sha1)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}

		fileMap[kmID] = sha1
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Iterate rows")
	}
	err = rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Close rows")
	}

	log.Info("Loaded files from MDB")
	return fileMap, nil
}

func analyzeSha1CSV(fileMap map[int64]string, mdb *sql.DB) error {
	input := "/home/edos/projects/kmedia/kmedia_files_sha1.csv"
	log.Infof("Processing data file: %s", input)
	f, err := os.Open(input)
	if err != nil {
		return errors.Wrap(err, "Open data file")
	}
	defer f.Close()

	i := 0
	csvReader := csv.NewReader(f)
	for {
		r, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrapf(err, "csv.Read() %d", i)
		}

		i++

		kmid, err := strconv.Atoi(r[0])
		if err != nil {
			return errors.Wrapf(err, "Bad KMID %s line %d", r[0], i)
		}
		sha1 := r[4]

		dbSha1, ok := fileMap[int64(kmid)]
		if !ok {
			//log.Infof("Unknown kmid %d", kmid)
			continue
		}

		if sha1 == dbSha1 {
			continue
		}

		log.Infof("SHA1 mismatch: %d, %s != %s", kmid, sha1, dbSha1)
		_, err = queries.Raw("UPDATE files set sha1=decode($1,'hex') where (properties->>'kmedia_id')::int = $2", sha1, kmid).Exec(mdb)
		if err != nil {
			log.Errorf("Update sha1 %d line %d: %s", kmid, i, err.Error())
		}
	}

	//i := 0
	//scanner := bufio.NewScanner(f)
	//for scanner.Scan() {
	//	i++
	//	line := strings.TrimSpace(scanner.Text())
	//	if line == "" {
	//		continue
	//	}
	//
	//	r := strings.Split(line, ",")
	//	kmid, err := strconv.Atoi(r[0])
	//	if err != nil {
	//		return errors.Wrapf(err, "Bad KMID %s", r[0])
	//	}
	//	sha1 := r[4]
	//
	//	dbSha1, ok := fileMap[int64(kmid)]
	//	if !ok {
	//		//log.Infof("Unknown kmid %d", kmid)
	//		continue
	//	}
	//
	//	if sha1 == dbSha1 {
	//		continue
	//	}
	//
	//	log.Infof("SHA1 mismatch: %d, %s != %s", kmid, sha1, dbSha1)
	//	_, err = queries.Raw(mdb, "UPDATE files set sha1=decode($1,'hex') where (properties->>'kmedia_id')::int = $2", sha1, kmid).Exec()
	//	if err != nil {
	//		return errors.Wrapf(err, "Update sha1 %d", kmid)
	//	}
	//}

	return nil
}

func analyzeMissingSha1CSV(fileMap map[int64]string, mdb *sql.DB) error {
	input := "/home/edos/projects/kmedia/kmedia_files_missing.csv"
	log.Infof("Processing data file: %s", input)
	f, err := os.Open(input)
	if err != nil {
		return errors.Wrap(err, "Open data file")
	}
	defer f.Close()

	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		r := strings.Split(line, ",")
		kmid, err := strconv.Atoi(r[0])
		if err != nil {
			return errors.Wrapf(err, "Bad KMID %s", r[0])
		}

		_, ok := fileMap[int64(kmid)]
		if ok {
			log.Infof("SHA1 should be null kmid %d", kmid)
			_, err = queries.Raw("UPDATE files SET sha1=NULL WHERE (properties->>'kmedia_id')::int = $1", kmid).Exec(mdb)
			if err != nil {
				return errors.Wrapf(err, "Update sha1 %d", kmid)
			}
		}
	}

	return nil
}
