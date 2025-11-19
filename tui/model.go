package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the main application state
type Model struct {
	// First row: Text input for search/commands
	textInput textinput.Model

	// Second row: Two viewports side by side
	leftViewport  viewport.Model // Left panel viewport
	rightViewport viewport.Model // Right panel viewport

	// Layout dimensions
	width  int
	height int
}

// InitialModel creates a new model with default values
func InitialModel() Model {
	// Initialize text input for the first row
	ti := textinput.New()
	ti.Placeholder = "Search or enter command..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Initialize left viewport for the second row
	leftVp := viewport.New(0, 0)
	leftVp.SetContent("Left Panel\n\nThis is the left viewport.\nIt will display file listings.")

	// Initialize right viewport for the second row
	rightVp := viewport.New(0, 0)
	rightVp.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	return Model{
		textInput:     ti,
		leftViewport:  leftVp,
		rightViewport: rightVp,
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
		// Text input row: Height(1) with borders = 3 total lines (1 content + 2 border)
		textInputRowHeight := 3
		// Viewport style height: remaining height after text input row
		viewportStyleHeight := msg.Height - textInputRowHeight
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

		return m, nil

	case tea.KeyMsg:
		// Handle quit
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"))):
			return m, tea.Quit
		}
	}

	// Update text input (first row)
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	// Update left viewport (second row, left column)
	m.leftViewport, cmd = m.leftViewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update right viewport (second row, right column)
	m.rightViewport, cmd = m.rightViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
