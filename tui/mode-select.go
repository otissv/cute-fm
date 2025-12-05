package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) SelectMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

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
	// Enter normal mode
	case bindings.Select.Matches(keyMsg.String()) ||
		bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = TuiModeNormal
		return m, nil
	}

	return m, nil
}
