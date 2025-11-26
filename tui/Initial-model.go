package tui

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/viewport"

	"cute/config"
	"cute/filesystem"
	"cute/theming"
)

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize left viewport for the second row
	leftVp := viewport.New()

	// Initialize right viewport for the second row
	rightVp := viewport.New()
	rightVp.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	// Load theme configuration.
	theme := theming.LoadTheme("cute.toml")

	// Load user-defined commands from the configuration file.
	cfgCommands := config.LoadCommands("cute.toml")

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

	// files, selected := loadDirectoryIntoView(&leftVp, theme, wd)
	_, selected := loadDirectoryIntoView(&leftVp, theme, wd)

	f := []filesystem.FileInfo{}

	m := Model{
		configDir:        cfgDir,
		fileListViewport: leftVp,
		previewViewport:  rightVp,
		allFiles:         f,
		files:            f,
		currentDir:       wd,
		selectedIndex:    selected,
		theme:            theme,
		isCommandBarOpen: false,
		commands:         cfgCommands,
		viewportHeight:   0,
		viewportWidth:    0,
		layoutRows:       []string{""},
		layout:           "",
		titleText:        "The Cute File Manager",
	}

	// Initialize the search input
	m.searchInput = m.SearchInput(">", "Filter...")

	// Initialize the command input
	m.commandInput = m.CommandInput("COMMAND: ", "")

	// Load command history for auto-complete
	m.commandHistory = m.LoadCommandHistory()
	m.historyMatches = []string{}
	m.historyIndex = -1

	m.CalcLayout()

	ActiveTuiMode = TuiModeNormal
	ActiveFileListMode = FileListModeDir

	return m
}
