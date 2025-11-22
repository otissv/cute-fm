package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"lsfm/command"
	"lsfm/config"
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
	searchBar  textinput.Model
	commandBar textinput.Model

	// Whether the command bar is currently active/visible.
	commandMode bool

	FileListViewport viewport.Model
	previewViewport  viewport.Model

	helpViewport viewport.Model

	// Data backing the file list viewport (directory listing).
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

	// commands holds user-defined commands from the config file.
	commands map[string]string

	// viewMode is the current logical file-list view mode (e.g. "ll", "ls").
	viewMode string

	// Layout dimensions
	width  int
	height int

	viewportHeight int
	viewportWidth  int

	layout     string
	layoutRows []string
}

// InitialModel creates a new model with default values.
// If startDir is non-empty, it will be used as the initial directory for the
// file list; otherwise the current working directory is used.
func InitialModel(startDir string) Model {
	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.Focus()
	searchInput.CharLimit = 256
	searchInput.Width = 50

	// Initialize command input
	commandInput := textinput.New()
	commandInput.Prompt = ":"
	commandInput.Placeholder = "command..."
	commandInput.CharLimit = 256
	commandInput.Width = 50

	// Initialize left viewport for the second row
	leftVp := viewport.New(0, 0)

	// Initialize right viewport for the second row
	rightVp := viewport.New(0, 0)
	rightVp.SetContent("Right Panel\n\nThis is the right viewport.\nIt will display file previews.")

	helpViewport := HelpViewport()

	// Load theme configuration.
	theme := theming.LoadTheme("lsfm.toml")

	// Load user-defined commands from the configuration file.
	cfgCommands := config.LoadCommands("lsfm.toml")

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

	return Model{
		searchBar:        searchInput,
		commandBar:       commandInput,
		FileListViewport: leftVp,
		previewViewport:  rightVp,
		helpViewport:     helpViewport,
		allFiles:         files,
		files:            files,
		currentDir:       wd,
		selectedIndex:    selected,
		theme:            theme,
		commandMode:      false,
		commands:         cfgCommands,
		viewMode:         "ll",
		viewportHeight:   0,
		viewportWidth:    0,
		layoutRows:       []string{""},
		layout:           "",
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

		// Set text input widths to full width (accounting for borders)
		m.searchBar.Width = msg.Width - 2
		m.commandBar.Width = msg.Width - 2

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

		// Recalculate layout and re-render the file table so that the last
		// column can pad to the new viewport width and the selected row
		// highlight reaches the edge.
		m = recalcLayout(m)

		return m, nil

	case tea.KeyMsg:
		// If a modal is active, handle its keys first.
		if m.activeModal != ModalNone {
			switch msg.String() {
			case "esc", "q", "?", "ctrl+c":
				// Close help modal.
				m.activeModal = ModalNone
				return m, nil
			}

			// For now, help modal is static; ignore other keys while open.
			return m, nil
		}

		// If the command bar is active, handle its control keys first.
		if m.commandMode {
			switch msg.String() {
			case "enter":
				// Execute the entered command.
				line := strings.TrimSpace(m.commandBar.Value())

				env := command.Environment{
					Cwd:            m.currentDir,
					ConfigCommands: m.commands,
				}

				res, err := command.Execute(env, line)

				// Apply environment changes.
				if res.Cwd != "" && res.Cwd != m.currentDir {
					m = changeDirectory(m, res.Cwd)
				} else if res.Refresh {
					// Re-list the current directory when requested by the command.
					m = changeDirectory(m, m.currentDir)
				}

				// Update view mode and re-apply filters so the file list view
				// actually changes when commands like "ll", "ls", "ld", "lf",
				// etc. are executed.
				if res.ViewMode != "" {
					m.viewMode = res.ViewMode
					m = applyFilter(m)
				}

				if res.OpenHelp {
					m.activeModal = ModalHelp
				}

				if res.Output != "" {
					m.previewViewport.SetContent(res.Output)
				}

				// On error, show it in the preview if there was no other output.
				if err != nil && res.Output == "" {
					m.previewViewport.SetContent(err.Error())
				}

				// Exit command mode and return focus to the search bar.
				m.commandMode = false
				m.commandBar.Blur()
				m.commandBar.SetValue("")
				m.searchBar.Focus()

				// Grow the file/preview viewports back to fill the freed space.
				m = recalcLayout(m)

				if res.Quit {
					return m, tea.Quit
				}

				return m, nil

			case "esc", "q", "ctrl+c":
				// Exit command mode and return focus to the search bar.
				m.commandMode = false
				m.commandBar.Blur()
				m.commandBar.SetValue("")
				m.searchBar.Focus()

				// Grow the file/preview viewports back to fill the freed space.
				m = recalcLayout(m)
				return m, nil
			}
		}

		// Navigate the file list with arrow keys (only when not in command mode).
		if !m.commandMode {
			switch msg.String() {
			case "up":
				m = moveSelection(m, -1)
				return m, nil
			case "down":
				m = moveSelection(m, 1)
				return m, nil
			}
		}

		// Open help modal with '?' when no modal is active and not in command mode.
		if !m.commandMode && msg.String() == "?" {
			m.activeModal = ModalHelp
			return m, nil
		}

		// Enter command mode with ':' when not already in command mode.
		if !m.commandMode && msg.String() == ":" {
			m.commandMode = true
			m.commandBar.SetValue("")
			m.commandBar.Focus()
			m.searchBar.Blur()

			// Shrink the file/preview viewports to make room for the command bar.
			m = recalcLayout(m)
			return m, nil
		}

		// Handle quit (only when not in command mode).
		if !m.commandMode {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"))):
				return m, tea.Quit
			}
		}
	}

	// Update the appropriate text input.
	if m.commandMode {
		// Command bar is active; update it instead of the search bar.
		m.commandBar, cmd = m.commandBar.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Update search input (first row) and apply filtering if the value changed.
		before := m.searchBar.Value()
		m.searchBar, cmd = m.searchBar.Update(msg)
		cmds = append(cmds, cmd)
		if m.searchBar.Value() != before {
			m = applyFilter(m)
		}
	}

	// Update left viewport (second row, left column)
	m.FileListViewport, cmd = m.FileListViewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update right viewport (second row, right column)
	m.previewViewport, cmd = m.previewViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// recalcLayout recalculates the viewport dimensions based on the current
