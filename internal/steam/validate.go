package steam

import (
	"encoding/hex"
	"fmt"
)

const steamAPIKeyLength = 32

func validateSteamAPIKey(apiKey string) error {
	if len(apiKey) != steamAPIKeyLength {
		return fmt.Errorf(
			"invalid Steam API key length: expected %d characters, got %d",
			steamAPIKeyLength,
			len(apiKey),
		)
	}

	_, err := hex.DecodeString(apiKey)
	if err != nil {
		return fmt.Errorf("invalid Steam API key format: %w", err)
	}

	return nil
}
