package components

import (
	"fmt"
	"strings"

	"cute/filesystem"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

type MenuModel struct {
	choices  []string
	cursor   int
	selected map[string]bool
}

func Menu(choices []string, cursor int, selected []filesystem.FileInfoColumn) MenuModel {
	selectedSet := make(map[string]bool, len(selected))
	for _, col := range selected {
		selectedSet[string(col)] = true
	}

	return MenuModel{
		choices:  choices,
		cursor:   cursor,
		selected: selectedSet,
	}
}

func (m MenuModel) View() string {
	var b strings.Builder
	for i, choice := range m.choices {
		iChoice := choice
		cursor := "[   ]"
		if m.selected != nil && m.selected[iChoice] {
			cursor = "[ x ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}

		if i == m.cursor {
			prompt := lipgloss.NewStyle().
				Bold(true).
				Render(">")

			cursor = "[ " + prompt + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, iChoice)
	}
	return strings.TrimRight(b.String(), "\n")
}

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
