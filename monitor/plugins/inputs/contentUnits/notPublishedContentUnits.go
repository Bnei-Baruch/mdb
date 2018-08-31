package contentUnits

import (
	"database/sql"
	"encoding/json"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/monitor/interfaces"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/inputs"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type NotPublishedContentUnits struct {
	mdb_url string
}

func (u *NotPublishedContentUnits) TryParseConfigurations(inputConfigs map[string]interface{}) error {
	log.Printf("Parsing not published content units input configurations: %s", inputConfigs)
	// TODO: In other place mapstructure works buthere not working - needs to be consistent
	// mapstructure.Decode(inputConfigs, &u)
	// log.Printf("Using mdb_url: %s", u.mdb_url)
	u.mdb_url = inputConfigs["mdb_url"].(string)
	log.Printf("Using mdb_url: %s", u.mdb_url)
	return nil
}

func (s *NotPublishedContentUnits) Description() string {
	return "Read metrics about not published content units"
}

func (s *NotPublishedContentUnits) SampleConfig() string {
	return `
  ## Postgres connection url
  mdb_url = "connection-url"
`
}

func (s *NotPublishedContentUnits) Gather(acc interfaces.Accumulator) error {
	log.Println("Gathering not published content units statistics...")
	log.Infof("Setting up connection to MDB using %s", s.mdb_url)
	db, err := sql.Open("postgres", s.mdb_url)
	utils.Must(err)
	totalNotPublishedContentUnits, err := models.ContentUnits(db, qm.Where("published = false")).All()
	utils.Must(err)
	var meals []*models.ContentUnit
	var lessonParts []*models.ContentUnit
	for _, contentUnit := range totalNotPublishedContentUnits {
		if contentUnit.TypeID == 19 { // Meals
			meals = append(meals, contentUnit)
		} else if contentUnit.TypeID == 11 { // Lesson parts
			lessonParts = append(lessonParts, contentUnit)
		}
	}
	defer db.Close()
	mealsJSON, err := json.Marshal(meals)
	utils.Must(err)
	lessonPartsJSON, err := json.Marshal(lessonParts)
	utils.Must(err)
	fields := map[string]interface{}{
		"total":            len(totalNotPublishedContentUnits),
		"totalMeals":       len(meals),
		"totalLessonParts": len(lessonParts),
		"meals":            mealsJSON,
		"lessonParts":      lessonPartsJSON,
	}
	acc.AddFields("not_published_content_units", fields, nil)

	return nil
}

func init() {
	inputs.Add("notpublishedcontentunits", func() interfaces.Input {
		return &NotPublishedContentUnits{}
	})
}
