package tui

import (
	"path/filepath"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// NormalMode handles keybindings when the TUI is in its default ("normal")
// browsing mode.
func (m Model) NormalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()

	// Vim-style numeric prefix: accumulate digits which will be applied
	// to the next navigation command in the file list.
	//
	// NOTE: if a digit key is also configured as a dedicated "Goto" binding,
	// we must *not* swallow it here; otherwise the Goto handler below would
	// never see the key press and the mode would never change.
	//
	// Examples:
	//   "5↓"  -> move 5 entries down
	//   "10↑" -> move 10 entries up
	if len(key) == 1 && unicode.IsDigit(rune(key[0])) && !bindings.Goto.Matches(key) {
		d := int(key[0] - '0')
		m.countPrefix = m.countPrefix*10 + d
		return m, nil
	}

	// Capture the current count before we reset it. A zero prefix means
	// "no explicit count", which we treat as 1.
	count := 1
	if m.countPrefix > 0 {
		count = m.countPrefix
	}

	// For any non-digit key we process, clear the prefix afterwards so
	// it only applies to a single command, like in Vim.
	defer func() {
		m.countPrefix = 0
	}()

	switch {
	// Add file
	case bindings.AddFile.Matches(key):
		if ActiveTuiMode != TuiModeAddFile {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeAddFile

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Add change (cd)
	case bindings.Cd.Matches(key):
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
	case bindings.Command.Matches(key):
		if ActiveTuiMode != TuiModeCommand {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeCommand

			m.commandInput.SetValue("")
			m.commandInput.Focus()
			m.historyMatches = []string{}
			m.historyIndex = -1

			return m, nil
		}

	// Copy file or folder
	case bindings.Copy.Matches(key):
		if ActiveTuiMode != TuiModeCopy {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeCopy

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Change file list to directories-only view
	case bindings.Directories.Matches(key):
		ActiveFileListMode = "ld"
		m.ApplyFilter()
		return m, nil

	// Move cursor down in file list (with optional count)
	case bindings.Down.Matches(key):
		// Arrow keys should move one row at a time, nano-style. Ignore any
		// numeric count prefix so that a stray digit doesn't cause the cursor
		// to "jump" or effectively page.
		m.fileList.CursorDown()
		m.UpdatePreview()
		return m, nil

	// Navigate into the selected directory.
	case bindings.Enter.Matches(key):
		selectedIdx := m.fileList.Index()
		if selectedIdx >= 0 && selectedIdx < len(m.files) {
			fi := m.files[selectedIdx]
			if fi.IsDir {
				m.ChangeDirectory(fi.Path)
				return m, nil
			}
		}

	// Change file list to files-only view
	case bindings.Files.Matches(key):
		ActiveFileListMode = "lf"
		m.ApplyFilter()
		return m, nil

	// Enter filter mode
	case bindings.Filter.Matches(key):
		if ActiveTuiMode != TuiModeFilter {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeFilter

			m.searchInput.Focus()

			return m, nil
		}

	// Enter Goto mode
	case bindings.Goto.Matches(key):
		if ActiveTuiMode != TuiModeGoto {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeGoto

			m.jumpTo = key
			m.commandInput.SetValue(key)
			m.commandInput.Focus()
			return m, nil
		}

		// Move cursor to end of file list
	case bindings.GoToEnd.Matches(keyMsg.String()):
		m.fileList.GoToEnd()
		m.UpdatePreview()
		return m, nil

		// Move cursor to start of file list
	case bindings.GoToStart.Matches(keyMsg.String()):
		m.fileList.GoToStart()
		m.UpdatePreview()
		return m, nil

	// Move file or folder
	case bindings.Move.Matches(key):
		if ActiveTuiMode != TuiModeMove {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeMove

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Open help modal
	case bindings.Help.Matches(key):
		if ActiveTuiMode != TuiModeHelp {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeHelp
			return m, nil
		}

	// Change file list to list-all-items view
	case bindings.List.Matches(key):
		ActiveFileListMode = "ll"
		m.ApplyFilter()
		return m, nil

	// Navigate to the parent directory.
	case bindings.Parent.Matches(key):
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
	case bindings.Preview.Matches(key):
		m.previewEnabled = !m.previewEnabled
		m.UpdatePreview()
		return m, nil

	// Quit application
	case bindings.Quit.Matches(key):
		SetQuitMode()
		return m, nil

	// Remove file or directory
	case bindings.Remove.Matches(key):
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
	case bindings.Select.Matches(key):
		if ActiveTuiMode != TuiModeSelect {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = TuiModeSelect
			return m, nil
		}

	// Move cursor up in file list (with optional count)
	case bindings.Up.Matches(key):
		// Arrow keys should move one row at a time, nano-style. Ignore any
		// numeric count prefix so that a stray digit doesn't cause the cursor
		// to "jump" or effectively page.
		m.fileList.CursorUp()
		m.UpdatePreview()
		return m, nil

		// Page down in file list (with optional count)
	case bindings.PageDown.Matches(key):
		for i := 0; i < count; i++ {
			m.fileList.NextPage()
		}
		m.UpdatePreview()
		return m, nil

	// Page up in file list (with optional count)
	case bindings.PageUp.Matches(key):
		for i := 0; i < count; i++ {
			m.fileList.PrevPage()
		}
		m.UpdatePreview()
		return m, nil
	}

	return m, nil
}
