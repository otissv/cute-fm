package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) CommandMode(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	beforeValue := m.commandInput.Value()
	m.commandInput, cmd = m.commandInput.Update(msg)
	cmds = append(cmds, cmd)

	// Update right viewport (second row, right column)
	m.rightViewport, cmd = m.rightViewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update history matches when input changes
	if m.commandInput.Value() != beforeValue {
		m.updateHistoryMatches()
	}

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Enter normal mode
	case bindings.Command.Matches(keyMsg.String()) ||
		bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Auto complete command
	case bindings.AutoComplete.Matches(keyMsg.String()):
		m.completeCommand()
		return m, nil

	// Get previous command
	case bindings.Up.Matches(keyMsg.String()):
		if len(m.commandHistory) > 0 {
			m.navigateHistory(-1)
			return m, nil
		}

	// Get next command
	case bindings.Down.Matches(keyMsg.String()):
		if len(m.commandHistory) > 0 {
			m.navigateHistory(1)
			return m, nil
		}

	// Execute command
	case bindings.Enter.Matches(keyMsg.String()):
		line := strings.TrimSpace(m.commandInput.Value())

		res, err := m.ExecuteCommand(line)

		// Apply environment changes.
		if res.Cwd != "" && res.Cwd != m.currentDir {
			m.ChangeDirectory(res.Cwd)
		} else if res.Refresh {
			// Re-list the current directory when requested by the command.
			m.ChangeDirectory(m.currentDir)
		}

		// Update view mode and re-apply filters so the file list view
		// actually changes when commands like "ll", "ls", "ld", "lf",
		// etc. are executed.
		if res.ViewMode != "" {
			ActiveFileListMode = FileListMode(res.ViewMode)
			m.ApplyFilter()
		}

		if res.OpenHelp {
			m.activeModal = ModalHelp
		}

		if res.Output != "" {
			m.rightViewport.SetContent(res.Output)
		}

		if err != nil && res.Output == "" {
			m.rightViewport.SetContent(err.Error())
		}

		m.commandInput.Blur()
		m.commandInput.SetValue("")
		m.searchInput.Focus()

		m.CalcLayout()

		if res.Quit {
			return m, tea.Quit
		}

		ActiveTuiMode = PreviousTuiMode

		return m, nil

	}

	return m, tea.Batch(cmds...)
}
