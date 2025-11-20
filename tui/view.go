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

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.theme.BorderColor))

	searchBarStyle := borderStyle.
		Width(m.width).
		Height(1).
		BorderBottom(true)

	searchBar := searchBarStyle.Render(
		m.searchBar.View(),
	)

	// Calculate viewport width (half of available width, accounting for borders)
	searchBarRowHeight := 3
	statusRowHeight := 1
	// Only reserve vertical space for the command bar when it is visible.
	commandRowHeight := 0
	if m.commandMode {
		commandRowHeight = 3
	}
	viewportWidth := m.width / 2
	viewportStyleHeight := m.height - (searchBarRowHeight + statusRowHeight + commandRowHeight)
	if viewportStyleHeight < 3 {
		viewportStyleHeight = 3 // Minimum height (1 content + 2 borders)
	}

	fileListViewportStyle := borderStyle.
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

	fileListViewportView := fileListViewportStyle.Render(
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
		fileListViewportView,
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

	// Command bar (only shown when command mode is active).
	commandBarStyle := borderStyle.
		Width(m.width).
		Background(lipgloss.Color(m.theme.CommandBar.Background)).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground)).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.CommandBar.PaddingBottom).
		PaddingLeft(m.theme.CommandBar.PaddingLeft).
		PaddingRight(m.theme.CommandBar.PaddingRight).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)

	rows := []string{
		searchBar,
		viewports,
		statusView,
	}

	if m.commandMode {
		commandBar := commandBarStyle.Render(
			m.commandBar.View(),
		)
		rows = append(rows, commandBar)
	}

	// Combine all rows vertically
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
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
		Style:   DefaultFloatingStyle(m.theme),
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
