package tui

import (
	"lsfm/theming"

	"github.com/charmbracelet/lipgloss"
)

type ViewPrimitive interface {
	View() string
}

type FloatingWindow struct {
	Content ViewPrimitive
	Width   int
	Height  int
	Title   string
	Style   lipgloss.Style
}

func DefaultFloatingStyle(theme theming.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.DefaultDialog.Background)).
		Foreground(lipgloss.Color(theme.DefaultDialog.Foreground)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.DefaultDialog.BorderColor)).
		PaddingTop(theme.DefaultDialog.PaddingTop).
		PaddingBottom(theme.DefaultDialog.PaddingBottom).
		PaddingLeft(theme.DefaultDialog.PaddingLeft).
		PaddingRight(theme.DefaultDialog.PaddingRight)
}

func (fw FloatingWindow) View(outerWidth, outerHeight int) string {
	if fw.Content == nil {
		return ""
	}

	contentView := fw.Content.View()

	style := fw.Style

	if fw.Width > 0 {
		style = style.Width(fw.Width)
	}
	if fw.Height > 0 {
		style = style.Height(fw.Height)
	}

	box := style.Render(contentView)

	return lipgloss.Place(
		outerWidth,
		outerHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars("///"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#303030")),
	)
}
