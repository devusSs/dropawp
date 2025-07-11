package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	ProjectName               string        `json:"project_name"`
	CreatedAt                 time.Time     `json:"created_at"`
	UpdatedAt                 time.Time     `json:"updated_at"`
	CooldownDuration          time.Duration `json:"cooldown_duration"`
	SkipSteamServicesCheck    bool          `json:"skip_steam_services_check"`
	SteamID64                 uint64        `json:"steam_id_64"`
	SkipSteamUserCheck        bool          `json:"skip_steam_user_check"`
	SkipFilterUntradableItems bool          `json:"skip_filter_untradable_items"`
	AdditionalItemsFile       string        `json:"additional_items_file"`

	filePath string
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

func Write(c *Config) error {
	if c == nil {
		return errors.New("config cannot be nil")
	}

	err := c.validate()
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	var f *os.File
	f, err = createConfigFile()
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

func Read() (*Config, error) {
	f, err := openConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	c := &Config{}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	var dir string
	dir, err = setupConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config directory: %w", err)
	}

	c.filePath = filepath.Join(dir, "config.json")

	err = c.validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return c, nil
}

func Delete(c *Config) error {
	if c == nil {
		return errors.New("config cannot be nil")
	}

	if c.filePath == "" {
		return errors.New("config file path is not set")
	}

	err := os.Remove(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	c.filePath = ""

	return nil
}

//nolint:funlen // This function has to be that long unfortunately, I am too lazy to refactor it.
func FromInput() (*Config, error) {
	projectName, err := getInput("Enter a project name (empty for random uuid)")
	if err != nil {
		return nil, fmt.Errorf("failed to get project name: %w", err)
	}

	if projectName == "" {
		projectName = uuid.New().String()
	}

	var cooldownDurationStr string
	cooldownDurationStr, err = getInput("Enter cooldown duration (e.g., 5m, 1h, 1d, 1w etc.)")
	if err != nil {
		return nil, fmt.Errorf("failed to get cooldown duration: %w", err)
	}

	if cooldownDurationStr == "" {
		cooldownDurationStr = minCooldownDuration.String()
	}

	var cooldownDuration time.Duration
	cooldownDuration, err = parseExtendedDuration(cooldownDurationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid cooldown duration: %w", err)
	}

	var skipSteamServicesCheckStr string
	skipSteamServicesCheckStr, err = getInput("Skip Steam services check? (true/false)")
	if err != nil {
		return nil, fmt.Errorf("failed to get skip Steam services check: %w", err)
	}

	if skipSteamServicesCheckStr == "" {
		//nolint:goconst // I will not create a constant for false lol.
		skipSteamServicesCheckStr = "false"
	}

	var skipSteamServicesCheck bool
	skipSteamServicesCheck, err = parseBool(skipSteamServicesCheckStr)
	if err != nil {
		return nil, fmt.Errorf("invalid skip Steam services check value: %w", err)
	}

	var steamID64Str string
	steamID64Str, err = getInput("Enter Steam ID64")
	if err != nil {
		return nil, fmt.Errorf("failed to get Steam ID64: %w", err)
	}

	if steamID64Str == "" {
		return nil, errors.New("steam id64 cannot be empty")
	}

	var steamID64 uint64
	steamID64, err = parseUint64(steamID64Str)
	if err != nil {
		return nil, fmt.Errorf("invalid Steam ID64: %w", err)
	}

	var skipSteamUserCheckStr string
	skipSteamUserCheckStr, err = getInput("Skip Steam user check? (true/false)")
	if err != nil {
		return nil, fmt.Errorf("failed to get skip Steam user check: %w", err)
	}

	if skipSteamUserCheckStr == "" {
		skipSteamUserCheckStr = "false"
	}

	var skipSteamUserCheck bool
	skipSteamUserCheck, err = parseBool(skipSteamUserCheckStr)
	if err != nil {
		return nil, fmt.Errorf("invalid skip Steam user check value: %w", err)
	}

	var skipFilterUntradableItemsStr string
	skipFilterUntradableItemsStr, err = getInput("Skip filter untradable items? (true/false)")
	if err != nil {
		return nil, fmt.Errorf("failed to get skip filter untradable items: %w", err)
	}

	if skipFilterUntradableItemsStr == "" {
		skipFilterUntradableItemsStr = "false"
	}

	var skipFilterUntradableItems bool
	skipFilterUntradableItems, err = parseBool(skipFilterUntradableItemsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid skip filter untradable items value: %w", err)
	}

	var additionalItemsFile string
	additionalItemsFile, err = getInput("Enter additional items file path (optional)")
	if err != nil {
		return nil, fmt.Errorf("failed to get additional items file path: %w", err)
	}

	c := &Config{
		ProjectName:               projectName,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
		CooldownDuration:          cooldownDuration,
		SkipSteamServicesCheck:    skipSteamServicesCheck,
		SteamID64:                 steamID64,
		SkipSteamUserCheck:        skipSteamUserCheck,
		SkipFilterUntradableItems: skipFilterUntradableItems,
		AdditionalItemsFile:       additionalItemsFile,
	}

	var dir string
	dir, err = setupConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config directory: %w", err)
	}

	c.filePath = filepath.Join(dir, "config.json")

	err = c.validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return c, nil
}

func SetFile(f string) {
	file = f
}

func FromFile() (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	c := &Config{}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
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

var file string
