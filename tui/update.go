package tui

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"cute/command"
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		m.width = msg.Width
		m.height = msg.Height + 2

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
		if ActiveTuiMode == TuiModeAddFile {
			return m.UtilityMode(msg, "touch")
		}

		if ActiveTuiMode == TuiModeCd {
			return m.UtilityMode(msg, "cd")
		}

		if ActiveTuiMode == TuiModeCommand {
			return m.CommandMode(msg)
		}

		if ActiveTuiMode == TuiModeCopy {
			return m.UtilityMode(msg, "cp")
		}

		if ActiveTuiMode == TuiModeFilter {
			return m.FilterMode(msg)
		}

		if ActiveTuiMode == TuiModeHelp {
			return m.HelpMode(msg)
		}

		if ActiveTuiMode == TuiModeMkdir {
			return m.UtilityMode(msg, "mkdir")
		}
		if ActiveTuiMode == TuiModeMove {
			return m.UtilityMode(msg, "mv")
		}

		if ActiveTuiMode == TuiModeNormal {
			return m.NormalMode(msg)
		}

		if ActiveTuiMode == TuiModeQuit {
			return m.QuitMode(msg)
		}

		if ActiveTuiMode == TuiModeRemove {
			return m.ConfirmMode(msg, "rm -r")
		}

	}

	return m, nil
}

func (m *Model) ExecuteCommand(line string) (command.Result, error) {
	if line != "" {
		m.AppendCommandHistory(line)
		// Reload history to include the new command
		m.commandHistory = m.LoadCommandHistory()
	}

	env := m.GetCommandEnvironment()

	res, err := command.Execute(env, line)

	return res, err
}

func (m *Model) GetSelectedEntry() *command.SelectedEntry {
	selectedIdx := m.fileList.Index()
	if selectedIdx < 0 || selectedIdx >= len(m.files) {
		return nil
	}

	fi := m.files[selectedIdx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(m.currentDir, fi.Name)
	}

	return &command.SelectedEntry{
		Name:  fi.Name,
		Path:  path,
		IsDir: fi.IsDir,
		Type:  fi.Type,
	}
}

// GetCommandEnvironment builds the command execution environment using the
// current model state, including the currently selected entry (if any).
func (m *Model) GetCommandEnvironment() command.Environment {
	return command.Environment{
		Cwd:      m.currentDir,
		Config:   m.runtimeConfig,
		Selected: m.GetSelectedEntry(),
	}
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

	// Update preview for the new selection after filtering.
	m.UpdatePreview()
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

	// And recompute the preview for the new directory/selection.
	m.UpdatePreview()
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
