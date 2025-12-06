package tui

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) ConfirmMode(msg tea.Msg, command string) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case bindings.Quit.Matches(keyMsg.String()):
		return m, tea.Quit

	case keyMsg.String() == "n" || bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = PreviousTuiMode
		return m, nil

	case keyMsg.String() == "y":
		selected := m.GetSelectedEntry().Path

		res, _ := m.ExecuteCommand(command + " " + selected)

		if res.Cwd != "" && res.Cwd != m.currentDir {
			m.ChangeDirectory(res.Cwd)
		} else if res.Refresh {
			m.ChangeDirectory(m.currentDir)
		}

		ActiveTuiMode = PreviousTuiMode
		return m, nil
	}

	return m, nil
}
