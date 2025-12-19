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

	isLeftViewportActive := m.activeViewport == LeftViewportType

	tuiMode := TuiMode(m, TuiModeComponentArgs{
		Height: 1,
		Width:  20,
	})

	viewModeText := ViewModeText(
		m, ViewModeTextComponentArgs{
			Height: 1,
			Width:  20,
		})

	header := Header(m, HeaderComponentArgs{
		Height: 1,
		Width:  m.width - 40,
	})

	headerView := lipgloss.NewStyle().
		PaddingBottom(1).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, tuiMode, viewModeText, header))

	searchBar := SearchBar(
		m, SearchBarComponentArgs{
			Width:  m.viewportWidth,
			Height: 1,
		})

	// sudoMode := SudoMode(m, ComponentArgs{
	// 	Height: 1,
	// })

	// if m.isSudo {
	// 	leftStatusBarItem = append([]string{sudoMode}, leftStatusBarItem...)
	// }

	fileInfoViewportView := ViewPort(
		m, ViewportComponentArgs{
			Width:    m.viewportWidth,
			Height:   m.viewportHeight + 1,
			Viewport: m.GetFileInfoViewport(),
		})

	deviceInfoViewportView := ViewPort(
		m, ViewportComponentArgs{
			Width:    m.viewportWidth,
			Height:   m.viewportHeight + 1,
			Viewport: m.GetDeviceInfoViewport(),
		})

	leftCurrentDir := CurrentDir(m, CurrentDirComponentArgs{
		Height:     1,
		CurrentDir: m.GetLeftPaneCurrentDir(),
	})

	filePane1StatusBar := StatusBar(
		m, StatusBarComponentArgs{
			Height: 1,
		},
		leftCurrentDir,
	)

	fileListView1 := FileList(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: LeftViewportType,
		})

	fileListView2 := FileList(
		m, FileListComponentArgs{
			Width:         m.viewportWidth,
			Height:        m.viewportHeight,
			SplitPaneType: RightViewportType,
		})

	placeholder := lipgloss.NewStyle().Render("")
	leftPaneHeader := SearchText(m, LeftViewportType)
	rightPaneHeader := placeholder

	if isLeftViewportActive {
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

	if ActiveFileListMode == FileListModeDevice {
		leftPaneItems = []string{
			leftPaneHeader,
			DeviceList(m, DeviceListComponentArgs{
				Width:  m.viewportWidth,
				Height: m.viewportHeight,
			}),
			"",
		}
	}

	rightPaneItems := []string{}

	if m.showRightPane {

		console.Log("%s View", m.activeSplitPane)

		switch m.activeSplitPane {

		case DeviceInfoSplitPaneType:
			rightPaneItems = []string{
				rightPaneHeader,
				deviceInfoViewportView,
			}

		case FileInfoSplitPaneType:
			rightPaneItems = []string{
				rightPaneHeader,
				fileInfoViewportView,
			}

		case FileListSplitPaneType:

			rightCurrentDir := CurrentDir(m, CurrentDirComponentArgs{
				Height:     1,
				CurrentDir: m.GetRightPaneCurrentDir(),
			})

			rightPaneHeader = SearchText(m, RightViewportType)

			if !isLeftViewportActive {
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

	sidePanel := ""

	if m.isSidePanelOpen {
		sidePanel = SidePanel(m)
	}

	viewports := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidePanel,
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
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Add New File",
			Placeholder: "Enter file name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCd:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Change Directory",
			Placeholder: "Enter directory...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeColumnVisibility:
		windowLayer := ColumnWindow(m, ColumnWindowArgs{
			Title: "Column Visibility",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeCommand:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Command",
			Placeholder: "Enter command..",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeCopy:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Copy",
			Placeholder: "Enter destination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeHelp:
		windowLayer := HelpWindow(m)
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeMkdir:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Add Directory",
			Placeholder: "Enter directory name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeMove:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Move",
			Placeholder: "Enter destination...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeQuit:
		windowLayer := QuitWindow(m)
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeRemove:
		windowLayer := DialogWindow(m, DialogWindowArgs{
			Title:   "Remove",
			Content: "Are you sure you want to remove\n\nYes (y) No (n)",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeRename:
		commandLayer := CommandWindow(m, CommandWindowArgs{
			Title:       "Rename",
			Placeholder: "New name...",
		})
		canvas = lipgloss.NewCanvas(baseLayer, commandLayer)

	case ModeSettings:
		windowLayer := SettingsWindow(m)
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	case ModeSort:
		windowLayer := ColumnWindow(m, ColumnWindowArgs{
			Title: "Sort Columns",
		})
		canvas = lipgloss.NewCanvas(baseLayer, windowLayer)

	default:
		canvas = lipgloss.NewCanvas(baseLayer)
	}

	v := tea.NewView(canvas)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion // Enable mouse support
	return v
}
