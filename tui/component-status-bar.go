package tui

import "charm.land/lipgloss/v2"

type StatusBarComponentArgs struct {
	Width  int
	Height int
}

func StatusBar(m Model, args StatusBarComponentArgs, items ...string) string {
	theme := m.GetTheme()

	statusStyle := lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(theme.Background)).
		Foreground(lipgloss.Color(theme.Foreground)).
		Height(args.Height).
		Width(args.Width)

	var flatItems []string
	flatItems = append(flatItems, items...)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, flatItems...)
	return statusStyle.Render(statusBar)
}
