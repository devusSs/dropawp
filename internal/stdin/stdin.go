package stdin

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	return reader.ReadString('\n')
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
