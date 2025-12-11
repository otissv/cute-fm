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
func InitialModel(startDir string) Model {
	fileInfoViewport := viewport.New()
	fileInfoViewport.SetContent("Right Pane\n\nThis is the right viewport.\nIt will display file previews.")

	// Determine initial directory for the file list.
	leftCurrentDir := startDir
	if leftCurrentDir == "" {
		var err error
		leftCurrentDir, err = os.Getwd()
		if err != nil {
			leftCurrentDir = "."
		}
	}
	cfgDir := getConfigDir()

	// Load Lua-based runtime configuration (theme + commands).
	runtimeCfg := config.LoadRuntimeConfig(cfgDir)

	// Load the initial directory.
	files := loadDirectory(leftCurrentDir)

	// Create the bubbles lists with file items for both panes.
	delegate := NewFileItemDelegate(runtimeCfg.Theme, 0, filesystem.ColumnNames)
	items := FileInfosToItems(files, nil)

	newList := func() list.Model {
		l := list.New(items, delegate, 0, 0)
		// Configure the list appearance - hide built-in UI elements since we have custom ones.
		l.SetShowTitle(false)
		l.SetShowStatusBar(false)
		l.SetShowFilter(false)
		l.SetShowHelp(false)
		l.SetShowPagination(false)
		l.SetFilteringEnabled(false)
		l.DisableQuitKeybindings()
		// Use a simple style for the list.
		l.Styles.NoItems = l.Styles.NoItems.Foreground(nil)
		return l
	}

	leftList := newList()
	rightList := newList()

	// Default visible columns for new panes.
	defaultColumns := []filesystem.FileInfoColumn{
		filesystem.ColumnPermissions,
		filesystem.ColumnUser,
		filesystem.ColumnGroup,
		filesystem.ColumnDateModified,
		filesystem.ColumnName,
	}

	m := Model{
		activeSplitPane:  FileInfoSplitPaneType,
		activeViewport:   LeftViewportType,
		configDir:        cfgDir,
		fileInfoViewport: fileInfoViewport,
		historyIndex:     -1,
		historyMatches:   []string{},
		isSplitPaneOpen:  false,
		isSudo:           false,
		jumpTo:           "",
		layout:           "",
		layoutRows:       []string{""},
		leftPane: filePane{
			currentDir:  leftCurrentDir,
			allFiles:    files,
			files:       files,
			fileList:    leftList,
			filterQuery: "",
			columns:     defaultColumns,
			marked:      make(map[string]bool),
		},
		menuCursor: 0,
		rightPane: filePane{
			currentDir:  leftCurrentDir,
			allFiles:    files,
			files:       files,
			fileList:    rightList,
			filterQuery: "",
			columns:     defaultColumns,
			marked:      make(map[string]bool),
		},
		runtimeConfig: runtimeCfg,
		showRightPane: true,
		sortColumnBy: SortColumnBy{
			column:    filesystem.ColumnName,
			direction: SortingAsc,
		},
		terminalType:   string(detectTerminalType()),
		theme:          runtimeCfg.Theme,
		titleText:      "Cute File Manager",
		viewportHeight: 0,
		viewportWidth:  0,
	}

	m.settings = Settings{
		StartDir:        leftCurrentDir,
		SortColumnBy:    filesystem.ColumnName,
		ColumnVisibiliy: m.leftPane.columns,
		SplitPane:       FileInfoSplitPaneType,
		FileListMode:    FileListModeList,
	}

	m.searchInput = m.SearchInput("> ", "Filter...")

	m.commandInput = m.CommandInput("", "")

	m.commandHistory = m.LoadCommandHistory()

	m.CalcLayout()

	ActiveTuiMode = ModeNormal
	ActiveFileListMode = FileListModeList

	m.UpdateFileInfoPane()

	return m
}

func loadDirectory(dir string) []filesystem.FileInfo {
	files, err := filesystem.ListDirectory(dir)
	if err != nil {
		return nil
	}
	return files
}

func (m *Model) UpdateFileListDelegate(width int) {
	leftDelegate := NewFileItemDelegate(m.theme, width, m.leftPane.columns)
	m.leftPane.fileList.SetDelegate(leftDelegate)

	rightDelegate := NewFileItemDelegate(m.theme, width, m.rightPane.columns)
	m.rightPane.fileList.SetDelegate(rightDelegate)
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
