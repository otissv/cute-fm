package tui

import (
	"charm.land/lipgloss/v2"
)

func SidePanel(m Model) string {
	theme := m.GetTheme()

	alignLeft := lipgloss.Left

	content := lipgloss.JoinVertical(lipgloss.Left, NewButton(m, ButtonArgs{
		label: "Computer",
		width: 20,
		align: &alignLeft,
	}).View())

	return lipgloss.NewStyle().
		Width(20).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		Padding(0, 1).
		Render(content)
}
