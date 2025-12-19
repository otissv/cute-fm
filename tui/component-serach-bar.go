package tui

import (
	"charm.land/lipgloss/v2"
)

type SearchBarComponentArgs struct {
	Width  int
	Height int
}

func SearchBar(m Model, args SearchBarComponentArgs) string {
	theme := m.GetTheme()
	view := m.GetSearchInputView()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Background(lipgloss.Color(theme.Background)).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(lipgloss.NormalBorder()).
		Height(args.Height).
		Width(args.Width).
		Render(view)
}
