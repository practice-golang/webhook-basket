package main // import "webhook-basket"

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

func WriteContent(c *gin.Context) {
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

	// r := gin.Default()
	r := gin.New()

	// r.Use(gin.Logger())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.Use(middleware.RequestLoggerMiddleware())

	r.Use(gin.Recovery())

	r.GET("/health", HealthCheck)
	r.POST("/write", WriteContent)

	r.Run("localhost:7749")
}
