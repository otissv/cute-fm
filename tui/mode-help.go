package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) HelpMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	m.commandInput.Blur()
	m.searchInput.Focus()

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Close help window
	case bindings.Help.Matches(keyMsg.String()) ||
		bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Scroll help content up
	case bindings.Up.Matches(keyMsg.String()):
		if m.helpScrollOffset > 0 {
			m.helpScrollOffset--
		}
		return m, nil

	// Scroll help content down
	case bindings.Down.Matches(keyMsg.String()):
		m.helpScrollOffset++
		return m, nil
	}

	return m, nil
}
