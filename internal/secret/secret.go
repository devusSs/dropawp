package secret

import (
	"errors"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

func GetInput(prompt string) (string, error) {
	if !isTerminal() {
		return "", errors.New("stdin is not a readable terminal")
	}

	if len(prompt) > 0 && prompt[len(prompt)-1] != ':' {
		prompt += ":"
		prompt += " "
	}

	fmt.Print(prompt)

	byteInput, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return string(byteInput), nil
}

type Key string

//nolint:gosec // These are not sensitive keys, just identifiers.
const (
	SteamAPIKey   Key = "steam_api_key"
	CSFloatAPIKey Key = "csfloat_api_key"
)

func Exists(key Key) (bool, error) {
	_, err := keyring.Get("dropawp", string(key))
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("failed to check if key exists: %w", err)
	}

	return true, nil
}

func Load(key Key) (string, error) {
	value, err := keyring.Get("dropawp", string(key))
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", fmt.Errorf("key not found: %s", key)
		}

		return "", fmt.Errorf("failed to load key: %w", err)
	}

	return value, nil
}

func Save(key Key, value string) error {
	if value == "" {
		return errors.New("value cannot be empty")
	}

	err := keyring.Set("dropawp", string(key), value)
	if err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	return nil
}

func DeleteAll() error {
	err := keyring.DeleteAll("dropawp")
	if err != nil {
		return fmt.Errorf("failed to delete all keys: %w", err)
	}

	return nil
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
