package ftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"webhook-basket/uploader/config"

	"github.com/secsy/goftp"

	gi "github.com/sabhiram/go-gitignore"
)

// Upload file to ftp server
func UploadFile(fc *goftp.Client, localFile, remoteFile string) (err error) {
	// log.Printf("Uploading '%s' to '%s' ..", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("unable to open local file: %v", err)
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		path = strings.ReplaceAll(path, "\\", "/")
		fc.Mkdir(path)
	}

	err = fc.Store(remoteFile, srcFile)
	if err != nil {
		return fmt.Errorf("unable to upload local file: %v", err)
	}
	// log.Printf("%d bytes copied", bytes)

	return nil
}

func ProcUploadMain(host config.Host) (err error) {
	srcBase := config.ReplacerSlash.Replace(host.SrcBase)
	srcRoot := filepath.Base(srcBase)
	srcCutPath := config.ReplacerSlash.Replace(strings.TrimSuffix(srcBase, srcRoot))

	var wbIgnore *gi.GitIgnore
	wbIgnorePath := filepath.Join(srcBase, ".wbignore")
	if _, err = os.Stat(wbIgnorePath); err == nil {
		wbIgnore, err = gi.CompileIgnoreFile(wbIgnorePath)
		if err != nil {
			if !strings.Contains(err.Error(), "The system cannot find the file specified") {
				return
			}
		}
	}

	ftpConfig := goftp.Config{
		User:            host.Username,
		Password:        host.Password,
		ActiveTransfers: !host.Passive,
	}

	fc, err := goftp.DialConfig(ftpConfig, host.Hostname+":"+host.Port)
	if err != nil {
		return
	}
	defer fc.Close()
	// log.Println("Connected.")

	ques := []config.QueSheet{}

	// err = filepath.Walk(srcBase, util.WalkDIR)
	err = filepath.Walk(srcBase, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ques = append(ques, config.QueSheet{Name: path, IsDIR: info.IsDir()})
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}

	dstBase := filepath.Join(host.DstBase, host.DstName)

	for _, q := range ques {
		srcPath := filepath.Join("", q.Name)
		dstPath := ""
		if host.DstName != "" {
			dstPath = strings.TrimPrefix(q.Name, srcCutPath)
			dstPath = filepath.Join(host.DstName, strings.TrimPrefix(dstPath, host.SrcName))
			dstPath = filepath.Join(host.DstBase, dstPath)
		} else {
			dstPath = filepath.Join(host.DstBase, strings.TrimPrefix(q.Name, srcCutPath))
		}
		dstPath = strings.ReplaceAll(dstPath, "\\", "/")

		relPathForIgnoreCheck := strings.TrimPrefix(dstPath, dstBase)
		if wbIgnore != nil && wbIgnore.MatchesPath(relPathForIgnoreCheck) {
			continue
		}

		switch q.IsDIR {
		case true:
			_, err = fc.Mkdir(dstPath)
			if err != nil {
				if err.Error() == "unexpected response: 550-Directory already exists" {
					continue
				}
				if err.Error() == "unexpected response: 550-Directory with same name already exists" {
					continue
				}
				if err.Error() == "failed parsing directory name: Directory created successfully" {
					continue
				}

				log.Println("DIR: ", dstPath)
				log.Println("DIR: ", err)
			}
		case false:
			err = UploadFile(fc, srcPath, dstPath)
			if err != nil {
				log.Fatalf("could not upload file: %v", err)
			}
		}
	}

	return nil
}
