package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"short-url-server/repo"
	"strings"
	"time"
)

type NewAddrRequest struct {
	Url string `json:"url"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var log = logrus.New()
	db := repo.New(log)

	server := gin.New()
	server.Use(ginlogrus.Logger(log), gin.Recovery())
	server.GET("/*path", func(c *gin.Context) {
		path := stripPath(c.Request.RequestURI, "/")
		if len(path) <= 1 {
			path = "index.html"
		}
		if url, err := db.FindUrlById(path); err != nil || url == nil || len(url.Url) < 1 {
			c.File("./assets/" + path)
		} else {
			c.Redirect(http.StatusMovedPermanently, url.Url)
		}
	})

	server.DELETE("/*path", func(c *gin.Context) {
		var path = c.Param("path")
		if path == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": "path should be not empty",
			})
		} else if err := db.DeleteUrl(path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "deleted successfully!",
			})
		}
	})

	server.POST("/", func(c *gin.Context) {
		var req NewAddrRequest
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("while unmarshal: %v", err),
			})
		} else if id, err := insertUrlUntilSuccess(db, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("while insert to db: %v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"id": id,
			})
		}
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func stripPath(uri string, s string) string {
	if strings.HasPrefix(uri, s) {
		return uri[len(s):]
	}
	return uri
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
