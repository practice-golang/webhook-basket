package util

import "os"

func DeleteDirectory(destination string) error {
	err := os.RemoveAll(destination)
	if err != nil {
		return err
	}

	return nil
}
