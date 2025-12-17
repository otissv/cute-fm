package tui

import (
	"cute/filesystem"

	"charm.land/lipgloss/v2"
)

type ColumnWindowArgs struct {
	Title      string
	Selected   string
	Unselected string
	Prompt     string
}

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

	columnNames := filesystem.FileInfoColumnNames
	menuChoices := make([]MenuChoice, len(columnNames))
	for i, col := range columnNames {
		menuChoices[i] = MenuChoice{
			Label: string(col),
			Type:  CHOICE_TYPE,
		}
	}

	menuCursorIndex := m.GetMenuCursorIndex()
	if menuCursorIndex < 0 {
		menuCursorIndex = 0
	}
	if menuCursorIndex >= len(menuChoices) {
		menuCursorIndex = len(menuChoices) - 1
	}

	var selectedColumns []filesystem.FileInfoColumn

	if ActiveTuiMode == ModeSort {
		sortBy := m.GetSortColumnBy()
		if sortByColumn := sortBy.Column(); sortByColumn != "" {
			selectedColumns = []filesystem.FileInfoColumn{sortByColumn}
		} else if m.settings.SortColumnBy != "" {
			// Fall back to settings if current sort column is not set
			selectedColumns = []filesystem.FileInfoColumn{m.settings.SortColumnBy}
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
		Choices:     menuChoices,
		CursorIndex: menuCursorIndex,
		CursorTypes: MenuCursor{
			Selected:   args.Selected,
			Unselected: args.Unselected,
		},
		Selected: selectedMap,
		Theme:    m.theme,
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
