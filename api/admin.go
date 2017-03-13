package api

import (
	"encoding/json"
    "encoding/hex"
	"net/http"

	"github.com/Bnei-Baruch/mdb/models"
	log "github.com/Sirupsen/logrus"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/gin-gonic/gin.v1"
)

type MarshableFile models.File

func (f *MarshableFile) MarshalJSON() ([]byte, error) {
    type Copy MarshableFile
    var b []byte
    if f.Sha1.Valid {
        b = f.Sha1.Bytes
    }
    return json.Marshal(&struct {
        Sha1 string `json:"sha1"`
        *Copy
    }{
        Sha1: hex.EncodeToString(b),
        Copy: (*Copy)(f),
    })
}

// AdminFilesHandler should also support following:
// query - text string to search in file properties.
// limit - nax number of files to fetch.
// first - offset (for pagination)
// http://.../files?sort=X,query=Y,first=Z,limit=W

// Gets list of all files.
func AdminFilesHandler(c *gin.Context) {

    sort := c.DefaultQuery("sort", "date")
    query := c.DefaultQuery("query", "*")
    first := c.DefaultQuery("first", "0")
    limit := c.DefaultQuery("limit", "30")


	tx, err := boil.Begin()
    var filesSlice []*models.File
	if err == nil {
        filesSlice, err = getFiles(tx)
		if err == nil {
			tx.Commit()
		} else {
			log.Error("Error handling admin files: ", err)
			if txErr := tx.Rollback(); txErr != nil {
				log.Error("Error rolling back DB transaction: ", txErr)
			}
		}
	}

	if err == nil {
        marshableFiles := make([]*MarshableFile, len(filesSlice))
        for i, f := range filesSlice {
            marshableFiles[i] = (*MarshableFile)(f)
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok", "files": marshableFiles})
	} else {
        c.Error(err).SetType(gin.ErrorTypePrivate)
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
	}
}

func getFiles(exec boil.Executor) ([]*models.File, error) {
	log.Info("Looking up files")
	f, err := models.Files(exec).All()
	if err == nil {
		return f, nil
	} else {
        return nil, err
	}
}
