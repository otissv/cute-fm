package tui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"cute/config"
	"cute/filesystem"
	"cute/theming"
)

type (
	ModalKind          string
	SplitPaneType      string
	TUIMode            string
	ActiveViewportType string
)

type TUIModes struct {
	TuiModeAddFile         TUIMode
	TuiModeCd              TUIMode
	TuiModeColumnVisibiliy TUIMode
	TuiModeCommand         TUIMode
	TuiModeCopy            TUIMode
	TuiModeFilter          TUIMode
	TuiModeGoto            TUIMode
	TuiModeHelp            TUIMode
	TuiModeMkdir           TUIMode
	TuiModeMove            TUIMode
	TuiModeNextDir         TUIMode
	TuiModeNormal          TUIMode
	TuiModeQuit            TUIMode
	TuiModeRemove          TUIMode
	TuiModeRename          TUIMode
	TuiModeSelect          TUIMode
	TuiModeSort            TUIMode
}

type (
	FileListMode  string
	FileListModes struct {
		TuiModeNormal  FileListMode
		TuiModeCommand FileListMode
		TuiModeFilter  FileListMode
		TuiModeHelp    FileListMode
	}
)

type ViewPrimitive interface {
	View() string
}

type ComponentArgs struct {
	Width  int
	Height int
}

type FileListComponentArgs struct {
	Width         int
	Height        int
	SplitPaneType ActiveViewportType
}

type CommandModalArgs struct {
	Title       string
	Prompt      string
	Placeholder string
}

type MenuCursor struct {
	Selected   string
	Unselected string
	Prompt     string
}
type MenuArgs struct {
	Choices     []string
	Cursor      int
	Selected    map[string]string
	CursorTypes MenuCursor
}

type DialogArgs struct {
	X int
	Y int
}

type DialogModalArgs struct {
	Title   string
	Content string
}

type ColumnModelArgs struct {
	Title      string
	Selected   string
	Unselected string
	Prompt     string
}

type SelectedEntry struct {
	// Name is the base name of the entry.
	Name string
	// Path is the full filesystem path to the entry.
	Path string
	// IsDir indicates whether the entry is a directory.
	IsDir bool
	// Type is the classified file type string ("directory", "regular", ...).
	Type string
}

type SortColumnByDirection string

type SortColumnBy struct {
	column    filesystem.FileInfoColumn
	direction SortColumnByDirection
}

// filePane holds all file-list-related state for a single pane (left or right).
// Keeping this in a separate struct lets us support independent panes while
// reusing the same logic for navigation, filtering, and directory history.
type filePane struct {
	currentDir      string
	allFiles        []filesystem.FileInfo
	files           []filesystem.FileInfo
	fileList        list.Model
	dirBackStack    []string
	dirForwardStack []string
	columns         []filesystem.FileInfoColumn
}

// Column returns the currently selected sort column.
func (s SortColumnBy) Column() filesystem.FileInfoColumn {
	return s.column
}

// Direction returns the current sort direction for the column.
func (s SortColumnBy) Direction() SortColumnByDirection {
	return s.direction
}

const (
	TuiModeAddFile         TUIMode = "ADD_FILE"
	TuiModeAutoComplete    TUIMode = "AUTOCOMPLETE"
	TuiModeCd              TUIMode = "CD"
	TuiModeColumnVisibiliy TUIMode = "COLUMN_VISIBILIY"
	TuiModeCommand         TUIMode = "COMMAND"
	TuiModeCopy            TUIMode = "COPY"
	TuiModeFilter          TUIMode = "FILTER"
	TuiModeGoto            TUIMode = "GOTO"
	TuiModeHelp            TUIMode = "HELP"
	TuiModeMkdir           TUIMode = "MKDIR"
	TuiModeMove            TUIMode = "MOVE"
	TuiModeNormal          TUIMode = "NORMAL"
	TuiModeParent          TUIMode = "PARENT"
	TuiModeQuit            TUIMode = "QUIT"
	TuiModeRemove          TUIMode = "REMOVE"
	TuiModeRename          TUIMode = "RENAME"
	TuiModeSelect          TUIMode = "SELECT"
	TuiModeSort            TUIMode = "SORT"

	ModalNone ModalKind = "None"
	ModalHelp ModalKind = "Help"

	SortingAsc  SortColumnByDirection = "ASC"
	SortingDesc SortColumnByDirection = "DESC"

	PreviewPanelType      SplitPaneType = "PREVIEW"
	FileInfoSplitPaneType SplitPaneType = "FILE_INFO"
	FileListSplitPaneType SplitPaneType = "FILE_LIST"

	LeftViewportType  ActiveViewportType = "LEFT"
	RightViewportType ActiveViewportType = "RIGHT"

	FileListModeList FileListMode = "ll"
	FileListModeFile FileListMode = "lf"
	FileListModeDir  FileListMode = "ld"
)

