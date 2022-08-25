package sftp

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"webhook-basket/uploader/config"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Upload file to sftp server
func UploadFile(sc *sftp.Client, localFile, remoteFile string) (err error) {
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
		sc.Mkdir(path)
	}

	// Note: SFTP Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		return fmt.Errorf("unable to open remote file: %v", err)
	}
	defer dstFile.Close()

	// bytes, err := io.Copy(dstFile, srcFile)
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("unable to upload local file: %v", err)
	}
	// log.Printf("%d bytes copied", bytes)

	return nil
}

func ProcMain(host config.Host) {
	srcBase := config.ReplacerSlash.Replace(host.SrcBase)
	srcRoot := filepath.Base(srcBase)
	srcCutPath := config.ReplacerSlash.Replace(strings.TrimSuffix(srcBase, srcRoot))

	var sshConfig = &ssh.ClientConfig{
		User:            host.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(host.Password)},
	}

	client, err := ssh.Dial("tcp", host.Hostname+":"+host.Port, sshConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	// log.Println("Connected.")

	// open an SFTP session over an existing ssh connection.
	sc, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

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
			err = sc.MkdirAll(dstPath)
			if err != nil {
				panic(err)
			}
		case false:
			err = UploadFile(sc, srcPath, dstPath)
			if err != nil {
				log.Fatalf("could not upload file: %v", err)
			}
		}
	}
}
