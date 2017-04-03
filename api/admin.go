package api

import (
	"encoding/json"
    "encoding/hex"
    "fmt"
	"net/http"
    "strconv"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/vattle/sqlboiler/boil"
    "github.com/vattle/sqlboiler/queries/qm"
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

type Params struct {
    sort string
    query string
    offset int
    limit int
}

func (p *Params) read(c *gin.Context) error {
    p.sort = c.DefaultQuery("sort", "date")
    p.query = c.DefaultQuery("query", "")
    var err1, err2 error
    p.offset, err1 = strconv.Atoi(c.DefaultQuery("offset", "0"))
    p.limit, err2 = strconv.Atoi(c.DefaultQuery("limit", "30"))
    return utils.CombineErr(err1, err2)
}

func InternalError(c *gin.Context, err error) {
    c.Error(err).SetType(gin.ErrorTypePrivate)
    c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
}

// Gets list of all files.
// AdminFilesHandler should also support following:
// query - text string to search in file properties.
// limit - nax number of files to fetch.
// first - offset (for pagination)
// http://.../files?sort=X,query=Y,first=Z,limit=W
func AdminFilesHandler(c *gin.Context) {
    var p Params
    err := p.read(c)
    if err != nil {
        InternalError(c, err)
        return
    }

    tx, err := boil.Begin()
    if err != nil {
        InternalError(c, err)
        return
    }
    var filesSlice []*models.File
    filesSlice, err = getFiles(tx, p)
    if err != nil {
        log.Error("Error handling admin files: ", err)
        if txErr := tx.Rollback(); txErr != nil {
            log.Error("Error rolling back DB transaction: ", txErr)
        }
        InternalError(c, err)
        return
    }

    marshableFiles := make([]*MarshableFile, len(filesSlice))
    for i, f := range filesSlice {
        marshableFiles[i] = (*MarshableFile)(f)
    }
    searchCount, err := getFilesCount(tx, p.query)
    if err != nil {
        InternalError(c, err)
        return
    }
    count, err := getFilesCount(tx, "")
    if err != nil {
        InternalError(c, err)
        return
    } else {
        // Commit transaction when all queries are done.
        tx.Commit()
    }
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "files": marshableFiles,
        "total": count,
        "matching": searchCount,
    })
}

func getFiles(exec boil.Executor, p Params) ([]*models.File, error) {
	log.Info("Looking up files")
    if (p.limit == 0) {
        return make([]*models.File, 0), nil
    }
    f, err := models.Files(exec,
        qm.Where("name like ?", []byte(fmt.Sprintf("%%%s%%", p.query))),
        qm.Limit(p.limit),
        qm.Offset(p.offset),
        qm.OrderBy("id")).All()
	if err == nil {
		return f, nil
	} else {
        return nil, err
	}
}

func getFilesCount(exec boil.Executor, query string) (int64, error) {
    c, err := models.Files(exec, qm.Where("name like ?", []byte(fmt.Sprintf("%%%s%%", query)))).Count()
	if err == nil {
		return c, nil
	} else {
        return -1, err
	}
}
