package tui

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"

	"cute/command"
)

func (m *Model) CommandBar() string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.CommandBar.Background)).
		BorderBottom(false).
		BorderForeground(lipgloss.Color(m.theme.BorderColor)).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(false).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground)).
		PaddingBottom(m.theme.CommandBar.PaddingBottom).
		PaddingLeft(m.theme.CommandBar.PaddingLeft).
		PaddingRight(m.theme.CommandBar.PaddingRight).
		PaddingTop(m.theme.Preview.PaddingTop).
		Width(m.width).
		Render(m.commandInput.View())
}

func (m *Model) CurrentDir() string {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Background(lipgloss.Color(m.theme.CurrentDir.Background)).
		Foreground(lipgloss.Color(m.theme.CurrentDir.Foreground)).
		MarginRight(2).
		PaddingBottom(1).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(0).
		Render(m.currentDir)
}

func (m *Model) FileList() string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.FileList.Background)).
		BorderBackground(lipgloss.Color(m.theme.FileList.Background)).
		BorderForeground(lipgloss.Color(m.theme.FileList.Border)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
		Foreground(lipgloss.Color(m.theme.FileList.Foreground)).
		Height(m.viewportHeight).
		PaddingBottom(m.theme.FileList.PaddingBottom).
		PaddingLeft(m.theme.FileList.PaddingLeft).
		PaddingRight(m.theme.FileList.PaddingRight).
		PaddingTop(m.theme.FileList.PaddingTop).
		Width(m.viewportWidth).
		Render(m.FileListViewport.View())
}

func (m *Model) Header() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(m.theme.Header.Background)).
		PaddingBottom(1).
		Width(m.width).
		Render(rainbowText(
			lipgloss.NewStyle().
				Background(lipgloss.Color(m.theme.Header.Background)),
			m.titleText,
			blends(m.theme.Primary, m.theme.Secondary),
		))
}

func (m *Model) Preview() string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.Preview.Background)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBackground(lipgloss.Color(m.theme.Background)).
		BorderForeground(lipgloss.Color(m.theme.Preview.Border)).
		Foreground(lipgloss.Color(m.theme.Preview.Foreground)).
		Height(m.viewportHeight).
		PaddingBottom(m.theme.Preview.PaddingBottom).
		PaddingLeft(m.theme.Preview.PaddingLeft).
		PaddingRight(m.theme.Preview.PaddingRight).
		PaddingTop(m.theme.Preview.PaddingTop).
		Width(m.viewportWidth).
		Render(m.previewViewport.View())
}

func (m *Model) PreviewTabs() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.theme.SearchBar.Foreground)).
		Background(lipgloss.Color(m.theme.SearchBar.Background)).
		BorderBackground(lipgloss.Color(m.theme.Background)).
		BorderForeground(lipgloss.Color(m.theme.SearchBar.Border)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Height(1).
		Width(m.viewportWidth).
		Render("Tabs")
}

func (m *Model) SearchBar() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.theme.SearchBar.Foreground)).
		Background(lipgloss.Color(m.theme.SearchBar.Background)).
		BorderBackground(lipgloss.Color(m.theme.SearchBar.Background)).
		BorderForeground(lipgloss.Color(m.theme.SearchBar.Border)).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(lipgloss.NormalBorder()).
		Height(1).
		Width(m.viewportWidth).
		Render(m.searchInput.View())
}

func (m *Model) StatusBar(items ...string) string {
	statusStyle := lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		Background(lipgloss.Color(m.theme.StatusBar.Background)).
		PaddingBottom(m.theme.StatusBar.PaddingBottom).
		PaddingLeft(m.theme.StatusBar.PaddingLeft).
		PaddingRight(m.theme.StatusBar.PaddingRight).
		PaddingTop(m.theme.Preview.PaddingTop).
		Width(m.width)

	var flatItems []string
	flatItems = append(flatItems, items...)

	statusBar := lipgloss.JoinHorizontal(lipgloss.Left, flatItems...)
	return statusStyle.Render(statusBar)
}

func (m *Model) ViewText() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(m.theme.ViewMode.Background)).
		Foreground(lipgloss.Color(m.theme.ViewMode.Foreground)).
		Width(10).
		Render(command.CmdViewModeStatus(m.viewMode))
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

func blends(colo1 string, color2 string) []color.Color {
	return gamut.Blends(lipgloss.Color(colo1), lipgloss.Color(color2), 50)
}

func rainbowText(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str += base.
			Background(lipgloss.Color(color.Hex())). // use blend as background
			Foreground(lipgloss.Color("#D4D4D4")).   // fixed text color
			Render(string(ss))
	}
	return str
}
