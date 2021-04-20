package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"math/rand"
	"net/http"
	"short-url-server/handlers"
	"short-url-server/repo"
	"time"
)


func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var log = logrus.New()
	db := repo.New(log)

	server := gin.New()
	server.Use(ginlogrus.Logger(log), gin.Recovery())

	server.GET("/", func(c *gin.Context) {
		c.File("./assets/index.html")
	})

	urlGroup := server.Group("/urls")
	{
		urlGroup.GET("/", handlers.UrlHandler(db))
		urlGroup.POST("/", handlers.AddUrlHandler(db))
		urlGroup.GET("/:id", handlers.SingleUrlHandler(db))
		urlGroup.PATCH("/:id", handlers.SingleUrlUpdateHandler(db))
		urlGroup.DELETE("/:id", handlers.SingleUrlDeleteHandler(db))
	}

	statGroup := server.Group("/stats")
	{
		statGroup.GET("/", handlers.StatHandler(db))
		statGroup.GET("/:id", handlers.SingleStatHandler(db))
	}

	server.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if url, err := db.FindUrlById(id); err != nil || url == nil || len(url.Url) < 1 {
			fmt.Println(id)
			c.File("./assets/" + id)
		} else {
			c.Redirect(http.StatusMovedPermanently, url.Url)
		}
	})

	server.NoRoute(static.ServeRoot("/", "./assets/"))

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
