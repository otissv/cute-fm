package tui

import "charm.land/lipgloss/v2"

type SudoModeComponentArgs struct {
	Width  int
	Height int
}

func SudoMode(m Model, args SudoModeComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Background(lipgloss.Color(theme.SudoMode.Background)).
		Foreground(lipgloss.Color(theme.SudoMode.Foreground)).
		PaddingBottom(0).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(0).
		Height(args.Height).
		Render(string("SUDO"))
}
