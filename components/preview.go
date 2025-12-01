package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func Preview(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()
	previewViewport := m.GetPreviewViewport()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Preview.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBackground(lipgloss.Color(theme.Preview.BorderBackground)).
		BorderForeground(lipgloss.Color(theme.Preview.Border)).
		Foreground(lipgloss.Color(theme.Preview.Foreground)).
		Height(args.Height).
		Width(args.Width).
		Render(previewViewport.View())
}
