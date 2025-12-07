package tui

import (
	"os"
	"path/filepath"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"

	"cute/config"
	"cute/filesystem"
)

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize right viewport for the second row
	fileInfoViewport := viewport.New()
	fileInfoViewport.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	// Determine initial directory for the file list.
	currentDir := startDir
	if currentDir == "" {
		var err error
		currentDir, err = os.Getwd()
		if err != nil {
			currentDir = "."
		}
	}
	cfgDir := getConfigDir()

	// Load Lua-based runtime configuration (theme + commands).
	runtimeCfg := config.LoadRuntimeConfig(cfgDir)

	// Load the initial directory.
	files := loadDirectory(currentDir)

	// Create the bubbles list with file items.
	delegate := NewFileItemDelegate(runtimeCfg.Theme, 0, filesystem.ColumnNames)
	items := FileInfosToItems(files)
	fileList := list.New(items, delegate, 0, 0)

	// Configure the list appearance - hide built-in UI elements since we have custom ones.
	fileList.SetShowTitle(false)
	fileList.SetShowStatusBar(false)
	fileList.SetShowFilter(false)
	fileList.SetShowHelp(false)
	fileList.SetShowPagination(false)
	fileList.SetFilteringEnabled(false)
	fileList.DisableQuitKeybindings()

	// Use a simple style for the list.
	fileList.Styles.NoItems = fileList.Styles.NoItems.Foreground(nil)

	m := Model{
		allFiles:           files,
		configDir:          cfgDir,
		currentDir:         currentDir,
		fileInfoViewport:   fileInfoViewport,
		fileList:           fileList,
		files:              files,
		imagePreviewActive: false,
		jumpTo:             "",
		lastPreviewedPath:  "",
		layout:             "",
		layoutRows:         []string{""},
		previewEnabled:     false,
		runtimeConfig:      runtimeCfg,
		showRightPanel:     true,
		isSudo:             false,
		terminalType:       string(detectTerminalType()),
		theme:              runtimeCfg.Theme,
		titleText:          "The Cute File Manager",
		viewportHeight:     0,
		viewportWidth:      0,
		columnVisibility: []filesystem.FileInfoColumn{
			filesystem.ColumnPermissions,
			filesystem.ColumnUser,
			filesystem.ColumnGroup,
			filesystem.ColumnDateModified,
			filesystem.ColumnName,
		},
		columnVisibilityCursor: 0,
	}

	// Initialize the search input
	m.searchInput = m.SearchInput(">", "Filter...")

	// Initialize the command input
	m.commandInput = m.CommandInput("", "")

	// Load command history for auto-complete
	m.commandHistory = m.LoadCommandHistory()
	m.historyMatches = []string{}
	m.historyIndex = -1

	m.CalcLayout()

	ActiveTuiMode = TuiModeNormal
	// Start in the default "list all" view mode so navigation (enter/backspace)
	// does not implicitly force the directories-only view ("ld").
	ActiveFileListMode = FileListModeList

	// Initialize preview for the initial selection, if any.
	m.UpdateFileInfoPanel()

	return m
}

// loadDirectory lists the given directory and returns the file list.
func loadDirectory(dir string) []filesystem.FileInfo {
	files, err := filesystem.ListDirectory(dir)
	if err != nil {
		return nil
	}
	return files
}

// UpdateFileListDelegate updates the delegate with a new width.
func (m *Model) UpdateFileListDelegate(width int) {
	delegate := NewFileItemDelegate(m.theme, width, m.columnVisibility)
	m.fileList.SetDelegate(delegate)
}

func getConfigDir() string {
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

	return cfgDir
}
