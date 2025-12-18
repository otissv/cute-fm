package tui

import (
	"unicode"

	tea "charm.land/bubbletea/v2"

	"cute/filesystem"
)

func (m Model) DeviceMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	bindings := GetKeyBindings()

	// Update the device list to handle scrolling
	m.deviceList, cmd = m.deviceList.Update(msg)
	cmds = append(cmds, cmd)

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, tea.Batch(cmds...)
	}

	key := keyMsg.String()

	// Vim-style numeric prefix: accumulate digits which will be applied
	// to the next navigation command in the device list.
	//
	// NOTE: if a digit key is also configured as a dedicated "Goto" binding,
	// we must *not* swallow it here; otherwise the Goto handler below would
	// never see the key press and the mode would never change.
	if len(key) == 1 && unicode.IsDigit(rune(key[0])) && !bindings.Goto.Matches(key) {
		d := int(key[0] - '0')
		m.countPrefix = m.countPrefix*10 + d
		return m, tea.Batch(cmds...)
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

	m.commandInput.Blur()
	m.searchInput.Focus()

	switch {
	// Quit application
	case bindings.Quit.Matches(key):
		SetQuitMode()
		return m, tea.Batch(cmds...)

	// Add change (cd)
	case bindings.Cd.Matches(key):
		if ActiveTuiMode != ModeCd {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeCd

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		}
		return m, nil

	// Open column visibility window
	case bindings.ColumnVisibilityFileList.Matches(key):
		if ActiveTuiMode != ModeColumnVisibility {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeColumnVisibility
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Enter command mode
	case bindings.Command.Matches(key):
		if ActiveTuiMode != ModeCommand {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeCommand

			m.commandInput.SetValue("")
			m.commandInput.Focus()
			m.historyMatches = []string{}
			m.historyIndex = -1
		}
		return m, nil

	// Move cursor down in device list (with optional count)
	case bindings.CursorDown.Matches(key):
		for i := 0; i < count; i++ {
			m.deviceList.CursorDown()
		}
		return m, tea.Batch(cmds...)

	// Move cursor up in device list (with optional count)
	case bindings.CursorUp.Matches(key):
		for i := 0; i < count; i++ {
			m.deviceList.CursorUp()
		}
		return m, tea.Batch(cmds...)

	// Change file list to directories-only view
	case bindings.Directories.Matches(key):
		ActiveFileListMode = "ld"
		ActiveTuiMode = ModeNormal
		m.ApplyFileListFilter()
		return m, nil

		// Change file list to files-only view
	case bindings.Files.Matches(key):
		ActiveFileListMode = "lf"
		ActiveTuiMode = ModeNormal
		m.ApplyFileListFilter()

		return m, nil

	// Goto home directory
	case bindings.Home.Matches(keyMsg.String()):
		res, _ := m.ExecuteCommand("cd ~/")
		ActiveTuiMode = ModeNormal
		ActiveFileListMode = PreviousFileListMode

		pane := m.GetActivePane()
		if res.Cwd != "" && res.Cwd != pane.currentDir {
			m.ChangeDirectory(res.Cwd)
		} else if res.Refresh {
			m.ReloadDirectory()
		}
		return m, nil

		// Open help window
	case bindings.Help.Matches(key):
		if ActiveTuiMode != ModeHelp {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeHelp
		}
		return m, nil

		// Change file list to list-all-items view
	case bindings.List.Matches(key):
		ActiveFileListMode = "ll"
		ActiveTuiMode = ModeNormal
		m.ApplyFileListFilter()
		return m, nil

		// Open settings window
	case bindings.Settings.Matches(key):
		if ActiveTuiMode != ModeSettings {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeSettings
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Move cursor to start of device list
	case bindings.GoToStart.Matches(key):
		m.deviceList.Select(0)
		return m, tea.Batch(cmds...)

	// Move cursor to end of device list
	case bindings.GoToEnd.Matches(key):
		items := m.deviceList.Items()
		if len(items) > 0 {
			m.deviceList.Select(len(items) - 1)
		}
		return m, tea.Batch(cmds...)

	// Enter Goto mode
	case bindings.Goto.Matches(key):
		if ActiveTuiMode != ModeGoto {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeGoto

			m.jumpTo = key
			m.commandInput.SetValue(key)
			m.commandInput.Focus()
		}
		return m, tea.Batch(cmds...)

	case bindings.Sort.Matches(keyMsg.String()):
		ActiveTuiMode = ModeSort
		// Initialize cursor to currently selected sort column
		sortBy := m.GetSortDeviceColumnBy()
		if sortBy.column != "" {
			for i, col := range filesystem.DeviceInfoColumnNames {
				if col == sortBy.column {
					m.menuCursorIndex = i
					break
				}
			}
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
