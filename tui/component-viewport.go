package tui

import (
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

type ViewportComponentArgs struct {
	Height   int
	Width    int
	Viewport viewport.Model
}

func ViewPort(m Model, args ViewportComponentArgs) string {
	theme := m.GetTheme()

	style := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		Foreground(lipgloss.Color(theme.Foreground)).
		Height(args.Height).
		Width(args.Width)

	return style.Render(args.Viewport.View())
}