var (
	ActiveFileListMode         = FileListModeList
	ActiveTuiMode      TUIMode = "NORMAL"
	PreviousTuiMode    TUIMode = "NORMAL"

	TuiModes = TUIModes{
		TuiModeAddFile:         TuiModeAddFile,
		TuiModeCd:              TuiModeCd,
		TuiModeColumnVisibiliy: TuiModeColumnVisibiliy,
		TuiModeCommand:         TuiModeCommand,
		TuiModeCopy:            TuiModeCopy,
		TuiModeFilter:          TuiModeFilter,
		TuiModeGoto:            TuiModeGoto,
		TuiModeHelp:            TuiModeHelp,
		TuiModeMkdir:           TuiModeMkdir,
		TuiModeNormal:          TuiModeNormal,
		TuiModeRemove:          TuiModeRemove,
		TuiModeRename:          TuiModeRename,
		TuiModeSelect:          TuiModeSelect,
		TuiModeSort:            TuiModeSort,
	}
)

type Model struct {
	configDir       string
	activeModal     ModalKind
	activeSplitPane SplitPaneType
	activeViewport  ActiveViewportType
	showRightPanel  bool
	isSudo          bool
	jumpTo          string
	isSearchBarOpen bool
	isSplitPaneOpen bool

	searchInput  textinput.Model
	commandInput textinput.Model

	// countPrefix stores a pending numeric prefix for Vim-style
	// navigation (e.g. "10j" / "3↓" in the file list).
	// A value of 0 means "no active prefix".
	countPrefix int

	// runtimeConfig holds the Lua-backed configuration (theme and commands).
	runtimeConfig *config.RuntimeConfig

	fileInfoViewport viewport.Model

	// Independent state for each file-list pane.
	leftPane  filePane
	rightPane filePane

	theme theming.Theme

	// Command history for auto-complete
	commandHistory []string
	historyMatches []string // Filtered matches based on current input
	historyIndex   int      // Current index in historyMatches for navigation

	width          int
	height         int
	viewportHeight int
	viewportWidth  int
	layout         string
	layoutRows     []string
	titleText      string

	menuCursor   int
	sortColumnBy SortColumnBy

	// Help modal scroll state
	helpScrollOffset int

	// Terminal / preview state
	terminalType string

	// Components
	CurrentDir   func(m Model, args ComponentArgs) string
	FileListView func(m Model, args FileListComponentArgs) string
	Header       func(m Model, args ComponentArgs) string
	Preview      func(m Model, args ComponentArgs) string
	PreviewTabs  func(m Model, args ComponentArgs) string
	SearchBar    func(m Model, args ComponentArgs) string
	StatusBar    func(m Model, args ComponentArgs, items ...string) string
	SudoMode     func(m Model, args ComponentArgs) string
	TuiMode      func(m Model, args ComponentArgs) string
	ViewModeText func(m Model, args ComponentArgs) string

	// Modals
	ColumnModal  func(m Model, args ColumnModelArgs) *lipgloss.Layer
	CommandModal func(m Model, args CommandModalArgs) *lipgloss.Layer
	DialogModal  func(m Model, args DialogModalArgs) *lipgloss.Layer
	HelpModal    func(m Model) *lipgloss.Layer
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) GetActiveModal() ModalKind {
	return m.activeModal
}

