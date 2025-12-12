package tui

import (
	"charm.land/lipgloss/v2"
)

func SettingsWindow(m Model) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	settings := GetSettings()
	choices := settings

	selected := map[string]string{
		choices[0].Label: choices[0].Label,
	}

	menu := NewMenu(MenuArgs{
		Choices:     choices,
		CursorIndex: SettingCursorIndex,
		CursorTypes: MenuCursor{
			Numbered: true,
		},
		Selected: selected,
		Theme:    m.theme,
	})

	contentItems := []string{}
	contentItems = append(contentItems, menu.View())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentItems...,
	)

	fw := FloatingWindow{
		Content: ViewPrimitive(content),
		Width:   50,
		Height:  50,
		Style:   DefaultFloatingStyle(theme),
		Title:   "Settings",
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}
