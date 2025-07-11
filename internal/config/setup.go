package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func createConfigFile() (*os.File, error) {
	configDir, err := setupConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config directory: %w", err)
	}

	configFilePath := filepath.Join(configDir, "config.json")

	file, err := os.Create(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create config file %s: %w", configFilePath, err)
	}

	return file, nil
}

func openConfigFile() (*os.File, error) {
	configDir, err := setupConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config directory: %w", err)
	}

	configFilePath := filepath.Join(configDir, "config.json")

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configFilePath, err)
	}

	return file, nil
}

func setupConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".dropawp", "config")

	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	return configDir, nil
}
