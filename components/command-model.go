package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

// CommandBar renders the bottom command bar using only the public TUI model
// interface, so this component can live outside the tui package.
func CommandModal(m tui.Model, args tui.CommandModalArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	commandInput := m.GetCommandInput()

	title := ""
	if args.Title != "" {
		title = args.Title
	}
	commandInput.Prompt = ""
	if args.Prompt != "" {
		commandInput.Prompt = args.Prompt
	}
	commandInput.Placeholder = ""
	if args.Placeholder != "" {
		commandInput.Placeholder = args.Placeholder
	}

	commandInputView := commandInput.View()

	// Choose a dialog-sized window, not full-screen.
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	fw := FloatingWindow{
		Content: viewPrimitive(commandInputView),
		Width:   modalWidth,
		Height:  4,
		Style:   DefaultFloatingStyle(theme),
		Title:   title,
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}
