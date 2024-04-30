package str

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var q = `
SELECT cu.id, cu.uid, cu.type_id, json_object_agg(f.language, f.*) as byLang 
	FROM content_units cu
	INNER JOIN files f ON f.content_unit_id = cu.id
	WHERE cu.id > $1 
	AND f.name LIKE '%%.mp3'
	AND f.type = 'audio'
	AND f.language = $2
	AND cu.published = TRUE
	AND cu.type_id IN (%s)
	AND cu.secure = 0
	AND f.published = TRUE
	AND f.secure = 0
GROUP BY cu.id, cu.uid
ORDER BY cu.id
LIMIT 1;
`

type BuildSrt struct {
	mdb      *sql.DB
	cuDir    string
	lastId   int64
	language string
	cts      string
}

func (s *BuildSrt) Run(cts, language string) {
	s.cts = cts
	s.language = language

	s.mdb = s.openDB()
	utils.Must(common.InitTypeRegistries(s.mdb))
	var err error
	s.lastId, err = fetchLastId()
	if err != nil {
		panic(err)
	}
	s.runLoop()
}
func (s *BuildSrt) runLoop() {
	if err := s.loop(); err != nil {
		utils.Must(s.saveBroken(err))
		s.lastId += 1
		s.runLoop()
	}
}

func (s *BuildSrt) loop() error {
	cu, byLang, err := s.getCuByPrevId()
	if err != nil {
		return err
	}

	if err := s.saveCuInfo(cu); err != nil {
		return err
	}
	for lang, f := range byLang {
		err = s.buildSrt(f, lang)
		if err != nil {
			return err
		}
	}

	if saveLastId(cu.ID) != nil {
		return err
	}
	s.lastId = cu.ID
	return s.loop()
}

func (s *BuildSrt) getCuByPrevId() (*models.ContentUnit, map[string]*models.File, error) {
	boil.DebugMode = true
	rows, err := queries.Raw(fmt.Sprintf(q, s.cts), s.lastId, s.language).Query(s.mdb)
	utils.Must(err)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var cu *models.ContentUnit
	var byLang map[string]*models.File
	for rows.Next() {
		var cuId int64
		var cuUid string
		var cuCt int64
		var _byLang []uint8
		if err = rows.Scan(&cuId, &cuUid, &cuCt, &_byLang); err != nil {
			return nil, nil, err
		}
		cu = &models.ContentUnit{
			ID:     cuId,
			UID:    cuUid,
			TypeID: cuCt,
		}
		if err := json.Unmarshal(_byLang, &byLang); err != nil {
			return nil, nil, err
		}
	}
	return cu, byLang, nil
}

func (s *BuildSrt) saveCuInfo(cu *models.ContentUnit) error {
	s.cuDir = fmt.Sprintf("%s/%s", viper.GetString("srt.res_dir"), cu.UID)
	if _, err := os.Stat(s.cuDir); !os.IsNotExist(err) {
		s.lastId = cu.ID
		return err
	}
	if err := os.Mkdir(s.cuDir, 0755); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", s.cuDir, "meta.csv")
	data := [][]string{
		{"unit uid", "unit content type"},
		{cu.UID, common.CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name},
	}

	return saveCsv(path, data)
}

func (s *BuildSrt) buildSrt(file *models.File, lang string) error {
	fUrl := fmt.Sprintf("https://cdn.kabbalahmedia.info/%s.mp3", file.UID)
	//_url := fmt.Sprintf("%s/stt?lang=%s&url=%s", viper.GetString("srt.ai_url"), lang, fUrl)
	resp, err := http.PostForm(fmt.Sprintf("%s/stt", viper.GetString("srt.ai_url")),
		url.Values{"lang": {lang}, "url": {fUrl}})
	defer func(resp *http.Response) {
		if resp == nil || resp.Body == nil {
			utils.Must(s.saveBroken(errors.Errorf("response have no body: %v", resp)))
		} else if resp.Body.Close() != nil {
			utils.Must(s.saveBroken(errors.Errorf("response have no body: %v", resp)))
		}
	}(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	if err = s.saveSrt(resp.Body, strings.Split(file.Name, ".")[0]); err != nil {
		return err
	}

	return err
}
func (s *BuildSrt) saveBroken(errLog error) error {
	path := fmt.Sprintf("%s/broken.scv", viper.GetString("srt.res_dir"))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	data := [][]string{
		{strconv.FormatInt(s.lastId, 10), fmt.Sprintf("%v", errLog)},
	}
	return writer.WriteAll(data)
}

func (s *BuildSrt) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

// helpers
func (s *BuildSrt) saveSrt(body io.ReadCloser, name string) error {
	out, err := os.Create(fmt.Sprintf("%s/%s.srt", s.cuDir, name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, body)
	return err
}
func fetchLastId() (int64, error) {
	path := fmt.Sprintf("%s/%s", viper.GetString("srt.res_dir"), "last_id.csv")
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return -1, errors.Wrap(err, "error reading CSV")
	}
	id, err := strconv.Atoi(records[1][0])
	if err != nil {
		return -1, err
	}
	return int64(id), nil
}

func saveLastId(id int64) error {
	path := fmt.Sprintf("%s/%s", viper.GetString("srt.res_dir"), "last_id.csv")
	data := [][]string{{"last_id"}, {fmt.Sprintf("%d", id)}}

	return saveCsv(path, data)
}

func saveCsv(path string, data [][]string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(data)
}
