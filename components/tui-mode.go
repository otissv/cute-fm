package components

import (
	"cute/tui"

	"charm.land/lipgloss/v2"
)

func TuiMode(m tui.Model, args tui.ComponentArgs) string {
	theme := m.GetTheme()
	activeTuiMode := tui.ActiveTuiMode

	foreground := ""
	background := ""

	switch activeTuiMode {
	case tui.TuiModeNormal:
		background = theme.TuiMode.NormalModeBackground
		foreground = theme.TuiMode.NormalModeForeground
	case tui.TuiModeCommand:
		background = theme.TuiMode.CommandModeBackground
		foreground = theme.TuiMode.CommandModeForeground
	case tui.TuiModeFilter:
		background = theme.TuiMode.FilterModeBackground
		foreground = theme.TuiMode.FilterModeForeground
	case tui.TuiModeHelp:
		background = theme.TuiMode.HelpModeBackground
		foreground = theme.TuiMode.HelpModeForeground
	case tui.TuiModeQuit:
		background = theme.TuiMode.QuitModeBackground
		foreground = theme.TuiMode.QuitModeForeground
	}

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Background(lipgloss.Color(background)).
		Foreground(lipgloss.Color(foreground)).
		PaddingBottom(0).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(0).
		Height(args.Height).
		Render(string(activeTuiMode))
}
