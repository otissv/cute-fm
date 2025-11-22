package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Calculate viewport width (half of available width, accounting for borders)
	searchBarRowHeight := 3
	statusRowHeight := 3
	headerRowHeight := 3
	// Only reserve vertical space for the command bar when it is visible.
	commandRowHeight := 0
	if m.commandMode {
		commandRowHeight = 3
	}
	m.viewportWidth = m.width / 2
	m.viewportHeight = m.height - (searchBarRowHeight + statusRowHeight + commandRowHeight + headerRowHeight)
	if m.viewportHeight < 3 {
		m.viewportHeight = 3 // Minimum height (1 content + 2 borders)
	}

	commandBar := m.CommandBar()
	fileListViewportView := m.FileList()
	previewViewportView := m.Preview()
	searchBar := m.SearchBar()
	statusView := m.StatusBar()
	headerView := m.Header()

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fileListViewportView,
		previewViewportView,
	)

	m.layoutRows = []string{
		headerView,
		searchBar,
		viewports,
		statusView,
	}

	if m.commandMode {
		m.layoutRows = append(m.layoutRows, commandBar)
	}

	m.layout = lipgloss.JoinVertical(
		lipgloss.Left,
		m.layoutRows...,
	)

	if m.activeModal == ModalNone {
		return m.layout
	}

	return m.Modal()
}
