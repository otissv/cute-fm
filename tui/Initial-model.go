package tui

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/viewport"

	"cute/config"
	"cute/filesystem"
)

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize left viewport for the second row
	leftViewport := viewport.New()

	// Initialize right viewport for the second row
	rightViewport := viewport.New()
	rightViewport.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	// Determine initial directory for the file list.
	wd := startDir
	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			wd = "."
		}
	}

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
	cfgDir := filepath.Join(userConfigDir, "cute")
	// Best-effort creation; ignore error so the TUI can still start.
	_ = os.MkdirAll(cfgDir, 0o755)

	// Load Lua-based runtime configuration (theme + commands).
	runtimeCfg := config.LoadRuntimeConfig(cfgDir)

	// Load the initial directory into the left viewport using the configured theme.
	_, selected := loadDirectoryIntoView(&leftViewport, runtimeCfg.Theme, wd)

	f := []filesystem.FileInfo{}

	m := Model{
		configDir:     cfgDir,
		runtimeConfig: runtimeCfg,

		leftViewport:   leftViewport,
		rightViewport:  rightViewport,
		allFiles:       f,
		files:          f,
		currentDir:     wd,
		selectedIndex:  selected,
		theme:          runtimeCfg.Theme,
		viewportHeight: 0,
		viewportWidth:  0,
		layoutRows:     []string{""},
		layout:         "",
		titleText:      "The Cute File Manager",
	}

	// Initialize the search input
	m.searchInput = m.SearchInput(">", "Filter...")

	// Initialize the command input
	m.commandInput = m.CommandInput("", "Enter prompt...")

	// Load command history for auto-complete
	m.commandHistory = m.LoadCommandHistory()
	m.historyMatches = []string{}
	m.historyIndex = -1

	m.CalcLayout()

	ActiveTuiMode = TuiModeNormal
	ActiveFileListMode = FileListModeDir

	return m
}
