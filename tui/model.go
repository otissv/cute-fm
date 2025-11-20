package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"lsfm/filesystem"
	"lsfm/theming"
)

// ModalKind represents the type of modal currently shown, if any.
type ModalKind int

const (
	ModalNone ModalKind = iota
	ModalHelp
)

// Model represents the main application state
type Model struct {
	// First row: Text input for search/commands
	textInput textinput.Model

	// Second row: Two viewports side by side
	leftViewport  viewport.Model // Left panel viewport
	rightViewport viewport.Model // Right panel viewport

	// Help modal content (rendered inside a floating window when active).
	helpViewport viewport.Model

	// Data backing the left viewport (directory listing).
	// allFiles contains the complete directory listing; files is the
	// currently visible (possibly filtered) subset.
	allFiles   []filesystem.FileInfo
	files      []filesystem.FileInfo
	currentDir string

	// Index of the currently selected file in the list (0-based).
	// -1 indicates "no selection".
	selectedIndex int

	// Currently active modal, if any.
	activeModal ModalKind

	// Theme configuration loaded from lsfm.toml.
	theme theming.Theme

	// Layout dimensions
	width  int
	height int
}

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize text input for the first row
	ti := textinput.New()
	ti.Placeholder = "Search or enter command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Initialize left viewport for the second row
	leftVp := viewport.New(0, 0)

	// Initialize right viewport for the second row
	rightVp := viewport.New(0, 0)
	rightVp.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	// Initialize help viewport content for the help modal.
	helpContent := `
Help
----

Navigation:
  Up/Down arrows   Move selection in file list
  Scroll wheel     Scroll file list

Search:
  Type in the search bar to filter files by name

General:
  ?                Toggle this help
  ctrl+c / ctrl+q  Quit
`
	helpVp := viewport.New(0, 0)
	helpVp.SetContent(strings.TrimSpace(helpContent))

	// Load theme configuration.
	theme := theming.LoadTheme("lsfm.toml")

	// Determine initial directory for the file list.
	wd := startDir
	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			wd = "."
		}
	}

	files, err := filesystem.ListDirectory(wd)
	selected := -1
	if err != nil {
		// If we can't list the directory, show a simple error message.
		leftVp.SetContent("Error reading directory:\n" + err.Error())
	} else {
		if len(files) > 0 {
			selected = 0
		}
		// Fill the left viewport with a table of directory contents.
		leftVp.SetContent(renderFileTable(theme, files, selected))
	}

	return Model{
		textInput:     ti,
		leftViewport:  leftVp,
		rightViewport: rightVp,
		helpViewport:  helpVp,
		allFiles:      files,
		files:         files,
		currentDir:    wd,
		selectedIndex: selected,
		theme:         theme,
	}
}

