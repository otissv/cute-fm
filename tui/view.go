package tui

import (
	"github.com/charmbracelet/lipgloss/v2"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	commandBar := m.CommandBar()
	currentDir := m.CurrentDir()
	fileListViewportView := m.FileList()
	headerView := m.Header()
	previewTabs := m.PreviewTabs()
	previewViewportView := m.Preview()
	searchBar := m.SearchBar()
	viewModeText := m.ViewText()

	statusBar := m.StatusBar(viewModeText, currentDir)

	filePanelRows := []string{
		searchBar,
		fileListViewportView,
	}

	if !m.isSearchBarOpen {
		// Only show searchBar if it's open
		filePanelRows = filePanelRows[1:]
	}

	filePanel := lipgloss.JoinVertical(
		lipgloss.Left,
		filePanelRows...,
	)

	previewPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		previewTabs,
		previewViewportView,
	)

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Center,
		filePanel,
		previewPanel,
	)

	m.layoutRows = []string{
		headerView,
		viewports,
		statusBar,
	}

	if m.isCommandBarOpen {
		m.layoutRows = append(m.layoutRows, commandBar)
	}

	layoutStyle := lipgloss.NewStyle().Background(lipgloss.Color(m.theme.Background))

	m.layout = lipgloss.JoinVertical(
		lipgloss.Center,
		m.layoutRows...,
	)

	if m.activeModal == ModalNone {
		return layoutStyle.Render(m.layout)
	}

	return m.Modal()
}
