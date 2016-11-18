package utils

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"net/http"
	"github.com/stvp/rollbar"
	"github.com/pkg/errors"
)

func MdbLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path  // some evil middleware modify this values

		c.Next()

		entry := logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"latency":    time.Now().Sub(start),
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}

func RollbarRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				if err, ok := rval.(error); ok {
					rollbar.RequestError(rollbar.CRIT, c.Request, err)
				} else {
					rollbar.RequestError(rollbar.CRIT, c.Request, errors.Errorf("%s", rval))

				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

