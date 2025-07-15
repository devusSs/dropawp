package lastrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var ErrLastRunNotExist = errors.New("last run file does not exist")

type LastRun struct {
	ProjectName string    `json:"project_name"`
	LastRun     time.Time `json:"last_run"`
}

func (l *LastRun) String() string {
	return fmt.Sprintf("%+v", *l)
}

func (l *LastRun) validate() error {
	if l.ProjectName == "" {
		return errors.New("project name cannot be empty")
	}

	if l.LastRun.IsZero() {
		return errors.New("last run time cannot be zero")
	}

	return nil
}

func Write(projectName string) error {
	lr := &LastRun{
		ProjectName: projectName,
		LastRun:     time.Now(),
	}

	err := lr.validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	var file *os.File
	file, err = createLastRunFile()
	if err != nil {
		return fmt.Errorf("failed to create last run file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(lr)
	if err != nil {
		return fmt.Errorf("failed to write last run data: %w", err)
	}

	return nil
}

func Read() (*LastRun, error) {
	file, err := openLastRunFile()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrLastRunNotExist
		}
		return nil, fmt.Errorf("failed to read last run file: %w", err)
	}
	defer file.Close()

	var lr LastRun
	err = json.NewDecoder(file).Decode(&lr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode last run data: %w", err)
	}

	err = lr.validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return &lr, nil
}

func createLastRunFile() (*os.File, error) {
	lastRunDir, err := setupLastRunDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup last run directory: %w", err)
	}

	lastRunFilePath := filepath.Join(lastRunDir, "lastrun.json")

	file, err := os.Create(lastRunFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create last run file %s: %w", lastRunFilePath, err)
	}

	return file, nil
}

func openLastRunFile() (*os.File, error) {
	lastRunDir, err := setupLastRunDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup last run directory: %w", err)
	}

	lastRunFilePath := filepath.Join(lastRunDir, "lastrun.json")

	file, err := os.Open(lastRunFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open last run file %s: %w", lastRunFilePath, err)
	}

	return file, nil
}

func setupLastRunDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	lastRunDir := filepath.Join(home, ".dropawp", "lastrun")

	err = os.MkdirAll(lastRunDir, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create last run directory %s: %w", lastRunDir, err)
	}

	return lastRunDir, nil
}
