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

	isLeftViewportActivce := m.activeViewport == LeftViewportType

	tuiMode := m.Components.TuiMode(m, TuiModeComponentArgs{
		Height: 1,
		Width:  20,
	})

	viewModeText := m.Components.ViewModeText(
		m, ViewModeTextComponentArgs{
			Height: 1,
			Width:  20,
		})

	header := m.Components.Header(m, HeaderComponentArgs{
		Height: 1,
		Width:  m.width - 40,
	})

	headerView := lipgloss.NewStyle().
		PaddingBottom(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, tuiMode, viewModeText, header))

	searchBar := m.Components.SearchBar(
		m, SearchBarComponentArgs{
			Width:  m.viewportWidth,
			Height: 1,
		})

	// sudoMode := m.Components.SudoMode(m, ComponentArgs{
	// 	Height: 1,
	// })

	// if m.isSudo {
	// 	leftStatusBarItem = append([]string{sudoMode}, leftStatusBarItem...)
	// }

	fileInfoViewportView := m.Components.FileInfo(
		m, FileInfoComponentArgs{
			Width:  m.viewportWidth,
			Height: m.viewportHeight + 1,
		})

	leftCurrentDir := m.Components.CurrentDir(m, CurrentDirComponentArgs{
		Height:     1,
		CurrentDir: m.GetLeftPaneCurrentDir(),
	})

	filePane1StatusBar := m.Components.StatusBar(
		m, StatusBarComponentArgs{
			Height: 1,
		},
		leftCurrentDir,
	)

	fileListView1 := m.Components.FileListView(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: LeftViewportType,
		})

	fileListView2 := m.Components.FileListView(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: RightViewportType,
		})

	placeholder := lipgloss.NewStyle().Render("")
	leftPaneHeader := m.Components.SearchText(m, LeftViewportType)
	rightPaneHeader := placeholder

	if isLeftViewportActivce {
		if ActiveTuiMode == ModeFilter {
			leftPaneHeader = searchBar
		}

		if ActiveTuiMode == ModeGoto {
			leftPaneHeader = "Jump to row: " + m.jumpTo
			rightPaneHeader = placeholder
		}
	}

	leftPaneItems := []string{
		leftPaneHeader,
		fileListView1,
		filePane1StatusBar,
	}

	rightPaneItems := []string{}

	if m.showRightPane {
		switch m.activeSplitPane {

		case FileInfoSplitPaneType:

			rightPaneItems = []string{
				rightPaneHeader,
				fileInfoViewportView,
			}

		case FileListSplitPaneType:
			rightCurrentDir := m.Components.CurrentDir(m, CurrentDirComponentArgs{
				Height:     1,
				CurrentDir: m.GetRightPaneCurrentDir(),
			})

			rightPaneHeader = m.Components.SearchText(m, RightViewportType)

			if !isLeftViewportActivce {
				if ActiveTuiMode == ModeGoto {
					rightPaneHeader = "Jump to row: " + m.jumpTo
				}

				if ActiveTuiMode == ModeFilter {
					rightPaneHeader = searchBar
				}
			}

			rightPaneItems = []string{
				rightPaneHeader,
				fileListView2,
				rightCurrentDir,
			}
		}
	}

	rightPane := lipgloss.JoinVertical(
		lipgloss.Left,
		rightPaneItems...,
	)

	leftPane := lipgloss.JoinVertical(
		lipgloss.Left,
		leftPaneItems...,
	)

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
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
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Add New File",
			Placeholder: "Enter file name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCd:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Change Directory",
			Placeholder: "Enter directory...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeColumnVisibility:
		windowLayer := m.Windows.Column(m, ColumnWindowArgs{
			Title: "Column Visibilty",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeCommand:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Command",
			Placeholder: "Enter commnad..",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCopy:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Copy",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeHelp:
		windowLayer := m.Windows.Help(m)
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeMkdir:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Add Directory",
			Placeholder: "Enter directory name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeMove:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Move",
			Placeholder: "Enter desination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeQuit:
		windowLayer := m.Windows.Dialog(m, DialogWindowArgs{
			Title:   "Quit",
			Content: "Press q to quit\n\nor\n\n press ESC to cancel",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeRemove:
		windowLayer := m.Windows.Dialog(m, DialogWindowArgs{
			Title:   "Remove",
			Content: "Are you sure you want to remove\n\nYes (y) No (n)",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeRename:
		commandLayer := m.Windows.Command(m, CommandWindowArgs{
			Title:       "Remane",
			Placeholder: "New name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeSettings:
		windowLayer := m.Windows.Settings(m)
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeSort:
		windowLayer := m.Windows.Column(m, ColumnWindowArgs{
			Title: "Sort Columns",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	return v
}
