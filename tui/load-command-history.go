package tui

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadCommandHistory reads the command history file and returns a slice of
// command lines, with the most recent commands at the end. It is best-effort
// and returns an empty slice if the file doesn't exist or can't be read.
func (m *Model) LoadCommandHistory() []string {
	if m.configDir == "" {
		return []string{}
	}

	historyPath := filepath.Join(m.configDir, "history")
	f, err := os.Open(historyPath)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	var history []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			history = append(history, line)
		}
	}

	return history
}
