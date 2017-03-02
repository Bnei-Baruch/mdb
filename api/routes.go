package api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/spf13/viper"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/operations/capture_start", CaptureStartHandler)
	router.POST("/operations/capture_stop", CaptureStopHandler)
	router.POST("/operations/demux", DemuxHandler)
	router.POST("/operations/send", SendHandler)
	router.POST("/operations/upload", UploadHandler)

	// Serve the log file.
	admin := router.Group("admin")
	admin.StaticFile("/log", viper.GetString("server.log"))

	// Serve the auto generated docs.
	router.StaticFile("/docs", viper.GetString("server.docs"))

	collections := router.Group("collections")
	collections.POST("/", CollectionsCreateHandler)

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})
}
