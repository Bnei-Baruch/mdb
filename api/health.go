package api

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/gin-gonic/gin.v1"
)

func HealthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	mdb := c.MustGet("MDB").(*sql.DB)

	// Uncomment once this lib/pq PR is merged
	// https://github.com/lib/pq/pull/737
	//err := mdb.PingContext(ctx)
	err := PingDB(ctx, mdb)

	if err == nil {
		err = ctx.Err()
	}

	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func PingDB(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "select 1")
	if err != nil {
		return driver.ErrBadConn // https://golang.org/pkg/database/sql/driver/#Pinger
	}
	defer rows.Close()
	return nil
}
