package tui

import (
	"os"

	"charm.land/lipgloss/v2"
)

func SettingsWindow(m Model) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	settings := GetSettingChoices()
	choices := settings
	selected := map[string]string{}

	// Map SplitPane setting to menu label
	switch m.settings.SplitPane {
	case PreviewPaneType:
		selected["Preview"] = "Preview"
	case FileInfoSplitPaneType:
		selected["File Info"] = "File Info"
	case FileListSplitPaneType:
		selected["File List"] = "File List"
	default:
		selected["None"] = "None"
	}

	// Map StartDir setting to menu label
	if m.settings.StartDir != "" {
		if homeDir, err := os.UserHomeDir(); err == nil && m.settings.StartDir == homeDir {
			selected["Home directory"] = "Home directory"
		} else {
			selected["Current directory"] = "Current directory"
		}
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
