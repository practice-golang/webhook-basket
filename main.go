package main // import "webhook-basket"

import (
	"io"
	"log"
	"net/http"
	"os"
	"webhook-basket/middleware"
	"webhook-basket/model"

	"github.com/gin-gonic/gin"
)

type Content struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func HealthCheck(c *gin.Context) {
	result := gin.H{"status": "ok"}

	c.JSON(http.StatusOK, result)
}

func PostSample(c *gin.Context) {
	queries := c.Request.URL.Query()
	for k, v := range queries {
		log.Println(k, v)
	}

	content := Content{}
	c.BindJSON(&content)

	c.JSON(http.StatusOK, content)
}

func DeployRepository(c *gin.Context) {
	content := Content{}
	c.BindJSON(&content)

	c.JSON(http.StatusOK, content)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	// Logging to a file.
	model.FileRequests, _ = os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	model.FileConnections, _ = os.Create("connections.log")
	defer model.FileConnections.Close()
	defer model.FileRequests.Close()

	gin.DefaultWriter = io.MultiWriter(model.FileConnections)

	r := gin.New()

	r.Use(gin.LoggerWithFormatter(middleware.LogFormatter))
	r.Use(middleware.RequestLoggerMiddleware())

	r.Use(gin.Recovery())

	r.GET("/health", HealthCheck)
	r.POST("/post-sample", PostSample)
	r.POST("/deploy", DeployRepository)

	r.Run("127.0.0.1:7749")
}
