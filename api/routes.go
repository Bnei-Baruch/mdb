package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health_check", HealthCheckHandler)

	operations := router.Group("operations")
	operations.POST("/capture_start", CaptureStartHandler)
	operations.POST("/capture_stop", CaptureStopHandler)
	operations.POST("/demux", DemuxHandler)
	operations.POST("/trim", TrimHandler)
	operations.POST("/send", SendHandler)
	operations.POST("/convert", ConvertHandler)
	operations.POST("/upload", UploadHandler)
	operations.POST("/sirtutim", SirtutimHandler)
	operations.POST("/insert", InsertHandler)
	operations.POST("/transcode", TranscodeHandler)
	operations.POST("/join", JoinHandler)
	operations.POST("/replace", ReplaceHLS)
	operations.GET("/descendant_units/:sha1", DescendantUnitsHandler)

	rest := router.Group("rest")
	rest.GET("/collections/", CollectionsListHandler)
	rest.POST("/collections/", CollectionsListHandler)
	rest.GET("/collections/:id/", CollectionHandler)
	rest.PUT("/collections/:id/", CollectionHandler)
	rest.DELETE("/collections/:id/", CollectionHandler)
	rest.PUT("/collections/:id/i18n/", CollectionI18nHandler)
	rest.GET("/collections/:id/content_units/", CollectionContentUnitsHandler)
	rest.POST("/collections/:id/order_positions", CollectionContentUnitsPositionHandler)
	rest.POST("/collections/:id/content_units/", CollectionContentUnitsHandler)
	rest.PUT("/collections/:id/content_units/:cuID", CollectionContentUnitsHandler)
	rest.DELETE("/collections/:id/content_units/:cuID", CollectionContentUnitsHandler)
	rest.POST("/collections/:id/activate", CollectionActivateHandler)
	rest.GET("/content_units/", ContentUnitsListHandler)
	rest.POST("/content_unit/autoname", ContentUnitAutoname)
	rest.POST("/content_units/", ContentUnitsListHandler)
	rest.GET("/content_units/:id/", ContentUnitHandler)
	rest.PUT("/content_units/:id/", ContentUnitHandler)
	rest.PUT("/content_units/:id/i18n/", ContentUnitI18nHandler)
	rest.GET("/content_units/:id/files/", ContentUnitFilesHandler)
	rest.POST("/content_units/:id/files/", ContentUnitFilesHandler)
	rest.GET("/content_units/:id/collections/", ContentUnitCollectionsHandler)
	rest.GET("/content_units/:id/derivatives/", ContentUnitDerivativesHandler)
	rest.POST("/content_units/:id/derivatives/", ContentUnitDerivativesHandler)
	rest.PUT("/content_units/:id/derivatives/:duID", ContentUnitDerivativesHandler)
	rest.DELETE("/content_units/:id/derivatives/:duID", ContentUnitDerivativesHandler)
	rest.GET("/content_units/:id/origins/", ContentUnitOriginsHandler)
	rest.GET("/content_units/:id/sources/", ContentUnitSourcesHandler)
	rest.POST("/content_units/:id/sources/", ContentUnitSourcesHandler)
	rest.DELETE("/content_units/:id/sources/:sourceID", ContentUnitSourcesHandler)
	rest.GET("/content_units/:id/tags/", ContentUnitTagsHandler)
	rest.POST("/content_units/:id/tags/", ContentUnitTagsHandler)
	rest.DELETE("/content_units/:id/tags/:tagID", ContentUnitTagsHandler)
	rest.GET("/content_units/:id/persons/", ContentUnitPersonsHandler)
	rest.POST("/content_units/:id/persons/", ContentUnitPersonsHandler)
	rest.DELETE("/content_units/:id/persons/:personID", ContentUnitPersonsHandler)
	rest.GET("/content_units/:id/publishers/", ContentUnitPublishersHandler)
	rest.POST("/content_units/:id/publishers/", ContentUnitPublishersHandler)
	rest.DELETE("/content_units/:id/publishers/:publisherID", ContentUnitPublishersHandler)
	rest.POST("/content_units/:id/merge", ContentUnitMergeHandler)
	rest.GET("/files/", FilesListHandler)
	rest.GET("/files/:id/", FileHandler)
	rest.PUT("/files/:id/", FileHandler)
	rest.GET("/files/:id/storages/", FileStoragesHandler)
	rest.GET("/files/:id/tree/", FilesWithOperationsTreeHandler)
	rest.GET("/operations/", OperationsListHandler)
	rest.GET("/operations/:id/", OperationItemHandler)
	rest.GET("/operations/:id/files/", OperationFilesHandler)
	rest.GET("/authors/", AuthorsHandler)
	rest.GET("/sources/", SourcesHandler)
	rest.POST("/sources/", SourcesHandler)
	rest.GET("/sources/:id/", SourceHandler)
	rest.PUT("/sources/:id/", SourceHandler)
	rest.PUT("/sources/:id/i18n/", SourceI18nHandler)
	rest.GET("/tags/", TagsHandler)
	rest.POST("/tags/", TagsHandler)
	rest.GET("/tags/:id/", TagHandler)
	rest.PUT("/tags/:id/", TagHandler)
	rest.PUT("/tags/:id/i18n/", TagI18nHandler)
	rest.GET("/persons/", PersonsListHandler)
	rest.POST("/persons/", PersonsListHandler)
	rest.GET("/persons/:id/", PersonHandler)
	rest.PUT("/persons/:id/", PersonHandler)
	rest.DELETE("/persons/:id/", PersonHandler)
	rest.PUT("/persons/:id/i18n/", PersonI18nHandler)
	rest.GET("/storages/", StoragesHandler)
	rest.GET("/publishers/", PublishersHandler)
	rest.POST("/publishers/", PublishersHandler)
	rest.GET("/publishers/:id/", PublisherHandler)
	rest.PUT("/publishers/:id/", PublisherHandler)
	rest.PUT("/publishers/:id/i18n/", PublisherI18nHandler)
	rest.GET("/labels/", LabelListHandler)
	rest.POST("/labels/", LabelListHandler)
	rest.GET("/labels/:id/", LabelHandler)
	rest.PUT("/labels/:id/", LabelHandler)
	rest.DELETE("/labels/:id/", LabelHandler)
	rest.PUT("/labels/:id/i18n/", LabelI18nHandler)
	rest.POST("/labels/:uid/i18n/", LabelAddI18nHandler)

	hierarchy := router.Group("hierarchy")
	hierarchy.GET("/sources/", SourcesHierarchyHandler)
	hierarchy.GET("/tags/", TagsHierarchyHandler)

	//router.GET("/recover", func(c *gin.Context) {
	//	panic("test recover")
	//})
	//router.GET("/error", func(c *gin.Context) {
	//	c.AbortWithError(500,
	//		errors.Wrap(errors.New("test error with stack"), "wrap msg")).
	//		SetType(gin.ErrorTypePrivate)
	//})
}
