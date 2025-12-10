package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func DialogModal(m tui.Model, args tui.DialogModalArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	fw := FloatingWindow{
		Content: viewPrimitive(args.Content),
		Width:   40,
		Height:  6,
		Style:   DefaultFloatingStyle(theme).Align(lipgloss.Center),
		Title:   args.Title,
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
