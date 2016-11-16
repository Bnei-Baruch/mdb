package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/spf13/viper"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/stvp/rollbar"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "MDB api server",
	Run: serverFn,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverDefaults() {
	viper.SetDefault("server", map[string]interface{}{
		"bind-address": ":8080",
		"mode": "debug",
		"rollbar-token": "",
		"rollbar-environment": "development",
	})
}

func serverFn(cmd *cobra.Command, args []string) {
	serverDefaults()

	// Setup logging
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Setup rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()

	var recovery gin.HandlerFunc
	if len(rollbar.Token) > 0 {
		recovery = utils.RollbarRecovery()
	} else {
		recovery = gin.Recovery()
	}

	router.Use(utils.MdbLoggerMiddleware(log.StandardLogger()), recovery)

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.GET("/recover", func(c *gin.Context) {
		panic("test recover")
	})

	log.Infoln("Running application")
	router.Run(viper.GetString("server.bind-address"))

	// This would be reasonable once we'll have graceful shutdown implemented
	//if len(rollbar.Token) > 0 {
	//	rollbar.Wait()
	//}
}
