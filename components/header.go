package components

import (
	"cute/theming"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func Header(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(theme.Header.Background)).
		PaddingBottom(1).
		Width(args.Width).
		Height(1).
		Render(theming.RainbowText(
			lipgloss.NewStyle().
				Background(lipgloss.Color(theme.Header.Background)),
			m.GetTitleText(),
			theming.Blends(theme.Primary, theme.Secondary),
		))
}
