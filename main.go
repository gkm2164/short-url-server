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

	server.GET("/", func(c *gin.Context) {
		c.File("./assets/index.html")
	})
	server.GET("/index.html", func(c *gin.Context) {
		c.File("./assets/index.html")
	})
	server.GET("/main.js", func(c *gin.Context) {
		c.File("./assets/main.js")
	})
	server.POST("/", func(c *gin.Context) {
		var req NewAddrRequest
		id := randomString()
		if value, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": fmt.Sprintf("while read: %v", err),
			})
		} else if err := json.Unmarshal(value, &req); err != nil {
			fmt.Println(string(value))
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

	server.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if url, err := db.FindUrlById(id); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.Redirect(http.StatusMovedPermanently, url.Url)
		}
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func randomString() string {
	size := 10
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	retStr := make([]uint8, size)
	for i := 0; i < size; i++ {
		retStr[i] = letters[rand.Int() % len(letters)]
	}
	return string(retStr)
}