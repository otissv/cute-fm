package tui

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

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
