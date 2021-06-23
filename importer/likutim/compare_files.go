package likutim

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

type Compare struct {
	storageDir     string
	allDocs        []string
	wordsByFileUid map[string][]map[string]int64 // file name -> list of paragraphs -> word -> counter
	result         []*Double
	errors         [][]interface{}
}

type Double struct {
	Save    string   `json:"save"`
	Doubles []string `json:"doubles"`
}

func (c *Compare) Run() []*Double {
	mdb := c.openDB()
	defer mdb.Close()

	err := os.MkdirAll(viper.GetString("likutim.os-dir"), os.ModePerm)
	if err != nil {
		log.Errorf("Can't create directory: %s", err)
	}
	dir, err := ioutil.TempDir(viper.GetString("likutim.os-dir"), "temp_kitvei_makor")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	c.storageDir = dir
	c.allDocs = []string{}

	c.wordsByFileUid = map[string][]map[string]int64{}

	//uids := []string{"MaJ1IBmg", "M4PYOzU5", "BD50WrEh", "EXoM7HY2", "gL4A2ajU", "YSjPO87t"}

	//שיעור בנושא "חוק הערבות"
	//uids1 := []string{"NGkBHYOc", "a0pql7A1", "QO0OKJVP", "0RWO6E45", "AYt7j3G6", "78yw7gN9", "qJvUzL2Y", "I94ixRdw", "CZKgZCho", "sVU30k0X", "9IemJTQK", "CsV4o1gV", "Z7l8O0cz", "gdfXW1Zv", "wLGlkTrI", "TcgWykjw", "ZDsPPWYf", "FfYqEnCe", "5qnzoah3", "N49WznDQ", "uRyVagJW", "AL3c18uY", "6NzSW1JD", "z3LrZmDf", "jl3FYSVM"}

	//שיעור בנושא "על כל פשעים תכסה אהבה"
	//uids2 := []string{"Kka7a73l", "le2bsaJc", "v18qCMuq", "no3cDNHM", "6u8IUtDc", "l4QMeRcV", "aUwVh0Z5", "y4o44oNV", "duBYBqbO"}
	//שיעור בנושא "על כל פשעים תכסה אהבה"  (הכנה לכנס וירטואלי " העולם החדש" 2020)
	//uids3 := []string{"ZERLFXAg", "vyB3Wx4o", "7cEvTDvN", "FLVxXiql"}
	//   שיעור בנושא "מתפילת העשירייה לתפילה של העשירייה עבור העולם"  (הכנה לכנס וירטואלי " העולם החדש" 2020)
	//uids4 := []string{"MBEAFhIh", "EUIRZA8s"}
	// problems files
	//uids := []string{"91e7mma7", "ruvfZyA4"}

	//uids := append(uids1, uids2...)
	//uids = append(uids, uids3...)
	//uids = append(uids, uids4...)

	cus, err := models.ContentUnits(mdb,
		qm.Select("distinct on (\"content_units\".id) \"content_units\".*"),
		qm.InnerJoin("files f ON f.content_unit_id = \"content_units\".id"),
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_KITEI_MAKOR].ID),
		qm.Load("Files", "Tags", "Tags.TagI18ns"),
		//qm.Offset(20),
		//qm.WhereIn("\"content_units\".uid in ?", utils.ConvertArgsString(uids)...),
		//qm.Limit(10),
	).All()
	if err != nil {
		log.Fatalf("Can't load units: %s", err)
		return c.result
	}

	c.countWordsPerFIle(cus)
	//order by file size - will be ease on recursive function
	sort.Slice(c.allDocs, func(i, j int) bool {
		return len(c.wordsByFileUid[c.allDocs[i]]) >= len(c.wordsByFileUid[c.allDocs[j]])
	})
	c.clearDuplicates(c.allDocs)
	log.Debugf("Result - uniq files: %v", c.result)
	c.printToCSV()
	err = c.saveJSON()
	if err != nil {
		log.Error(err)
	}
	return c.result
}

