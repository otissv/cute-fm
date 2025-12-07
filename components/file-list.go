package components

import (
	"charm.land/lipgloss/v2"

	"cute/tui"
)

func FileList(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()
	fileList := m.GetFileList()

	// Content width is viewport width minus left/right borders.
	contentWidth := args.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	header := tui.RenderFileHeaderRow(theme, contentWidth, m.GetColumnVisibility())
	body := fileList.View()

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileList.Background)).
		BorderBackground(lipgloss.Color(theme.FileList.Background)).
		BorderForeground(lipgloss.Color(theme.FileList.Border)).
		Foreground(lipgloss.Color(theme.FileList.Foreground)).
		Height(args.Height).
		Width(args.Width).
		Render(inner)
}
