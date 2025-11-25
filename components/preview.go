package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func Preview(m tui.Model) string {
	theme := m.GetTheme()
	previewViewport := m.GetPreviewViewport()
	viewportHeight := m.GetViewportHeight()
	viewportWidth := m.GetViewportHeight()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Preview.Background)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color(theme.Preview.Border)).
		Foreground(lipgloss.Color(theme.Preview.Foreground)).
		Height(viewportHeight).
		Width(viewportWidth).
		Render(previewViewport.View())
}
