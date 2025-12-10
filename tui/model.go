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
	ModeAddFile           TUIMode
	ModeCd                TUIMode
	ModeColumnVisibiliy   TUIMode
	ModeCommand           TUIMode
	ModeCopy              TUIMode
	ModeFilter            TUIMode
	ModeGoto              TUIMode
	ModeHelp              TUIMode
	ModeMkdir             TUIMode
	ModeMove              TUIMode
	ModeNormal            TUIMode
	ModeQuit              TUIMode
	ModeRemove            TUIMode
	ModeRename            TUIMode
	ModeSelect            TUIMode
	ModeSort              TUIMode
	ModeFileListSplitPane TUIMode
}

type (
	FileListMode  string
	FileListModes struct {
		ModeNormal  FileListMode
		ModeCommand FileListMode
		ModeFilter  FileListMode
		ModeHelp    FileListMode
	}
)

type ViewPrimitive interface {
	View() string
}

type ComponentArgs struct {
	Width  int
	Height int
}

type CurrentDirComponentArgs struct {
	Width      int
	Height     int
	CurrentDir string
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
	currentDir string
	allFiles   []filesystem.FileInfo
	files      []filesystem.FileInfo
	fileList   list.Model
	// filterQuery stores the current filter string for this pane only.
	// Each pane maintains its own independent filter so that split panes
	// do not share search state.
	filterQuery     string
	dirBackStack    []string
	dirForwardStack []string
	columns         []filesystem.FileInfoColumn
	// marked tracks which file paths are currently marked for actions (multi-select)
	// independently of the cursor row. The key is the full filesystem path.
	marked map[string]bool
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
	ModeAddFile           TUIMode = "ADD_FILE"
	ModeCd                TUIMode = "CD"
	ModeColumnVisibiliy   TUIMode = "COLUMN VISIBILIY"
	ModeCommand           TUIMode = "COMMAND"
	ModeCopy              TUIMode = "COPY"
	ModeFilter            TUIMode = "FILTER"
	ModeGoto              TUIMode = "GOTO"
	ModeHelp              TUIMode = "HELP"
	ModeMkdir             TUIMode = "MKDIR"
	ModeMove              TUIMode = "MOVE"
	ModeNormal            TUIMode = "NORMAL"
	ModeParent            TUIMode = "PARENT"
	ModeQuit              TUIMode = "QUIT"
	ModeRemove            TUIMode = "REMOVE"
	ModeRename            TUIMode = "RENAME"
	ModeSelect            TUIMode = "SELECT"
	ModeSort              TUIMode = "SORT"
	ModeFileListSplitPane TUIMode = "SPLIT PANE"

	ModalNone ModalKind = "None"
	ModalHelp ModalKind = "Help"

	SortingAsc  SortColumnByDirection = "ASC"
	SortingDesc SortColumnByDirection = "DESC"

	PreviewPaneType       SplitPaneType = "PREVIEW"
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
		ModeAddFile:         ModeAddFile,
		ModeCd:              ModeCd,
		ModeColumnVisibiliy: ModeColumnVisibiliy,
		ModeCommand:         ModeCommand,
		ModeCopy:            ModeCopy,
		ModeFilter:          ModeFilter,
		ModeGoto:            ModeGoto,
		ModeHelp:            ModeHelp,
		ModeMkdir:           ModeMkdir,
		ModeNormal:          ModeNormal,
		ModeRemove:          ModeRemove,
		ModeRename:          ModeRename,
		ModeSelect:          ModeSelect,
		ModeSort:            ModeSort,
	}
)

