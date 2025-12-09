package tui

import (
	"cute/console"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m Model) View() tea.View {
	if m.width == 0 {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	tuiMode := m.TuiMode(m, ComponentArgs{
		Height: 1,
		Width:  20,
	})

	viewModeText := m.ViewModeText(
		m, ComponentArgs{
			Height: 1,
			Width:  20,
		})

	header := m.Header(m, ComponentArgs{
		Height: 1,
		Width:  m.width - 40,
	})

	headerView := lipgloss.NewStyle().
		PaddingBottom(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, tuiMode, viewModeText, header))

	searchBar := m.SearchBar(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: 1,
		})

	// sudoMode := m.SudoMode(m, ComponentArgs{
	// 	Height: 1,
	// })

	// if ActiveTuiMode == ModeGoto {
	// 	leftStatusBarItem = []string{tuiMode, m.jumpTo}
	// }

	// if m.isSudo {
	// 	leftStatusBarItem = append([]string{sudoMode}, leftStatusBarItem...)
	// }

	leftCurrentDir := m.CurrentDir(m, CurrentDirComponentArgs{
		Height:     1,
		CurrentDir: m.GetLeftPaneCurrentDir(),
	})

	filePanel1StatusBar := m.StatusBar(
		m, ComponentArgs{
			Height: 1,
		},
		leftCurrentDir,
	)

	fileListView1 := m.FileListView(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: LeftViewportType,
		})

	fileListView2 := m.FileListView(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: RightViewportType,
		})

	leftPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		searchBar,
		fileListView1,
		filePanel1StatusBar,
	)

	fileInfoViewportView := m.FileInfo(
		m, ComponentArgs{
			Width:  m.viewportWidth,
			Height: m.viewportHeight + 1,
		})

	rightPanel := fileInfoViewportView

	if m.showRightPanel {
		console.Log("ActiveTuiMode: %s", ActiveTuiMode)

		if m.isSplitPaneOpen {
			console.Log("m.isSplitPaneOpen: %s", "true")
		} else {
			console.Log("m.isSplitPaneOpen: %s", "false")
		}

		switch m.activeSplitPane {

		case FileInfoSplitPaneType:

			rightPanel = lipgloss.JoinVertical(
				lipgloss.Left,
				// Placeholder to align with left panel
				lipgloss.NewStyle().
					Height(1).
					Render(""),
				fileInfoViewportView,
			)

		case FileListSplitPaneType:
			rightCurrentDir := m.CurrentDir(m, CurrentDirComponentArgs{
				Height:     1,
				CurrentDir: m.GetRightPaneCurrentDir(),
			})

			rightPanel = lipgloss.JoinVertical(
				lipgloss.Left,
				searchBar,
				fileListView2,
				rightCurrentDir,
			)
		}
	}

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	layoutStyle := lipgloss.NewStyle()

	m.layout = lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		viewports,
	)

	baseContent := layoutStyle.Render(m.layout)
	baseLayer := lipgloss.NewLayer(baseContent)

	var canvas *lipgloss.Canvas
	switch ActiveTuiMode {

	case ModeAddFile:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Add New File",
			Placeholder: "Enter file name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCd:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Change Directory",
			Placeholder: "Enter directory...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeColumnVisibiliy:
		modalLayer := m.ColumnModal(m, ColumnModelArgs{
			Title: "Column Visibilty",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case ModeCommand:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Command",
			Placeholder: "Enter commnad..",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCopy:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Copy",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeHelp:
		modalLayer := m.HelpModal(m)
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case ModeMkdir:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Add Directory",
			Placeholder: "Enter directory name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeMove:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Move",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeQuit:
		modalLayer := m.DialogModal(m, DialogModalArgs{
			Title:   "Quit",
			Content: "Press q to quit\n\nor\n\n press ESC to cancel",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case ModeRemove:
		modalLayer := m.DialogModal(m, DialogModalArgs{
			Title:   "Remove",
			Content: "Are you sure you want to remove\n\nYes (y) No (n)",
		})
		canvas = lipgloss.NewCanvas(baseLayer, modalLayer)

	case ModeRename:
		commandLayer := m.CommandModal(m, CommandModalArgs{
			Title:       "Remane",
			Placeholder: "New name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeSort:
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
