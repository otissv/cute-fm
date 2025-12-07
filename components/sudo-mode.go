package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func SudoMode(m tui.Model, args tui.ComponentArgs) string {
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
