package tui

import (
	"path/filepath"
	"unicode"

	tea "charm.land/bubbletea/v2"

	"cute/filesystem"
)

func (m Model) NormalMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()

	// Vim-style numeric prefix: accumulate digits which will be applied
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
	// Quit application
	case bindings.Quit.Matches(key):
		SetQuitMode()
		return m, nil

	// Change to sudo mode
	case bindings.Sudo.Matches(key):
		m.isSudo = !m.isSudo
		return m, nil

	// Add file
	case bindings.AddFile.Matches(key):
		if ActiveTuiMode != ModeAddFile {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeAddFile

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Add change (cd)
	case bindings.Cd.Matches(key):
		if ActiveTuiMode != ModeCd {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeCd

			m.commandInput.SetValue("")
			m.commandInput.Focus()
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

		// Switch to device view
	case bindings.Device.Matches(key):
		if ActiveTuiMode != ModeDevice {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeDevice
			ActiveFileListMode = FileListModeDevice
			m.previousSplitPane = m.activeSplitPane
			m.activeSplitPane = DeviceInfoSplitPaneType

			// Refresh device list when entering device mode
			devices, err := filesystem.ListDevices()
			if err == nil {
				m.lastDevices = devices
				m.applyDeviceSorting()
				m.updateDeviceListItems()
			}
		} else {
			ActiveTuiMode = PreviousTuiMode
			ActiveFileListMode = PreviousFileListMode
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

	// Move cursor down in file list (with optional count)
	case bindings.CursorDown.Matches(key):
		pane := m.GetActivePane()
		pane.fileList.CursorDown()
		m.UpdateFileInfoPane()
		return m, nil

	// Move cursor up in file list (with optional count)
	case bindings.CursorUp.Matches(key):
		pane := m.GetActivePane()
		pane.fileList.CursorUp()
		m.UpdateFileInfoPane()
		return m, nil

	// Copy file or folder
	case bindings.Copy.Matches(key):
		if ActiveTuiMode != ModeCopy {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeCopy

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Change file list to directories-only view
	case bindings.Directories.Matches(key):
		ActiveFileListMode = "ld"
		m.ApplyFileListFilter()
		return m, nil

	// Navigate into the selected directory.
	case bindings.Enter.Matches(key):
		pane := m.GetActivePane()
		selectedIdx := pane.fileList.Index()
		if selectedIdx >= 0 && selectedIdx < len(pane.files) {
			fi := pane.files[selectedIdx]
			if fi.IsDir {
				m.ChangeDirectory(fi.Path)
			}
		}
		return m, nil

	// Open file info split pane;
	case bindings.FileInfoPane.Matches(key):
		ActiveTuiMode = ModeNormal
		m.previousSplitPane = m.activeSplitPane
		m.activeSplitPane = FileInfoSplitPaneType
		m.isSplitPaneOpen = false
		m.activeViewport = LeftViewportType
		return m, nil

	// Change file list to files-only view
	case bindings.Files.Matches(key):
		ActiveFileListMode = "lf"
		m.ApplyFileListFilter()
		return m, nil

	// Enter filter mode
	case bindings.Filter.Matches(key):
		if ActiveTuiMode != ModeFilter {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeFilter

			pane := m.GetActivePane()
			m.searchInput.SetValue(pane.filterQuery)
			m.searchInput.Focus()
		}
		return m, nil

	// Enter Goto mode
	case bindings.Goto.Matches(key):
		if ActiveTuiMode != ModeGoto {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeGoto

			m.jumpTo = key
			m.commandInput.SetValue(key)
			m.commandInput.Focus()
			return m, nil
		}

		// Move cursor to end of file list
	case bindings.GoToEnd.Matches(keyMsg.String()):
		pane := m.GetActivePane()
		pane.fileList.GoToEnd()
		m.UpdateFileInfoPane()
		return m, nil

		// Move cursor to start of file list
	case bindings.GoToStart.Matches(keyMsg.String()):
		pane := m.GetActivePane()
		pane.fileList.GoToStart()
		m.UpdateFileInfoPane()
		return m, nil

		// Goto home directory
	case bindings.Home.Matches(keyMsg.String()):
		res, _ := m.ExecuteCommand("cd ~/")

		pane := m.GetActivePane()
		if res.Cwd != "" && res.Cwd != pane.currentDir {
			m.ChangeDirectory(res.Cwd)
		} else if res.Refresh {
			m.ReloadDirectory()
		}

		return m, nil

	// Make directory
	case bindings.Mkdir.Matches(key):
		if ActiveTuiMode != ModeMkdir {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeMkdir

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Move file or folder
	case bindings.Move.Matches(key):
		if ActiveTuiMode != ModeMove {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeMove

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Open help window
	case bindings.Help.Matches(key):
		if ActiveTuiMode != ModeHelp {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeHelp
			return m, nil
		}

	// Change file list to list-all-items view
	case bindings.List.Matches(key):
		ActiveFileListMode = "ll"
		m.ApplyFileListFilter()
		return m, nil

	// Navigate to the parent directory.
	case bindings.Parent.Matches(key):
		pane := m.GetActivePane()
		parent := filepath.Dir(pane.currentDir)
		if parent != "" && parent != pane.currentDir {
			m.ChangeDirectory(parent)
		} else {
			m.ReloadDirectory()
		}
		return m, nil

	// Navigate backwards/forwards through directory history.
	case bindings.PreviousDir.Matches(key):
		m.NavigatePreviousDir()
		return m, nil

	case bindings.NextDir.Matches(key):
		m.NavigateNextDir()
		return m, nil

	// Open preview split pane
	case bindings.PreviewPane.Matches(key):
		m.previousSplitPane = m.activeSplitPane
		m.activeSplitPane = PreviewPaneType
		m.isSplitPaneOpen = false
		return m, nil

	// Remove file or directory
	case bindings.Remove.Matches(key):
		if ActiveTuiMode != ModeRemove {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeRemove

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Move file or folder
	case bindings.Rename.Matches(key):
		if ActiveTuiMode != ModeRename {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeRename

			m.commandInput.SetValue("")
			m.commandInput.Focus()
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

	// Enter select mode and mark the current row.
	case bindings.Select.Matches(key):
		if ActiveTuiMode != ModeSelect {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeSelect
			(&m).toggleCurrentSelection()
		}
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

		// Open settings window
	case bindings.SidePanel.Matches(key):
		m.isSidePanelOpen = !m.isSidePanelOpen
		m.CalcLayout()
		return m, nil

		// Open sort window
	case bindings.Sort.Matches(key):
		if ActiveTuiMode != ModeSort {
			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeSort
			if ActiveFileListMode == FileListModeDevice {
				sortBy := m.GetSortDeviceColumnBy()
				if sortBy.column != "" {
					for i, col := range filesystem.DeviceInfoColumnNames {
						if col == sortBy.column {
							m.menuCursorIndex = i
							break
						}
					}
				}
			} else {
				sortBy := m.GetSortFileListColumnBy()
				if sortBy.column != "" {
					for i, col := range filesystem.FileInfoColumnNames {
						if col == sortBy.column {
							m.menuCursorIndex = i
							break
						}
					}
				}
			}
		} else {
			ActiveTuiMode = PreviousTuiMode
		}
		return m, nil

		// Open file list split pane
	case bindings.Tab.Matches(key) && !m.isSplitPaneOpen:
		if ActiveTuiMode != ModeFileListSplitPane {

			PreviousTuiMode = ActiveTuiMode
			ActiveTuiMode = ModeFileListSplitPane

			m.previousSplitPane = m.activeSplitPane
			m.activeSplitPane = FileListSplitPaneType
			m.isSplitPaneOpen = true
		}

		return m, nil

		// Switch panes in file list slipt mode
	case bindings.Tab.Matches(key):
		m.previousSplitPane = m.activeSplitPane
		m.activeSplitPane = FileListSplitPaneType
		m.isSplitPaneOpen = true

		if m.isSplitPaneOpen {
			if m.activeViewport == LeftViewportType {
				m.activeViewport = RightViewportType
			} else {
				m.activeViewport = LeftViewportType
			}
		}
		return m, nil

	case bindings.ToggleRightPane.Matches(key):
		m.showRightPane = !m.showRightPane
		m.CalcLayout()
		return m, nil

		// Page down in file list (with optional count)
	case bindings.PageDown.Matches(key):
		pane := m.GetActivePane()
		for i := 0; i < count; i++ {
			pane.fileList.NextPage()
		}
		m.UpdateFileInfoPane()
		return m, nil

	// Page up in file list (with optional count)
	case bindings.PageUp.Matches(key):
		pane := m.GetActivePane()
		for i := 0; i < count; i++ {
			pane.fileList.PrevPage()
		}
		m.UpdateFileInfoPane()
		return m, nil
	}

	return m, nil
}
