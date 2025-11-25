package components

import (
	"strings"

	"cute/tui"

	"charm.land/bubbles/v2/viewport"
)

// HelpModal renders the help dialog as a floating window using information
// exposed by the public TUI model.
func HelpModal(m tui.Model) string {
	width, height := m.GetSize()
	theme := m.GetTheme()

	helpContent := `
	Help
	----
	
	Navigation:
		Up/Down arrows   Move selection in file list
		Scroll wheel     Scroll file list
	
	Search:
		Type in the search bar to filter files by name
	
	General:
		?                Toggle this help
		ctrl+c / ctrl+q  Quit
	`

	helpViewport := viewport.New()
	helpViewport.SetContent(strings.TrimSpace(helpContent))

	// Choose a dialog-sized window, not full-screen.
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	modalHeight := height / 2
	if modalHeight > 16 {
		modalHeight = 16
	}
	if modalHeight < 6 {
		modalHeight = 6
	}

	fw := FloatingWindow{
		Content: helpViewport,
		Width:   modalWidth,
		Height:  modalHeight,
		Style:   DefaultFloatingStyle(theme),
	}

	// Return just the dialog view; compositing with the base layout is handled
	// at the View() layer via Lip Gloss layers/canvas.
	return fw.View(width, height)
}
