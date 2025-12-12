package tui

import "charm.land/lipgloss/v2"

type CurrentDirComponentArgs struct {
	Width      int
	Height     int
	CurrentDir string
}

func CurrentDir(m Model, args CurrentDirComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Background(lipgloss.Color(theme.CurrentDir.Background)).
		Foreground(lipgloss.Color(theme.CurrentDir.Foreground)).
		PaddingBottom(0).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(0).
		Height(args.Height).
		Render(args.CurrentDir)
}
