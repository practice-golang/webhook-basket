package util

import "os"

func DeleteDirectory(target string) error {
	err := os.RemoveAll(target)
	if err != nil {
		return err
	}

	return nil
}