// window size and whether command mode is active, then re-renders the file
// table and ensures the selection is visible.
func recalcLayout(m Model) Model {
	if m.width <= 0 || m.height <= 0 {
		return m
	}

	statusRowHeight := 1    // Status row at the bottom: 1 content line
	searchBarRowHeight := 3 // Search bar: 1 content line + 2 border lines

	// Only reserve vertical space for the command bar when it is visible.
	commandRowHeight := 0
	if m.commandMode {
		commandRowHeight = 3 // Command bar: 1 content line + 2 border lines
	}

	// Viewport style height: remaining height after the top and bottom rows.
	viewportHeight := m.height - (statusRowHeight + searchBarRowHeight + commandRowHeight)
	if viewportHeight < 3 {
		viewportHeight = 3 // Minimum: 1 content + 2 borders
	}
	// Viewport content height (scrollable area): style height - 2 border lines.
	viewportContentHeight := viewportHeight - 2
	if viewportContentHeight < 1 {
		viewportContentHeight = 1 // Minimum content height
	}

	// Calculate viewport width (half of available width, accounting for borders).
	viewportWidth := m.width / 2

	// Update left viewport dimensions (height is the content height).
	m.FileListViewport.Width = viewportWidth
	m.FileListViewport.Height = viewportContentHeight

	// Update right viewport dimensions (height is the content height).
	m.previewViewport.Width = viewportWidth
	m.previewViewport.Height = viewportContentHeight

	// Re-render the file table for the new width and ensure the selection is
	// still visible.
	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
	m = ensureSelectionVisible(m)

	return m
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
	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
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
	viewHeight := m.FileListViewport.Height
	if viewHeight <= 0 {
		return m
	}

	// Current scroll offset (top visible line).
	y := m.FileListViewport.YOffset

	// If the selected line is above the viewport, scroll up.
	if line < y+1 {
		m.FileListViewport.SetYOffset(line - 1)
		return m
	}

	// If the selected line is below the viewport, scroll down so it becomes
	// the last visible line.
	if line > y+viewHeight-1 {
		m.FileListViewport.SetYOffset(line - viewHeight + 1)
	}

	return m
}