type Model struct {
	activeModal        ModalKind
	activeSplitPane    SplitPaneType
	activeViewport     ActiveViewportType
	commandHistory     []string // Command history for auto-complete
	commandInput       textinput.Model
	configDir          string
	countPrefix        int            // countPrefix stores a pending numeric prefix for Vim-style navigation (e.g. "10j" / "3↓" in the file list). A value of 0 means "no active prefix".
	fileInfoViewport   viewport.Model // Independent state for each file-list pane.
	height             int
	helpScrollOffset   int      // Help modal scroll state
	historyIndex       int      // Current index in historyMatches for navigation
	historyMatches     []string // Filtered matches based on current input
	isActionInProgress bool
	isSplitPaneOpen    bool
	isSudo             bool
	jumpTo             string
	layout             string
	layoutRows         []string
	leftPane           filePane
	menuCursor         int
	rightPane          filePane
	runtimeConfig      *config.RuntimeConfig // runtimeConfig holds the Lua-backed configuration (theme and commands).
	// searchInput is the interactive text input used when the user is in
	// filter mode. It is shared visually, but each pane keeps its own
	// filterQuery string so that split panes can have independent filters.
	searchInput    textinput.Model
	showRightPane  bool
	sortColumnBy   SortColumnBy
	terminalType   string // Terminal / preview state
	theme          theming.Theme
	titleText      string
	viewportHeight int
	viewportWidth  int
	width          int

	// Components
	CurrentDir   func(m Model, args CurrentDirComponentArgs) string
	FileListView func(m Model, args FileListComponentArgs) string
	FileInfo     func(m Model, args ComponentArgs) string
	Header       func(m Model, args ComponentArgs) string
	PreviewTabs  func(m Model, args ComponentArgs) string
	SearchBar    func(m Model, args ComponentArgs) string
	SearchText   func(m Model, view ActiveViewportType) string
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

// activePane returns the filePane corresponding to the currently active
// viewport. When the file-list split pane is not open, this is always the
// left pane.
func (m *Model) GetActivePane() *filePane {
	if m.activeViewport == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return &m.rightPane
	}
	return &m.leftPane
}

func (m Model) GetActiveViewport() ActiveViewportType {
	return m.activeViewport
}

func (m Model) GetColumnVisibility() []filesystem.FileInfoColumn {
	return m.leftPane.columns
}

func (m Model) GetColumnVisibilityForViewport(view ActiveViewportType) []filesystem.FileInfoColumn {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.columns
	}
	return m.leftPane.columns
}

func (m Model) GetCommandInput() textinput.Model {
	return m.commandInput
}

func (m Model) GetCommandInputView() string {
	return m.commandInput.View()
}

func (m Model) GetConfigDir() string {
	return m.configDir
}

func (m Model) GetIsActionInProgress() bool {
	return m.isSplitPaneOpen
}

func (m Model) GetIsSplitPaneOpen() bool {
	return m.isSplitPaneOpen
}

func (m Model) GetLeftPaneCurrentDir() string {
	return m.leftPane.currentDir
}

func (m Model) GetLeftPaneFileList() list.Model {
	return m.leftPane.fileList
}

// GetLeftPaneFileListForViewport returns the file list model for the given viewport
// (left or right). In split file-list mode, each pane has its own list;
// otherwise, only the left list is active.
func (m Model) GetLeftPaneFileListForViewport(view ActiveViewportType) list.Model {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.fileList
	}
	return m.leftPane.fileList
}

func (m Model) GetMenuCursor() int {
	return m.menuCursor
}

func (m Model) GetPreviewViewport() viewport.Model {
	return m.fileInfoViewport
}

func (m Model) GetRightCurrentDir() string {
	return m.rightPane.currentDir
}

func (m Model) GetRightPaneCurrentDir() string {
	return m.rightPane.currentDir
}

func (m Model) GetSearchInputText() string {
	// Default to the active viewport's filter text.
	if m.activeViewport == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.filterQuery
	}
	return m.leftPane.filterQuery
}

func (m Model) GetSearchInputView() string {
	return m.searchInput.View()
}

// GetSearchInputTextForViewport returns the filter text associated with the
// given viewport, allowing split panes to display their own filter headers.
func (m Model) GetSearchInputTextForViewport(view ActiveViewportType) string {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.filterQuery
	}
	return m.leftPane.filterQuery
}

func (m Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m Model) GetSortColumnBy() SortColumnBy {
	return m.sortColumnBy
}

func (m Model) GetTheme() theming.Theme {
	return m.theme
}

func (m Model) GetTerminalType() string {
	return m.terminalType
}

func (m Model) GetTitleText() string {
	return m.titleText
}

func (m Model) GetHelpScrollOffset() int {
	return m.helpScrollOffset
}
