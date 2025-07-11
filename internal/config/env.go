package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func SetEnvFile(file string) {
	envFile = file
}

func FromEnv() (*Config, error) {
	err := loadEnvFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment file: %w", err)
	}

	c := &Config{}
	err = parseEnv(c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	var dir string
	dir, err = setupConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config directory: %w", err)
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.filePath = filepath.Join(dir, "config.json")

	err = c.validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return c, nil
}

var envFile string

func loadEnvFile() error {
	if envFile == "" {
		return nil
	}

	err := godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("failed to load environment file %s: %w", envFile, err)
	}

	return nil
}

func parseEnv(v any) error {
	if v == nil {
		return errors.New("nil value provided for parsing")
	}

	opts := env.Options{Prefix: "DROPAWP_", RequiredIfNoDef: true}
	err := env.ParseWithOptions(v, opts)
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return nil
}
