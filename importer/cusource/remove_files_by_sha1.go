package cusource

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
)

func RemoveFilesBySHA1() {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	defer mdb.Close()
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	run(mdb)
}

func run(exec boil.Executor) {
	names := getNames()
	// actual removal
	err := models.Files(exec,
		qm.WhereIn("name in ?", utils.ConvertArgsString(names)...),
	).UpdateAll(models.M{"removed_at": null.TimeFrom(time.Now().UTC())})

	utils.Must(err)
	/*
			// file removed events
			removedFiles, err := models.Files(exec,
				qm.Select("id", "uid", "published"),
				qm.WhereIn("id in ?", utils.ConvertArgsInt64(fIDs)...)).
				All()
			if err != nil {
				return nil, errors.Wrap(err, "Refresh files to remove")
			}

			for i := range removedFiles {
				evnts = append(evnts, events.FileRemoveEvent(removedFiles[i]))
				wasPublished = wasPublished || removedFiles[i].Published
			}


		utils.Must(err)

	*/
}

func getNames() []string {
	path := viper.GetString("source-import.source-dir")
	r, err := os.Open(path)
	utils.Must(err)
	gzr, err := gzip.NewReader(r)
	utils.Must(err)
	defer utils.Must(gzr.Close())

	tr := tar.NewReader(gzr)
	result := make([]string, 1)

	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.Must(err)
		}

		if isDoc := strings.Contains(header.Name, ".doc"); header.Typeflag == tar.TypeReg && isDoc {
			spl := strings.Split(header.Name, "/")
			result = append(result, spl[len(spl)-1])
		}
	}
	return result
}
