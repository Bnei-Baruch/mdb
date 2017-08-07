package cmd

import (
	"database/sql"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "MDB api server",
	Run:   serverFn,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Infof("Starting MDB API server version %s", version.Version)

	log.Info("Setting up connection to MDB")
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	defer db.Close()
	boil.SetDB(db)
	boil.DebugMode = viper.GetString("server.mode") == "debug"

	log.Info("Initializing type registries")
	utils.Must(api.InitTypeRegistries(db))

	// Setup Rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")
	rollbar.CodeVersion = version.Version

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = append(corsConfig.AllowMethods, http.MethodDelete)
	corsConfig.AllowAllOrigins = true

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()
	router.Use(
		utils.MdbLoggerMiddleware(),
		utils.ErrorHandlingMiddleware(),
		cors.New(corsConfig),
		utils.RecoveryMiddleware())

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
