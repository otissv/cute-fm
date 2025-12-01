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

	leftViewportView := m.FileListView(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: m.viewportHeight,
		})
	headerView := m.Header(m, ComponentArgs{
		Width: m.width,
	})
	previewTabs := m.PreviewTabs(m, ComponentArgs{
		Width:  m.viewportWidth,
		Height: 1,
	})
	rightViewportView := m.Preview(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: m.viewportHeight,
		})
	searchBar := m.SearchBar(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: 1,
		})

	currentDir := m.CurrentDir(m, ComponentArgs{
		Height: 1,
	})
	tuiMode := m.TuiMode(m, ComponentArgs{
		Height: 1,
	})
	viewModeText := m.ViewModeText(
		m, ComponentArgs{
			Width:  10,
			Height: 1,
		})

	statusBar := m.StatusBar(
		m, ComponentArgs{
			Width:  m.width,
			Height: 1,
		},
		tuiMode,
		viewModeText,
		currentDir,
	)

	filePanelRows := []string{
		searchBar,
		leftViewportView,
	}

	filePanel := lipgloss.JoinVertical(
		lipgloss.Left,
		filePanelRows...,
	)

	previewPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		previewTabs,
		rightViewportView,
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
	}

	layoutStyle := lipgloss.NewStyle()

	m.layout = lipgloss.JoinVertical(
		lipgloss.Center,
		m.layoutRows...,
	)

	baseContent := layoutStyle.Render(m.layout)
	baseLayer := lipgloss.NewLayer(baseContent)

	var canvas *lipgloss.Canvas
	switch ActiveTuiMode {

	case TuiModeCommand:
		commandLayer := m.CommandModal(m)
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeHelp:
		modalLayer := m.HelpModal(m)
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case TuiModeQuit:
		modalLayer := m.QuitModal(m)
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	return v
}
