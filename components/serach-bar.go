package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func SearchBar(m tui.Model) string {
	theme := m.GetTheme()
	width := m.GetViewportWidth()
	view := m.GetSearchInputView()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SearchBar.Foreground)).
		Background(lipgloss.Color(theme.SearchBar.Background)).
		BorderBackground(lipgloss.Color(theme.SearchBar.Background)).
		BorderForeground(lipgloss.Color(theme.SearchBar.Border)).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		Height(1).
		Width(width).
		Render(view)
}
