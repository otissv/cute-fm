package tui

import (
	"charm.land/lipgloss/v2"
)

type DialogWindowArgs struct {
	Title   string
	Content string
}

func DialogWindow(m Model, args DialogWindowArgs) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	fw := FloatingWindow{
		Content: ViewPrimitive(args.Content),
		Width:   40,
		Height:  6,
		Style:   DefaultFloatingStyle(theme).Align(lipgloss.Center),
		Title:   args.Title,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
