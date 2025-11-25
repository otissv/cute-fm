package tui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"cute/filesystem"
	"cute/theming"
)

type ModalKind string

const (
	ModalNone = "none"
	ModalHelp = "help"
)

type ViewPrimitive interface {
	View() string
}

type Model struct {
	configDir string

	searchInput  textinput.Model
	commandInput textinput.Model

	commands map[string]string

	fileListViewport viewport.Model
	previewViewport  viewport.Model

	allFiles      []filesystem.FileInfo
	files         []filesystem.FileInfo
	currentDir    string
	selectedIndex int // Index of the currently selected file in the list (0-based). -1 indicates "no selection".
	viewMode      string

	activeModal ModalKind

	theme theming.Theme

	isCommandBarOpen bool
	isSearchBarOpen  bool

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
	HelpModal   func(m Model) string
	CommandBar  func(m Model) string
	SearchBar   func(m Model) string
	CurrentDir  func(m Model) string
	Header      func(m Model) string
	StatusBar   func(m Model, items ...string) string
	ViewText    func(m Model) string
	PreviewTabs func(m Model) string
	Preview     func(m Model) string
	FileList    func(m Model) string
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) GetCommandInput() textinput.Model {
	return m.commandInput
}

func (m Model) GetCommandInputView() string {
	return m.commandInput.View()
}

func (m Model) GetCurrentDir() string {
	return m.currentDir
}

func (m Model) GetSearchInputView() string {
	return m.searchInput.View()
}

func (m Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m Model) GetTheme() theming.Theme {
	return m.theme
}

func (m Model) GetViewportWidth() int {
	return m.viewportWidth
}

// Getters for additional model fields so components outside the tui package
// can access read-only state without touching unexported fields directly.

func (m Model) GetConfigDir() string {
	return m.configDir
}

func (m Model) GetSearchInput() textinput.Model {
	return m.searchInput
}

func (m Model) GetCommands() map[string]string {
	return m.commands
}

func (m Model) GetFileListViewport() viewport.Model {
	return m.fileListViewport
}

func (m Model) GetPreviewViewport() viewport.Model {
	return m.previewViewport
}

func (m Model) GetAllFiles() []filesystem.FileInfo {
	return m.allFiles
}

func (m Model) GetFiles() []filesystem.FileInfo {
	return m.files
}

func (m Model) GetSelectedIndex() int {
	return m.selectedIndex
}

func (m Model) GetViewMode() string {
	return m.viewMode
}

func (m Model) GetActiveModal() ModalKind {
	return m.activeModal
}

func (m Model) IsCommandBarOpen() bool {
	return m.isCommandBarOpen
}

func (m Model) IsSearchBarOpen() bool {
	return m.isSearchBarOpen
}

func (m Model) GetCommandHistory() []string {
	return m.commandHistory
}

func (m Model) GetHistoryMatches() []string {
	return m.historyMatches
}

func (m Model) GetHistoryIndex() int {
	return m.historyIndex
}

func (m Model) GetViewportHeight() int {
	return m.viewportHeight
}

func (m Model) GetLayout() string {
	return m.layout
}

func (m Model) GetLayoutRows() []string {
	return m.layoutRows
}

func (m Model) GetTitleText() string {
	return m.titleText
}
