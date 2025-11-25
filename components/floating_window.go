package components

import (
	"cute/theming"

	"cute/tui"

	"charm.land/lipgloss/v2"
)

type FloatingWindow struct {
	Content tui.ViewPrimitive
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
		BorderForeground(lipgloss.Color(theme.DefaultDialog.Border)).
		BorderBackground(lipgloss.Color(theme.DefaultDialog.Background)).
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

	return box
}
