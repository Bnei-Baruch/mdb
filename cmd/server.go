package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/spf13/viper"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
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
	})
}

func serverFn(cmd *cobra.Command, args []string) {
	serverDefaults()

	// Setup logging
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Setup routes
	router := gin.New()
	router.Use(utils.MdbLoggerMiddleware(log.StandardLogger()), gin.Recovery())

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

	router.Run(viper.GetString("server.bind-address"))
}
