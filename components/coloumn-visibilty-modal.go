package components

import (
	"cute/filesystem"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ColoumnVisibiltyModal(m tui.Model) *lipgloss.Layer {
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
	// Convert the typed column identifiers into plain strings for the Menu.
	columnNames := filesystem.ColumnNames
	menuChoices := make([]string, len(columnNames))
	for i, col := range columnNames {
		menuChoices[i] = string(col)
	}

	// Use the cursor stored on the TUI model so navigation in ColumnVisibiliyMode
	// is reflected visually in the modal.
	menuCursor := m.GetColumnVisibilityCursor()
	if menuCursor < 0 {
		menuCursor = 0
	}
	if menuCursor >= len(menuChoices) {
		menuCursor = len(menuChoices) - 1
	}

	// Pass the currently selected columns so the menu can display [x] markers.
	menu := Menu(menuChoices, menuCursor, m.GetColumnVisibility())

	fw := FloatingWindow{
		Content: menu, // Menu implements tui.ViewPrimitive via its View method.
		Width:   modalWidth,
		Height:  10,
		Style:   DefaultFloatingStyle(theme),
		Title:   "Column Visibilty",
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
