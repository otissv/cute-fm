package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func SortModal(m tui.Model) *lipgloss.Layer {
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

	fw := FloatingWindow{
		// Content: viewPrimitive(commandInputView),
		Content: viewPrimitive(""),
		Width:   modalWidth,
		Height:  4,
		Style:   DefaultFloatingStyle(theme),
		Title:   "Sort",
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
