package tui

import "charm.land/lipgloss/v2"

func TuiMode(m Model, args ComponentArgs) string {
	theme := m.GetTheme()
	activeTuiMode := ActiveTuiMode

	foreground := ""
	background := ""

	switch activeTuiMode {
	case ModeNormal:
		background = theme.TuiMode.NormalModeBackground
		foreground = theme.TuiMode.NormalModeForeground
	case ModeCommand:
		background = theme.TuiMode.CommandModeBackground
		foreground = theme.TuiMode.CommandModeForeground
	case ModeFilter:
		background = theme.TuiMode.FilterModeBackground
		foreground = theme.TuiMode.FilterModeForeground
	case ModeHelp:
		background = theme.TuiMode.HelpModeBackground
		foreground = theme.TuiMode.HelpModeForeground
	case ModeQuit:
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
