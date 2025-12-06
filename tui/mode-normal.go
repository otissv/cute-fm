package tui

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func (m Model) NormalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {

	// Add file
	case bindings.AddFile.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeAddFile {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeAddFile

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

		// Add change
	case bindings.Cd.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeCd {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeCd

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Enter command mode
	case bindings.Command.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeCommand {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeCommand

			m.commandInput.SetValue("")
			m.commandInput.Focus()
			m.historyMatches = []string{}
			m.historyIndex = -1

			return m, nil
		}

	// Change file list to directoties only view
	case bindings.Directories.Matches(keyMsg.String()):
		ActiveFileListMode = "ld"
		m.ApplyFilter()
		return m, nil

	// Move cursor down in file list
	case bindings.Down.Matches(keyMsg.String()):
		m.fileList.CursorDown()
		m.UpdatePreview()
		return m, nil

		// Navigate into the selected directory.
	case bindings.Enter.Matches(keyMsg.String()):
		selectedIdx := m.fileList.Index()
		if selectedIdx >= 0 && selectedIdx < len(m.files) {
			fi := m.files[selectedIdx]
			if fi.IsDir {
				m.ChangeDirectory(fi.Path)
				return m, nil
			}
		}

	// Change file list to files only view
	case bindings.Files.Matches(keyMsg.String()):
		ActiveFileListMode = "lf"
		m.ApplyFilter()
		return m, nil

	// Enter filter mode
	case bindings.Filter.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeFilter {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeFilter

			m.searchInput.Focus()

			return m, nil
		}

	// Move move cursor to end of file list
	case bindings.GoToEnd.Matches(keyMsg.String()):
		m.fileList.GoToEnd()
		m.UpdatePreview()
		return m, nil

		// Move move cursor to start of file list
	case bindings.GoToStart.Matches(keyMsg.String()):
		m.fileList.GoToStart()
		m.UpdatePreview()
		return m, nil

		// Add new directory
	case bindings.Mkdir.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeMkdir {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeMkdir

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Open help modal
	case bindings.Help.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeHelp {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeHelp
			return m, nil
		}

		// Change file list to list all items view
	case bindings.List.Matches(keyMsg.String()):
		ActiveFileListMode = "ll"
		m.ApplyFilter()
		return m, nil

	// Navigate to the parent directory.
	case bindings.Parent.Matches(keyMsg.String()):
		parent := filepath.Dir(m.currentDir)
		if parent != "" && parent != m.currentDir {
			m.ChangeDirectory(parent)
		} else {
			// Even if we're at the root (Dir("/") == "/"), attempt to
			// reload so the listing stays fresh.
			m.ChangeDirectory(m.currentDir)
		}
		return m, nil

	// Toggle preview
	case bindings.Preview.Matches(keyMsg.String()):
		m.previewEnabled = !m.previewEnabled
		m.UpdatePreview()
		return m, nil

	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

		// Remove file or direcory
	case bindings.Remove.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeRemove {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeRemove

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Enter select mode
	case bindings.Select.Matches(keyMsg.String()):
		if ActiveTuiMode != TuiModeSelect {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeSelect
			return m, nil
		}

	// Move cursor up in file list
	case bindings.Up.Matches(keyMsg.String()):
		m.fileList.CursorUp()
		m.UpdatePreview()
		return m, nil
	}

	return m, nil
}
