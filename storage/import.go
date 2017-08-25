package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"
)

const (
	STATUS_ONLINE   = "online"
	STATUS_NEARLINE = "nearline"
	STATUS_OFFLINE  = "offline"
)

var EMPTY_IDS = make([]int64, 0)

type StorageDevice struct {
	ID       string `json:"id"`
	Country  string `json:"country"`
	Location string `json:"location"`
	Status   string `json:"status"`
	Access   string `json:"access"`
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
// Once we have these, we compute the status diff for each file.
// We create new mappings and delete no longer existing ones.
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

	err = syncStorages(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Sync storages"))
	}

	dataFile, err := downloadDataFile()
	if err != nil {
		panic(errors.Wrap(err, "Download data file"))
	}

	fileMap, err := loadMDBFiles(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Load files from MDB"))
	}

	fileStorages, err := loadMDBFilesMappings(mdb)
	if err != nil {
		panic(errors.Wrap(err, "Load files storages from MDB"))
	}

	err = processDataFile(dataFile, mdb, fileMap, fileStorages)
	if err != nil {
		panic(errors.Wrap(err, "Process data file"))
	}

	err = clearStatusForMissing(mdb, fileMap, fileStorages)
	if err != nil {
		panic(errors.Wrap(err, "Clear status for missing"))
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

	url := fmt.Sprintf("%scatalog", viper.GetString("storage.api-url"))

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

// get storages from API
// get storages from MDB
// Process diff: create new, update existing and remove non existing
func syncStorages(db *sql.DB) error {
	// Get storages from API
	url := fmt.Sprintf("%sstorages", viper.GetString("storage.api-url"))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	utils.Must(err)

	req.Header.Set("User-Agent", fmt.Sprintf("MDB_%s", version.Version))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Do http request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Read http response body")
	}

	var apiStorages []StorageDevice
	err = json.Unmarshal(body, &apiStorages)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}
	log.Infof("Got %d storages from API", len(apiStorages))

	apiStoragesMap := make(map[string]*StorageDevice, len(apiStorages))
	for i := range apiStorages {
		apiStoragesMap[apiStorages[i].ID] = &apiStorages[i]
	}

	// Fetch storages from MDB
	mdbStoragesMap, err := getMDBStorageMap(db)
	if err != nil {
		return errors.Wrap(err, "Load storages from MDB")
	}
	log.Infof("Got %d storages from MDB", len(mdbStoragesMap))

	// Sync API and MDB storage devices

	// create or update MDB models
	tx, err := db.Begin()
	utils.Must(err)
	for i := range apiStorages {
		s := apiStorages[i]
		mdbS, ok := mdbStoragesMap[s.ID]
		if ok {
			// update
			mdbS.Country = s.Country
			mdbS.Location = s.Location
			mdbS.Status = s.Status
			mdbS.Access = s.Access
			err = mdbS.Update(tx, "country", "location", "status", "access")
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Update MDB storage %d %s", mdbS.ID, s.ID)
			}

			delete(mdbStoragesMap, s.ID)
		} else {
			// create new
			mdbS = &models.Storage{
				Name:     s.ID,
				Country:  s.Country,
				Location: s.Location,
				Status:   s.Status,
				Access:   s.Access,
			}
			err = mdbS.Insert(tx)
			if err != nil {
				utils.Must(tx.Rollback())
				return errors.Wrapf(err, "Create new MDB storage %s", s.ID)
			}
		}
	}

	utils.Must(tx.Commit())

	// remove mdb models not found in API anymore
	if len(mdbStoragesMap) > 0 {
		ids := make([]int64, len(mdbStoragesMap))
		i := 0
		for _, v := range mdbStoragesMap {
			ids[i] = v.ID
			i++
		}

		log.Infof("Deleting %d storages from MDB", len(ids))

		tx, err = db.Begin()
		utils.Must(err)

		// TODO: this is buggy, seems like a bug in sqlboiler...
		// need to dig in and fix asap
		err = models.Storages(tx, qm.WhereIn("id in ?", ids)).DeleteAll()
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrap(err, "Delete storages from MDB")
		}

		// Note: the ON DELETE CASCADE in files_storages join table

		utils.Must(tx.Commit())
	}

	return nil
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
		var fID int64
		var sha1 string
		err = rows.Scan(&fID, &sha1)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}

		fileMap[sha1] = fID
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

func loadMDBFilesMappings(db *sql.DB) (m map[int64][]int64, err error) {
	rows, err := queries.Raw(db, "SELECT * FROM files_storages").Query()
	if err != nil {
		err = errors.Wrap(err, "Load files mappings from MDB")
		return
	}

	m = make(map[int64][]int64)
	for rows.Next() {
		var fID int64
		var sID int64
		err = rows.Scan(&fID, &sID)
		if err != nil {
			err = errors.Wrap(err, "sql.Scan")
			return
		}

		if v, ok := m[fID]; ok {
			m[fID] = append(v, sID)
		} else {
			m[fID] = []int64{sID}
		}
	}

	err = rows.Err()
	if err != nil {
		err = errors.Wrap(err, "Iterate rows")
		return
	}
	err = rows.Close()
	if err != nil {
		err = errors.Wrap(err, "Close rows")
		return
	}

	return
}

type fileDiffUOW struct {
	fID     int64
	current []int64
	next    map[int64]bool
}

