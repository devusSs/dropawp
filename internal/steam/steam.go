package steam

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

var ErrContextNil = errors.New("context cannot be nil")

func applyHeaders(req *http.Request) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.113 Safari/537.36",
	)

	return nil
}

const steamID64Regex = `^7656119\d{10}$`

func validateSteamID64(steamID64 uint64) error {
	re := regexp.MustCompile(steamID64Regex)
	s := strconv.FormatUint(steamID64, 10)

	if !re.MatchString(s) {
		return errors.New("invalid steamID64 format")
	}

	return nil
}

const steamAPIKeyLength = 32

func validateSteamAPIKey(apiKey string) error {
	if apiKey == "" {
		return errors.New("api key cannot be empty")
	}

	if len(apiKey) != steamAPIKeyLength {
		return fmt.Errorf("api key must be exactly %d characters long", steamAPIKeyLength)
	}

	_, err := hex.DecodeString(apiKey)
	if err != nil {
		return errors.New("api key must be a valid hexadecimal string")
	}

	return nil
}
