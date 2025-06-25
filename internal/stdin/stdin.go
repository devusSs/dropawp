package stdin

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ReadFromStdin(prompt string) (string, error) {
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	if !isTerminal(os.Stdin) {
		return "", errors.New("stdin is not a terminal")
	}

	if prompt[len(prompt)-1] != ':' {
		prompt += ": "
	}

	if prompt[len(prompt)-1] != ' ' {
		prompt += " "
	}

	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}

	info, err := f.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
