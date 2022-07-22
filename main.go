package main // import "webhook-basket"

import (
	"io"
	"log"
	"net/http"
	"os"
	"webhook-basket/downloader"
	"webhook-basket/middleware"
	"webhook-basket/model"
	"webhook-basket/util"

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
	destination := c.Query("destination")
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deployment destination is required"})
		return
	}

	request := model.Request{}
	c.BindJSON(&request)

	request.Destination = destination

	err := downloader.CloneRepository(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

func DeleteReposRoot(c *gin.Context) {
	destination := model.CloneRepoRoot

	err := util.DeleteDirectory(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
	r.DELETE("/repos-root", DeleteReposRoot)

	r.Run("127.0.0.1:7749")
}
