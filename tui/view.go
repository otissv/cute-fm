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

	defaultStatus := []string{tuiMode, viewModeText, currentDir}

	statusBarItem := defaultStatus

	if ActiveTuiMode == TuiModeGoto {
		statusBarItem = []string{tuiMode, m.jumpTo}
	}

	statusBar := m.StatusBar(
		m, ComponentArgs{
			Width:  m.width,
			Height: 1,
		},
		statusBarItem...,
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

	case TuiModeAddFile:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Add New File",
			Placeholder: "Enter file name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeCd:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Change Directory",
			Placeholder: "Enter directory...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeCommand:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Command",
			Placeholder: "Enter commnad..",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeCopy:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Copy",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeHelp:
		modalLayer := m.HelpModal(m)
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case TuiModeMkdir:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Add Directory",
			Placeholder: "Enter directory name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeMove:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Move",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeQuit:
		modalLayer := m.DialogModal(m, DialogModalArgs{
			Title:   "Quit",
			Content: "Press q to quit\n\nor\n\n press ESC to cancel",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)
	case TuiModeRemove:
		modalLayer := m.DialogModal(m, DialogModalArgs{
			Title:   "Remove",
			Content: "Are you sure you want to remove\n\nYes (y) No (n)",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	return v
}
