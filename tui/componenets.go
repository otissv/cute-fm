package tui

import (
	"image/color"
	"strings"

	"lsfm/command"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

var blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)

func (m *Model) CommandBar() string {
	commandBarStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground)).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.CommandBar.PaddingBottom).
		PaddingLeft(m.theme.CommandBar.PaddingLeft).
		PaddingRight(m.theme.CommandBar.PaddingRight).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)

	return commandBarStyle.Render(
		m.commandBar.View(),
	)
}

func (m *Model) FileList() string {
	fileListViewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		Width(m.viewportWidth).
		Height(m.viewportHeight).
		Background(lipgloss.Color(m.theme.FileList.Background)).
		Foreground(lipgloss.Color(m.theme.FileList.Foreground)).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		BorderRight(true).
		PaddingTop(m.theme.FileList.PaddingTop).
		PaddingBottom(m.theme.FileList.PaddingBottom).
		PaddingLeft(m.theme.FileList.PaddingLeft).
		PaddingRight(m.theme.FileList.PaddingRight)

	return fileListViewportStyle.Render(
		m.FileListViewport.View(),
	)
}

func (m *Model) Preview() string {
	previewViewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		Width(m.viewportWidth).
		Height(m.viewportHeight).
		Background(lipgloss.Color(m.theme.Preview.Background)).
		Foreground(lipgloss.Color(m.theme.Preview.Foreground)).
		BorderTop(false).
		BorderRight(false).
		BorderBottom(false).
		BorderLeft(false).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.Preview.PaddingBottom).
		PaddingLeft(m.theme.Preview.PaddingLeft).
		PaddingRight(m.theme.Preview.PaddingRight)

	return previewViewportStyle.Render(
		m.previewViewport.View(),
	)
}

func (m *Model) SearchBar() string {
	searchBarStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		Width(m.width).
		Height(1).
		BorderBottom(true)

	return searchBarStyle.Render(
		m.searchBar.View(),
	)
}

func (m *Model) StatusBar() string {
	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(m.theme.StatusBar.Background)).
		Foreground(lipgloss.Color(m.theme.StatusBar.Foreground)).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.StatusBar.PaddingBottom).
		PaddingLeft(m.theme.StatusBar.PaddingLeft).
		PaddingRight(m.theme.StatusBar.PaddingRight)

	statusText := m.currentDir
	if statusText == "" {
		statusText = "."
	}

	viewModeText := command.CmdViewModeStatus(m.viewMode)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, viewModeText, statusText)
	return statusStyle.Render(statusBar)
}

func (m *Model) Header() string {
	headerStyle := lipgloss.NewStyle().
		Width(m.width).
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(m.theme.StatusBar.Background)).
		Foreground(lipgloss.Color(m.theme.StatusBar.Foreground)).
		PaddingTop(m.theme.Preview.PaddingTop).
		PaddingBottom(m.theme.StatusBar.PaddingBottom).
		PaddingLeft(m.theme.StatusBar.PaddingLeft).
		PaddingRight(m.theme.StatusBar.PaddingRight)

	headerText := m.currentDir
	if headerText == "" {
		headerText = "."
	}

	viewModeText := command.CmdViewModeStatus(m.viewMode)

	header := lipgloss.JoinHorizontal(lipgloss.Left, viewModeText, headerText)
	return headerStyle.Render(rainbow(lipgloss.NewStyle(), header, blends))
}

func HelpViewport() viewport.Model {
	helpContent := `
	Help
	----
	
	Navigation:
		Up/Down arrows   Move selection in file list
		Scroll wheel     Scroll file list
	
	Search:
		Type in the search bar to filter files by name
	
	General:
		?                Toggle this help
		ctrl+c / ctrl+q  Quit
	`
	helpViewport := viewport.New(0, 0)
	helpViewport.SetContent(strings.TrimSpace(helpContent))

	return helpViewport
}

func (m *Model) Modal() string {
	// Decide which content to show inside the floating window.
	var content ViewPrimitive
	switch m.activeModal {
	case ModalHelp:
		content = m.helpViewport
	}
	// Choose a dialog-sized window, not full-screen.
	modalWidth := m.width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	modalHeight := m.height / 2
	if modalHeight > 16 {
		modalHeight = 16
	}
	if modalHeight < 6 {
		modalHeight = 6
	}

	fw := FloatingWindow{
		Content: content,
		Width:   modalWidth,
		Height:  modalHeight,
		Style:   DefaultFloatingStyle(m.theme),
	}

	dialogView := fw.View(m.width, m.height)
	return overlayModel(m.layout, dialogView)
}

// overlayDialog composes the base layout and the dialog view.
func overlayModel(layout, dialog string) string {
	baseLines := strings.Split(layout, "\n")
	dialogLines := strings.Split(dialog, "\n")

	maxLines := len(baseLines)
	if len(dialogLines) > maxLines {
		maxLines = len(dialogLines)
	}

	out := make([]string, maxLines)

	for i := 0; i < maxLines; i++ {
		var baseLine, dlgLine string
		if i < len(baseLines) {
			baseLine = baseLines[i]
		}
		if i < len(dialogLines) {
			dlgLine = dialogLines[i]
		}

		if strings.TrimSpace(dlgLine) != "" {
			out[i] = dlgLine
		} else {
			out[i] = baseLine
		}
	}

	return strings.Join(out, "\n")
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
