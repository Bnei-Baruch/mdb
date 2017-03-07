package migrations

import (
	_ "github.com/Bnei-Baruch/mdb/gmodels"
	"github.com/Bnei-Baruch/mdb/kmodels"
	. "github.com/vattle/sqlboiler/queries/qm"
	"database/sql"
	"os"
	. "github.com/y0ssar1an/q"
	"github.com/Bnei-Baruch/mdb/gmodels"
	"github.com/vattle/sqlboiler/boil"
	"strconv"
	"fmt"
)

var (
	MDB      *sql.DB
	KMediaDB *sql.DB
)

func openDB() {
	var err error

	MDB, err = sql.Open("postgres", "host=localhost user=postgres dbname=mdb sslmode=disable")
	if err != nil {
		Q("MDB open:", err)
		os.Exit(1)
	}
	KMediaDB, _ = sql.Open("postgres", "host=localhost user=postgres dbname=kmedia sslmode=disable")
	if err != nil {
		Q("KMediaDB open:", err)
		os.Exit(1)
	}
	boil.DebugMode = true
}

func closeDB() {
	MDB.Close()
	KMediaDB.Close()
}

var languageTranslator = map[string]string{
	"HEB": "he",
	"RUS": "ru",
	"ENG": "en",
	"SPA": "es",
	"TRK": "",
	"MKD": "",
	"LAV": "",
	"POR": "",
	"RON": "",
	"FAR": "",
	"GEO": "",
	"LIT": "",
}

func ImportOldKmedia() {
	var err error

	openDB()
	defer closeDB()

	// Load Content Types
	lessonCT, err := kmodels.ContentTypes(KMediaDB, Where("name = ?", "Lesson")).One()
	if err != nil {
		Q("Unable to load content type 'Lesson'")
	}
	dailyLessonCT, err := gmodels.ContentTypes(MDB, Where("name = ?", "DAILY_LESSON")).One()
	if err != nil {
		Q("Unable to load content type 'DAILY_LESSON'")
	}
	saturdayLessonCT, err := gmodels.ContentTypes(MDB, Where("name = ?", "SATURDAY_LESSON")).One()
	if err != nil {
		Q("Unable to load content type 'SATURDAY_LESSON'")
	}
	Q(lessonCT, dailyLessonCT, saturdayLessonCT)

	//currentTranslationId, _ := gmodels.StringTranslations(MDB, Select("id"), OrderBy("id DESC"), Limit(1)).One()
	if err != nil {
		Q("Unable to load current Translation ID")
	}

	// virual lesson (virtual_lessons) -> collection
	// lesson part (containers) -> content unit
	// files (file_assets) -> files

	// VirtualLesson.limit(1).each do |vl|
	// 	name = "Morning lesson"
	// 	cl = Collection.create
	// 	cl.set_name('ENG',name)
	// 	vl.containers.each do |con|
	// 		cu = ContentUnit.create(name: con.name, description)
	// 		cu -> cl
	// 		con.file_assets.each do |fa|
	// 			f = File.create(fa)
	// 			f -> cu
	// 			op = OpAddFiles(f)
	// 		end
	// 	end
	// end

	Q("Virtual Lessons")
	vls, err := kmodels.VirtualLessons(KMediaDB, Limit(10), Load("Containers", "ContainerDescriptions")).All()
	if err != nil {
		Q("Unable to load virtual lessons")
		fmt.Errorf("%v\n", err)
		os.Exit(2)
	}
	for _, vl := range vls {

		Q("VL " + strconv.Itoa(vl.ID))
		containers, err := vl.Containers(KMediaDB, Limit(10), Load("FileAssets")).All()
		if err != nil {
			Q("Unable to load containers")
			os.Exit(1)
		}
		LoadOrCreateVL(vl, containers)

		for _, container := range containers {
			// LoadOrCreate(container)
			// Select new content type DAILY_LESSON or SATURDAY_LESSON
			// LoadOrCreate(container.descriptions)
			// Go over files
			Q(container)
		}
	}

	// TODO: Cleanup:
	// - remove collections without content units

	//var containers []gmodels.Container
	//	KMediaDB.Limit(10).Preload("FileAssets").Preload("Descriptions").Find(&containers)
	//	contentType := gmodels.ContentType{Name: "program"}
	//	MDB.FirstOrCreate(&contentType)
	//	for _, container := range containers {
	//		q.Q(container)
	//		var currentID uint64
	//		var descriptionExist bool
	//		MDB.DB().QueryRow("select id from string_translations order by id desc limit 1").Scan(&currentID)
	//		q.Q(currentID)
	//		for _, description := range container.Descriptions {
	//			if description.ContainerDesc != "" {
	//				MDB.Create(&models.StringTranslation{
	//					ID:               currentID + 1,
	//					Language:         strings.ToLower(description.Lang),
	//					CreatedAt:        time.Now(),
	//					OriginalLanguage: container.Lang,
	//					Text:             description.ContainerDesc,
	//				})
	//				descriptionExist = true
	//			}
	//		}
	//		if !descriptionExist {
	//			continue
	//		}
	//
	//		MDB.Create(&models.ContentUnit{
	//			ContentType: contentType,
	//			TranslatedContent: models.TranslatedContent{
	//				NameID:        currentID + 1,
	//				DescriptionID: currentID + 1,
	//			},
	//			//ID            int
	//			//UID: uuid.NewV4().String(),//           string
	//			//TypeID        string
	//			//TranslatedContent
	//			CreatedAt: container.CreatedAt, //     time.Time
	//			//Properties    JsonB
	//			//Files         []File
	//			//Collections   []Collection `gorm:"many2many:collections_content_units;"`
	//		})
	//		//Q(container)
	//		//newFile := File{}
	//		//Q(file)
	//		//KMediaDB.Model(&file).Association("Containers").Find(&file.Containers)
	//	}
	//	//pretty.Println(files)
}

func LoadOrCreateVL(vl *kmodels.VirtualLesson, containers []*kmodels.Container) *gmodels.Collection {
	// Search in collections for existing virtual lesson by film_date
	// In case there are many collections with the same film_date --> compare them to have the same containers
	// or at least the first one
	collections := gmodels.Collections(MDB, Where())

	// If it does not exist we have to create a new one; otherwise
}
