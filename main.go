package main // import "github.com/practice-golang/webhook-basket"

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/practice-golang/webhook-basket/config"
	"github.com/practice-golang/webhook-basket/copier"
)

var (
	file *os.File

	//go:embed ini/webhook-basket.ini
	sampleINI string

	cfg *ini.File
)

func serverSetup() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(
		middleware.CORS(),
		middleware.Recover(),
	)

	e.POST("/webhook", copier.CopyRepository)

	return e
}

func dumpHandler(c echo.Context, reqBody, resBody []byte) {
	header := time.Now().Format("2006-01-02 15:04:05") + " - "
	body := string(reqBody)
	body = strings.Replace(body, "\r\n", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	data := header + body + "\n"

	f, err := os.OpenFile(
		"request-body.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Println(err)
		return
	}
}

func readINI() (err error) {
	// cfg, err = ini.Load("webhook-basket.ini")
	cfg, err = ini.LoadSources(ini.LoadOptions{SpaceBeforeInlineComment: true}, "webhook-basket.ini")
	if err != nil {
		log.Print("Fail to read ini. ")

		f, err := os.Create("webhook-basket.ini")
		if err != nil {
			log.Println("Create INI: ", err)
			return err
		}
		defer f.Close()

		_, err = f.WriteString(sampleINI + "\n")
		if err != nil {
			log.Println("Create INI: ", err)
			return err
		}

		log.Println("webhook-basket.ini is created")
	}

	if cfg != nil {
		config.Repos = map[string]config.Repository{}
		repoSections := cfg.SectionStrings()

		for _, section := range repoSections {
			switch section {
			case "DEFAULT":
				config.ServerInfo = config.ServerPreference{
					AppName:     cfg.Section(section).Key("APP_NAME").Value(),
					PrepareRoot: cfg.Section(section).Key("PREPARE_ROOT").Value(),
					ListenAddr:  cfg.Section(section).Key("LISTEN_ADDR").Value(),
					ListenPort:  cfg.Section(section).Key("LISTEN_PORT").Value(),
				}
			case "do_not_use":
				continue
			default:
				repo := config.Repository{
					Name:       section,
					CloneURI:   cfg.Section(section).Key("CLONE_URI").Value(),
					SSHos:      cfg.Section(section).Key("SSH_OS").Value(),
					SSHuri:     cfg.Section(section).Key("SSH_ADDRESS").Value(),
					SSHid:      cfg.Section(section).Key("SSH_ID").Value(),
					SSHpasswd:  cfg.Section(section).Key("SSH_PASSWD").Value(),
					DeployRoot: cfg.Section(section).Key("DEPLOY_ROOT").Value(),
					RepoRoot:   cfg.Section(section).Key("DESTINATION_PATH").Value(),
				}

				config.Repos[repo.Name] = repo
			}
		}
	}

	return
}

func main() {
	var err error

	err = readINI()
	if err != nil {
		log.Fatalln(err)
		log.Fatalln(err)
	}

	e := serverSetup()

	file, err = os.OpenFile(
		"connection.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} - remote_ip:${remote_ip}, host:${host}, ` +
			`method:${method}, uri:${uri},status:${status}, error:${error}, ` +
			`${header:Authorization}, query:${query:property}, form:${form}, ` + "\n",
		Output: file,
	}))

	e.Use(middleware.BodyDump(dumpHandler))

	e.Logger.Fatal(e.Start(config.ServerInfo.ListenAddr + ":" + config.ServerInfo.ListenPort))
}
