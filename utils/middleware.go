package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/go-playground/validator.v8"
)

func MdbLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path // some evil middleware modify this values

		c.Next()

		log.WithFields(log.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"latency":    time.Now().Sub(start),
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}).Info()
	}
}

// Recover with error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				err, ok := rval.(error)
				if !ok {
					err = errors.Errorf("panic: %s", rval)
				}
				c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
			}
		}()

		c.Next()
	}
}

func ValidationErrorMessage(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return "required"
	case "max":
		return fmt.Sprintf("cannot be longer than %s", e.Param)
	case "min":
		return fmt.Sprintf("must be longer than %s", e.Param)
	case "len":
		return fmt.Sprintf("must be %s characters long", e.Param)
	case "email":
		return "invalid email format"
	case "hexadecimal":
		return "invalid hexadecimal value"
	default:
		return "invalid value"
	}
}

func BindErrorMessage(err error) string {
	switch err.(type) {
	case *json.SyntaxError:
		e := err.(*json.SyntaxError)
		return fmt.Sprintf("json: %s [offset: %d]", e.Error(), e.Offset)
	case *json.UnmarshalTypeError:
		e := err.(*json.UnmarshalTypeError)
		return fmt.Sprintf("json: expecting %s got %s [offset: %d]", e.Type.String(), e.Value, e.Offset)
	default:
		return err.Error()
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypePublic:
					log.Warn("Public error: %s", e.Error())

					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"status": "error", "error": e.Error()})
					}

				case gin.ErrorTypeBind:
					// Keep the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}

					switch e.Err.(type) {
					case validator.ValidationErrors:
						errs := e.Err.(validator.ValidationErrors)
						errMap := make(map[string]string)
						for field, err := range errs {
							msg := ValidationErrorMessage(err)
							log.WithFields(log.Fields{
								"field": field,
								"error": msg,
							}).Warn("Validation error")
							errMap[err.Field] = msg
						}
						c.JSON(status, gin.H{"status": "error", "errors": errMap})
					default:
						log.WithFields(log.Fields{
							"error": e.Err.Error(),
						}).Warn("Bind error")
						c.JSON(status, gin.H{
							"status": "error",
							"error":  BindErrorMessage(e.Err),
						})
					}

				default:
					// Log all other errors
					log.Error(e.Err)
					if len(rollbar.Token) != 0 {
						rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
					}
				}
			}

			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError,
					gin.H{"status": "error", "error": "Internal Server Error"})
			}
		}
	}
}