func (m Model) GetAllFiles() []filesystem.FileInfo {
	return m.leftPane.allFiles
}

func (m Model) GetCommandHistory() []string {
	return m.commandHistory
}

func (m Model) GetCommandInput() textinput.Model {
	return m.commandInput
}

func (m Model) GetCommandInputView() string {
	return m.commandInput.View()
}

func (m Model) GetCommands() map[string]string {
	// Deprecated: commands are now defined in Lua and executed through the
	// runtimeConfig; this method remains only to satisfy any existing callers.
	return nil
}

func (m Model) GetConfigDir() string {
	return m.configDir
}

func (m Model) GetCurrentDir() string {
	return m.leftPane.currentDir
}

func (m Model) GetRightCurrentDir() string {
	return m.rightPane.currentDir
}

func (m Model) GetFileList() list.Model {
	// For backward compatibility, return the left pane's file list.
	return m.leftPane.fileList
}

// GetFileListForViewport returns the file list model for the given viewport
// (left or right). In split file-list mode, each pane has its own list;
// otherwise, only the left list is active.
func (m Model) GetFileListForViewport(view ActiveViewportType) list.Model {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.fileList
	}
	return m.leftPane.fileList
}

func (m Model) GetFiles() []filesystem.FileInfo {
	return m.leftPane.files
}

func (m Model) GetHistoryIndex() int {
	return m.historyIndex
}

func (m Model) GetHistoryMatches() []string {
	return m.historyMatches
}

func (m Model) GetLayout() string {
	return m.layout
}

func (m Model) GetLayoutRows() []string {
	return m.layoutRows
}

func (m Model) GetPreviewViewport() viewport.Model {
	return m.fileInfoViewport
}

func (m Model) GetSearchInput() textinput.Model {
	return m.searchInput
}

func (m Model) GetSearchInputView() string {
	return m.searchInput.View()
}

func (m Model) GetSelectedIndex() int {
	return m.leftPane.fileList.Index()
}

func (m Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m Model) GetTheme() theming.Theme {
	return m.theme
}

// GetTerminalType returns the detected terminal type (e.g. "kitty").
func (m Model) GetTerminalType() string {
	return m.terminalType
}

func (m Model) GetTitleText() string {
	return m.titleText
}

func (m Model) GetViewportHeight() int {
	return m.viewportHeight
}

func (m Model) GetViewportWidth() int {
	return m.viewportWidth
}

func (m Model) GetHelpScrollOffset() int {
	return m.helpScrollOffset
}

// GetColumnVisibility returns the currently selected columns for the *left*
// file list. This is kept for backward compatibility with code that assumes a
// single file list.
func (m Model) GetColumnVisibility() []filesystem.FileInfoColumn {
	return m.leftPane.columns
}

// GetColumnVisibilityForViewport returns the selected columns for the given
// viewport (left or right). When the right file-list pane is active, this
// allows it to have independent column visibility.
func (m Model) GetColumnVisibilityForViewport(view ActiveViewportType) []filesystem.FileInfoColumn {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.columns
	}
	return m.leftPane.columns
}

// GetMenuCursor returns the current cursor index in the
// column-visibility modal.
func (m Model) GetMenuCursor() int {
	return m.menuCursor
}

// GetSortColumnBy returns the current sort column and direction.
func (m Model) GetSortColumnBy() SortColumnBy {
	return m.sortColumnBy
}

func (m Model) IsSearchBarOpen() bool {
	return m.isSearchBarOpen
}

func (m Model) GetActiveSplitPane() SplitPaneType {
	return m.activeSplitPane
}

func (m Model) GetActiveViewport() ActiveViewportType {
	return m.activeViewport
}

func (m Model) GetIsSplitPaneOpen() bool {
	return m.isSplitPaneOpen
}

// activePane returns the filePane corresponding to the currently active
// viewport. When the file-list split pane is not open, this is always the
// left pane.
func (m *Model) activePane() *filePane {
	if m.activeViewport == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return &m.rightPane
	}
	return &m.leftPane
}
