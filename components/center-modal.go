package components

import (
	"charm.land/lipgloss/v2"
)

func CenterModal(modalContent string, width, height int) *lipgloss.Layer {
	dialogWidth := lipgloss.Width(modalContent)
	dialogHeight := lipgloss.Height(modalContent)
	x := 10
	y := 10
	if width > dialogWidth {
		x = (width - dialogWidth) / 2
	}
	if height > dialogHeight {
		y = (height - dialogHeight) / 2
	}

	return lipgloss.NewLayer(modalContent).X(x).Y(y)
}
