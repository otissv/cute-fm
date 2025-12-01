package components

import (
	"charm.land/lipgloss/v2"

	"cute/tui"
)

func FileList(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()
	fileList := m.GetFileList()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileList.Background)).
		BorderBackground(lipgloss.Color(theme.FileList.Background)).
		BorderForeground(lipgloss.Color(theme.FileList.Border)).
		BorderStyle(lipgloss.NormalBorder()).
		Foreground(lipgloss.Color(theme.FileList.Foreground)).
		Height(args.Height).
		Width(args.Width).
		Render(fileList.View())
}
