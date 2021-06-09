package cusource

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type KiteiMakorCompare struct {
	storageDir     string
	allDocs        []string
	wordsByFileUid map[string][]map[string]int64 // file name -> list of paragraphs -> word -> counter
	result         []double
	errors         [][]interface{}
}

type double struct {
	save string
	same []string
}

func (c *KiteiMakorCompare) Run() {
	mdb := c.openDB()
	defer mdb.Close()

	dir, err := ioutil.TempDir("/home/david", "temp_kitvei_makor")
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
	uids := []string{"91e7mma7", "ruvfZyA4"}

	//uids := append(uids1, uids2...)
	//uids = append(uids, uids3...)
	//uids = append(uids, uids4...)

	cus, err := models.ContentUnits(mdb,
		qm.Select("distinct on (\"content_units\".id) \"content_units\".*"),
		qm.InnerJoin("files f ON f.content_unit_id = \"content_units\".id"),
		qm.Where("type_id = ?", common.CONTENT_TYPE_REGISTRY.ByName[common.CT_KITEI_MAKOR].ID),
		qm.Load("Files", "Tags", "Tags.TagI18ns"),
		qm.WhereIn("\"content_units\".uid in ?", utils.ConvertArgsString(uids)...),
		//qm.Limit(100),
	).All()
	if err != nil {
		log.Errorf("can't load units: %s", err)
	}

	c.mapFiles(cus)
	c.clearDuplicates(c.allDocs)
	log.Debugf("Result - uniq files: %v", c.result)
	c.printToCSV()
}

func (c *KiteiMakorCompare) mapFiles(units []*models.ContentUnit) {
	for _, u := range units {
		for _, f := range u.R.Files {
			if f.Language.String != "he" || !strings.Contains(f.Name, ".doc") {
				continue
			}
			log.Debugf("Start count words. unit id: %d, file id: %d", u.ID, f.ID)
			url := viper.GetString("source-import.word-counter-url") + "?f=" + f.UID
			resp, err := http.Get(url)
			if err != nil {
				log.Errorf("Can't count words per paragraph in file unit id: %d, file id: %d", u.ID, f.ID, err)
				e := make([]interface{}, 2)
				e[0] = f
				e[1] = err
				continue
			}

			ps := []map[string]int64{}
			if err = json.NewDecoder(resp.Body).Decode(&ps); err != nil {
				log.Errorf("Error on decode response", err)
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

func (c *KiteiMakorCompare) clearDuplicates(docs []string) {
	log.Debugf("Start recursive function docs length: ", len(docs))
	if len(docs) == 0 {
		return
	}

	d := double{
		save: docs[0],
		same: []string{},
	}
	var nextDocs []string

	for i, n := range docs {
		if i == 0 {
			continue
		}
		eq, b := c.compareFiles(docs[0], n)
		if eq {
			d.save = b
			d.same = append(d.same, n)
		} else {
			nextDocs = append(nextDocs, n)
		}
	}
	c.result = append(c.result, d)
	log.Debugf("Call recursive function docs length: ", len(nextDocs))
	c.clearDuplicates(nextDocs)
}

func (c *KiteiMakorCompare) compareFiles(n1, n2 string) (bool, string) {
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

func (c *KiteiMakorCompare) compareParagraph(p1, p2 map[string]int64) bool {
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

func (c *KiteiMakorCompare) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("source-import.test-url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

func (c *KiteiMakorCompare) printToCSV() {
	lines := []string{"File that stay", "All duplicates"}
	for _, r := range c.result {
		l := fmt.Sprintf("\n%s, %v", r.save, r.same)
		lines = append(lines, l)
	}
	lines = append(lines, "\nFiles  with exceptions")
	for _, e := range c.errors {
		l := fmt.Sprintf("\n%v, %v", e[0], e[1])
		lines = append(lines, l)
	}
	b := []byte(strings.Join(lines, ","))
	err := ioutil.WriteFile(viper.GetString("source-import.kitvei-makor-duplicates"), b, 0644)
	utils.Must(err)
}