// applyFilter recomputes the visible file list based on the current value of
// the text input. The filter is a case-insensitive substring match on the file
// name. When the filter changes, the selection is clamped to the new list and
// the table is re-rendered.
func applyFilter(m Model) Model {
	query := strings.TrimSpace(m.searchBar.Value())

	// If there is no backing data yet, nothing to do.
	if len(m.allFiles) == 0 {
		return m
	}

	// First apply the current view mode (ll, ls, ld, lf, etc.).
	base := filterByViewMode(m.allFiles, m.viewMode)

	// Then apply the search query on top.
	if query == "" {
		m.files = base
	} else {
		lq := strings.ToLower(query)
		var filtered []filesystem.FileInfo
		for _, fi := range base {
			if strings.Contains(strings.ToLower(fi.Name), lq) {
				filtered = append(filtered, fi)
			}
		}
		m.files = filtered
	}

	// Adjust selection for the new list.
	if len(m.files) == 0 {
		m.selectedIndex = -1
		m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
		return m
	}

	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	if m.selectedIndex >= len(m.files) {
		m.selectedIndex = len(m.files) - 1
	}

	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
	m = ensureSelectionVisible(m)

	return m
}

// filterByViewMode filters the given file list according to the current view
// mode. It does not modify the original slice.
func filterByViewMode(files []filesystem.FileInfo, mode string) []filesystem.FileInfo {
	if mode == "" {
		mode = "lll"
	}

	out := make([]filesystem.FileInfo, 0, len(files))

	switch mode {
	case "ls":
		// Hide dotfiles (roughly emulating eza/ls without -a).
		for _, fi := range files {
			if strings.HasPrefix(fi.Name, ".") {
				continue
			}
			out = append(out, fi)
		}
	case "ld":
		// Only directories (hidden or not). Prefer the IsDir flag, but also
		// fall back to the classified file type to be defensive.
		for _, fi := range files {
			if fi.IsDir || fi.Type == "directory" {
				out = append(out, fi)
			}
		}
	case "lf":
		// Only files (hidden or not).
		for _, fi := range files {
			if !fi.IsDir {
				out = append(out, fi)
			}
		}
	default:
		// "ll"  or any unknown mode: show everything.
		out = append(out, files...)
	}

	return out
}

// loadDirectoryIntoView lists the given directory and loads it into the left
// viewport using the provided theme. It returns the file list and the selected
// index (0-based, or -1 if there is no selection).
func loadDirectoryIntoView(vp *viewport.Model, theme theming.Theme, dir string) ([]filesystem.FileInfo, int) {
	files, err := filesystem.ListDirectory(dir)
	selected := -1
	if err != nil {
		vp.SetContent("Error reading directory:\n" + err.Error())
		return nil, selected
	}

	if len(files) > 0 {
		selected = 0
	}

	// At this point we don't yet know the viewport width, so pass 0 for the
	// totalWidth and let a later resize re-render with the proper width.
	vp.SetContent(renderFileTable(theme, files, selected, 0))

	return files, selected
}

// changeDirectory updates the model to point at a new current directory and
// reloads the file list.
func changeDirectory(m Model, dir string) Model {
	files, selected := loadDirectoryIntoView(&m.FileListViewport, m.theme, dir)
	m.currentDir = dir
	m.allFiles = files
	m.files = files
	m.selectedIndex = selected
	// Re-apply search/view filters for the new directory.
	m = applyFilter(m)
	return m
}
