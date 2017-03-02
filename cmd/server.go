package cmd

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"
	"database/sql"
	"github.com/vattle/sqlboiler/boil"
	"github.com/Gurpartap/logrus-stack"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "MDB api server",
	Run:   serverFn,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverDefaults() {
	viper.SetDefault("server", map[string]interface{}{
		"bind-address":        ":8080",
		"mode":                "debug",
		"rollbar-token":       "",
		"rollbar-environment": "development",
		"log":                 "./logs/mdb.log",
		"docs":                "./docs.html",
	})
}

var router *gin.Engine

func serverFn(cmd *cobra.Command, args []string) {
	rand.Seed(time.Now().UTC().UnixNano())
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.AddHook(logrus_stack.StandardHook())
	serverDefaults()

	log.Info("Setting up connection to MDB")
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	boil.SetDB(db)
	boil.DebugMode = true

	log.Info("Initializing type registries")
	if err := api.CONTENT_TYPE_REGISTRY.Init(); err != nil {
		log.Fatal(err)
	}
	if err := api.OPERATION_TYPE_REGISTRY.Init(); err != nil {
		log.Fatal(err)
	}

	// Setup Rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")
	rollbar.CodeVersion = version.Version

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router = gin.New()

	var recovery gin.HandlerFunc
	if len(rollbar.Token) > 0 {
		recovery = utils.RollbarRecoveryMiddleware()
	} else {
		recovery = gin.Recovery()
	}

	router.Use(
		utils.MdbLoggerMiddleware(log.StandardLogger()),
		utils.ErrorHandlingMiddleware(),
		utils.GinBodyLogMiddleware,
		cors.Default(),
		recovery)

	router.POST("/operations/capture_start", api.CaptureStartHandler)
	router.POST("/operations/capture_stop", api.CaptureStopHandler)
	router.POST("/operations/demux", api.DemuxHandler)
	router.POST("/operations/send", api.SendHandler)
	router.POST("/operations/upload", api.UploadHandler)

	// Serve the log file.
	admin := router.Group("admin")
	admin.StaticFile("/log", viper.GetString("server.log"))

	// Serve the auto generated docs.
	router.StaticFile("/docs", viper.GetString("server.docs"))

	collections := router.Group("collections")
	collections.POST("/", api.CollectionsCreateHandler)

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})

	log.Infoln("Running application")
	if cmd != nil {
		router.Run(viper.GetString("server.bind-address"))
	}

	// This would be reasonable once we'll have graceful shutdown implemented
	//if len(rollbar.Token) > 0 {
	//	rollbar.Wait()
	//}
}
