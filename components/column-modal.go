package components

import (
	"cute/filesystem"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ColumnModal(m tui.Model, args tui.ColumnModelArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	// Choose a dialog-sized window, not full-screen.
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	// Column identifiers are strongly typed so we don't pass around raw strings.
	// Convert the typed column identifiers into plain strings for the NewMenu.
	columnNames := filesystem.ColumnNames
	menuChoices := make([]string, len(columnNames))
	for i, col := range columnNames {
		menuChoices[i] = string(col)
	}

	// Use the cursor stored on the TUI model so navigation in ColumnVisibiliyMode
	// is reflected visually in the modal.
	menuCursor := m.GetMenuCursor()
	if menuCursor < 0 {
		menuCursor = 0
	}
	if menuCursor >= len(menuChoices) {
		menuCursor = len(menuChoices) - 1
	}

	// Pass the currently selected columns so the menu can display markers.
	// In column-visibility mode, this is the set of visible columns.
	// In sort mode, this is the single column currently used for sorting.
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
		Content: menu, // NewMenu implements tui.ViewPrimitive via its View method.
		Width:   modalWidth,
		Height:  10,
		Style:   DefaultFloatingStyle(theme),
		Title:   args.Title,
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
