package tui

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) GotoMode(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	// Keep the jumpTo status text in sync with the current input value.
	m.jumpTo = m.commandInput.Value()

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Confirm goto on Enter.
	case bindings.Enter.Matches(keyMsg.String()):
		inputValue := strings.TrimSpace(m.commandInput.Value())

		if inputValue != "" && len(m.files) > 0 {
			// If the input ends with "-", move the selection up by that many rows.
			// Otherwise, move the selection down by that many rows.
			moveBackward := false
			if strings.HasSuffix(inputValue, "-") {
				moveBackward = true
				inputValue = strings.TrimSpace(strings.TrimSuffix(inputValue, "-"))
			}

			if inputValue != "" {
				if n, err := strconv.Atoi(inputValue); err == nil {
					if n < 0 {
						n = -n
					}
					if n > 0 {
						current := m.fileList.Index()
						if current < 0 {
							current = 0
						}

						var target int
						if moveBackward {
							target = current - n
						} else {
							target = current + n
						}

						if target < 0 {
							target = 0
						}
						if target >= len(m.files) {
							target = len(m.files) - 1
						}

						m.fileList.Select(target)
						m.UpdatePreview()
					}
				}
			}
		}

		m.commandInput.Blur()
		m.commandInput.SetValue("")
		m.jumpTo = ""

		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Cancel goto mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = TuiModeNormal
		m.commandInput.Blur()
		m.commandInput.SetValue("")
		m.jumpTo = ""
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
