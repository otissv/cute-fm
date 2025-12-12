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
		Background(lipgloss.Color(theme.StatusBar.Background)).
		Height(args.Height).
		PaddingBottom(theme.StatusBar.PaddingBottom).
		PaddingLeft(theme.StatusBar.PaddingLeft).
		PaddingRight(theme.StatusBar.PaddingRight).
		PaddingTop(theme.StatusBar.PaddingTop).
		Width(args.Width)

	var flatItems []string
	flatItems = append(flatItems, items...)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, flatItems...)
	return statusStyle.Render(statusBar)
}
