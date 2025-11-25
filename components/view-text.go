package components

import (
	"cute/command"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func ViewText(m tui.Model) string {
	theme := m.GetTheme()

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(lipgloss.Color(theme.ViewMode.Background)).
		Foreground(lipgloss.Color(theme.ViewMode.Foreground)).
		Width(10).
		Render(command.CmdViewModeStatus(m.GetViewMode()))
}
