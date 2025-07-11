package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

func (c *Config) validate() error {
	err := validateProjectName(c.ProjectName)
	if err != nil {
		return fmt.Errorf("invalid project_name: %w", err)
	}

	err = validateProjectTimes(c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("invalid project times: %w", err)
	}

	err = validateCooldownDuration(c.CooldownDuration)
	if err != nil {
		return fmt.Errorf("invalid cooldown_duration: %w", err)
	}

	err = validateSteamID64(c.SteamID64)
	if err != nil {
		return fmt.Errorf("invalid steam_id_64: %w", err)
	}

	err = validateAdditionalItemsFile(c.AdditionalItemsFile)
	if err != nil {
		return fmt.Errorf("invalid additional_items_file: %w", err)
	}

	return nil
}

const (
	minProjectNameLength = 4
	maxProjectNameLength = 16
)

func validateProjectName(name string) error {
	if name == "" {
		return errors.New("project_name cannot be empty")
	}

	if len(name) < minProjectNameLength {
		return fmt.Errorf("project_name must be at least %d characters long", minProjectNameLength)
	}

	if len(name) > maxProjectNameLength {
		return fmt.Errorf("project_name must be at most %d characters long", maxProjectNameLength)
	}

	for _, char := range name {
		if !unicode.IsLetter(char) {
			return fmt.Errorf("project_name can only contain letters, got '%c'", char)
		}

		if !unicode.IsLower(char) {
			return fmt.Errorf("project_name must be lowercase, got '%c'", char)
		}
	}

	return nil
}

func validateProjectTimes(createdAt time.Time, updatedAt time.Time) error {
	if createdAt.IsZero() {
		return errors.New("created_at cannot be zero")
	}

	if updatedAt.IsZero() {
		return errors.New("updated_at cannot be zero")
	}

	if updatedAt.Before(createdAt) {
		return errors.New("updated_at cannot be before created_at")
	}

	return nil
}

const minCooldownDuration = 10 * time.Minute

func validateCooldownDuration(duration time.Duration) error {
	if duration < minCooldownDuration {
		return fmt.Errorf("cooldown_duration must be at least %s", minCooldownDuration)
	}

	return nil
}

const steamID64Regex = `^7656119\d{10}$`

func validateSteamID64(steamID64 uint64) error {
	s := strconv.FormatUint(steamID64, 10)
	regex := regexp.MustCompile(steamID64Regex)

	if !regex.MatchString(s) {
		return fmt.Errorf("steam_id_64 must match regex %s, got '%s'", steamID64Regex, s)
	}

	return nil
}

func validateAdditionalItemsFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("additional_items_file does not exist: %s", filePath)
		}

		return fmt.Errorf("error checking additional_items_file: %w", err)
	}

	return nil
}
