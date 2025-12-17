package tui

import (
	"charm.land/lipgloss/v2"
)

type DialogWindowArgs struct {
	Title   string
	Content string
	Width   int
	Height  int
}

func DialogWindow(m Model, args DialogWindowArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	windowWidth := 40
	windowHeight := 6

	if args.Width != 0 {
		windowWidth = args.Width
	}

	if args.Height != 0 {
		windowHeight = args.Height
	}

	fw := FloatingWindow{
		Content: ViewPrimitive(args.Content),
		Width:   windowWidth,
		Height:  windowHeight,
		Style:   DefaultFloatingStyle(theme).Align(lipgloss.Center),
		Title:   args.Title,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
