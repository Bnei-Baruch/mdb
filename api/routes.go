package api

import (
	"github.com/pkg/errors"
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
	rest.GET("/collections/:id/content_units/", CollectionContentUnitsHandler)
	rest.POST("/collections/:id/activate", CollectionActivateHandler)
	rest.GET("/content_units/", ContentUnitsListHandler)
	rest.GET("/content_units/:id/", ContentUnitItemHandler)
	rest.GET("/content_units/:id/files/", ContentUnitFilesHandler)
	rest.GET("/content_units/:id/collections/", ContentUnitCollectionsHandler)
	rest.GET("/files/", FilesListHandler)
	rest.GET("/files/:id/", FileItemHandler)
	rest.GET("/operations/", OperationsListHandler)
	rest.GET("/operations/:id/", OperationItemHandler)
	rest.GET("/operations/:id/files/", OperationFilesHandler)
	rest.GET("/tags/", TagsHandler)
	rest.POST("/tags/", TagsHandler)
	rest.GET("/tags/:id/", TagHandler)
	rest.PUT("/tags/:id/", TagHandler)
	rest.PUT("/tags/:id/i18n/", TagI18nHandler)

	hierarchy := router.Group("hierarchy")
	hierarchy.GET("/sources/", SourcesHierarchyHandler)
	hierarchy.GET("/tags/", TagsHierarchyHandler)

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})
	router.GET("/error", func(c *gin.Context) {
		c.AbortWithError(500,
			errors.Wrap(errors.New("test error with stack"), "wrap msg")).
			SetType(gin.ErrorTypePrivate)
	})
}
