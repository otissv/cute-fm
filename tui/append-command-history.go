package tui

import (
	"os"
	"path/filepath"
)

// AppendCommandHistory appends a command line to the history file
// inside the config directory. It is best-effort and silently ignores
// any filesystem errors so as not to disrupt the TUI.
func (m *Model) AppendCommandHistory(line string) {
	if line == "" || m.configDir == "" {
		return
	}

	historyPath := filepath.Join(m.configDir, "history")

	f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(line + "\n")
}
