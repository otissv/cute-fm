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
		commandBar = m.CommandBar(
			m, ComponentArgs{
				Width:  m.width,
				Height: 1,
			})
	}

	fileListViewportView := m.FileList(
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
	previewViewportView := m.Preview(
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

	layoutStyle := lipgloss.NewStyle()

	m.layout = lipgloss.JoinVertical(
		lipgloss.Center,
		m.layoutRows...,
	)

	baseContent := layoutStyle.Render(m.layout)
	baseLayer := lipgloss.NewLayer(baseContent)

	var canvas *lipgloss.Canvas
	switch ActiveTuiMode {

	case TuiModeHelp:
		if m.HelpModal != nil {

			modalLayer := m.HelpModal(m)
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
