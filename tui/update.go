package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"cute/command"
	"cute/filesystem"
	"cute/theming"
)

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
		m.CalcLayout()

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
		if m.isCommandBarOpen {
			switch msg.String() {
			case "enter":
				// Execute the entered command.
				line := strings.TrimSpace(m.commandInput.Value())

				env := command.Environment{
					Cwd:            m.currentDir,
					ConfigCommands: m.commands,
				}

				res, err := command.Execute(env, line)

				// Apply environment changes.
				if res.Cwd != "" && res.Cwd != m.currentDir {
					m.ChangeDirectory(res.Cwd)
				} else if res.Refresh {
					// Re-list the current directory when requested by the command.
					m.ChangeDirectory(m.currentDir)
				}

				// Update view mode and re-apply filters so the file list view
				// actually changes when commands like "ll", "ls", "ld", "lf",
				// etc. are executed.
				if res.ViewMode != "" {
					m.viewMode = res.ViewMode
					m.ApplyFilter()
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
				m.isCommandBarOpen = false
				m.commandInput.Blur()
				m.commandInput.SetValue("")
				m.searchInput.Focus()

				// Grow the file/preview viewports back to fill the freed space.
				m.CalcLayout()

				if res.Quit {
					return m, tea.Quit
				}

				return m, nil

			case "esc", "q", "ctrl+c":
				// Exit command mode and return focus to the search bar.
				m.isCommandBarOpen = false
				m.commandInput.Blur()
				m.commandInput.SetValue("")
				m.searchInput.Focus()

				// Grow the file/preview viewports back to fill the freed space.
				m.CalcLayout()
				return m, nil
			}
		}

		// Navigate the file list with arrow keys (only when not in command mode).
		switch msg.String() {
		case "?":
			m.activeModal = ModalHelp
			return m, nil
		case ":":
			m.isCommandBarOpen = true
			m.commandInput.Blur()
			m.commandInput.SetValue("")
			m.searchInput.Focus()
			m.CalcLayout()
			return m, nil

		case "ctrl+f":
			if !m.isSearchBarOpen {
				m.searchInput.Focus()
			} else {
				m.searchInput.Blur()
			}
			m.isSearchBarOpen = !m.isSearchBarOpen
		case "up":
			m.moveSelection(-1)
			return m, nil
		case "down":
			m.moveSelection(1)
			return m, nil
		}

		// Handle quit (only when not in command mode).
		if !m.isCommandBarOpen {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "ctrl+q"))):
				return m, tea.Quit
			}
		}
	}

	// Update the appropriate text input.
	if m.isCommandBarOpen {
		// Command bar is active; update it instead of the search bar.
		m.commandInput, cmd = m.commandInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Update search input (first row) and apply filtering if the value changed.
		before := m.searchInput.Value()
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
		if m.searchInput.Value() != before {
			m.ApplyFilter()
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

// moveSelection updates the selection by delta rows (negative for up,
// positive for down). The selection is clamped to the valid range and the
// table is re-rendered.
func (m *Model) moveSelection(delta int) {
	if len(m.files) == 0 {
		return
	}

	newIndex := m.selectedIndex + delta
	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= len(m.files) {
		newIndex = len(m.files) - 1
	}
	if newIndex == m.selectedIndex {
		return
	}

	m.selectedIndex = newIndex
	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
	m.EnsureSelectionVisible()
}

// EnsureSelectionVisible adjusts the left viewport's scroll offset so that the
// selected row remains visible.
func (m *Model) EnsureSelectionVisible() {
	if m.selectedIndex < 0 {
		return
	}

	// Header row is at line 0; first file row is at line 1.
	line := 1 + m.selectedIndex
	viewHeight := m.FileListViewport.Height
	if viewHeight <= 0 {
		return
	}

	// Current scroll offset (top visible line).
	y := m.FileListViewport.YOffset

	// If the selected line is above the viewport, scroll up.
	if line < y+1 {
		m.FileListViewport.SetYOffset(line - 1)
		return
	}

	// If the selected line is below the viewport, scroll down so it becomes
	// the last visible line.
	if line > y+viewHeight-1 {
		m.FileListViewport.SetYOffset(line - viewHeight + 1)
	}
}

// ApplyFilter recomputes the visible file list based on the current value of
// the text input. The filter is a case-insensitive substring match on the file
// name. When the filter changes, the selection is clamped to the new list and
// the table is re-rendered.
func (m *Model) ApplyFilter() {
	query := strings.TrimSpace(m.searchInput.Value())

	// If there is no backing data yet, nothing to do.
	if len(m.allFiles) == 0 {
		return
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
	}

	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	if m.selectedIndex >= len(m.files) {
		m.selectedIndex = len(m.files) - 1
	}

	m.FileListViewport.SetContent(renderFileTable(m.theme, m.files, m.selectedIndex, m.FileListViewport.Width))
	m.EnsureSelectionVisible()
}

// ChangeDirectory updates the model to point at a new current directory and
// reloads the file list.
func (m *Model) ChangeDirectory(dir string) {
	files, selected := loadDirectoryIntoView(&m.FileListViewport, m.theme, dir)
	m.currentDir = dir
	m.allFiles = files
	m.files = files
	m.selectedIndex = selected
	// Re-apply search/view filters for the new directory.
	m.ApplyFilter()
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
