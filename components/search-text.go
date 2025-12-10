package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

// SearchText renders the current filter text for a specific viewport (left or
// right). Each split pane can therefore display its own filter independently.
func SearchText(m tui.Model, view tui.ActiveViewportType) string {
	return lipgloss.NewStyle().
		Render("> " + m.GetSearchInputTextForViewport(view))
}
