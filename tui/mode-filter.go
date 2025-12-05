package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) FilterMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	// Update search input (first row) and apply filtering if the value changed.
	before := m.searchInput.Value()
	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)
	if m.searchInput.Value() != before {
		m.ApplyFilter()
	}

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Enter normal mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = TuiModeNormal
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
