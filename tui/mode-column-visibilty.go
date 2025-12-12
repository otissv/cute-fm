package tui

import (
	"cute/filesystem"

	tea "charm.land/bubbletea/v2"
)

func (m Model) ColumnVisibilityMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	// Move cursor up within the column list.
	case bindings.Up.Matches(keyMsg.String()):
		if m.menuCursorIndex > 0 {
			m.menuCursorIndex--
		}
		return m, nil

	// Move cursor down within the column list.
	case bindings.Down.Matches(keyMsg.String()):
		maxIdx := len(filesystem.ColumnNames) - 1
		if m.menuCursorIndex < maxIdx {
			m.menuCursorIndex++
		}
		return m, nil

	// Toggle the currently focused column and stay in this mode.
	case bindings.Select.Matches(keyMsg.String()):
		if len(filesystem.ColumnNames) == 0 {
			return m, nil
		}

		cur := m.menuCursorIndex
		if cur < 0 {
			cur = 0
		}
		if cur >= len(filesystem.ColumnNames) {
			cur = len(filesystem.ColumnNames) - 1
		}

		col := filesystem.ColumnNames[cur]

		// Work on the active pane's column visibility.
		pane := m.GetActivePane()

		// Toggle presence of col in the columnVisibility set, but always rebuild
		visible := make(map[filesystem.FileInfoColumn]bool, len(filesystem.ColumnNames))
		for _, c := range pane.columns {
			visible[c] = true
		}

		if visible[col] {
			delete(visible, col)
		} else {
			visible[col] = true
		}

		// Rebuild in canonical order.
		newCols := make([]filesystem.FileInfoColumn, 0, len(visible))
		for _, c := range filesystem.ColumnNames {
			if visible[c] {
				newCols = append(newCols, c)
			}
		}
		pane.columns = newCols

		// Rebuild the file list delegate so the visible columns update
		// immediately to reflect the new selection.
		listContentWidth := m.viewportWidth - 2
		if listContentWidth < 1 {
			listContentWidth = 1
		}
		m.UpdateFileListDelegate(listContentWidth)

		return m, nil

	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil
	// Enter normal mode
	case bindings.Select.Matches(keyMsg.String()) ||
		bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		m.menuCursorIndex = 0
		return m, nil

	case bindings.Sort.Matches(keyMsg.String()):
		ActiveTuiMode = ModeSort
		return m, nil
	}

	return m, nil
}
