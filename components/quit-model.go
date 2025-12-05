package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

// CommandBar renders the bottom command bar using only the public TUI model
// interface, so this component can live outside the tui package.
func QuitModal(m tui.Model) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()
	content := "Press q to quit\n\nor\n\n press ESC to cancel"

	fw := FloatingWindow{
		Content: textView(content),
		Width:   40,
		Height:  6,
		Style:   DefaultFloatingStyle(theme).Align(lipgloss.Center),
		Title:   "Quit",
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
