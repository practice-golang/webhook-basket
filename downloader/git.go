package downloader

import (
	"errors"
	"log"
	"os"
	"webhook-basket/model"
	"webhook-basket/uploader/config"
	"webhook-basket/uploader/ftp"
	"webhook-basket/uploader/sftp"

	"github.com/go-git/go-git/v5"
)

func CloneAndUploadRepository(request model.Request) error {
	var err error

	if request.Destination == "" {
		return errors.New("Deployment destination is required")
	}

	request.Ftp = model.FtpServerInfo

	repoName := request.Repository.Name
	cloneURI := request.Repository.CloneURL

	repoPath := model.TempClonedRepoRoot + "/" + repoName

	auth := model.AuthInfo

	// Clone
	if _, err = os.Stat(repoPath + "/.git"); os.IsNotExist(err) {
		// clone new
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{URL: cloneURI, Progress: os.Stdout, Auth: auth})
		if err != nil {
			log.Println("Git clone error: ", err)
		}
	} else {
		// pull existing
		r, err := git.PlainOpen(repoPath)
		if err != nil {
			log.Println("Git pull error: ", err)
		}

		w, err := r.Worktree()
		if err != nil {
			log.Println("Git pull error: ", err)
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth})
		if err != nil {
			log.Println("Git pull error: ", err)
		}
	}

	host := config.Host{
		Type:     request.Ftp.Type,
		Hostname: request.Ftp.Host,
		Port:     request.Ftp.Port,
		Username: request.Ftp.Username,
		Password: request.Ftp.Password,
		SrcBase:  repoPath,
		DstBase:  request.Destination,
	}

	switch host.Type {
	case "ftp":
		ftp.ProcMain(host)
	case "sftp":
		sftp.ProcMain(host)
	}

	return err
}
