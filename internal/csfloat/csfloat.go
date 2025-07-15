package csfloat

import (
	"errors"
	"net/http"
)

var ErrContextNil = errors.New("context cannot be nil")

func applyHeaders(req *http.Request, apiKey string) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if apiKey == "" {
		return errors.New("api key cannot be empty")
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return nil
}
