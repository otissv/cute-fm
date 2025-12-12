package tui

import (
	"charm.land/lipgloss/v2"
)

type CommandWindowArgs struct {
	Title       string
	Prompt      string
	Placeholder string
}

func CommandWindow(m Model, args CommandWindowArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	commandInput := m.GetCommandInput()

	title := ""
	if args.Title != "" {
		title = args.Title
	}
	commandInput.Prompt = ""
	if args.Prompt != "" {
		commandInput.Prompt = args.Prompt
	}
	commandInput.Placeholder = ""
	if args.Placeholder != "" {
		commandInput.Placeholder = args.Placeholder
	}

	commandInputView := commandInput.View()

	// Dialog-sized window
	windowWidth := width / 2
	if windowWidth > 60 {
		windowWidth = 60
	}
	if windowWidth < 30 {
		windowWidth = 30
	}

	fw := FloatingWindow{
		Content: ViewPrimitive(commandInputView),
		Width:   windowWidth,
		Height:  4,
		Style:   DefaultFloatingStyle(theme),
		Title:   title,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
