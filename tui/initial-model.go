package tui

import (
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"

	"cute/config"
	"cute/filesystem"
	"cute/theming"
)

func InitialModel(startDir string) Model {
	fileInfoViewport := viewport.New()
	fileInfoViewport.SetContent("Right Pane\n\nThis is the right viewport.\nIt will display file previews.")

	configDir := config.GetConfigDir()

	if err := SaveDefaultSettings(configDir); err != nil {
		_ = err
	}

	tomlSettings, _ := LoadSettings(configDir)

	// Determine initial directory for the file list.
	// Priority: command line argument > settings.toml > current directory
	leftCurrentDir := startDir
	if leftCurrentDir == "" {
		// Try to load from settings
		if tomlSettings != nil && tomlSettings.StartDir != "" {
			if tomlSettings.StartDir == "Home" {
				if homeDir, err := os.UserHomeDir(); err == nil {
					leftCurrentDir = homeDir
				}
			} else {
				leftCurrentDir = tomlSettings.StartDir
			}
		}
		// Fallback to current directory
		if leftCurrentDir == "" {
			var err error
			leftCurrentDir, err = os.Getwd()
			if err != nil {
				leftCurrentDir = "."
			}
		}
	}

	files := loadDirectory(leftCurrentDir)
	theme := theming.GetTheme()

	// Create the bubbles lists with file items for both panes.
	delegate := NewFileItemDelegate(theme, 0, filesystem.ColumnNames)
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
		configDir:        configDir,
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
		menuCursorIndex: 0,
		rightPane: filePane{
			currentDir:  leftCurrentDir,
			allFiles:    files,
			files:       files,
			fileList:    rightList,
			filterQuery: "",
			columns:     defaultColumns,
			marked:      make(map[string]bool),
		},
		showRightPane: true,
		sortColumnBy: SortColumnBy{
			column:    filesystem.ColumnName,
			direction: SortingAsc,
		},
		terminalType:   string(detectTerminalType()),
		theme:          theme,
		titleText:      "Cute File Manager",
		viewportHeight: 0,
		viewportWidth:  0,
	}

	// Initialize default settings
	defaultSettings := Settings{
		StartDir:            leftCurrentDir,
		SortColumnBy:        filesystem.ColumnName,
		SortColumnDirection: SortingAsc,
		ColumnVisibility:    m.leftPane.columns,
		SplitPane:           FileInfoSplitPaneType,
		FileListMode:        FileListModeList,
	}

	m.settings = MergeSettings(defaultSettings, tomlSettings)
	m.settings.StartDir = leftCurrentDir

	m.activeSplitPane = m.settings.SplitPane
	if m.settings.SplitPane != "" {
		m.isSplitPaneOpen = m.settings.SplitPane != ""
	}

	if m.settings.FileListMode != "" {
		ActiveFileListMode = m.settings.FileListMode
	}

	if len(m.settings.ColumnVisibility) > 0 {
		m.leftPane.columns = m.settings.ColumnVisibility
		m.rightPane.columns = m.settings.ColumnVisibility
	}

	if m.settings.SortColumnBy != "" {
		m.sortColumnBy.column = m.settings.SortColumnBy
	}
	if m.settings.SortColumnDirection != "" {
		m.sortColumnBy.direction = m.settings.SortColumnDirection
	}

	m.searchInput = m.SearchInput("> ", "Filter...")
	m.commandInput = m.CommandInput("", "")
	m.commandHistory = m.LoadCommandHistory()

	m.CalcLayout()

	ActiveTuiMode = ModeNormal

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
