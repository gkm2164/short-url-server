package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"short-url-server/common"
	"short-url-server/repo"
)

const generalQuery = `SELECT REGEXP_REPLACE(REGEXP_REPLACE(url, 'https?://', ''), '/.*', '') AS HOSTNAME, COUNT(*) AS CNT FROM urls GROUP BY HOSTNAME`
const queryByHostName = `
SELECT * FROM (
    SELECT REGEXP_REPLACE(REGEXP_REPLACE(url, 'https?://', ''), '/.*', '') AS HOSTNAME, COUNT(*) AS CNT
    FROM urls
    GROUP BY HOSTNAME) A
WHERE HOSTNAME='%s';`

func StatHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c * gin.Context) {
		var rows *sql.Rows
		var err error
		type RetStruct struct {
			Hostname string `json:"hostname"`
			Count    int    `json:"count"`
		}

		var rets []RetStruct
		hostname := c.Query("hostname")
		if hostname == "" {
			if rows, err = db.DB().Raw(generalQuery).Rows(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": common.Failed,
					"message": "Error while digging from DB",
				})
				goto ret
			}
		} else {
			if rows, err = db.DB().Raw(fmt.Sprintf(queryByHostName, hostname)).Rows(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": common.Failed,
					"message": "Error while digging from DB",
				})
				goto ret
			}
		}

		defer rows.Close()

		for rows.Next() {
			var hostname string
			var count int

			_ = rows.Scan(&hostname, &count)
			rets = append(rets, RetStruct{
				Hostname: hostname,
				Count:    count,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": common.Succeed,
			"stats":  rets,
		})
	ret:
	}
}

func SingleStatHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithError(http.StatusNotImplemented, errors.New("thinking..."))
	}
}