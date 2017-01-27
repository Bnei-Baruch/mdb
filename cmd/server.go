package cmd

import (
	"github.com/Bnei-Baruch/mdb/dal"
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"math/rand"
	"net/http"
	"time"
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
	})
}

var router *gin.Engine

func serverFn(cmd *cobra.Command, args []string) {
	rand.Seed(time.Now().UTC().UnixNano())
	serverDefaults()

	// Setup logging
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Setup rollbar
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

	router.Use(utils.MdbLoggerMiddleware(log.StandardLogger()),
		utils.ErrorHandlingMiddleware(),
		utils.GinBodyLogMiddleware,
		cors.Default(),
		recovery)

	dal.Init()
	router.POST("/operations/capture_start", CaptureStartHandler)
	router.POST("/operations/capture_stop", CaptureStopHandler)
	router.POST("/operations/demux", DemuxHandler)
	router.POST("/operations/send", SendHandler)
	router.POST("/operations/upload", UploadHandler)

    // Serve the log file.
	admin := router.Group("admin")
	admin.StaticFile("/log", viper.GetString("server.log"))

    // Serve the auto generated docs.
	router.StaticFile("/docs", viper.GetString("server.docs"))

	collections := router.Group("collections")
	collections.POST("/", rest.CollectionsCreateHandler)

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

// Starts capturing file, i.e., morning lesson or other program.
func CaptureStartHandler(c *gin.Context) {
	var cs rest.CaptureStart
	if c.BindJSON(&cs) == nil {
        if err := dal.CaptureStart(cs); err == nil {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  err.Error(),
            })
        }
    }
}

// Stops capturing file, i.e., morning lesson or other program.
func CaptureStopHandler(c *gin.Context) {
	var cs rest.CaptureStop
	if c.BindJSON(&cs) == nil {
        if err := dal.CaptureStop(cs); err == nil {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  err.Error(),
            })
        }
	}
}

// Demux
func DemuxHandler(c *gin.Context) {
	var demux rest.Demux
	if c.BindJSON(&demux) == nil {
        if err := dal.Demux(demux); err == nil {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  err.Error(),
            })
        }
	}
}

// Moves file from capture machine to other storage.
func SendHandler(c *gin.Context) {
	var send rest.Send
	if c.BindJSON(&send) == nil {
        if err := dal.Send(send); err == nil {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  err.Error(),
            })
        }
	}
}

// Enabled file to be accessible from URL.
func UploadHandler(c *gin.Context) {
	var upload rest.Upload
	if c.BindJSON(&upload) == nil {
        if err := dal.Upload(upload); err == nil {
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  err.Error(),
            })
        }
	}
}
