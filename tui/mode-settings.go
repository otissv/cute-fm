package tui

import (
	"cute/console"

	tea "charm.land/bubbletea/v2"
)

var (
	SettingTabIndex    = 0
	SettingCursorIndex = 1

	SettingTabs = []ActiveSetting{
		SETTING_START,
		SETTING_SPLIT_PANE,
		SETTING_FILE_LIST_MODE,
		SETTING_COLUMN_VISIBILITY,
		SETTING_SORT_BY_COLUMN,
		SETTING_SORT_COLUMN_DIRECTION,
	}
)

func (m Model) SettingsMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()
	settings := GetSettings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Enter normal mode-
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		return m, nil

	case bindings.Down.Matches(keyMsg.String()):

		SettingCursorIndex += 1

		if SettingCursorIndex == len(GetSettings()) {
			SettingCursorIndex = 1
		}

		if settings[SettingCursorIndex].Type == HEADING_CHOICE_TYPE {
			SettingCursorIndex += 1
		}

		console.Log("%v", SettingCursorIndex)

		return m, nil

	case bindings.Up.Matches(keyMsg.String()):
		SettingCursorIndex -= 1

		if SettingCursorIndex == 0 {
			SettingCursorIndex = len(settings) - 1
		}

		if settings[SettingCursorIndex].Type == HEADING_CHOICE_TYPE {
			if SettingCursorIndex == 0 {
				SettingCursorIndex = len(settings) - 1
			} else {
				SettingCursorIndex -= 1
			}
		}

		if SettingCursorIndex == -1 {
			SettingCursorIndex = 1
		}

		return m, nil
	}

	return m, nil
}
