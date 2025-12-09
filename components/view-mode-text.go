package components

import (
	"cute/command"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ViewModeText(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color(theme.ViewMode.Background)).
		Foreground(lipgloss.Color(theme.ViewMode.Foreground)).
		Height(args.Height).
		PaddingRight(1).
		Width(args.Width).
		Render(command.CmdViewModeStatus(string(tui.ActiveFileListMode)))
}
