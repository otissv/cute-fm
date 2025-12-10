package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func SearchText(m tui.Model, view tui.ActiveViewportType) string {
	return lipgloss.NewStyle().
		Render("> " + m.GetSearchInputTextForViewport(view))
}
