package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"database/sql"
	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"
	"github.com/vattle/sqlboiler/boil"
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
		"log":  "./logs/mdb.log",
		"docs": "./docs.html",
	})
}

func serverFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	serverDefaults()

	log.Infof("Starting MDB API server version %s", version.Version)

	log.Info("Setting up connection to MDB")
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	boil.SetDB(db)
	boil.DebugMode = viper.GetString("server.mode") == "debug"

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
	router := gin.New()

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

	api.SetupRoutes(router)

	log.Infoln("Running application")
	if cmd != nil {
		router.Run(viper.GetString("server.bind-address"))
	}

	// This would be reasonable once we'll have graceful shutdown implemented
	//if len(rollbar.Token) > 0 {
	//	rollbar.Wait()
	//}
}
