package hebcal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	H_ROSH_HASHANA     = "rosh-hashana"
	H_TZOM_GEDALIAH    = "tzom-gedaliah"
	H_YOM_KIPPUR       = "yom-kippur"
	H_SUKKOT           = "sukkot"
	H_SHMINI_ATZERET   = "shmini-atzeret"
	H_YOM_HAALIYAH     = "yom-haaliyah"
	H_SIGD             = "sigd"
	H_CHANUKAH         = "chanukah"
	H_ASARA_BTEVET     = "asara-btevet"
	H_TU_BISHVAT       = "tu-bishvat"
	H_PURIM_KATAN      = "purim-katan"
	H_TAANIT_ESTHER    = "taanit-esther"
	H_PURIM            = "purim"
	H_SHUSHAN_PURIM    = "shushan-purim"
	H_TAANIT_BECHOROT  = "taanit-bechorot"
	H_PESACH           = "pesach"
	H_YOM_HASHOAH      = "yom-hashoah"
	H_YOM_HAZIKARON    = "yom-hazikaron"
	H_YOM_HAATZMAUT    = "yom-haatzmaut"
	H_PESACH_SHENI     = "pesach-sheni"
	H_LAG_BAOMER       = "lag-baomer"
	H_YOM_YERUSHALAYIM = "yom-yerushalayim"
	H_SHAVUOT          = "shavuot"
	H_TZOM_TAMMUZ      = "tzom-tammuz"
	H_TISHA_BAV        = "tisha-bav"
	H_TU_BAV           = "tu-bav"
	H_LEIL_SELICHOT    = "leil-selichot"
)

type HebcalItems struct {
	Items []HebcalItem `json:"items"`
}

type HebcalItem struct {
	Link          string `json:"link"`
	Memo          string `json:"memo"`
	Category      string `json:"category"`
	Subcategory   string `json:"subcat"`
	Title         string `json:"title"`
	OriginalTitle string `json:"title_orig"`
	Hebrew        string `json:"hebrew"`
	Date          string `json:"date"`
}

type Hebcal struct {
	all map[string]map[int][]HebcalItem
}

func (h *Hebcal) Load() error {
	allItems := make([]HebcalItems, 0)
	re := regexp.MustCompile("^hebcal_\\d{4}\\.json$")
	err := filepath.Walk("hebcal/data", func(path string, info os.FileInfo, err error) error {
		if re.MatchString(info.Name()) {
			f, err := os.Open(path)
			if err != nil {
				return errors.Wrapf(err, "os.Open %s", info.Name())
			}

			var items HebcalItems
			if err := json.NewDecoder(f).Decode(&items); err != nil {
				return errors.Wrapf(err, "json.Decode %s", info.Name())
			}
			allItems = append(allItems, items)
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "Load hebcal data")
	}

	// group by link and year
	h.all = make(map[string]map[int][]HebcalItem, 50)
	for i := range allItems {
		x := allItems[i]
		for j := range x.Items {
			y := x.Items[j]

			pos := strings.LastIndex(y.Link, "/")
			k := y.Link[pos+1:]
			v, ok := h.all[k]
			if !ok {
				v = make(map[int][]HebcalItem)
				h.all[k] = v
			}

			year, err := strconv.Atoi(y.Date[:4])
			if err != nil {
				return errors.Wrapf(err, "Bad date [%d,%d] %s: %s", i, j, y.Date[:4], err.Error())
			}

			vv, ok := v[year]
			if !ok {
				vv = make([]HebcalItem, 0)
			}
			v[year] = append(vv, y)
		}
	}

	// sort by date
	for _, v := range h.all {
		for _, vv := range v {
			sort.Slice(vv, func(i, j int) bool {
				return vv[i].Date < vv[j].Date
			})
		}
	}

	// Special treatment for chanukah which might fall on a year change like in 2017 and 2005.
	// We look for consecutive days of the holiday spanning a georgian year change
	// and move them under the appropriate year (first day of holiday rules.)
	chanukah := h.all[H_CHANUKAH]
	for year := range chanukah {
		w, ok := chanukah[year+1]
		if !ok {
			continue
		}

		v := chanukah[year]
		s, err := time.Parse("2006-01-02", v[len(v)-1].Date)
		if err != nil {
			return errors.Wrapf(err, "time.Parse Date %s", v[len(v)-1].Date)
		}

		for {
			if len(w) == 0 {
				break
			}

			e, err := time.Parse("2006-01-02", w[0].Date)
			if err != nil {
				return errors.Wrapf(err, "time.Parse Date %s", w[0].Date)
			}

			if s.Sub(e) == -time.Hour*24 {
				s = e
				chanukah[year] = append(chanukah[year], w[0])
				w = w[1:]
				chanukah[year+1] = w
			} else {
				break
			}
		}
	}

	return nil
}

func (h *Hebcal) Print() {
	for k, v := range h.all {
		fmt.Printf("%s\n", k)
		for year, vv := range v {
			fmt.Printf("\t%d\n", year)
			for i := range vv {
				item := vv[i]
				fmt.Printf("\t\t%s\t%s\n", item.OriginalTitle, item.Date)
			}
		}
	}

	for k := range h.all {
		fmt.Println(k)
	}
}

func (h *Hebcal) GetPeriod(holiday string, year int) (start string, end string) {
	if x, ok := h.all[holiday]; ok {
		if y, ok := x[year]; ok {
			start = y[0].Date
			end = y[len(y)-1].Date
		}
	}
	return
}