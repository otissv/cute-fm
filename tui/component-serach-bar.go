package tui

import (
	"charm.land/lipgloss/v2"
)

func SearchBar(m Model, args ComponentArgs) string {
	theme := m.GetTheme()
	view := m.GetSearchInputView()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SearchBar.Foreground)).
		Background(lipgloss.Color(theme.SearchBar.Background)).
		BorderBackground(lipgloss.Color(theme.SearchBar.Background)).
		BorderForeground(lipgloss.Color(theme.SearchBar.Border)).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(lipgloss.NormalBorder()).
		Height(args.Height).
		Width(args.Width).
		Render(view)
}
