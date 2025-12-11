package tui

import (
	"charm.land/lipgloss/v2"
)

func CenterWindow(windowContent string, width, height int) *lipgloss.Layer {
	dialogWidth := lipgloss.Width(windowContent)
	dialogHeight := lipgloss.Height(windowContent)
	x := 10
	y := 10
	if width > dialogWidth {
		x = (width - dialogWidth) / 2
	}
	if height > dialogHeight {
		y = (height - dialogHeight) / 2
	}

	return lipgloss.NewLayer(windowContent).X(x).Y(y)
}
