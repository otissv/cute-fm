package tui

import (
	"charm.land/lipgloss/v2"
)

type ComputerListComponentArgs struct {
	Width  int
	Height int
}

func ComputerList(m Model, args ComputerListComponentArgs) string {
	theme := m.GetTheme()

	// Calculate content dimensions (subtract borders)
	contentWidth := args.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Update delegate with current width
	delegate := NewDeviceItemDelegate(theme, contentWidth, m.deviceColumns)
	m.computerList.SetDelegate(delegate)
	m.computerList.SetWidth(contentWidth)
	m.computerList.SetHeight(args.Height - 3) // -3 for header and borders

	// Render header
	header := RenderDeviceHeaderRow(DeviceHeaderRowArgs{
		Theme:      theme,
		TotalWidth: contentWidth,
		columns:    m.deviceColumns,
	})

	body := m.computerList.View()

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileList.Background)).
		BorderBackground(lipgloss.Color(theme.FileList.Background)).
		BorderForeground(lipgloss.Color("#1E1E1E")).
		BorderStyle(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color(theme.FileList.Foreground)).
		Height(args.Height).
		Width(args.Width).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	return baseStyle.
		Render(inner)
}
