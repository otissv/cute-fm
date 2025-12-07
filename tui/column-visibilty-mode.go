package tui

import (
	"cute/filesystem"

	tea "charm.land/bubbletea/v2"
)

func (m Model) ColumnVisibiliyMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	// Move cursor up within the column list.
	case bindings.Up.Matches(keyMsg.String()):
		if m.columnVisibilityCursor > 0 {
			m.columnVisibilityCursor--
		}
		return m, nil

	// Move cursor down within the column list.
	case bindings.Down.Matches(keyMsg.String()):
		maxIdx := len(filesystem.ColumnNames) - 1
		if m.columnVisibilityCursor < maxIdx {
			m.columnVisibilityCursor++
		}
		return m, nil

	// Toggle the currently focused column and stay in this mode.
	case bindings.Select.Matches(keyMsg.String()):
		if len(filesystem.ColumnNames) == 0 {
			return m, nil
		}

		cur := m.columnVisibilityCursor
		if cur < 0 {
			cur = 0
		}
		if cur >= len(filesystem.ColumnNames) {
			cur = len(filesystem.ColumnNames) - 1
		}

		col := filesystem.ColumnNames[cur]

		// Toggle presence of col in the columnVisibility slice.
		found := false
		newCols := make([]filesystem.FileInfoColumn, 0, len(m.columnVisibility))
		for _, c := range m.columnVisibility {
			if c == col {
				found = true
				continue // drop to "unselect"
			}
			newCols = append(newCols, c)
		}
		if !found {
			newCols = append(newCols, col)
		}
		m.columnVisibility = newCols

		return m, nil

	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil
	// Enter normal mode
	case bindings.Select.Matches(keyMsg.String()) ||
		bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = TuiModeNormal
		return m, nil
	}

	return m, nil
}
