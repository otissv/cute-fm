package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func StatusBar(m tui.Model, items ...string) string {
	theme := m.GetTheme()
	width, _ := m.GetSize()

	statusStyle := lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(theme.StatusBar.Background)).
		PaddingBottom(theme.StatusBar.PaddingBottom).
		PaddingLeft(theme.StatusBar.PaddingLeft).
		PaddingRight(theme.StatusBar.PaddingRight).
		PaddingTop(theme.Preview.PaddingTop).
		Width(width)

	var flatItems []string
	flatItems = append(flatItems, items...)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, flatItems...)
	return statusStyle.Render(statusBar)
}
