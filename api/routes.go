package api

import (
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func SetupRoutes(router *gin.Engine) {
	operations := router.Group("operations")
	operations.POST("/capture_start", CaptureStartHandler)
	operations.POST("/capture_stop", CaptureStopHandler)
	operations.POST("/demux", DemuxHandler)
	operations.POST("/trim", TrimHandler)
	operations.POST("/send", SendHandler)
	operations.POST("/convert", ConvertHandler)
	operations.POST("/upload", UploadHandler)

	rest := router.Group("rest")
	rest.GET("/collections/", CollectionsListHandler)
	rest.GET("/collections/:id/", CollectionItemHandler)
	rest.POST("/collections/:id/activate", CollectionActivateHandler)
	rest.GET("/content_units/", ContentUnitsListHandler)
	rest.GET("/content_units/:id/", ContentUnitItemHandler)
	rest.GET("/files/", FilesListHandler)
	rest.GET("/files/:id/", FileItemHandler)

	hierarchy := router.Group("hierarchy")
	hierarchy.GET("/sources/", SourcesHierarchyHandler)
	hierarchy.GET("/tags/", TagsHierarchyHandler)

	// Serve the log file.
	admin := router.Group("admin")
	admin.StaticFile("/log", viper.GetString("server.log"))

	// Serve the auto generated docs.
	router.StaticFile("/docs", viper.GetString("server.docs"))

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})
}
