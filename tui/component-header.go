package tui

import (
	"cute/theming"

	"charm.land/lipgloss/v2"
)

func Header(m Model, args ComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Right).
		Background(lipgloss.Color(theme.Header.Background)).
		Height(args.Height).
		PaddingBottom(0).
		PaddingRight(1).
		PaddingTop(0).
		Width(args.Width).
		Render(theming.RainbowText(
			lipgloss.NewStyle().
				Background(lipgloss.Color(theme.Header.Background)),
			m.GetTitleText(),
			theming.Blends(theme.Primary, theme.Secondary),
		))
}
