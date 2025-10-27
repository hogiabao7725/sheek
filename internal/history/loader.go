package history

import (
	"bufio"
	"os"
	"path/filepath"
)

func LoadZshHistory() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Build full path to .zsh_history
	historyFile := filepath.Join(homeDir, ".zsh_history")

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
