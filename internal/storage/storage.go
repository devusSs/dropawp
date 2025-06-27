package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func CreateDir(dir string) error {
	if dir == "" {
		return errors.New("directory path cannot be empty")
	}

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return errors.New("failed to create directory: " + err.Error())
		}
	}

	return nil
}

func CleanupOld(dir string, deadline time.Duration) error {
	if dir == "" {
		return errors.New("directory path cannot be empty")
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.New("failed to read directory: " + err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			return fmt.Errorf("directory %s contains subdirectories, expected only files", dir)
		}

		filePath := filepath.Join(dir, file.Name())

		var info os.FileInfo
		info, err = os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("failed to get info for file %s: %w", filePath, err)
		}

		if time.Since(info.ModTime()) > deadline {
			err = os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to remove old storage file %s: %w", filePath, err)
			}
		}
	}

	return nil
}
