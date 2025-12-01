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

type ModalKind string

const (
	ModalNone ModalKind = "None"
	ModalHelp ModalKind = "Help"
)

type TUIMode string

type TUIModes struct {
	TuiModeNormal  TUIMode
	TuiModeCommand TUIMode
	TuiModeFilter  TUIMode
	TuiModeHelp    TUIMode
	TuiModeSelect  TUIMode
	TuiModeQuit    TUIMode
}

const (
	TuiModeNormal  TUIMode = "NORMAL"
	TuiModeCommand TUIMode = "COMMAND"
	TuiModeFilter  TUIMode = "FILTER"
	TuiModeHelp    TUIMode = "HELP"
	TuiModeSelect  TUIMode = "SELECT"
	TuiModeQuit    TUIMode = "QUIT"
)

var TuiModes = TUIModes{
	TuiModeNormal:  TuiModeNormal,
	TuiModeCommand: TuiModeCommand,
	TuiModeFilter:  TuiModeFilter,
	TuiModeHelp:    TuiModeHelp,
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

const (
	FileListModeList FileListMode = "ll"
	FileListModeFile FileListMode = "lf"
	FileListModeDir  FileListMode = "ld"
)

type ViewPrimitive interface {
	View() string
}

type ComponentArgs struct {
	Width  int
	Height int
}

type DialogArgs struct {
	X int
	Y int
}

var (
	ActiveFileListMode = FileListModeList
	ActiveTuiMode      TUIMode
	PreviousTuiMode    TUIMode
)

type Model struct {
	configDir string

	searchInput  textinput.Model
	commandInput textinput.Model

	// runtimeConfig holds the Lua-backed configuration (theme and commands).
	runtimeConfig *config.RuntimeConfig

	fileList      list.Model // Bubbles list for file listing
	rightViewport viewport.Model

	allFiles   []filesystem.FileInfo
	files      []filesystem.FileInfo
	currentDir string

	activeModal ModalKind

	theme theming.Theme

	isSearchBarOpen bool

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

	// Components
	CurrentDir   func(m Model, args ComponentArgs) string
	FileListView func(m Model, args ComponentArgs) string
	Header       func(m Model, args ComponentArgs) string
	Preview      func(m Model, args ComponentArgs) string
	PreviewTabs  func(m Model, args ComponentArgs) string
	SearchBar    func(m Model, args ComponentArgs) string
	StatusBar    func(m Model, args ComponentArgs, items ...string) string
	TuiMode      func(m Model, args ComponentArgs) string
	ViewModeText func(m Model, args ComponentArgs) string

	// Modals
	HelpModal    func(m Model) *lipgloss.Layer
	CommandModal func(m Model) *lipgloss.Layer
	QuitModal    func(m Model) *lipgloss.Layer
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
	return m.currentDir
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
	return m.rightViewport
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

func (m Model) GetTitleText() string {
	return m.titleText
}

func (m Model) GetViewportHeight() int {
	return m.viewportHeight
}

func (m Model) GetViewportWidth() int {
	return m.viewportWidth
}

func (m Model) IsSearchBarOpen() bool {
	return m.isSearchBarOpen
}
