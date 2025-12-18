package tui

import (
	"cute/filesystem"
	"cute/utils"

	tea "charm.land/bubbletea/v2"
)

func (m Model) ColumnVisibilityMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	var columnNames []string
	var columns []string

	pane := m.GetActivePane()

	if ActiveFileListMode == FileListModeComputer {
		columnNames = utils.ToStringSlice(filesystem.DeviceInfoColumnNames)
		columns = utils.ToStringSlice(m.deviceColumns)
	} else {
		columnNames = utils.ToStringSlice(filesystem.FileInfoColumnNames)
		columns = utils.ToStringSlice(pane.columns)
	}

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Move cursor up within the column list.
	case bindings.CursorUp.Matches(keyMsg.String()):
		if m.menuCursorIndex > 0 {
			m.menuCursorIndex--
		}
		return m, nil

	// Move cursor down within the column list.
	case bindings.CursorDown.Matches(keyMsg.String()):
		maxIdx := len(columnNames) - 1
		if m.menuCursorIndex < maxIdx {
			m.menuCursorIndex++
		}
		return m, nil

	// Toggle the currently focused column and stay in this mode.
	case bindings.Select.Matches(keyMsg.String()):
		if len(columnNames) == 0 {
			return m, nil
		}

		cur := m.menuCursorIndex
		if cur < 0 {
			cur = 0
		}
		if cur >= len(columnNames) {
			cur = len(columnNames) - 1
		}

		col := columnNames[cur]

		// Toggle presence of col in the columnVisibility set, but always rebuild
		visible := make(map[string]bool, len(columnNames))
		for _, c := range columns {
			visible[string(c)] = true
		}

		if visible[col] {
			delete(visible, col)
		} else {
			visible[col] = true
		}

		if len(visible) == 0 {
			return m, nil
		}

		// Rebuild in canonical order.
		if ActiveFileListMode == FileListModeComputer {
			newCols := make([]filesystem.DeviceInfoColumn, 0, len(visible))
			for _, c := range columnNames {
				if visible[c] {
					newCols = append(newCols, filesystem.DeviceInfoColumn(c))
				}
			}
			m.deviceColumns = newCols
		} else {
			newCols := make([]filesystem.FileInfoColumn, 0, len(visible))
			for _, c := range columnNames {
				if visible[c] {
					newCols = append(newCols, filesystem.FileInfoColumn(c))
				}
			}
			pane.columns = newCols

			// Rebuild the file list delegate so the visible columns update
			listContentWidth := m.viewportWidth - 2
			if listContentWidth < 1 {
				listContentWidth = 1
			}
			m.UpdateFileListDelegate(listContentWidth)
		}

		return m, nil

	// Enter normal mode
	case bindings.Cancel.Matches(keyMsg.String()) || bindings.Select.Matches(keyMsg.String()):
		pane := m.GetActivePane()
		m.settings.ColumnVisibilityFileList = pane.columns

		if err := SaveSettings(m); err != nil {
			_ = err
		}

		ActiveTuiMode = ModeNormal
		m.menuCursorIndex = 0
		return m, nil

	case bindings.Sort.Matches(keyMsg.String()):
		ActiveTuiMode = ModeSort
		return m, nil
	}

	return m, nil
}
