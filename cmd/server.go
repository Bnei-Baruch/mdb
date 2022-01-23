package cmd

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-oidc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"github.com/volatiletech/sqlboiler/boil"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/permissions"
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
	//boil.SetDB(db)
	boil.DebugMode = viper.GetString("server.mode") == "debug"

	log.Info("Initializing type registries")
	utils.Must(common.InitTypeRegistries(db))

	emitter, err := events.InitEmitter()
	utils.Must(err)

	// Setup Rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")
	rollbar.CodeVersion = version.Version

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = append(corsConfig.AllowMethods, http.MethodDelete)
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true

	// Authentication
	var oidcIDTokenVerifiers []*oidc.IDTokenVerifier
	if viper.GetBool("authentication.enable") {
		for _, issuer := range viper.GetStringSlice("authentication.issuers") {
			oidcProvider, err := oidc.NewProvider(context.TODO(), issuer)
			utils.Must(err)
			oidcIDTokenVerifiers = append(oidcIDTokenVerifiers, oidcProvider.Verifier(&oidc.Config{
				SkipClientIDCheck: true,
			}))
		}
	}

	// casbin
	enforcer := permissions.NewEnforcer()
	enforcer.EnableEnforce(viper.GetBool("permissions.enable"))
	enforcer.EnableLog(viper.GetBool("permissions.log"))

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()
	router.Use(
		utils.MdbLoggerMiddleware(),
		utils.EnvMiddleware(db, emitter, enforcer, oidcIDTokenVerifiers),
		utils.ErrorHandlingMiddleware(),
		permissions.AuthenticationMiddleware(),
		cors.New(corsConfig),
		utils.RecoveryMiddleware())

	api.SetupRoutes(router)

	srv := &http.Server{
		Addr:    viper.GetString("server.bind-address"),
		Handler: router,
	}

	go func() {
		// service connections
		log.Infoln("Running application")
		if err := srv.ListenAndServe(); err != nil {
			log.Infof("Server listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infof("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown:", err)
	}

	log.Infof("Close events emitter ...")
	events.CloseEmitter(ctx)

	if len(rollbar.Token) > 0 {
		log.Infof("Wait for rollbar ...")
		rollbar.Wait()
	}

	log.Infof("Shutdown complete")
}
