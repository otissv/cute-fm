package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) UtilityMode(msg tea.Msg, command string) (tea.Model, tea.Cmd) {
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

	m.searchInput.Blur()

	m.commandInput, cmd = m.commandInput.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Confirm command on Enter.
	case bindings.Enter.Matches(keyMsg.String()):

		inputValue := strings.TrimSpace(m.commandInput.Value())

		if inputValue != "" {
			line := command + " " + inputValue
			selectedEntry := m.GetSelectedEntry()

			// For copy and move operations, automatically include the selected path
			// as the source and use the input as the destination.
			if (command == "cp" || command == "mv") && selectedEntry != nil {
				line = command + " " + selectedEntry.Path + " " + inputValue
			}

			res, _ := m.ExecuteCommand(line)

			if res.Cwd != "" && res.Cwd != m.currentDir {
				m.ChangeDirectory(res.Cwd)
			} else if res.Refresh {
				m.ChangeDirectory(m.currentDir)
			}
		}

		m.commandInput.Blur()
		m.commandInput.SetValue("")

		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Cancel add-file mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = TuiModeNormal
		m.commandInput.Blur()
		m.commandInput.SetValue("")
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
