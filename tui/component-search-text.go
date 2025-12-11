package tui

import (
	"charm.land/lipgloss/v2"
)

func SearchText(m Model, view ActiveViewportType) string {
	return lipgloss.NewStyle().
		Render("> " + m.GetSearchInputTextForViewport(view))
}
