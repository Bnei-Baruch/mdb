package utils

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"

	"bytes"
	"fmt"
	"io"
    "io/ioutil"
	"net/http"
    "regexp"
    "strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stvp/rollbar"
)

func MdbLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path // some evil middleware modify this values

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
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  be.Err.Error(),
			})
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type bodyLogReadCloser struct {
	body   *bytes.Buffer
	reader io.Reader
	closer io.Closer
}

func (blr *bodyLogReadCloser) Close() error {
	return blr.closer.Close()
}

func (blr bodyLogReadCloser) Read(p []byte) (n int, err error) {
	return blr.reader.Read(p)
}

func GinBodyLogMiddleware(c *gin.Context) {
	blr := &bodyLogReadCloser{
		body:   bytes.NewBufferString(""),
		closer: c.Request.Body,
	}
	blw := &bodyLogWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: c.Writer,
	}
	if "/admin/log" != c.Request.URL.Path {
		blr.reader = io.TeeReader(c.Request.Body, blr.body)
		c.Request.Body = blr
		c.Writer = blw
	}

	c.Next()

	if "/admin/log" != c.Request.URL.Path {
		fmt.Println("Request body: " + blr.body.String())
		fmt.Println("Response body: " + blw.body.String())
	}
}

func FixDoubleQuotesInInput(c *gin.Context) {
    // TODO: Remove when input request fixed.
    b, _ := ioutil.ReadAll(c.Request.Body)
    body := string(b)
    fmt.Printf("body:%s\n", body)
    re := regexp.MustCompile("\"(\"[[:alnum:]]*\")\"")
    if matches := re.FindStringSubmatch(body); len(matches) > 1 {
        fmt.Printf("matches:%v\n", matches)
        for _, m := range matches[1:] {
            body = strings.Replace(
                body,
                fmt.Sprintf("\"%s\"", m),
                fmt.Sprintf("%s", m),
                1)
        }
    }
    fmt.Printf("fixed:%s\n", body)
    c.Request.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))

    c.Next()
}

