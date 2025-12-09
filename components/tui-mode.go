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
	case tui.ModeNormal:
		background = theme.TuiMode.NormalModeBackground
		foreground = theme.TuiMode.NormalModeForeground
	case tui.ModeCommand:
		background = theme.TuiMode.CommandModeBackground
		foreground = theme.TuiMode.CommandModeForeground
	case tui.ModeFilter:
		background = theme.TuiMode.FilterModeBackground
		foreground = theme.TuiMode.FilterModeForeground
	case tui.ModeHelp:
		background = theme.TuiMode.HelpModeBackground
		foreground = theme.TuiMode.HelpModeForeground
	case tui.ModeQuit:
		background = theme.TuiMode.QuitModeBackground
		foreground = theme.TuiMode.QuitModeForeground
	}

	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Left).
		Background(lipgloss.Color(background)).
		Foreground(lipgloss.Color(foreground)).
		Height(args.Height).
		PaddingBottom(0).
		PaddingRight(1).
		PaddingTop(0).
		Width(args.Width).
		Height(args.Height).
		Render(string(activeTuiMode))
}
