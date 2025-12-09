package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func CurrentDir(m tui.Model, args tui.CurrentDirComponentArgs) string {
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
