package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Define styles for the layout.
	// Border color is driven by the theme loaded from lsfm.toml.
	borderColor := ""
	if m.theme.BorderColor != "" {
		borderColor = m.theme.BorderColor
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor))

	// Search Text input
	textInputStyle := borderStyle.
		Width(m.width).
		Height(1).
		BorderBottom(true)

	textInputView := textInputStyle.Render(
		m.textInput.View(),
	)

	// Text input row: Height(1) with borders = 3 total lines (1 content + 2 border)
	textInputRowHeight := 3 // 1 content line + 2 border lines

	// Status row at the bottom: 1 content line
	statusRowHeight := 1
	// Viewport style height = total height - (text input row + status row)

	// Calculate viewport width (half of available width, accounting for borders)
	viewportWidth := m.width / 2

	// Calculate viewport style height (content height + borders)
	viewportStyleHeight := m.height - (textInputRowHeight + statusRowHeight)
	if viewportStyleHeight < 3 {
		viewportStyleHeight = 3 // Minimum height (1 content + 2 borders)
	}

	FileListViewportStyle := borderStyle.
		Width(viewportWidth).
		Height(viewportStyleHeight).
		Background(lipgloss.Color(m.theme.FileList.Background)).
		Foreground(lipgloss.Color(m.theme.FileList.Foreground)).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		BorderRight(true).
		PaddingTop(m.theme.FileList.PaddingTop).
		PaddingBottom(m.theme.FileList.PaddingBottom).
		PaddingLeft(m.theme.FileList.PaddingLeft).
		PaddingRight(m.theme.FileList.PaddingRight)

	FileListViewportView := FileListViewportStyle.Render(
		m.FileListViewport.View(),
	)

	previewViewportStyle := borderStyle.
		Width(viewportWidth).
		Height(viewportStyleHeight).
		Background(lipgloss.Color(m.theme.Preview.Background)).
		Foreground(lipgloss.Color(m.theme.Preview.Foreground)).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		BorderLeft(false).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.Preview.PaddingBottom).
		PaddingLeft(m.theme.Preview.PaddingLeft).
		PaddingRight(m.theme.Preview.PaddingRight)

	previewViewportView := previewViewportStyle.Render(
		m.previewViewport.View(),
	)

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Left,
		FileListViewportView,
		previewViewportView,
	)

	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(m.theme.StatusBar.Background)).
		Foreground(lipgloss.Color(m.theme.StatusBar.Foreground)).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.StatusBar.PaddingBottom).
		PaddingLeft(m.theme.StatusBar.PaddingLeft).
		PaddingRight(m.theme.StatusBar.PaddingRight)

	statusText := m.currentDir
	if statusText == "" {
		statusText = "."
	}
	statusView := statusStyle.Render(statusText)

	// Combine all rows vertically
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		textInputView,
		viewports,
		statusView,
	)

	// If no modal is active, show the base layout.
	if m.activeModal == ModalNone {
		return layout
	}

	// Decide which content to show inside the floating window.
	var content ViewPrimitive
	switch m.activeModal {
	case ModalHelp:
		content = m.helpViewport
	}

	// Choose a dialog-sized window, not full-screen.
	modalWidth := m.width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	modalHeight := m.height / 2
	if modalHeight > 16 {
		modalHeight = 16
	}
	if modalHeight < 6 {
		modalHeight = 6
	}

	fw := FloatingWindow{
		Content: content,
		Width:   modalWidth,
		Height:  modalHeight,
		Style:   DefaultFloatingStyle(),
	}

	dialogView := fw.View(m.width, m.height)
	return overlayDialog(layout, dialogView)
}

// overlayDialog composes the base layout and the dialog view.
func overlayDialog(layout, dialog string) string {
	baseLines := strings.Split(layout, "\n")
	dialogLines := strings.Split(dialog, "\n")

	maxLines := len(baseLines)
	if len(dialogLines) > maxLines {
		maxLines = len(dialogLines)
	}

	out := make([]string, maxLines)

	for i := 0; i < maxLines; i++ {
		var baseLine, dlgLine string
		if i < len(baseLines) {
			baseLine = baseLines[i]
		}
		if i < len(dialogLines) {
			dlgLine = dialogLines[i]
		}

		if strings.TrimSpace(dlgLine) != "" {
			out[i] = dlgLine
		} else {
			out[i] = baseLine
		}
	}

	return strings.Join(out, "\n")
}
