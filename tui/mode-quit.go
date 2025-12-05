package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) QuitMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case bindings.Quit.Matches(keyMsg.String()):
		return m, tea.Quit

	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = PreviousTuiMode
		return m, nil
	}

	return m, nil
}
