package downloader

import (
	"log"
	"os"
	"webhook-basket/model"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepository(request model.Request) error {
	var err error

	repoName := request.Repository.Name
	cloneURI := request.Repository.CloneURL

	repoPath := model.CloneRepoRoot + "/" + repoName

	auth := &http.BasicAuth{
		Username: "edp1096",
		Password: "01026473448",
	}

	// Clone
	if _, err = os.Stat(repoPath + "/.git"); os.IsNotExist(err) {
		// clone new
		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{URL: cloneURI, Progress: os.Stdout, Auth: auth})
		if err != nil {
			log.Println("Git clone error: ", err)
		}
	} else {
		files, err := os.ReadDir(repoPath)
		if err != nil {
			log.Println("dir: ", err)
		}

		for _, f := range files {
			log.Println(f.Name(), f.IsDir())
		}

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

	return err
}
