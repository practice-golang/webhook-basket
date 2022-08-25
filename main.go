package main // import "webhook-basket"

import (
	_ "embed"
	"flag"
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

type Content struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

//go:embed webhook-basket.ini
var sampleINI string

var listen string

func createINI(iniPath string) {
	if _, err := os.Stat(iniPath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", iniPath)
		os.Exit(1)
	}

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

func setupINI() {
	iniPath := "webhook-basket.ini"

	if len(os.Args) > 1 {
		FlagSetINI := flag.String("ini", "[filename.ini]", " Run server with the ini file")
		FlagGetINI := flag.Bool("getini", false, " Get sample ini file")

		flag.Usage = func() {
			flagSet := flag.CommandLine
			fmt.Printf("Usage of %s:\n", "webhook-basket")
			fmt.Printf("  %-19s Run server\n", "without option")

			order := []string{"ini", "getini"}
			for _, name := range order {
				flag := flagSet.Lookup(name)
				switch name {
				case "ini":
					fmt.Printf("  -%-18s%s\n", flag.Name+" "+flag.Value.String(), flag.Usage)
				case "getini":
					fmt.Printf("  -%-18s%s\n", flag.Name, flag.Usage)
				}
			}
		}

		flag.Parse()

		if *FlagGetINI {
			createINI(iniPath)
		}

		if *FlagSetINI != "" && *FlagSetINI != "[filename.ini]" {
			iniPath = *FlagSetINI

			if _, err := os.Stat(iniPath); os.IsNotExist(err) {
				fmt.Printf("File %s not found\n", iniPath)
				os.Exit(1)
			}
		}
	}

	cfg, err := ini.Load(iniPath)
	if err != nil {
		switch iniPath {
		case "webhook-basket.ini":
			createINI(iniPath)
		default:
			fmt.Println("Load INI: ", err)
			os.Exit(1)
		}
	}

	if cfg != nil {
		if cfg.Section("server").HasKey("LISTEN") {
			listen = cfg.Section("server").Key("LISTEN").String()
		}

		if cfg.Section("path").HasKey("CLONED_REPO_ROOT") {
			model.TempClonedRepoRoot = cfg.Section("path").Key("CLONED_REPO_ROOT").String()
		}
		if cfg.Section("path").HasKey("DEPLOYMENT_ROOT") {
			model.DeploymentRoot = cfg.Section("path").Key("DEPLOYMENT_ROOT").String()
		}

		if cfg.Section("git").HasKey("USERNAME") {
			model.AuthInfo.Username = cfg.Section("git").Key("USERNAME").String()
		}
		if cfg.Section("git").HasKey("PASSWORD") {
			model.AuthInfo.Password = cfg.Section("git").Key("PASSWORD").String()
		}

		if cfg.Section("ftp").HasKey("TYPE") {
			model.FtpServerInfo.Type = cfg.Section("ftp").Key("TYPE").String()
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
		if cfg.Section("ftp").HasKey("PASSIVE") {
			model.FtpServerInfo.Passive, err = cfg.Section("ftp").Key("PASSIVE").Bool()
			if err != nil {
				model.FtpServerInfo.Passive = true
			}
		}
	}
}

func HealthCheck(c *gin.Context) {
	result := gin.H{"status": "ok"}

	c.JSON(http.StatusOK, result)
}

func DeployRepository(c *gin.Context) {
	request := model.Request{}
	c.BindJSON(&request)

	request.DeployRoot = c.Query("deploy-root")
	request.DeployName = c.Query("deploy-name")

	if request.DeployRoot == "" {
		request.DeployRoot = model.DeploymentRoot
	}

	err := downloader.CloneAndUploadRepository(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DeleteReposRoot(c *gin.Context) {
	target := model.TempClonedRepoRoot

	err := util.DeleteDirectory(target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	listen = "127.0.0.1:7749"

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
	r.POST("/deploy", DeployRepository)
	r.DELETE("/repos-root", DeleteReposRoot)

	fmt.Println("Listen: ", listen)
	r.Run(listen)
}
