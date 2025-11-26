package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

// textView is a minimal implementation of tui.ViewPrimitive that just renders
// the given string. This lets us render the command input directly without an
// extra viewport layer.
type textView string

func (t textView) View() string {
	return string(t)
}

// CommandBar renders the bottom command bar using only the public TUI model
// interface, so this component can live outside the tui package.
func CommandModal(m tui.Model) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()
	commandInputView := m.GetCommandInputView()

	// Choose a dialog-sized window, not full-screen.
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	fw := FloatingWindow{
		Content: textView(commandInputView),
		Width:   modalWidth,
		Height:  4,
		Style:   DefaultFloatingStyle(theme),
		Title:   "Command",
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
