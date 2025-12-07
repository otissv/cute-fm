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
	SplitPanelType     string
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
	Width          int
	Height         int
	SplitPanelType ActiveViewportType
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

	PreviewPanelType       SplitPanelType = "PREVIEW"
	FileInfoSplitPanelType SplitPanelType = "FILE_INFO"
	FileListSplitPanelType SplitPanelType = "FILE_LIST"

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
	configDir        string
	activeModal      ModalKind
	activeSplitPanel SplitPanelType
	activeViewport   ActiveViewportType
	showRightPanel   bool
	isSudo           bool
	jumpTo           string
	isSearchBarOpen  bool
	isSplitPanelOpen bool

	searchInput  textinput.Model
	commandInput textinput.Model

	// countPrefix stores a pending numeric prefix for Vim-style
	// navigation (e.g. "10j" / "3↓" in the file list).
	// A value of 0 means "no active prefix".
	countPrefix int

	// runtimeConfig holds the Lua-backed configuration (theme and commands).
	runtimeConfig *config.RuntimeConfig

	fileList         list.Model // Bubbles list for file listing
	fileInfoViewport viewport.Model

	allFiles        []filesystem.FileInfo
	files           []filesystem.FileInfo
	leftCurrentDir  string
	rightCurrentDir string

	// Directory navigation history (similar to a web browser's back/forward).
	// dirBackStack holds previously visited directories; the most recent entry
	// is at the end of the slice.
	dirBackStack []string
	// dirForwardStack holds directories we can move forward to after going
	// back. The most recent "forward" target is at the end of the slice.
	dirForwardStack []string

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

	menuCursor       int
	columnVisibility []filesystem.FileInfoColumn
	sortColumnBy     SortColumnBy

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
	return m.allFiles
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
	return m.leftCurrentDir
}

func (m Model) GetRightCurrentDir() string {
	return m.rightCurrentDir
}

func (m Model) GetFileList() list.Model {
	return m.fileList
}

func (m Model) GetFiles() []filesystem.FileInfo {
	return m.files
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
	return m.fileList.Index()
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

// GetColumnVisibility returns the currently selected columns for the file list.
func (m Model) GetColumnVisibility() []filesystem.FileInfoColumn {
	return m.columnVisibility
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

func (m Model) GetActiveSplitPanel() SplitPanelType {
	return m.activeSplitPanel
}

func (m Model) GetActiveViewport() ActiveViewportType {
	return m.activeViewport
}

func (m Model) GetIsSplitPanelOpen() bool {
	return m.isSplitPanelOpen
}
