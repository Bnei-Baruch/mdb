package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/casbin/casbin"
	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"github.com/stvp/rollbar"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/go-playground/validator.v8"

	"github.com/Bnei-Baruch/mdb/events"
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

func EnvMiddleware(mdb *sql.DB, emitter events.EventEmitter, enforcer *casbin.Enforcer,
	tokenVerifiers []*oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MDB", mdb)
		c.Set("EVENTS_EMITTER", emitter)
		c.Set("PERMISSIONS_ENFORCER", enforcer)
		c.Set("TOKEN_VERIFIERS", tokenVerifiers)
		c.Next()
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
					if e.Err != nil {
						errMsg := e.Error()
						if !strings.Contains(errMsg, "oidc: token is expired ") {
							log.Warnf("Public error: %s", errMsg)
						}
						c.JSON(c.Writer.Status(), gin.H{"status": "error", "error": errMsg})
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
					st, ok := e.Err.(StackTracer)
					if ok {
						fmt.Printf("%s: %+v\n", st, st.StackTrace())
					}

					// Log to rollbar if we have a token setup
					if len(rollbar.Token) != 0 {
						if ok {
							rollbar.RequestErrorWithStack(rollbar.ERR, c.Request, e.Err,
								ErrorsToRollbarStack(st))
						} else {
							rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
						}
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
