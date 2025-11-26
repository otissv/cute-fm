package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

// CommandBar renders the bottom command bar using only the public TUI model
// interface, so this component can live outside the tui package.
func CommandBar(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()
	view := m.GetCommandInputView()

	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.CommandBar.Background)).
		BorderBottom(false).
		BorderForeground(lipgloss.Color(theme.BorderColor)).
		BorderLeft(false).
		BorderRight(false).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(false).
		Foreground(lipgloss.Color(theme.CommandBar.Foreground)).
		PaddingBottom(theme.CommandBar.PaddingBottom).
		PaddingLeft(theme.CommandBar.PaddingLeft).
		PaddingRight(theme.CommandBar.PaddingRight).
		PaddingTop(theme.Preview.PaddingTop).
		Width(args.Width).
		Height(args.Height).
		Render(view)
}
