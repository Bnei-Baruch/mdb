package rest

import (
    "bytes"
    "encoding/json"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type CaptureStart struct {
    Type        string  `json:"type"`
    Station     string  `json:"station"`
    User        string  `json:"user"`
    FileName    string  `json:"file_name"`
    CaptureID   string  `json:"capture_id"`
}

func CaptureStartHandler(c *gin.Context) {
    var cs CaptureStart
    example := &CaptureStart{
        Type:       "some_type",
        Station:    "1.2.3.4",
        User:       "username",
        FileName:   "this/is/file.name",
        CaptureID:  "13eA3b1341ff",
    }
    log.Infoln("Content Length:", c.Request.ContentLength);
    buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    body := buf.String()
    log.Infoln("Body:", body);
    if err := json.Unmarshal([]byte(body), &cs); err == nil {
        if cs.Type == "" || cs.Station == "" || cs.User == "" || cs.FileName == "" || cs.CaptureID == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "Error": "One or more required fields are empty.",
                "Input": body,
                "Example": example,
            })
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
        return
    } else {
        c.JSON(http.StatusBadRequest, gin.H{
            "Error": "Could not parse JSON from input paylod.",
            "Details": err.Error(),
            "Example": example,
        })
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }
}

