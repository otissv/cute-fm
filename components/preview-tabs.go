package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func PreviewTabs(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.SearchBar.Foreground)).
		Background(lipgloss.Color(theme.SearchBar.Background)).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color(theme.SearchBar.Border)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Height(args.Height).
		Width(args.Width).
		Render("Tabs")
}
