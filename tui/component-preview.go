package tui

import (
	"charm.land/lipgloss/v2"
)

func FileInfo(m Model, args ComponentArgs) string {
	theme := m.GetTheme()
	previewViewport := m.GetPreviewViewport()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileInfo.Background)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBackground(lipgloss.Color(theme.FileInfo.BorderBackground)).
		BorderForeground(lipgloss.Color(theme.FileInfo.Border)).
		Foreground(lipgloss.Color(theme.FileInfo.Foreground)).
		Height(args.Height).
		Width(args.Width).
		Render(previewViewport.View())
}
