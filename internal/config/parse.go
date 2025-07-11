package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

func getInput(prompt string) (string, error) {
	if !isTerminal() {
		return "", errors.New("stdin is not a readable terminal")
	}

	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	if prompt[len(prompt)-1] != ':' {
		prompt += ":"
	}

	if prompt[len(prompt)-1] != ' ' {
		prompt += " "
	}

	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	return strings.TrimSpace(text), nil
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func parseBool(s string) (bool, error) {
	if s == "" {
		return false, errors.New("missing bool value")
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("invalid bool value: %s", s)
	}

	return b, nil
}

func parseUint64(s string) (uint64, error) {
	if s == "" {
		return 0, errors.New("missing uint64 value")
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint64 value: %s", s)
	}

	return i, nil
}

const (
	timeUnitDay  = 24 * time.Hour
	timeUnitWeek = 7 * timeUnitDay
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
		"d":  timeUnitDay,
		"w":  timeUnitWeek,
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
