// @title URL Shortener API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Api-Key

package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toorop/gin-logrus"
	"math/rand"
	"net/http"
	"short-url-server/docs"
	_ "short-url-server/docs"
	"short-url-server/handlers"
	"short-url-server/repo"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var log = logrus.New()
	db := repo.New(log)

	server := gin.New()
	server.RedirectTrailingSlash = false
	server.Use(ginlogrus.Logger(log), gin.Recovery())

	server.GET("/", func(c *gin.Context) {
		c.File("./assets/index.html")
	})

	urlGroup := server.Group("/urls")
	{
		urlGroup.GET("", handlers.UrlHandler(db))
		urlGroup.POST("", handlers.AddUrlHandler(db))
		urlGroup.GET("/:id", handlers.SingleUrlHandler(db))
		urlGroup.PATCH("/:id", handlers.SingleUrlUpdateHandler(db))
		urlGroup.DELETE("/:id", handlers.SingleUrlDeleteHandler(db))
	}

	statGroup := server.Group("/stats")
	{
		statGroup.GET("", handlers.StatHandler(db))
		statGroup.GET("/:id", handlers.SingleStatHandler(db))
	}

	docInit()
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func docInit() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
