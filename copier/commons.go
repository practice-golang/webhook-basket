package copier

import (
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/labstack/echo/v4"
	"github.com/practice-golang/webhook-basket/commander"
	"github.com/practice-golang/webhook-basket/config"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v5"
)

func doCopyRepository(bodyBytes []byte) {
	// cref := gjson.Get(string(bodyBytes), "ref")
	// committerName := gjson.Get(string(bodyBytes), "commits.#.committer.name")
	// committerEmail := gjson.Get(string(bodyBytes), "commits.#.committer.email")
	repoName := gjson.Get(string(bodyBytes), "repository.name")
	cloneURI := config.Repos[repoName.String()].CloneURI

	// log.Println(cref, committerName, committerEmail, repoName, cloneURI)

	sshOS := config.Repos[repoName.String()].SSHos
	sshURI := config.Repos[repoName.String()].SSHuri
	sshID := config.Repos[repoName.String()].SSHid
	sshPASSWD := config.Repos[repoName.String()].SSHpasswd
	deployPath := config.Repos[repoName.String()].DeployRoot
	prepareRoot := config.ServerInfo.PrepareRoot

	// Clone
	if _, err := os.Stat(prepareRoot + "/" + repoName.String() + "/.git"); os.IsNotExist(err) {
		// clone new
		_, err := git.PlainClone(prepareRoot+"/"+repoName.String(), false, &git.CloneOptions{
			URL:      cloneURI,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Println("Git clone ERR: ", err)
		}
	} else {
		// pull existing
		r, err := git.PlainOpen(prepareRoot + "/" + repoName.String())
		if err != nil {
			log.Println("Git pull ERR: ", err)
		}

		w, err := r.Worktree()
		if err != nil {
			log.Println("Git pull ERR: ", err)
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			log.Println("Git pull ERR: ", err)
		}
	}

	// SCP
	basePaths := ""
	filepath.WalkDir(prepareRoot+"/"+repoName.String(), func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			f, err := os.Open("./" + s)
			if err != nil {
				log.Println("SCP open file ERR: ", err)
			}

			proot := strings.ReplaceAll(prepareRoot, "./", "")
			dpath := strings.ReplaceAll("./"+s, proot+"/", "")
			dpath = strings.ReplaceAll(dpath, proot+"\\", "")
			dpath = strings.ReplaceAll(dpath, "./", "")
			dpath = strings.ReplaceAll(deployPath+"/"+dpath, "/", "\\")

			println("workdir: ", s+" -> "+dpath)

			clientConfig, _ := auth.PasswordKey(sshID, sshPASSWD, ssh.InsecureIgnoreHostKey())
			client := scp.NewClient(sshURI, &clientConfig)
			defer client.Close()

			err = client.Connect()
			if err != nil {
				log.Println("SCP deploy ERR: ", err)
			}

			err = client.CopyFile(f, dpath, "0655")
			if err != nil {
				log.Println("SCP copy file ERR: ", err)
			}
		} else {
			// basePaths += "/" + d.Name()
			basePaths = strings.ReplaceAll(s, prepareRoot+"/", "")
			basePaths = strings.ReplaceAll(basePaths, prepareRoot+"\\", "")
			fullPath := deployPath + "/" + basePaths

			log.Println("Full path: ", fullPath+" ---- "+s)
			log.Println("id/pwd ", sshID, sshPASSWD)

			conn, err := commander.Connect(sshURI, sshID, sshPASSWD)
			if err != nil {
				log.Println("SSH mkdir connect ERR: ", err)
			}
			defer conn.Close()

			dpath := ""
			switch sshOS {
			case "WINDOWS":
				dpath = strings.ReplaceAll(fullPath, "/", "\\")
				_, err = conn.SendCommands("md " + dpath)
				if err != nil {
					log.Println("SSH command ERR: ", err)
				}
			case "LINUX":
				dpath = strings.ReplaceAll(fullPath, "\\", "/")
				_, err = conn.SendCommands("mkdir " + dpath)
				if err != nil {
					log.Println("SSH command ERR: ", err)
				}
			}
			log.Println("Make dir: ", dpath)
		}

		return nil
	})
}

func CopyRepository(c echo.Context) error {
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}

	go doCopyRepository(bodyBytes)

	r := make(map[string]interface{})
	r["status"] = "200"
	r["message"] = "OK"

	return c.JSON(http.StatusOK, r)
}
