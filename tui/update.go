package tui

import (
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"cute/command"
	"cute/console"
	"cute/filesystem"
	"cute/theming"
)

func SetQuitMode() {
	if ActiveTuiMode != TuiModeQuit {
		PreviousTuiMode = ActiveTuiMode
		ActiveTuiMode = TuiModeQuit
	}
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

		helpWidth := msg.Width / 2
		helpHeight := msg.Height / 2
		if helpWidth < 20 {
			helpWidth = 20
		}
		if helpHeight < 5 {
			helpHeight = 5
		}

		m.CalcLayout()

		return m, nil

	case tea.KeyMsg:

		if ActiveTuiMode == TuiModeQuit {
			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				return m, tea.Quit

			case "esc":
				ActiveTuiMode = PreviousTuiMode
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeNormal {
			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				SetQuitMode()
				return m, nil
			case "?":
				if ActiveTuiMode != TuiModeHelp {
					PreviousTuiMode = ActiveTuiMode
					ActiveTuiMode = TuiModeHelp
					return m, nil
				}

			case ":":
				if ActiveTuiMode != TuiModeCommand {
					PreviousTuiMode = ActiveTuiMode
					ActiveTuiMode = TuiModeCommand

					// Prepare the command input when entering command mode.
					m.commandInput.SetValue("")
					m.commandInput.Focus()
					m.historyMatches = []string{}
					m.historyIndex = -1

					return m, nil
				}
			case "s":
				if ActiveTuiMode != TuiModeSelect {
					PreviousTuiMode = ActiveTuiMode
					ActiveTuiMode = TuiModeSelect
					return m, nil
				}

			case "f":
				if ActiveTuiMode != TuiModeFilter {
					PreviousTuiMode = ActiveTuiMode
					ActiveTuiMode = TuiModeFilter
					return m, nil
				}
			}
		}

		if ActiveTuiMode == TuiModeHelp {
			m.commandInput.Blur()
			m.searchInput.Focus()

			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				SetQuitMode()
				return m, nil

			case "esc", "?":
				ActiveTuiMode = PreviousTuiMode
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeFilter {
			// Update search input (first row) and apply filtering if the value changed.
			before := m.searchInput.Value()
			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)
			if m.searchInput.Value() != before {
				m.ApplyFilter()
			}

			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				SetQuitMode()
				return m, nil
			case "esc":
				ActiveTuiMode = TuiModeNormal
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeSelect {
			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				SetQuitMode()
				return m, nil
			case "esc":
				ActiveTuiMode = TuiModeNormal
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeCommand {
			m.searchInput.Blur()

			// Command bar is active; update it instead of the search bar.
			beforeValue := m.commandInput.Value()
			m.commandInput, cmd = m.commandInput.Update(msg)
			cmds = append(cmds, cmd)

			// Update left viewport (second row, left column)
			m.leftViewport, cmd = m.leftViewport.Update(msg)
			cmds = append(cmds, cmd)

			// Update right viewport (second row, right column)
			m.rightViewport, cmd = m.rightViewport.Update(msg)
			cmds = append(cmds, cmd)

			// Update history matches when input changes
			if m.commandInput.Value() != beforeValue {
				m.updateHistoryMatches()
			}

			switch msg.String() {
			case "ctrl+c", "ctrl+q":
				SetQuitMode()
				return m, nil

			case "esc", ":":
				ActiveTuiMode = PreviousTuiMode
				return m, nil

			case "tab":
				m.completeCommand()
				return m, nil

			case "up":

				if len(m.commandHistory) > 0 {
					m.navigateHistory(-1)
					return m, nil
				}
				// If no history, fall through to let textinput handle it

			case "down":

				if len(m.commandHistory) > 0 {
					m.navigateHistory(1)
					return m, nil
				}
				// If no history, fall through to let textinput handle it

			case "enter":
				line := strings.TrimSpace(m.commandInput.Value())
				if line != "" {
					m.AppendCommandHistory(line)
					// Reload history to include the new command
					m.commandHistory = m.LoadCommandHistory()
				}

				var selected *command.SelectedEntry
				if m.selectedIndex >= 0 && m.selectedIndex < len(m.files) {
					fi := m.files[m.selectedIndex]
					selected = &command.SelectedEntry{
						Name:  fi.Name,
						Path:  filepath.Join(m.currentDir, fi.Name),
						IsDir: fi.IsDir,
						Type:  fi.Type,
					}
				}

				env := command.Environment{
					Cwd:      m.currentDir,
					Config:   m.runtimeConfig,
					Selected: selected,
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
					ActiveFileListMode = FileListMode(res.ViewMode)
					console.Log("%s %s", ActiveFileListMode, res.ViewMode)
					m.ApplyFilter()
				}

				if res.OpenHelp {
					m.activeModal = ModalHelp
				}

				if res.Output != "" {
					m.rightViewport.SetContent(res.Output)
				}

				if err != nil && res.Output == "" {
					m.rightViewport.SetContent(err.Error())
				}

				m.commandInput.Blur()
				m.commandInput.SetValue("")
				m.searchInput.Focus()

				m.CalcLayout()

				if res.Quit {
					return m, tea.Quit
				}

				ActiveTuiMode = PreviousTuiMode

				return m, nil

			}
		}

		// Navigate the file list with arrow keys (only when not in command mode).
		switch msg.String() {

		case "up":
			m.moveSelection(-1)
			return m, nil
		case "down":
			m.moveSelection(1)
			return m, nil
		}
	}

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
	m.leftViewport.SetContent(
		renderFileTable(m.theme, m.files, m.selectedIndex, m.leftViewport.Width()),
	)
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
	viewHeight := m.leftViewport.Height()
	if viewHeight <= 0 {
		return
	}

	// Current scroll offset (top visible line).
	y := m.leftViewport.YOffset()

	// If the selected line is above the viewport, scroll up.
	if line < y+1 {
		m.leftViewport.SetYOffset(line - 1)
		return
	}

	// If the selected line is below the viewport, scroll down so it becomes
	// the last visible line.
	if line > y+viewHeight-1 {
		m.leftViewport.SetYOffset(line - viewHeight + 1)
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

	base := filterByViewMode(m.allFiles)

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
		m.leftViewport.SetContent(
			renderFileTable(m.theme, m.files, m.selectedIndex, m.leftViewport.Width()),
		)
	}

	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	if m.selectedIndex >= len(m.files) {
		m.selectedIndex = len(m.files) - 1
	}

	m.leftViewport.SetContent(
		renderFileTable(m.theme, m.files, m.selectedIndex, m.leftViewport.Width()),
	)
	m.EnsureSelectionVisible()
}

// ChangeDirectory updates the model to point at a new current directory and
// reloads the file list.
func (m *Model) ChangeDirectory(dir string) {
	files, selected := loadDirectoryIntoView(&m.leftViewport, m.theme, dir)
	m.currentDir = dir
	m.allFiles = files
	m.files = files
	m.selectedIndex = selected
	// Re-apply search/view filters for the new directory.
	m.ApplyFilter()
}

// filterByViewMode filters the given file list according to the current view
// mode. It does not modify the original slice.
func filterByViewMode(files []filesystem.FileInfo) []filesystem.FileInfo {
	if ActiveFileListMode == "" {
		ActiveFileListMode = "ll"
	}

	out := make([]filesystem.FileInfo, 0, len(files))

	switch ActiveFileListMode {
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
