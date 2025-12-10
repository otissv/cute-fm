package tui

import (
	"strings"

	"cute/console"

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

			if (command == "cp" || command == "mv") && selectedEntry != nil {
				line = command + " " + selectedEntry.Path + " " + inputValue
			}

			if command == "rename" {
				line = "mv " + selectedEntry.Name + " " + inputValue

				console.Log("%s", line)
			}

			res, _ := m.ExecuteCommand(line)

			pane := m.GetActivePane()
			if res.Cwd != "" && res.Cwd != pane.currentDir {
				m.ChangeDirectory(res.Cwd)
			} else if res.Refresh {
				m.ReloadDirectory()
			}
		}

		m.commandInput.Blur()
		m.commandInput.SetValue("")

		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Cancel add-file mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		m.commandInput.Blur()
		m.commandInput.SetValue("")
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
