package storage

import (
	"bufio"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	//"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/lib/pq"
	"strings"
)

var INDEX_LINE_RE = regexp.MustCompile("\\[\"(.*)\",\"(.*)\",(\\d+),(\\d+)]")

type PhysicalFile struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
}

// A single storage location index
// 	key - SHA1 checksum
// 	value - list of physical copies of that file in this storage location
type LocationIndex map[string][]*PhysicalFile

// All physical copies of a single file in all storage locations
// 	key - storage location name
//	value - list of physical copies of that file in this storage location
type FileCopies map[string][]*PhysicalFile

// A combined index of all storage locations
//	key - SHA1 checksum
//	value - Map of storage name to physical copies
type MasterIndex map[string]FileCopies

func CreateMasterIndex() {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Starting index merge")

	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()
	boil.SetDB(mdb)
	//boil.DebugMode = true

	mIndex, err := loadMasterIndex(mdb)
	utils.Must(err)
	log.Infof("Master index size is %d", len(mIndex))

	log.Info("Merging local_archive index")
	//indexPath, _ := filepath.Abs("storage/.index/files-kabbalahmedia")
	//err = MergeIndex(mdb, mIndex, "kabbalahmedia", indexPath)
	indexPath, _ := filepath.Abs("storage/.index/files-archive")
	err = MergeIndex(mdb, mIndex, "local_archive", indexPath)
	utils.Must(err)

	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func loadMasterIndex(db *sql.DB) (MasterIndex, error) {
	total, err := models.Files(db, qm.Where("sha1 IS NOT NULL")).Count()
	if err != nil {
		return nil, errors.Wrap(err, "Get total files in MDB")
	}
	log.Infof("%d files with sha1 in MDB", total)

	master := make(MasterIndex, total)

	rows, err := queries.Raw(db,
		`SELECT id, encode(sha1, 'hex'), properties -> 'storage' FROM files WHERE sha1 IS NOT NULL`).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "Load files")
	}

	for rows.Next() {
		var fid int64
		var sha1 string
		var sProps null.JSON
		err = rows.Scan(&fid, &sha1, &sProps)
		if err != nil {
			return nil, errors.Wrap(err, "Scan row")
		}

		var s FileCopies
		if sProps.Valid {
			err = json.Unmarshal(sProps.JSON, &s)
			if err != nil {
				return nil, errors.Wrapf(err, "json.Unmarshal [%d]", fid)
			}
		}

		master[sha1] = s
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "Iter rows")
	}
	err = rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Close rows")
	}

	return master, nil
}

// Merges the index of the given storage location into the master index.
func MergeIndex(db *sql.DB, master MasterIndex, location string, indexPath string) error {
	log.Infof("Reading index at %s", indexPath)
	index, err := readIndex(indexPath)
	if err != nil {
		return errors.Wrap(err, "Read index")
	}

	log.Info("Computing diff with master index")
	diff := computeDiff(master, index, location)
	log.Infof("%d files changed in index", len(diff))

	log.Info("Saving diff results to DB")
	successful, err := saveDiff(db, master, diff, location)
	if err != nil {
		return errors.Wrap(err, "Save diff to DB")
	}
	log.Infof("Successfully updated %d out of %d files in diff", successful, len(diff))

	if successful == len(diff) {
		return nil
	} else {
		return errors.Errorf("%d files couldn't be saved to DB", len(diff)-successful)
	}
}

func readIndex(path string) (LocationIndex, error) {
	idxFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Open file")
	}
	defer idxFile.Close()

	index := make(LocationIndex, 500000)  // we estimate number of lines in index file to 0.5M
	i := 0
	badLines := 0
	scanner := bufio.NewScanner(idxFile)
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		sha1, pf, err := parseLineSplit(line)
		if err != nil {
			log.Warnf("Bad line [%d]: %s", i, err.Error())
			badLines++
			continue
		}

		fLocations, ok := index[sha1]
		if !ok {
			fLocations = make([]*PhysicalFile, 0)
		}
		index[sha1] = append(fLocations, pf)
	}
	log.Infof("Index has %d files in %d total lines, %d bad lines", len(index), i, badLines)

	return index, nil
}

