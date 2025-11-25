package components

import (
	"charm.land/lipgloss/v2"

	"cute/tui"
)

func FileList(m tui.Model) string {
	theme := m.GetTheme()
	viewportHeight := m.GetViewportHeight()
	viewportWidth := m.GetViewportHeight()
	fileListViewport := m.GetFileListViewport()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileList.Background)).
		BorderBackground(lipgloss.Color(theme.FileList.Background)).
		BorderForeground(lipgloss.Color(theme.FileList.Border)).
		BorderStyle(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color(theme.FileList.Foreground)).
		Height(viewportHeight).
		Width(viewportWidth).
		Render(fileListViewport.View())
}
