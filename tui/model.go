package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"

	"cute/config"
	"cute/filesystem"
	"cute/theming"
)

type ModalKind int

const (
	ModalNone ModalKind = iota
	ModalHelp
)

type Model struct {
	searchInput  textinput.Model
	commandInput textinput.Model

	commands map[string]string

	FileListViewport viewport.Model
	previewViewport  viewport.Model
	helpViewport     viewport.Model

	allFiles      []filesystem.FileInfo
	files         []filesystem.FileInfo
	currentDir    string
	selectedIndex int // Index of the currently selected file in the list (0-based). -1 indicates "no selection".
	viewMode      string

	activeModal ModalKind

	theme theming.Theme

	isCommandBarOpen bool
	isSearchBarOpen  bool

	width          int
	height         int
	viewportHeight int
	viewportWidth  int
	layout         string
	layoutRows     []string
	titleText      string
}

// CommandInput constructs the command-line text input styled to match the command bar.
func (m Model) CommandInput(placeholder string, prompt string) textinput.Model {
	commandInput := textinput.New()
	commandInput.Prompt = prompt
	commandInput.Placeholder = placeholder
	commandInput.CharLimit = 256
	commandInput.Width = 50

	// Style the inner text input so its background/foreground align with the command bar.
	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.CommandBar.Background)).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground))

	commandInput.PromptStyle = baseStyle
	commandInput.TextStyle = baseStyle
	commandInput.PlaceholderStyle = baseStyle.
		Foreground(lipgloss.Color(m.theme.CommandBar.Placeholder))
	commandInput.CursorStyle = baseStyle

	return commandInput
}

// SearchInput constructs the search text input styled to match the search bar.
func (m Model) SearchInput(placeholder string, focus bool) textinput.Model {
	searchInput := textinput.New()
	searchInput.Placeholder = placeholder

	searchInput.CharLimit = 256
	searchInput.Width = 50

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.SearchBar.Background)).
		Foreground(lipgloss.Color(m.theme.SearchBar.Foreground))

	searchInput.TextStyle = baseStyle
	searchInput.PlaceholderStyle = baseStyle.
		Foreground(lipgloss.Color(m.theme.SearchBar.Placeholder))
	searchInput.CursorStyle = baseStyle

	if focus {
		searchInput.Focus()
	}

	return searchInput
}

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize left viewport for the second row
	leftVp := viewport.New(0, 0)

	// Initialize right viewport for the second row
	rightVp := viewport.New(0, 0)
	rightVp.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	helpViewport := HelpViewport()

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

	files, selected := loadDirectoryIntoView(&leftVp, theme, wd)

	m := Model{
		FileListViewport: leftVp,
		previewViewport:  rightVp,
		helpViewport:     helpViewport,
		allFiles:         files,
		files:            files,
		currentDir:       wd,
		selectedIndex:    selected,
		theme:            theme,
		isCommandBarOpen: true,
		commands:         cfgCommands,
		viewMode:         "ll",
		viewportHeight:   0,
		viewportWidth:    0,
		layoutRows:       []string{""},
		layout:           "",
		titleText:        "The Cute File Manager",
	}

	m.searchInput = m.SearchInput("Search...", false)
	m.commandInput = m.CommandInput("", "COMMAND: ")

	m.CalcLayout()

	return m
}

// Init initializes the model (required by Bubble Tea)
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// CalcLayout recalculates the viewport dimensions based on the current
// window size and whether command mode is active, then re-renders the file
// table and ensures the selection is visible.
func (m *Model) CalcLayout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	statusRowHeight := 1    // Status row at the bottom: 1 content line
	searchBarRowHeight := 3 // Search bar: 1 content line + 2 border lines
	headerRowHeight := 3

	// Only reserve vertical space for the command bar when it is visible.
	commandRowHeight := 0
	if m.isCommandBarOpen {
		commandRowHeight = 3 // Command bar: 1 content line + 2 border lines
	}

	// Viewport style height: remaining height after the top and bottom rows.
	viewportHeight := m.height - (statusRowHeight + searchBarRowHeight + commandRowHeight + headerRowHeight)
	if viewportHeight < 3 {
		viewportHeight = 3 // Minimum: 1 content + 2 borders
	}
	// Persist the total viewport box height so Lip Gloss containers (FileList,
	// Preview) can render with a fixed height instead of expanding to fit
	// their content.
	m.viewportHeight = viewportHeight

	// Viewport content height (scrollable area): style height - 2 border lines.
	viewportContentHeight := viewportHeight - 2
	if viewportContentHeight < 1 {
		viewportContentHeight = 1 // Minimum content height
	}

	// Calculate viewport width (half of available width, accounting for borders).

	m.viewportWidth = (m.width / 2) - 2

	// Update left viewport dimensions (height is the content height).
	m.FileListViewport.Width = m.viewportWidth
	m.FileListViewport.Height = viewportContentHeight / 2

	// Update right viewport dimensions (height is the content height).
	m.previewViewport.Width = m.viewportWidth
	m.previewViewport.Height = viewportContentHeight / 2

	// Re-render the file table for the new width and ensure the selection is
	// still visible.
	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
	m.EnsureSelectionVisible()
}
