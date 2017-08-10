package storage

import (
	"bufio"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	STATUS_ONLINE   = "ONLINE"
	STATUS_NEARLINE = "NEARLINE"
	STATUS_OFFLINE  = "OFFLINE"
)

type Location struct {
	ID      string
	Country string
}

var DEFAULT_LOCATIONS = map[string]*Location{
	"il-merkaz": {ID: "il-merkaz", Country: "Israel"},
}

type FileCopy struct {
	LocationID string `json:"location"`
	StorageID  string `json:"storage"`
	Status     string `json:"status"`
}

type Statistics struct {
}

func ImportStorageStatus() {
	defer func() {
		if rval := recover(); rval != nil {
			debug.PrintStack()
			err, ok := rval.(error)
			if !ok {
				err = errors.Errorf("panic: %s", rval)
			}

			// TODO proper handling of panic
			log.WithError(err).Error("Panic")
		}
	}()

	doStuff()
}

// The way we do stuff:
// We download the data file and locations map from the storage api.
// Once we have these, we create a temp table and fill it with the data in the file.
// We update `files` table from that temp table.
// We clear storage status for files not found in data file
// Finally, we drop temp table and clean old local copies of data files
func doStuff() {
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

	// TODO: proper error handling

	dataFile, err := downloadDataFile()
	if err != nil {
		panic(errors.Wrap(err, "Download data file"))
	}

	locations, err := syncLocations()
	if err != nil {
		panic(errors.Wrap(err, "Get locations definitions"))
	}

	fileMap, err := loadMDBFiles(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Load files from MDB"))
	}

	tempTable, err := initTempTable(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Initializing temp update table"))
	}

	count, err := processDataFile(dataFile, mdb, fileMap, tempTable, locations)
	if err != nil {
		panic(errors.Wrap(err, "Process data file"))
	}

	rowsAffected, err := saveUpdates(mdb, tempTable)
	if err != nil {
		panic(errors.Wrap(err, "Save updates to MDB"))
	}

	// count equals rowsAffected ?
	if count == rowsAffected {
		log.Infof("%d files were updated in MDB", count)
	} else {
		log.Warnf("Count mismatch, count != rowsAffected %d != %d", count, rowsAffected)
	}

	err = clearStatusForMissing(mdb, fileMap)
	if err != nil {
		panic(errors.Wrap(err, "Clear status for missing"))
	}

	err = shutdownTempTable(mdb, tempTable)
	if err != nil {
		panic(errors.Wrap(err, "Shutdown temp update table"))
	}

	err = cleanupDataDir(dataFile)
	if err != nil {
		panic(errors.Wrap(err, "Cleanup data dir"))
	}

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func downloadDataFile() (string, error) {
	dir, err := getDataDir()
	if err != nil {
		return "", errors.Wrap(err, "Get data directory")
	}

	fileName := fmt.Sprintf("index_%s.txt", time.Now().Format("20060102150405"))
	output := path.Join(dir, fileName)
	log.Infof("Downloading storage status index to %s", output)

	url := fmt.Sprintf("%sstorage", viper.GetString("storage.api-url"))

	err = utils.DownloadUrl(url, output)
	if err != nil {
		return "", errors.Wrap(err, "Download file")
	}

	return output, nil
}

func cleanupDataDir(dataFile string) error {
	dir, err := getDataDir()
	if err != nil {
		return errors.Wrap(err, "Get data directory")
	}

	// Remove all previous data files in directory
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("Error walking %s: %s", path, err.Error())
			return nil
		}
		if info.IsDir() || path == dataFile {
			return nil
		}

		ex := os.Remove(path)
		if ex != nil {
			return errors.Wrapf(ex, "Remove file %s", path)
		}

		return nil
	})

	return nil
}

func syncLocations() (map[string]*Location, error) {
	return DEFAULT_LOCATIONS, nil
}

func loadMDBFiles(db *sql.DB) (map[string]int64, error) {
	total, err := models.Files(db, qm.Where("sha1 IS NOT NULL")).Count()
	if err != nil {
		return nil, errors.Wrap(err, "Get total files in MDB")
	}
	log.Infof("%d files with sha1 in MDB", total)

	fileMap := make(map[string]int64, total)

	rows, err := queries.Raw(db,
		`SELECT id, encode(sha1, 'hex') FROM files WHERE sha1 IS NOT NULL`).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "Load files")
	}

	for rows.Next() {
		var fid int64
		var sha1 string
		err = rows.Scan(&fid, &sha1)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}

		fileMap[sha1] = fid
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

func initTempTable(db *sql.DB) (string, error) {
	tmpTable := strings.ToLower(fmt.Sprintf("storage_diff_%s", utils.GenerateName(4)))
	log.Infof("Creating temp table: %s", tmpTable)
	q := fmt.Sprintf("CREATE TABLE %s (sha1 BYTEA, storage JSONB)", tmpTable)
	_, err := db.Exec(q)
	if err != nil {
		return "", errors.Wrap(err, "Create temp table")
	}
	return tmpTable, nil
}

