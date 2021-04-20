package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"short-url-server/repo"
	"strings"
	"time"
)

type NewAddrRequest struct {
	Url string `json:"url"`
}

const Succeed = "succeed"
const Failed = "failed"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var log = logrus.New()
	db := repo.New(log)

	server := gin.New()
	server.Use(ginlogrus.Logger(log), gin.Recovery())
	server.GET("/*path", func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.RequestURI, "/")
		if len(path) <= 1 {
			path = "index.html"
		}
		if strings.HasPrefix(path, "urls") {
			urlHandler(db)(c)
		} else if strings.HasPrefix(path, "stats") {
			urlStatHandler(db)(c)
		} else if url, err := db.FindUrlById(path); err != nil || url == nil || len(url.Url) < 1 {
			c.File("./assets/" + path)
		} else {
			c.Redirect(http.StatusMovedPermanently, url.Url)
		}
	})

	server.DELETE("/*path", func(c *gin.Context) {
		var path = c.Param("path")
		if path == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  Failed,
				"message": "path should be not empty",
			})
		} else if err := db.DeleteUrl(path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  Failed,
				"message": err,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"status":  Succeed,
				"message": "deleted successfully!",
			})
		}
	})

	server.POST("/", func(c *gin.Context) {
		var req NewAddrRequest
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status": Failed,
				"message": fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status": Failed,
				"message": fmt.Sprintf("while unmarshal: %v", err),
			})
		} else if id, err := insertUrlUntilSuccess(db, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": Failed,
				"message": fmt.Sprintf("while insert to db: %v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"status": Succeed,
				"id": id,
			})
		}
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

const generalQuery = `SELECT REGEXP_REPLACE(REGEXP_REPLACE(url, 'https?://', ''), '/.*', '') AS HOSTNAME, COUNT(*) AS CNT FROM urls GROUP BY HOSTNAME`
const queryByHostName = `
SELECT * FROM (
    SELECT REGEXP_REPLACE(REGEXP_REPLACE(url, 'https?://', ''), '/.*', '') AS HOSTNAME, COUNT(*) AS CNT
    FROM urls
    GROUP BY HOSTNAME) A
WHERE HOSTNAME='%s';`

var IdRegexUnderUrl = regexp.MustCompile("/urls/(?P<id>[a-zA-Z0-9]+)(\\?.*)?")
var IdRegexUnderStat = regexp.MustCompile("/stats/(?P<id>[a-zA-Z0-9]+)(\\?.*)?")
var idIdxUnderUrl = IdRegexUnderUrl.SubexpIndex("id")
var idIdxUnderStat = IdRegexUnderStat.SubexpIndex("id")

func urlHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := strings.TrimPrefix(c.Request.RequestURI, "/urls")
		if IdRegexUnderUrl.MatchString(uri) {
			idmatch := IdRegexUnderUrl.FindStringSubmatch(uri)[idIdxUnderUrl]
			if url, err := db.FindUrlById(idmatch); err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  Failed,
					"message": "id not found",
				})
			} else {
				c.JSON(http.StatusFound, gin.H{
					"status": Succeed,
					"url":    url,
				})
			}
		} else {
			if ret, err := db.FindAllUrls(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  Failed,
					"message": "retrieving urls failed",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": Succeed,
					"urls":   ret,
				})
			}
		}
	}
}

func urlStatHandler(db *repo.DDB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
					"status": Failed,
					"message": "Error while digging from DB",
				})
				goto ret
			}
		} else {
			if rows, err = db.DB().Raw(fmt.Sprintf(queryByHostName, hostname)).Rows(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": Failed,
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
			"status": Succeed,
			"stats":  rets,
		})
	ret:
	}
}

const UrlSize = 11

func insertUrlUntilSuccess(db *repo.DDB, url string) (*string, error) {
	var success = false
	for !success {
		id := randomString()
		if _, err := db.InsertUrl(id, url); err != nil {
			return nil, fmt.Errorf("errors when trying to put item: %v", err)
		} else {
			success = true
			return &id, nil
		}
	}

	return nil, errors.New("should not reach here")
}

func randomString() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	retStr := make([]uint8, UrlSize)
	for i := 0; i < UrlSize; i++ {
		retStr[i] = letters[rand.Uint32()%uint32(len(letters))]
	}
	return string(retStr)
}
