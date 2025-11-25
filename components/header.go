package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func Header(m tui.Model) string {
	theme := m.GetTheme()
	width, _ := m.GetSize()

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(theme.Header.Background)).
		PaddingBottom(1).
		Width(width).
		Render(tui.RainbowText(
			lipgloss.NewStyle().
				Background(lipgloss.Color(theme.Header.Background)),
			m.GetTitleText(),
			tui.Blends(theme.Primary, theme.Secondary),
		))
}
