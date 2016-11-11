package utils

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
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

