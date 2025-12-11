package tui

import (
	"cute/filesystem"

	"charm.land/lipgloss/v2"
)

func ColumnWindow(m Model, args ColumnWindowArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	// Dialog-sized window
	windowWidth := width / 2
	if windowWidth > 60 {
		windowWidth = 60
	}
	if windowWidth < 30 {
		windowWidth = 30
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

	if ActiveTuiMode == ModeSort {
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

	menu := NewMenu(MenuArgs{
		Choices:  menuChoices,
		Cursor:   menuCursor,
		Selected: selectedMap,
		CursorTypes: MenuCursor{
			Selected:   args.Selected,
			Unselected: args.Unselected,
		},
	})

	fw := FloatingWindow{
		Content: menu,
		Width:   windowWidth,
		Height:  10,
		Style:   DefaultFloatingStyle(theme),
		Title:   args.Title,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
