package components

import (
	"cute/filesystem"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ColumnModal(m tui.Model, args tui.ColumnModelArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	// Dialog-sized window
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	columnNames := filesystem.ColumnNames
	menuChoices := make([]string, len(columnNames))
	for i, col := range columnNames {
		menuChoices[i] = string(col)
	}

	menuCursor := m.GetMenuCursor()
	if menuCursor < 0 {
		menuCursor = 0
	}
	if menuCursor >= len(menuChoices) {
		menuCursor = len(menuChoices) - 1
	}

	var selectedColumns []filesystem.FileInfoColumn

	if tui.ActiveTuiMode == tui.ModeSort {
		sortBy := m.GetSortColumnBy()
		if sortByColumn := sortBy.Column(); sortByColumn != "" {
			selectedColumns = []filesystem.FileInfoColumn{sortByColumn}
		}
	} else {
		selectedColumns = m.GetColumnVisibility()
	}
	selectedMap := make(map[string]string, len(selectedColumns))
	for _, col := range selectedColumns {
		name := string(col)
		selectedMap[name] = name
	}

	menu := NewMenu(tui.MenuArgs{
		Choices:  menuChoices,
		Cursor:   menuCursor,
		Selected: selectedMap,
		CursorTypes: tui.MenuCursor{
			Selected:   args.Selected,
			Unselected: args.Unselected,
		},
	})

	fw := FloatingWindow{
		Content: menu,
		Width:   modalWidth,
		Height:  10,
		Style:   DefaultFloatingStyle(theme),
		Title:   args.Title,
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