// Parse a single line in an index file.
// Implementation is tight hard to current format for performance.
// Using regexps was an order of magnitude slower.
//
// The regex for a line is: ["(.*)<path>","(.*)<sha1>",(\d+),(\d+)]
func parseLineSplit(line string) (string, *PhysicalFile, error) {
	s := strings.Split(line[1:len(line)-1],"\"")
	if len(s) != 5 {
		return "", nil, errors.Errorf("bad format [len(s)=%d]: %s", len(s),line)
	}

	sha1 := s[3]
	if !utils.SHA1_RE.MatchString(sha1) {
		return "", nil, errors.Errorf("bad checksum %s", sha1)
	}

	pf := &PhysicalFile{
		Path: s[1],
	}

	s = strings.Split(s[4],",")
	if len(s) != 3 {
		return "", nil, errors.Errorf("bad format numbers [len(s)=%d]: %s", len(s),line)
	}
	var err error

	pf.Size, err = strconv.ParseInt(s[1], 10, 0)
	if err != nil {
		return "", nil, errors.Errorf("size not an integer, %s", s[2])
	}

	pf.ModTime, err = strconv.ParseInt(s[2], 10, 0)
	if err != nil {
		return "", nil, errors.Errorf("last_modified not an integer, %s", s[3])
	}

	return sha1, pf, nil
}

func computeDiff(master MasterIndex, index LocationIndex, location string) LocationIndex {
	diff := make(LocationIndex)

	for sha1, allCopies := range master {
		var prev []*PhysicalFile
		if allCopies != nil {
			prev = allCopies[location]
		}
		next := index[sha1]
		modified, _ := fileDiff(prev, next)
		//log.Infof("SHA1 %s: modified=%s, changes=%d", sha1, modified, changes)

		if modified {
			diff[sha1] = next
		}
	}

	return diff
}

// Return the high level diff statistics for the copies of a file in a single storage location.
// 	Returns:
//	modified (bool) - any changes at all ?
//	changes (int) - diff in total number of copies. Positive, means added, negative, removed.
func fileDiff(prev, next []*PhysicalFile) (modified bool, changes int) {
	changes = len(next) - len(prev)
	modified = changes != 0

	for i := 0; !modified && i < len(next); i++ {
		x := next[i]
		exists := false
		for j := 0; !exists && j < len(prev); j++ {
			exists = x.Path == prev[j].Path
		}
		modified = !exists
	}

	return
}

// Saving the diff in the DB may involve many updates.
// To increase performance we use the postgres's FROM clause in UPDATE statement.
// We first create a temp table and insert updated 'storage' property of each file.
// We then execute a single UPDATE statement.
// Finally we drop the temporary table.
// See http://tapoueh.org/blog/2013/03/15-batch-update
func saveDiff(db *sql.DB, master MasterIndex, diff LocationIndex, location string) (int, error) {
	if len(diff) == 0 {
		return 0, nil
	}

	tmpTable := strings.ToLower(fmt.Sprintf("storage_diff_%s", utils.GenerateName(4)))
	log.Infof("Creating temp table: %s", tmpTable)
	q := fmt.Sprintf("CREATE TABLE %s (sha1 BYTEA, storage JSONB)", tmpTable)
	log.Info(q)
	_, err := db.Exec(q)
	if err != nil {
		return 0, errors.Wrap(err, "Create temp table")
	}
	defer func() {
		log.Infof("Dropping temp table: %s", tmpTable)
		_, err = db.Exec(fmt.Sprintf("DROP TABLE %s", tmpTable))
		if err != nil {
			log.Errorf("Error dropping temp table [%s]: %s", tmpTable, err.Error())
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Begin transaction")
	}

	stmt, err := tx.Prepare(pq.CopyIn(tmpTable, "sha1", "storage"))
	if err != nil {
		return 0, errors.Wrap(err, "Prepare statement")
	}

	log.Info("Populating temp table")
	for sha1, copies := range diff {
		sb, err := hex.DecodeString(sha1)
		if err != nil {
			log.Errorf("hex Decode sha1=%s: %s", sha1, err.Error())
			continue
		}

		s, ok := master[sha1]
		if !ok {
			s = make(FileCopies)
		}
		delete(s, location)
		if copies != nil {
			s[location] = copies
		}

		sj, err := json.Marshal(s)
		if err != nil {
			log.Errorf("json Marshal copies sha1=%s, copies=%s: %s", sha1, copies, err.Error())
			continue
		}

		_, err = stmt.Exec(sb, string(sj))
		if err != nil {
			log.Errorf("Error saving diff to DB, sha1=%s, copies=%s: %s", sha1, copies, err.Error())
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return 0, errors.Wrap(err, "Final exec statement")
	}
	err = stmt.Close()
	if err != nil {
		return 0, errors.Wrap(err, "Close statement")
	}
	err = tx.Commit()
	if err != nil {
		return 0, errors.Wrap(err, "Commit transaction")
	}

	log.Info("Executing update")
	q = fmt.Sprintf(`UPDATE files f SET properties = jsonb_set(properties, '{storage}', t.storage)
	FROM %s t WHERE f.sha1 = t.sha1`, tmpTable)
	r, err := db.Exec(q)
	if err != nil {
		return 0, errors.Wrap(err, "Save diff to files table")
	}
	ra, err := r.RowsAffected()
	if err == nil {
		return int(ra), nil
	} else {
		return 0, errors.Wrap(err, "Get number of rows affected")
	}
}
