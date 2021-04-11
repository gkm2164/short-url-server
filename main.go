package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"short-url/repo"
	"strings"
)

type NewAddrRequest struct {
	Url string `json:"url"`
}

type DeleteAddrRequest struct {
	Id string `json:"id"`
}

func main() {
	var log = logrus.New()
	db := repo.New(log)

	server := gin.New()

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

	server.POST("/", func(c *gin.Context) {
		var req NewAddrRequest
		id := randomString()
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("while unmarshal: %v", err),
			})
		} else if err := db.InsertUrl(id, req.Url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("while insert to db: %v", err),
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"id": id,
			})
		}
	})

	server.DELETE("/", func(c *gin.Context) {
		var req DeleteAddrRequest
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": err,
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": err,
			})
		} else if err := db.DeleteUrl(req.Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "deleted successfully!",
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

func randomString() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	retStr := make([]uint8, UrlSize)
	for i := 0; i < UrlSize; i++ {
		retStr[i] = letters[rand.Uint32()%uint32(len(letters))]
	}
	return string(retStr)
}
