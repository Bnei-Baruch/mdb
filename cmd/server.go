package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "MDB api server",
	Run: execute,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func setDefaults() {
	viper.SetDefault("server", map[string]interface{}{
		"bind-address": ":8080",
	})
}

func execute(cmd *cobra.Command, args []string) {
	setDefaults()
	router := gin.Default()

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