func shutdownTempTable(db *sql.DB, tableName string) error {
	log.Infof("Dropping temp table: %s", tableName)
	_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", tableName))
	if err != nil {
		return errors.Wrapf(err, "Dropping temp table: %s", tableName)
	}
	return nil
}

func processDataFile(path string, db *sql.DB, fileMap map[string]int64, tempTable string,
	locations map[string]*Location) (int64, error) {

	log.Infof("Processing data file: %s", path)
	f, err := os.Open(path)
	if err != nil {
		return 0, errors.Wrap(err, "Open data file")
	}
	defer f.Close()

	tx, err := db.Begin()
	utils.Must(err)

	stmt, err := tx.Prepare(pq.CopyIn(tempTable, "sha1", "storage"))
	if err != nil {
		utils.Must(tx.Rollback())
		return 0, errors.Wrap(err, "Prepare statement")
	}

	log.Info("Populating temp table")
	i := 0
	count := int64(0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		sha1 := line[:40]
		var sha1B []byte
		sha1B, err = hex.DecodeString(sha1)
		if err != nil {
			err = errors.Wrapf(err, "hex.DecodeString %s", sha1)
			break
		}

		var val []FileCopy
		err = json.Unmarshal([]byte(line[41:]), &val)
		if err != nil {
			err = errors.Wrapf(err, "json.Unmarshal line [%d]", i)
			break
		}

		if len(val) == 0 {
			log.Warnf("Empty storage status line [%d]", i)
			continue
		}

		// Verify location IDs & statuses
		for j := range val {
			if _, ok := locations[val[j].LocationID]; !ok {
				log.Warnf("Unknown location %s line [%d]", val[j].LocationID, i)
				continue
			}

			switch strings.ToUpper(val[j].Status) {
			case STATUS_ONLINE, STATUS_NEARLINE, STATUS_OFFLINE:
				continue
			default:
				log.Warnf("Unknown status %s line [%d]", val[j].Status, i)
				continue
			}
		}

		// we Unmarshal, Marshal to keep data consistent with code
		var valB []byte
		valB, err = json.Marshal(val)
		if err != nil {
			err = errors.Wrapf(err, "json.Marshal line [%d]", i)
			break
		}

		_, err = stmt.Exec(sha1B, string(valB))
		if err != nil {
			err = errors.Wrapf(err, "Insert to temp table line [%d]", i)
			break
		}

		count++

		// we keep only files which have no storage
		// in next phase we clear their storage value in MDB
		delete(fileMap, sha1)
	}
	if err != nil {
		utils.Must(stmt.Close())
		utils.Must(tx.Rollback())
		return 0, err
	}

	_, err = stmt.Exec()
	if err != nil {
		utils.Must(tx.Rollback())
		return 0, errors.Wrap(err, "Final exec statement")
	}
	err = stmt.Close()
	if err != nil {
		utils.Must(tx.Rollback())
		return 0, errors.Wrap(err, "Close statement")
	}

	utils.Must(tx.Commit())

	return count, nil
}

func saveUpdates(db *sql.DB, tableName string) (int64, error) {
	log.Info("Updating files with storage status")

	tx, err := db.Begin()
	utils.Must(err)

	q := fmt.Sprintf(`UPDATE files f SET properties = jsonb_set(properties, '{storage}', t.storage)
	FROM %s t WHERE f.sha1 = t.sha1`, tableName)
	r, err := tx.Exec(q)
	if err != nil {
		utils.Must(tx.Rollback())
		return 0, errors.Wrap(err, "Save diff to files table")
	}

	ra, err := r.RowsAffected()
	if err != nil {
		utils.Must(tx.Rollback())
		return 0, errors.Wrap(err, "Get number of rows affected")
	}

	utils.Must(tx.Commit())

	return ra, nil
}

func clearStatusForMissing(db *sql.DB, missingFiles map[string]int64) error {
	log.Infof("Clearing storage status for %d missing files", len(missingFiles))

	i := 0
	ids := make([]string, len(missingFiles))
	for _, id := range missingFiles {
		ids[i] = strconv.FormatInt(id, 10)
		i++
	}

	pageSize := 1000
	i = 0
	for pageSize*i < len(ids) {
		start := pageSize * i
		end := utils.Min(start+pageSize, len(ids)-1)

		tx, err := db.Begin()
		utils.Must(err)

		_, err = tx.Exec(fmt.Sprintf(`UPDATE files SET properties = properties - 'storage' WHERE id IN (%s)`,
			strings.Join(ids[start:end], ",")))
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Clear storage page %d", i)
		}

		utils.Must(tx.Commit())
		i++
	}

	return nil
}

func getDataDir() (string, error) {
	dir := viper.GetString("storage.index-directory")
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "os.Getwd")
		}
		dir = path.Join(cwd, "storage", ".index")
	}

	return dir, nil
}
