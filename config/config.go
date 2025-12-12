package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	// Resolve and ensure the configuration directory exists.
	userConfigDir, err := os.UserConfigDir()
	if err != nil || userConfigDir == "" {
		// Fallback to $HOME/.config if UserConfigDir is unavailable.
		homeDir, herr := os.UserHomeDir()
		if herr != nil || homeDir == "" {
			userConfigDir = "."
		} else {
			userConfigDir = filepath.Join(homeDir, ".config")
		}
	}
	configDir := filepath.Join(userConfigDir, "cute")
	// Best-effort creation; ignore error so the TUI can still start.
	_ = os.MkdirAll(configDir, 0o755)

	return configDir
}
