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
	case bindings.Up.Matches(keyMsg.String()):
		if m.menuCursor > 0 {
			m.menuCursor--
		}
		return m, nil

	// Move cursor down within the column list.
	case bindings.Down.Matches(keyMsg.String()):
		maxIdx := len(filesystem.ColumnNames) - 1
		if m.menuCursor < maxIdx {
			m.menuCursor++
		}
		return m, nil

	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Apply sorting based on the currently focused column.
	case bindings.Select.Matches(keyMsg.String()):
		if len(filesystem.ColumnNames) == 0 {
			return m, nil
		}

		cur := m.menuCursor
		if cur < 0 {
			cur = 0
		}
		if cur >= len(filesystem.ColumnNames) {
			cur = len(filesystem.ColumnNames) - 1
		}

		col := filesystem.ColumnNames[cur]

		// Toggle sort direction when selecting the same column; otherwise start
		// with ascending order for a newly selected column.
		if m.sortColumnBy.column == col {
			if m.sortColumnBy.direction == SortingAsc {
				m.sortColumnBy.direction = SortingDesc
			} else {
				m.sortColumnBy.direction = SortingAsc
			}
		} else {
			m.sortColumnBy.column = col
			m.sortColumnBy.direction = SortingAsc
		}

		m.ApplyFilter()
		m.menuCursor = 0

		return m, nil

	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		m.menuCursor = 0
		return m, nil

	case bindings.ColumnVisibiliy.Matches(keyMsg.String()):
		ActiveTuiMode = ModeColumnVisibiliy
		return m, nil
	}

	return m, nil
}
