package components

import (
	"strings"

	"cute/tui"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

// HelpModal renders the help dialog as a floating window using information
// exposed by the public TUI model.
func HelpModal(m tui.Model) *lipgloss.Layer {
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
		Title:   "HELP",
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}

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
