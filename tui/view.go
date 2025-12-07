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

	headerView := m.Header(m, ComponentArgs{
		Width: m.width,
	})

	searchBar := m.SearchBar(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: 1,
		})

	leftCurrentDir := m.CurrentDir(m, ComponentArgs{
		Height: 1,
	})

	tuiMode := m.TuiMode(m, ComponentArgs{
		Height: 1,
	})

	sudoMode := m.SudoMode(m, ComponentArgs{
		Height: 1,
	})

	viewModeText := m.ViewModeText(
		m, ComponentArgs{
			Width:  10,
			Height: 1,
		})

	leftStatus := []string{tuiMode, viewModeText, leftCurrentDir}

	leftStatusBarItem := leftStatus

	if ActiveTuiMode == TuiModeGoto {
		leftStatusBarItem = []string{tuiMode, m.jumpTo}
	}

	if m.isSudo {
		leftStatusBarItem = append([]string{sudoMode}, leftStatusBarItem...)
	}

	statusBar := m.StatusBar(
		m, ComponentArgs{
			Width:  m.width,
			Height: 1,
		},
		leftStatusBarItem...,
	)

	fileListView1 := m.FileListView(
		m, FileListComponentArgs{
			Width:          m.viewportWidth,
			Height:         m.viewportHeight,
			SplitPanelType: LeftViewportType,
		})

	fileListView2 := m.FileListView(
		m, FileListComponentArgs{
			Width:          m.viewportWidth,
			Height:         m.viewportHeight,
			SplitPanelType: RightViewportType,
		})

	filePanel1Rows := []string{
		searchBar,
		fileListView1,
	}

	filePanel2Rows := []string{
		searchBar,
		fileListView2,
	}

	leftPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		filePanel1Rows...,
	)

	fileInfoViewportView := m.Preview(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: m.viewportHeight,
		})

	rightPanel := lipgloss.JoinVertical(
		lipgloss.Left,
	)

	if m.showRightPanel {
		switch m.activeSplitPanel {
		case FileInfoSplitPanelType:
			rightPanel = lipgloss.JoinVertical(
				lipgloss.Left,
				fileInfoViewportView,
			)
		case FileListSplitPanelType:
			rightPanel = lipgloss.JoinVertical(
				lipgloss.Left,
				filePanel2Rows...,
			)
		}
	}

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	m.layoutRows = []string{
		headerView,
		viewports,
		statusBar,
	}

	layoutStyle := lipgloss.NewStyle()

	m.layout = lipgloss.JoinVertical(
		lipgloss.Left,
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

	case TuiModeColumnVisibiliy:
		modalLayer := m.ColumnModal(m, ColumnModelArgs{
			Title: "Column Visibilty",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

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

	case TuiModeRename:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Remane",
			Placeholder: "New name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case TuiModeSort:
		modalLayer := m.ColumnModal(m, ColumnModelArgs{
			Title: "Sort Columns",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	return v
}
