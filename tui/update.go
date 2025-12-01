package tui

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"cute/command"
	"cute/console"
	"cute/filesystem"
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

	bindings := GetKeyBindings()

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
			switch {
			case bindings.Quit.Matches(msg.String()):
				return m, tea.Quit

			case bindings.Cancel.Matches(msg.String()):
				ActiveTuiMode = PreviousTuiMode
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeNormal {
			switch msg.String() {
			case "ctrl+c":
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

					// Prepare the search input when entering filter mode.
					m.searchInput.Focus()

					return m, nil
				}

			case "up", "k":
				m.fileList.CursorUp()
				return m, nil

			case "down", "j":
				m.fileList.CursorDown()
				return m, nil

			case "g":
				m.fileList.GoToStart()
				return m, nil

			case "G":
				m.fileList.GoToEnd()
				return m, nil

			case "ctrl+l":
				ActiveFileListMode = "ll"
				m.ApplyFilter()
				return m, nil

			case "ctrl+d":
				ActiveFileListMode = "ld"
				m.ApplyFilter()
				return m, nil

			case "ctrl+f":
				ActiveFileListMode = "lf"
				m.ApplyFilter()
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeHelp {
			m.commandInput.Blur()
			m.searchInput.Focus()

			switch msg.String() {
			case "ctrl+c":
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
			case "ctrl+c":
				SetQuitMode()
				return m, nil
			case "esc":
				ActiveTuiMode = TuiModeNormal
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeSelect {
			switch msg.String() {
			case "ctrl+c":
				SetQuitMode()
				return m, nil
			case "esc":
				ActiveTuiMode = TuiModeNormal
				return m, nil
			}
		}

		if ActiveTuiMode == TuiModeCommand {
			m.searchInput.Blur()

			// Command modal is active; update it instead of the search bar.
			beforeValue := m.commandInput.Value()
			m.commandInput, cmd = m.commandInput.Update(msg)
			cmds = append(cmds, cmd)

			// Update right viewport (second row, right column)
			m.rightViewport, cmd = m.rightViewport.Update(msg)
			cmds = append(cmds, cmd)

			// Update history matches when input changes
			if m.commandInput.Value() != beforeValue {
				m.updateHistoryMatches()
			}

			switch msg.String() {
			case "ctrl+c":
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

				res, err := m.ExecuteCommand(line)

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
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) ExecuteCommand(line string) (command.Result, error) {
	if line != "" {
		m.AppendCommandHistory(line)
		// Reload history to include the new command
		m.commandHistory = m.LoadCommandHistory()
	}

	var selected *command.SelectedEntry
	selectedIdx := m.fileList.Index()
	if selectedIdx >= 0 && selectedIdx < len(m.files) {
		fi := m.files[selectedIdx]
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

	return res, err
}

// ApplyFilter recomputes the visible file list based on the current value of
// the text input. The filter is a case-insensitive substring match on the file
// name. When the filter changes, the list is updated with the new items.
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

	// Update the list with new items.
	items := FileInfosToItems(m.files)
	m.fileList.SetItems(items)

	// Reset selection to first item if we have items.
	if len(m.files) > 0 {
		m.fileList.Select(0)
	}
}

// ChangeDirectory updates the model to point at a new current directory and
// reloads the file list.
func (m *Model) ChangeDirectory(dir string) {
	files, err := filesystem.ListDirectory(dir)
	if err != nil {
		m.rightViewport.SetContent("Error reading directory:\n" + err.Error())
		return
	}

	m.currentDir = dir
	m.allFiles = files
	m.files = files

	// Update the list with new items.
	items := FileInfosToItems(files)
	m.fileList.SetItems(items)

	// Select the first item if we have items.
	if len(files) > 0 {
		m.fileList.Select(0)
	}

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