// Init initializes the model (required by Bubble Tea)
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		m.width = msg.Width
		m.height = msg.Height

		// Calculate viewport dimensions
		// CWD row: 1 content line
		// Text input row: Height(1) with borders = 3 total lines (1 content + 2 border)
		cwdRowHeight := 1
		textInputRowHeight := 3
		// Viewport style height: remaining height after the top two rows
		viewportStyleHeight := msg.Height - (cwdRowHeight + textInputRowHeight)
		if viewportStyleHeight < 3 {
			viewportStyleHeight = 3 // Minimum: 1 content + 2 borders
		}
		// Viewport content height (scrollable area): style height - 2 border lines
		viewportContentHeight := viewportStyleHeight - 2
		if viewportContentHeight < 1 {
			viewportContentHeight = 1 // Minimum content height
		}

		// Calculate viewport width (half of available width, accounting for borders)
		viewportWidth := msg.Width / 2

		// Update left viewport dimensions
		// Height is the content height (viewport handles scrolling internally)
		m.leftViewport.Width = viewportWidth
		m.leftViewport.Height = viewportContentHeight

		// Update right viewport dimensions
		// Height is the content height (viewport handles scrolling internally)
		m.rightViewport.Width = viewportWidth
		m.rightViewport.Height = viewportContentHeight

		// Set text input width to full width (accounting for borders)
		m.textInput.Width = msg.Width - 2

		// Resize the help viewport to fit nicely in a floating window.
		helpWidth := msg.Width / 2
		helpHeight := msg.Height / 2
		if helpWidth < 20 {
			helpWidth = 20
		}
		if helpHeight < 5 {
			helpHeight = 5
		}
		m.helpViewport.Width = helpWidth
		m.helpViewport.Height = helpHeight - 2 // account for borders/padding

		// Ensure the selected row stays visible after resize.
		m = ensureSelectionVisible(m)

		return m, nil

	case tea.KeyMsg:
		// If a modal is active, handle its keys first.
		if m.activeModal != ModalNone {
			switch msg.String() {
			case "esc", "q", "?":
				// Close help modal.
				m.activeModal = ModalNone
				return m, nil
			}

			// For now, help modal is static; ignore other keys while open.
			return m, nil
		}

		// Navigate the file list with arrow keys.
		switch msg.String() {
		case "up":
			m = moveSelection(m, -1)
			return m, nil
		case "down":
			m = moveSelection(m, 1)
			return m, nil
		}

		// Open help modal with '?' when no modal is active.
		if msg.String() == "?" {
			m.activeModal = ModalHelp
			return m, nil
		}

		// Handle quit
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"))):
			return m, tea.Quit
		}
	}

	// Update text input (first row) and apply filtering if the value changed.
	before := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	if m.textInput.Value() != before {
		m = applyFilter(m)
	}

	// Update left viewport (second row, left column)
	m.leftViewport, cmd = m.leftViewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update right viewport (second row, right column)
	m.rightViewport, cmd = m.rightViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// moveSelection returns an updated model with the selection moved by delta
// rows (negative for up, positive for down). The selection is clamped to the
// valid range and the table is re-rendered.
func moveSelection(m Model, delta int) Model {
	if len(m.files) == 0 {
		return m
	}

	newIndex := m.selectedIndex + delta
	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= len(m.files) {
		newIndex = len(m.files) - 1
	}
	if newIndex == m.selectedIndex {
		return m
	}

	m.selectedIndex = newIndex
	m.leftViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex))
	m = ensureSelectionVisible(m)

	return m
}

// ensureSelectionVisible adjusts the left viewport's scroll offset so that the
// selected row remains visible.
func ensureSelectionVisible(m Model) Model {
	if m.selectedIndex < 0 {
		return m
	}

	// Header row is at line 0; first file row is at line 1.
	line := 1 + m.selectedIndex
	viewHeight := m.leftViewport.Height
	if viewHeight <= 0 {
		return m
	}

	// Current scroll offset (top visible line).
	y := m.leftViewport.YOffset

	// If the selected line is above the viewport, scroll up.
	if line < y+1 {
		m.leftViewport.SetYOffset(line - 1)
		return m
	}

	// If the selected line is below the viewport, scroll down so it becomes
	// the last visible line.
	if line > y+viewHeight-1 {
		m.leftViewport.SetYOffset(line - viewHeight + 1)
	}

	return m
}

// applyFilter recomputes the visible file list based on the current value of
// the text input. The filter is a case-insensitive substring match on the file
// name. When the filter changes, the selection is clamped to the new list and
// the table is re-rendered.
func applyFilter(m Model) Model {
	query := strings.TrimSpace(m.textInput.Value())

	// If there is no backing data yet, nothing to do.
	if len(m.allFiles) == 0 {
		return m
	}

	if query == "" {
		// Reset to full list.
		m.files = m.allFiles
	} else {
		lq := strings.ToLower(query)
		var filtered []filesystem.FileInfo
		for _, fi := range m.allFiles {
			if strings.Contains(strings.ToLower(fi.Name), lq) {
				filtered = append(filtered, fi)
			}
		}
		m.files = filtered
	}

	// Adjust selection for the new list.
	if len(m.files) == 0 {
		m.selectedIndex = -1
		m.leftViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex))
		return m
	}

	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	if m.selectedIndex >= len(m.files) {
		m.selectedIndex = len(m.files) - 1
	}

	m.leftViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex))
	m = ensureSelectionVisible(m)

	return m
}
