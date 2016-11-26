package rest

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type CaptureStart struct {
	Type      string  `json:"type"`
	Station   string  `json:"station"`
	User      string  `json:"user"`
	FileName  string  `json:"file_name" binding:"required,max=25"`
	CaptureID string  `json:"capture_id" binding:"required,max=255"`
}

func CaptureStartHandler(c *gin.Context) {
	var cs CaptureStart
	if c.BindJSON(&cs) == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

