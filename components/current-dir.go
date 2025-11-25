package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func CurrentDir(m tui.Model) string {
	theme := m.GetTheme()
	currentDir := m.GetCurrentDir()

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Background(lipgloss.Color(theme.CurrentDir.Background)).
		Foreground(lipgloss.Color(theme.CurrentDir.Foreground)).
		MarginRight(2).
		PaddingBottom(1).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(0).
		Render(currentDir)
}
