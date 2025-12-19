package tui

import (
	"charm.land/lipgloss/v2"
)

type DeviceListComponentArgs struct {
	Width  int
	Height int
}

func DeviceList(m Model, args DeviceListComponentArgs) string {
	theme := m.GetTheme()

	// Calculate content dimensions (subtract borders)
	contentWidth := args.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Update delegate with current width
	delegate := NewDeviceItemDelegate(theme, contentWidth, m.deviceColumns)
	m.deviceList.SetDelegate(delegate)
	m.deviceList.SetWidth(contentWidth)
	m.deviceList.SetHeight(args.Height - 3) // -3 for header and borders

	// Render header
	header := RenderDeviceHeaderRow(DeviceHeaderRowArgs{
		Theme:      theme,
		TotalWidth: contentWidth,
		columns:    m.deviceColumns,
	})

	body := m.deviceList.View()

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Background)).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color(theme.ActiveBorderColor)).
		BorderStyle(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color(theme.Foreground)).
		Height(args.Height).
		Width(args.Width).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	return baseStyle.
		Render(inner)
}
