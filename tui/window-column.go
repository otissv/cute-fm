package tui

import (
	"cute/filesystem"
	"cute/utils"

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

	var columnNames []string

	if ActiveFileListMode == FileListModeComputer {
		columnNames = utils.ToStringSlice(filesystem.DeviceInfoColumnNames)
	} else {
		columnNames = utils.ToStringSlice(filesystem.FileInfoColumnNames)
	}

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

	var selectedColumns []string

	// var selectedColumns []filesystem.FileInfoColumn

	if ActiveTuiMode == ModeSort {
		sortBy := m.GetSortColumnBy()
		if sortByColumn := sortBy.Column(); sortByColumn != "" {
			selectedColumns = []string{string(sortByColumn)}
		} else if m.settings.SortFileListColumnBy != "" {
			// Fall back to settings if current sort column is not set
			selectedColumns = []string{string(m.settings.SortFileListColumnBy)}
		}
	} else {
		if ActiveFileListMode == FileListModeComputer {
			selectedColumns = utils.ToStringSlice(m.deviceColumns)
		} else {
			selectedColumns = utils.ToStringSlice(m.GetFileListColumnVisibility())
		}
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

	style := DefaultFloatingStyle(theme)

	if ActiveTuiMode == ModeColumnVisibility || ActiveTuiMode == ModeSort {
		style = style.BorderForeground(lipgloss.Color(theme.ActiveBorderColor))
	}

	fw := FloatingWindow{
		Content: menu,
		Width:   windowWidth,
		Height:  10,
		Style:   style,
		Title:   args.Title,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
