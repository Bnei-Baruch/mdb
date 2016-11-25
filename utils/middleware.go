package utils

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"net/http"
	"github.com/stvp/rollbar"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v8"
)

func MdbLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path  // some evil middleware modify this values

		c.Next()

		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"latency":    time.Now().Sub(start),
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}).Info()

		if len(c.Errors) > 0 {
			logger.Error(c.Errors.String())
		}
	}
}

func RollbarRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// Log panics
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

		// Log context errors
		for _, err := range c.Errors.ByType(gin.ErrorTypePrivate) {
			rollbar.RequestError(rollbar.CRIT, c.Request, err.Err)
		}
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if be := c.Errors.ByType(gin.ErrorTypeBind).Last(); be != nil {
			var errorMessages []interface{}
			for _, err := range be.Err.(validator.ValidationErrors) {
				errorMessages = append(errorMessages, gin.H{err.Field: err.ActualTag})
			}
			c.JSON(-1, gin.H{"errors":errorMessages})
		}
	}
}

