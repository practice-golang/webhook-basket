package ftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"webhook-basket/uploader/config"

	"github.com/secsy/goftp"
)

// Upload file to ftp server
func UploadFile(fc *goftp.Client, localFile, remoteFile string) (err error) {
	// log.Printf("Uploading '%s' to '%s' ..", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("Unable to open local file: %v", err)
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
		return fmt.Errorf("Unable to upload local file: %v", err)
	}
	// log.Printf("%d bytes copied", bytes)

	return nil
}

func ProcMain(host config.Host) {
	ftpConfig := goftp.Config{
		User:            host.Username,
		Password:        host.Password,
		ActiveTransfers: !host.Passive,
	}

	fc, err := goftp.DialConfig(ftpConfig, host.Hostname+":"+host.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer fc.Close()
	// log.Println("Connected.")

	srcBase := config.ReplacerSlash.Replace(host.SrcBase)
	srcRoot := filepath.Base(srcBase)
	srcCutPath := config.ReplacerSlash.Replace(strings.TrimSuffix(srcBase, srcRoot))

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

		switch q.IsDIR {
		case true:
			_, err = fc.Mkdir(dstPath)
			if err != nil {
				if err.Error() == "unexpected response: 550-Directory already exists" {
					continue
				}
				if err.Error() == "failed parsing directory name: Directory created successfully" {
					continue
				}

				log.Println("DIR: ", err)
			}
		case false:
			err = UploadFile(fc, srcPath, dstPath)
			if err != nil {
				log.Fatalf("could not upload file: %v", err)
			}
		}
	}
}
