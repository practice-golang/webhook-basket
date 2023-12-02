package util

import (
	"os"

	"golang.org/x/crypto/ssh"
)

func DeleteDirectory(target string) error {
	err := os.RemoveAll(target)
	if err != nil {
		return err
	}

	return nil
}

// https://hwdhyeon.github.io/golang/how-to-use-ssh-in-golang
func ReadSshPemKey(file string) (authMethod ssh.AuthMethod, err error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
