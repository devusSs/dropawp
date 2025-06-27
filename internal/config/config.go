package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	LogKeepInterval     time.Duration `json:"log_keep_interval"`
	StorageKeepInterval time.Duration `json:"storage_keep_interval"`
	WantSteamCheck      bool          `json:"want_steam_check"`
	Username            string        `json:"username"`
	SteamID64           string        `json:"steam_id_64"`
	FilterUntradable    bool          `json:"filter_untradable"`
	AdditionalItemsFile string        `json:"additional_items_file"`
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

func (c *Config) Write(filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	err := os.MkdirAll(filepath.Dir(filePath), 0750)
	if err != nil {
		return fmt.Errorf("failed to create directory for config file: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	return nil
}

func Load(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("file path cannot be empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	c := &Config{}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return c, nil
}

const (
	MinLogKeepInterval     = 1 * time.Hour
	MinStorageKeepInterval = 24 * time.Hour
)

func NewConfig(
	logKeepInterval string,
	storageKeepInterval string,
	wantSteamCheck string,
	username string,
	steamID64 string,
	filterUntradable string,
	additionalItemsFile string,
) (*Config, error) {
	logKeepDur, err := parseExtendedDuration(logKeepInterval)
	if err != nil {
		return nil, err
	}

	if logKeepDur < MinLogKeepInterval {
		return nil, fmt.Errorf("log keep interval must be at least %s", MinLogKeepInterval)
	}

	storageKeepDur, err := parseExtendedDuration(storageKeepInterval)
	if err != nil {
		return nil, err
	}

	if storageKeepDur < MinStorageKeepInterval {
		return nil, fmt.Errorf("storage keep interval must be at least %s", MinStorageKeepInterval)
	}

	wantSteamCheckBool := true
	if wantSteamCheck != "" {
		wantSteamCheckBool, err = strconv.ParseBool(wantSteamCheck)
		if err != nil {
			return nil, err
		}
	}

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if steamID64 == "" {
		return nil, errors.New("steam_id_64 cannot be empty")
	}

	filterUntradableBool, err := strconv.ParseBool(filterUntradable)
	if err != nil {
		return nil, fmt.Errorf("invalid filter_untradable value: %w", err)
	}

	_, err = os.Stat(additionalItemsFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to check additional items file: %w", err)
	}

	return &Config{
		LogKeepInterval:     logKeepDur,
		StorageKeepInterval: storageKeepDur,
		WantSteamCheck:      wantSteamCheckBool,
		Username:            username,
		SteamID64:           steamID64,
		FilterUntradable:    filterUntradableBool,
		AdditionalItemsFile: additionalItemsFile,
	}, nil
}

const (
	dayFormatter  = 24 * time.Hour
	weekFormatter = 7 * dayFormatter
)

func parseExtendedDuration(s string) (time.Duration, error) {
	unitMap := map[string]time.Duration{
		"ns": time.Nanosecond,
		"us": time.Microsecond,
		"µs": time.Microsecond,
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  dayFormatter,
		"w":  weekFormatter,
	}

	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	for unit, dur := range unitMap {
		if strings.HasSuffix(s, unit) {
			numStr := strings.TrimSuffix(s, unit)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return time.Duration(float64(dur) * num), nil
		}
	}

	return 0, errors.New("unknown duration format")
}
