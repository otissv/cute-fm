package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) AddFileMode(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	// Confirm and create the file on Enter.
	case bindings.Enter.Matches(keyMsg.String()):
		line := strings.TrimSpace(m.commandInput.Value())
		if line != "" {
			res, _ := m.ExecuteCommand("touch " + line)

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
