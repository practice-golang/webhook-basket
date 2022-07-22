package main // import "webhook-basket"

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"webhook-basket/downloader"
	"webhook-basket/middleware"
	"webhook-basket/model"
	"webhook-basket/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

//go:embed wb.ini
var sampleINI string

type Content struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func setupINI() {
	iniPath := "wb.ini"

	cfg, err := ini.Load(iniPath)
	if err != nil {
		f, err := os.Create(iniPath)
		if err != nil {
			log.Fatalln("Create INI: ", err)
		}
		defer f.Close()

		_, err = f.WriteString(sampleINI + "\n")
		if err != nil {
			log.Fatalln("Create INI: ", err)
		}

		fmt.Println(iniPath + " is created")
		fmt.Println("Please modify " + iniPath + " then run again")

		os.Exit(1)
	}

	if cfg != nil {
		if cfg.Section("path").HasKey("CLONED_REPO_ROOT") {
			model.TempClonedRepoRoot = cfg.Section("path").Key("CLONED_REPO_ROOT").String()
		}
		if cfg.Section("path").HasKey("DEPLOYMENT_ROOT") {
			model.TempClonedRepoRoot = cfg.Section("path").Key("DEPLOYMENT_ROOT").String()
		}

		if cfg.Section("git").HasKey("USERNAME") {
			model.AuthInfo.Username = cfg.Section("git").Key("USERNAME").String()
		}
		if cfg.Section("git").HasKey("PASSWORD") {
			model.AuthInfo.Password = cfg.Section("git").Key("PASSWORD").String()
		}

		if cfg.Section("ftp").HasKey("HOST") {
			model.FtpServerInfo.Host = cfg.Section("ftp").Key("HOST").String()
		}
		if cfg.Section("ftp").HasKey("PORT") {
			model.FtpServerInfo.Port = cfg.Section("ftp").Key("PORT").String()
		}
		if cfg.Section("ftp").HasKey("USERNAME") {
			model.FtpServerInfo.Username = cfg.Section("ftp").Key("USERNAME").String()
		}
		if cfg.Section("ftp").HasKey("PASSWORD") {
			model.FtpServerInfo.Password = cfg.Section("ftp").Key("PASSWORD").String()
		}
	}

	log.Println("FTP server: ", model.FtpServerInfo.Host)
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
	destination := model.TempClonedRepoRoot

	err := util.DeleteDirectory(destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	setupINI()

	// Logging to a file.
	model.FileRequests, _ = os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	model.FileConnections, _ = os.Create("connections.log")
	defer model.FileConnections.Close()
	defer model.FileRequests.Close()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
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
