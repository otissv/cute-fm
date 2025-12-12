package tui

import (
	"cute/command"

	"charm.land/lipgloss/v2"
)

type ViewModeTextComponentArgs struct {
	Width  int
	Height int
}

func ViewModeText(m Model, args ViewModeTextComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color(theme.ViewMode.Background)).
		Foreground(lipgloss.Color(theme.ViewMode.Foreground)).
		Height(args.Height).
		PaddingRight(1).
		Width(args.Width).
		Render(command.CmdViewModeStatus(string(ActiveFileListMode)))
}
