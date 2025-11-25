package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if m.width == 0 {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	commandBar := ""
	if m.CommandBar != nil {
		commandBar = m.CommandBar(m)
	}
	currentDir := m.CurrentDir(m)
	fileListViewportView := m.FileList(m)
	headerView := m.Header(m)
	previewTabs := m.PreviewTabs(m)
	previewViewportView := m.Preview(m)
	searchBar := m.SearchBar(m)
	viewModeText := m.ViewText(m)

	statusBar := m.StatusBar(m, viewModeText, currentDir)

	filePanelRows := []string{
		searchBar,
		fileListViewportView,
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
		lipgloss.Top,
		filePanel,
		previewPanel,
	)

	m.layoutRows = []string{
		headerView,
		viewports,
		statusBar,
		commandBar,
	}

	layoutStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.Background)).
		Height(m.height).
		Width(m.width)

	m.layout = lipgloss.JoinVertical(
		lipgloss.Center,
		m.layoutRows...,
	)

	baseContent := layoutStyle.Render(m.layout)
	baseLayer := lipgloss.NewLayer(baseContent)

	var canvas *lipgloss.Canvas
	switch m.activeModal {

	case ModalHelp:
		if m.HelpModal != nil {
			modalContent := m.HelpModal(m)

			dialogWidth := lipgloss.Width(modalContent)
			dialogHeight := lipgloss.Height(modalContent)
			x := 0
			y := 0
			if m.width > dialogWidth {
				x = (m.width - dialogWidth) / 2
			}
			if m.height > dialogHeight {
				y = (m.height - dialogHeight) / 2
			}

			modalLayer := lipgloss.NewLayer(modalContent).X(x).Y(y)
			canvas = lipgloss.NewCanvas(baseLayer, modalLayer)
		} else {
			// No help modal configured; just render the base layout.
			canvas = lipgloss.NewCanvas(baseLayer)
		}
	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	return v
}
