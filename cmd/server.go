package cmd

import (
    "github.com/Bnei-Baruch/mdb/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
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
	router := gin.Default()

    router.POST("/operations/capture", func(c *gin.Context) {
        var op models.Operation
        if c.BindJSON(&op) == nil {
            c.JSON(http.StatusOK, gin.H{"status": "HHH"})
        }
    })

	router.Run(viper.GetString("server.bind-address"))
}
