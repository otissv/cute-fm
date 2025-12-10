package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func SearchText(m tui.Model) string {
	return lipgloss.NewStyle().
		Render("> " + m.GetSearchInputText())
}