// Note: this function modifies the contents of fileMap
func processDataFile(path string, db *sql.DB, fileMap map[string]int64, fileStorages map[int64][]int64) error {
	log.Infof("Processing data file: %s", path)
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "Open data file")
	}
	defer f.Close()

	// Fetch storages from MDB
	sMap, err := getMDBStorageMap(db)
	if err != nil {
		return errors.Wrap(err, "Load storages from MDB")
	}
	log.Infof("Got %d storages from MDB", len(sMap))

	log.Info("Setting up workers")
	jobs := make(chan fileDiffUOW, 100)
	var workersWG sync.WaitGroup
	for w := 1; w <= 5; w++ {
		go makeWorker(db, workersWG)(jobs)
	}

	log.Info("Queueing work")
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
		fID, ok := fileMap[sha1]
		if !ok {
			continue
		}

		count++
		// we keep only files which have no storage
		// in next phase we clear their storage value in MDB
		delete(fileMap, sha1)

		// dedup names since the API doesn't do that for us at the moment

		names := strings.Split(strings.Replace(line[42:len(line)-1], "\"", "", -1), ",")
		locations := make(map[int64]bool)
		for j := range names {
			if s, ok := sMap[names[j]]; ok {
				locations[s.ID] = true
			} else if names[j] != "" {
				log.Warnf("Unknown storage device %s line [%d]", names[j], i)
			}
		}

		current, ok := fileStorages[fID]
		if !ok {
			current = EMPTY_IDS
		}

		jobs <- fileDiffUOW{
			fID:     fID,
			current: current,
			next:    locations,
		}
	}

	log.Info("Closing jobs channel")
	close(jobs)

	log.Info("Waiting for workers to finish")
	workersWG.Wait()

	log.Infof("Good files %d, total lines %d", count, i)

	return nil
}

func makeWorker(db *sql.DB, wg sync.WaitGroup) func(jobs <-chan fileDiffUOW) {
	wg.Add(1)
	return func(jobs <-chan fileDiffUOW) {
		for uow := range jobs {
			if err := doFile(db, uow.fID, uow.current, uow.next); err != nil {
				log.Errorf("Process file %d: %s", uow.fID, err.Error())
			}
		}
		wg.Done()
	}
}

func doFile(db *sql.DB, fID int64, current []int64, next map[int64]bool) error {

	// we diff current vs next storage devices status and mark:
	// 1. which existing needs to be deleted since they're not in new state
	// 2. which new mappings needs to be created
	// unchanged mappings are untouched

	toDelete := make([]int64, 0)
	for j := range current {
		x := current[j]
		if _, ok := next[x]; ok {
			delete(next, x) // we only keep new mappings to be created
		} else {
			toDelete = append(toDelete, x)
		}
	}

	if len(next) == 0 && len(toDelete) == 0 {
		return nil // no op
	}

	tx, err := db.Begin()
	utils.Must(err)

	// create new mappings
	if len(next) > 0 {
		values := make([]string, 0)
		for sID := range next {
			values = append(values, fmt.Sprintf("(%d,%d)", fID, sID))
		}
		res, err := queries.Raw(tx,
			fmt.Sprintf("INSERT INTO files_storages (file_id, storage_id) VALUES %s",
				strings.Join(values, ","))).Exec()
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrap(err, "Insert mappings")
		}
		ra, err := res.RowsAffected()
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrap(err, "Insert mappings, retrieve rows affected")
		}
		if ra != int64(len(values)) {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Insert mappings, rows affected, expected %d got %d", len(values), ra)
		}
	}

	// delete old mappings
	if len(toDelete) > 0 {
		values := make([]string, 0)
		for sID := range toDelete {
			values = append(values, strconv.Itoa(sID))
		}
		res, err := queries.Raw(tx,
			fmt.Sprintf("DELETE FROM files_storages WHERE file_id=%d AND storage_id IN (%s)",
				fID, strings.Join(values, ","))).
			Exec()
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Delete mappings")
		}
		ra, err := res.RowsAffected()
		if err != nil {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Delete mappings, retrieve rows affected")
		}
		if ra != int64(len(values)) {
			utils.Must(tx.Rollback())
			return errors.Wrapf(err, "Delete mappings, rows affected, expected %d got %d", len(values), ra)
		}
	}

	utils.Must(tx.Commit())

	return nil
}

func clearStatusForMissing(db *sql.DB, missingFiles map[string]int64, fileStorages map[int64][]int64) error {
	// create a slice of missing files IDs
	ids := make([]string, 0)
	for _, id := range missingFiles {
		if _, ok := fileStorages[id]; ok {
			ids = append(ids, strconv.FormatInt(id, 10))
		}
	}

	log.Infof("Clearing storage status for %d missing files", len(ids))
	if len(ids) == 0 {
		return nil
	}

	// go over slice by pages
	pageSize := 1000
	i := 0
	for pageSize*i < len(ids) {
		start := pageSize * i
		end := utils.Min(start+pageSize, len(ids)-1)

		tx, err := db.Begin()
		utils.Must(err)

		_, err = tx.Exec(fmt.Sprintf(`DELETE FROM files_storages WHERE file_id IN (%s)`,
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

func getMDBStorageMap(db *sql.DB) (m map[string]*models.Storage, err error) {
	all, err := models.Storages(db).All()
	if err != nil {
		return
	}

	m = make(map[string]*models.Storage, len(all))
	for i := range all {
		m[all[i].Name] = all[i]
	}

	return
}