func (c *Compare) countWordsPerFIle(units []*models.ContentUnit) {
	for _, u := range units {
		for _, f := range u.R.Files {
			if f.Language.String != "he" || !strings.Contains(f.Name, ".doc") {
				continue
			}
			log.Debugf("Start count words. unit id: %d, file id: %d", u.ID, f.ID)
			url := viper.GetString("source-import.word-counter-url") + "?f=" + f.UID
			resp, err := http.Get(url)
			if err != nil {
				log.Errorf("Can't count words per paragraph in file unit id: %d, file id: %d. Error: %s", u.ID, f.ID, err)
				e := make([]interface{}, 2)
				e[0] = f
				e[1] = err
				continue
			}

			ps := []map[string]int64{}
			if err = json.NewDecoder(resp.Body).Decode(&ps); err != nil {
				log.Errorf("Error on decode response . Error: %s", err)
				e := make([]interface{}, 2)
				e[0] = f
				e[1] = err
				c.errors = append(c.errors, e)
				continue
			}
			log.Debug("Successfully count words.")
			c.wordsByFileUid[f.UID] = ps
			c.allDocs = append(c.allDocs, f.UID)
		}
	}
}

func (c *Compare) clearDuplicates(docs []string) {
	log.Debugf("Start recursive function docs length: %d ", len(docs))
	if len(docs) == 0 {
		return
	}

	d := Double{
		Save:    docs[0],
		Doubles: []string{docs[0]},
	}
	var nextDocs []string

	for i, n := range docs {
		if i == 0 {
			continue
		}
		eq, b := c.compareFiles(docs[0], n)
		if eq {
			d.Save = b
			d.Doubles = append(d.Doubles, n)
		} else {
			nextDocs = append(nextDocs, n)
		}
	}
	c.result = append(c.result, &d)
	log.Debugf("Call recursive function docs length: %d", len(nextDocs))
	c.clearDuplicates(nextDocs)
}

func (c *Compare) compareFiles(n1, n2 string) (bool, string) {
	f1 := c.wordsByFileUid[n1]
	f2 := c.wordsByFileUid[n2]
	biggerName := n1

	if len(f2) > len(f1) {
		f1 = c.wordsByFileUid[n2]
		f2 = c.wordsByFileUid[n1]
		biggerName = n2
	}

	eq := 0
	for _, p1 := range f1 {
		hasP1 := false
		for _, p2 := range f2 {
			if c.compareParagraph(p1, p2) {
				hasP1 = true
				break
			}
		}

		if hasP1 {
			eq = eq + 1
		}
	}

	if eq == len(f2) {
		return true, biggerName
	}
	percent := float64(eq) / float64(len(f2))

	return percent > 0.9, biggerName
}

func (c *Compare) compareParagraph(p1, p2 map[string]int64) bool {
	keys := map[string]float64{}
	wordsCount := float64(0)
	for k, c := range p1 {
		cf := float64(c)
		keys[k] = cf
		wordsCount += cf
	}

	for k, c := range p2 {
		cf := float64(c)
		wordsCount += cf
		v, ok := keys[k]
		if !ok {
			keys[k] = cf
			continue
		}
		diff := cf - v
		if diff > 0 {
			keys[k] = diff
		} else {
			keys[k] = -1 * diff
		}
	}

	diffCount := float64(0)
	for _, v := range keys {
		diffCount += v
	}
	if wordsCount == 0 {
		return false
	}

	return diffCount/wordsCount < 0.1
}

func (c *Compare) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

func (c *Compare) printToCSV() {
	addToFile("File that stay, All duplicates")
	for _, r := range c.result {
		addToFile(fmt.Sprintf(",\n%s, %v", r.Save, r.Doubles))
	}
	addToFile(",\nFiles  with exceptions")
	for _, e := range c.errors {
		addToFile(fmt.Sprintf(",\n%v, %v", e[0], e[1]))
	}
}

func addToFile(line string) {
	p := path.Join(viper.GetString("likutim.os-dir"), "kitvei-makor-duplicates.csv")
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Must(err)

	if _, err := f.Write([]byte(line)); err != nil {
		f.Close()
		log.Error(err)
	}

	utils.Must(f.Close())
}

func (c *Compare) saveJSON() error {
	j, err := json.Marshal(c.result)
	if err != nil {
		return errors.Wrapf(err, "Error on create json. Result: %v.", c.result)
	}
	p := path.Join(viper.GetString("likutim.os-dir"), "kitvei-makor-duplicates.json")
	err = ioutil.WriteFile(p, j, 0644)
	if err != nil {
		return errors.Wrapf(err, "Error on save json: %s, path: %s.", j, p)
	}
	return nil
}
