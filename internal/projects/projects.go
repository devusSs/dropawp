package projects

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Project struct {
	Name       string `json:"name"`
	LogsDir    string `json:"logs_dir"`
	ConfigFile string `json:"config_file"`
	StorageDir string `json:"storage_dir"`
}

func (p *Project) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *Project) JSON() string {
	m, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling project info: %v", err)
	}

	return string(m)
}

func Load(dir string) ([]*Project, error) {
	setup(dir)

	file, err := os.Open(projectFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open projects file: %w", err)
	}
	defer file.Close()

	projects := make([]*Project, 0)
	err = json.NewDecoder(file).Decode(&projects)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return projects, nil
		}

		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

func SaveProject(dir string, p *Project) error {
	setup(dir)

	if p == nil {
		return errors.New("project cannot be nil")
	}

	file, err := os.OpenFile(projectFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open projects file: %w", err)
	}
	defer file.Close()

	projects, err := Load(dir)
	if err != nil {
		return fmt.Errorf("failed to load existing projects: %w", err)
	}

	for _, project := range projects {
		if project.Name == p.Name {
			return fmt.Errorf("project with name '%s' already exists", p.Name)
		}
	}

	projects = append(projects, p)

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate projects file: %w", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to start of projects file: %w", err)
	}

	err = json.NewEncoder(file).Encode(projects)
	if err != nil {
		return fmt.Errorf("failed to encode projects: %w", err)
	}

	return nil
}

var (
	once        = &sync.Once{}
	projectFile = ""
)

func setup(dir string) {
	once.Do(func() {
		if dir == "" {
			fmt.Fprintf(os.Stderr, "Error: Projects directory cannot be empty\n")
			os.Exit(1)
		}

		projectFile = filepath.Join(dir, ".projects")

		err := os.MkdirAll(filepath.Dir(projectFile), 0750)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating projects directory: %v\n", err)
			os.Exit(1)
		}

		_, err = os.Stat(projectFile)
		if os.IsNotExist(err) {
			var file *os.File
			file, err = os.Create(projectFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating projects file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
		}
	})
}
