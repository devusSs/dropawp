package storage

import (
	"errors"
	"os"
)

func CreateDir(dir string) error {
	if dir == "" {
		return errors.New("directory path cannot be empty")
	}

	// TODO: implement me
	// storageDir = dir

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return errors.New("failed to create directory: " + err.Error())
		}
	}

	return nil
}

// var storageDir = ""
