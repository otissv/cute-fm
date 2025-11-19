package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Define styles for the layout
	// To change border color, modify the BorderForeground() value below
	// Options: color names ("blue", "red"), hex colors ("#874BFD"), or ANSI codes ("62")
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#636363")) // Current: purple/blue (ANSI code 62)

	// First row: Text input
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

	// Second row: Two viewports side by side
	// Calculate viewport width (half of available width, accounting for borders)
	viewportWidth := m.width / 2

	// Calculate viewport style height (content height + borders)
	// Text input row: Height(1) with borders = 3 total lines (1 content + 2 border)
	// Viewport style height = total height - text input row height
	textInputRowHeight := 3 // 1 content line + 2 border lines
	viewportStyleHeight := m.height - textInputRowHeight
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
	// First row: text input
	// Second row: two viewports
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		textInputView,
		secondRow, // No spacing string - borders will connect
	)

	return layout
}
