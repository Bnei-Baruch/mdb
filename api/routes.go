package api

import (
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func SetupRoutes(router *gin.Engine) {
    // Operations API
	router.POST("/operations/capture_start", CaptureStartHandler)
	router.POST("/operations/capture_stop", CaptureStopHandler)
	router.POST("/operations/demux", DemuxHandler)
	router.POST("/operations/trim", TrimHandler)
	router.POST("/operations/send", SendHandler)
	router.POST("/operations/upload", UploadHandler)

    // Admin API
	admin := router.Group("admin")

    // Serving admin UI
	admin.StaticFile("/", viper.GetString("server.admin-ui"))
	admin.Static("/build", "./admin-ui/build/")

    // Admin rest handlers.
    admin.GET("/rest/files", AdminFilesHandler)
	admin.StaticFile("/rest/log", viper.GetString("server.log"))


	// Serve the auto generated docs.
	router.StaticFile("/docs", viper.GetString("server.docs"))

	collections := router.Group("collections")
	collections.POST("/", CollectionsCreateHandler)

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})
}
