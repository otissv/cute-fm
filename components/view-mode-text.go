package components

import (
	"cute/command"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ViewModeText(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(theme.ViewMode.Background)).
		Foreground(lipgloss.Color(theme.ViewMode.Foreground)).
		Height(args.Height).
		Width(args.Width).
		// Use the current active file-list mode instead of always rendering
		// the default "list all" mode, so commands like "ls", "ld", "lf"
		// are reflected in the UI.
		Render(command.CmdViewModeStatus(string(tui.ActiveFileListMode)))
}
