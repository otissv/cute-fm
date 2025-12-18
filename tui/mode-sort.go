package tui

import (
	"cute/filesystem"

	tea "charm.land/bubbletea/v2"
)

func (m Model) SortMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	// Move cursor up within the column list.
	case bindings.CursorUp.Matches(keyMsg.String()):
		if m.menuCursorIndex > 0 {
			m.menuCursorIndex--
		}
		return m, nil

	// Move cursor down within the column list.
	case bindings.CursorDown.Matches(keyMsg.String()):
		var maxIdx int
		if ActiveFileListMode == FileListModeDevice {
			maxIdx = len(filesystem.DeviceInfoColumnNames) - 1
		} else {
			maxIdx = len(filesystem.FileInfoColumnNames) - 1
		}
		if m.menuCursorIndex < maxIdx {
			m.menuCursorIndex++
		}
		return m, nil

	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Apply sorting based on the currently focused column.
	case bindings.Select.Matches(keyMsg.String()) || bindings.Enter.Matches(keyMsg.String()):

		if ActiveFileListMode == FileListModeDevice {
			if len(filesystem.DeviceInfoColumnNames) == 0 {
				return m, nil
			}

			cur := m.menuCursorIndex
			if cur < 0 {
				cur = 0
			}
			if cur >= len(filesystem.DeviceInfoColumnNames) {
				cur = len(filesystem.DeviceInfoColumnNames) - 1
			}

			col := filesystem.DeviceInfoColumnNames[cur]

			// Toggle sort direction when selecting the same column
			if m.sortDeviceColumnBy.column == col {
				if m.sortDeviceColumnBy.direction == SortingAsc {
					m.sortDeviceColumnBy.direction = SortingDesc
				} else {
					m.sortDeviceColumnBy.direction = SortingAsc
				}
			} else {
				m.sortDeviceColumnBy.column = col
				m.sortDeviceColumnBy.direction = SortingAsc
			}

			m.ApplyDeviceSorting()
		} else {
			if len(filesystem.FileInfoColumnNames) == 0 {
				return m, nil
			}

			cur := m.menuCursorIndex
			if cur < 0 {
				cur = 0
			}
			if cur >= len(filesystem.FileInfoColumnNames) {
				cur = len(filesystem.FileInfoColumnNames) - 1
			}

			col := filesystem.FileInfoColumnNames[cur]

			// Toggle sort direction when selecting the same column
			if m.sortFileListColumnBy.column == col {
				if m.sortFileListColumnBy.direction == SortingAsc {
					m.sortFileListColumnBy.direction = SortingDesc
				} else {
					m.sortFileListColumnBy.direction = SortingAsc
				}
			} else {
				m.sortFileListColumnBy.column = col
				m.sortFileListColumnBy.direction = SortingAsc
			}

			m.ApplyFileListFilter()
		}

		return m, nil

	case bindings.Cancel.Matches(keyMsg.String()):
		if ActiveFileListMode == FileListModeDevice {
			m.settings.SortDeviceColumnBy = filesystem.FileInfoColumn(m.sortDeviceColumnBy.column)
			m.settings.SortDeviceColumnDirection = m.sortDeviceColumnBy.direction
		} else {
			m.settings.SortFileListColumnBy = m.sortFileListColumnBy.column
			m.settings.SortFileListColumnDirection = m.sortFileListColumnBy.direction
		}

		if err := SaveSettings(m); err != nil {
			_ = err
		}

		ActiveTuiMode = PreviousTuiMode
		m.menuCursorIndex = 0
		return m, nil

	case bindings.ColumnVisibilityFileList.Matches(keyMsg.String()):
		ActiveTuiMode = ModeColumnVisibility
		return m, nil
	}

	return m, nil
}
