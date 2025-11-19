package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Define styles for the layout.
	// Border color is driven by the theme loaded from lsfm.toml.
	borderColor := "#636363"
	if m.theme.BorderColor != "" {
		borderColor = m.theme.BorderColor
	}

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor))

	// Top row: current working directory (no border, just a single line)
	cwdStyle := lipgloss.NewStyle().
		Width(m.width)
	cwdText := m.currentDir
	if cwdText == "" {
		cwdText = "."
	}
	cwdView := cwdStyle.Render(cwdText)

	// Second row: Text input
	// Create a bordered container for the text input
	// To change the height, modify the Height() value below (e.g., Height(3) for taller)
	// Keep top, left, right borders; remove bottom border so it connects with viewports
	textInputStyle := borderStyle.
		Width(m.width).
		Height(1).
		BorderBottom(true)

	textInputView := textInputStyle.Render(
		m.textInput.View(),
	)

	// Third row: Two viewports side by side
	// Calculate viewport width (half of available width, accounting for borders)
	viewportWidth := m.width / 2

	// Calculate viewport style height (content height + borders)
	// CWD row: 1 content line
	// Text input row: Height(1) with borders = 3 total lines (1 content + 2 border)
	// Viewport style height = total height - (cwd row + text input row)
	cwdRowHeight := 1
	textInputRowHeight := 3 // 1 content line + 2 border lines
	viewportStyleHeight := m.height - (cwdRowHeight + textInputRowHeight)
	if viewportStyleHeight < 3 {
		viewportStyleHeight = 3 // Minimum height (1 content + 2 borders)
	}

	// Left viewport (second row, left column)
	// Keep top, left, bottom borders; remove right border so it connects with right viewport
	leftViewportStyle := borderStyle.
		Width(viewportWidth).
		Height(viewportStyleHeight). // Full height of the row (content + borders)
		BorderRight(true)

	leftViewportView := leftViewportStyle.Render(
		m.leftViewport.View(),
	)

	// Right viewport (second row, right column)
	// Keep top, right, bottom borders; remove left border so it connects with left viewport
	rightViewportStyle := borderStyle.
		Width(viewportWidth).
		Height(viewportStyleHeight). // Full height of the row (content + borders)
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		BorderLeft(false)

	rightViewportView := rightViewportStyle.Render(
		m.rightViewport.View(),
	)

	// Combine viewports horizontally without spacing (borders will connect)
	secondRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftViewportView,
		rightViewportView,
	)

	// Combine all rows vertically without spacing (borders will connect)
	// First row: current working directory
	// Second row: text input
	// Third row: two viewports
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		cwdView,
		textInputView,
		secondRow, // No spacing string - borders will connect
	)

	return layout
}
