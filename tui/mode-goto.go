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

		if PreviousTuiMode == ModeComputer {
			m.applyRelativeGotoComputerList(inputValue)
		} else {
			m.applyRelativeGoto(inputValue)
		}

		m.commandInput.Blur()
		m.commandInput.SetValue("")
		m.jumpTo = ""

		ActiveTuiMode = PreviousTuiMode
		return m, nil

	// Cancel goto mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		m.commandInput.Blur()
		m.commandInput.SetValue("")
		m.jumpTo = ""
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) applyRelativeGoto(inputValue string) bool {
	inputValue = strings.TrimSpace(inputValue)
	pane := m.GetActivePane()

	if inputValue == "" || len(pane.files) == 0 {
		return false
	}

	moveBackward := false

	// A trailing "-" or sytarts with "-" means "move up".
	if strings.HasSuffix(inputValue, "-") || strings.HasPrefix(inputValue, "-") {
		moveBackward = true
		inputValue = strings.TrimSpace(strings.TrimSuffix(inputValue, "-"))
		if inputValue == "" {
			return false
		}
	}

	n, err := strconv.Atoi(inputValue)
	if err != nil {
		return false
	}

	// A leading "-" also means "move up".
	if n < 0 {
		moveBackward = true
		n = -n
	}

	if n <= 0 {
		return false
	}

	current := pane.fileList.Index()
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
	if target >= len(pane.files) {
		target = len(pane.files) - 1
	}

	if target == current {
		return false
	}

	pane.fileList.Select(target)
	m.UpdateFileInfoPane()

	return true
}

func (m *Model) applyRelativeGotoComputerList(inputValue string) bool {
	inputValue = strings.TrimSpace(inputValue)
	items := m.computerList.Items()

	if inputValue == "" || len(items) == 0 {
		return false
	}

	moveBackward := false

	// A trailing "-" or starts with "-" means "move up".
	if strings.HasSuffix(inputValue, "-") || strings.HasPrefix(inputValue, "-") {
		moveBackward = true
		inputValue = strings.TrimSpace(strings.TrimSuffix(inputValue, "-"))
		if inputValue == "" {
			return false
		}
	}

	n, err := strconv.Atoi(inputValue)
	if err != nil {
		return false
	}

	// A leading "-" also means "move up".
	if n < 0 {
		moveBackward = true
		n = -n
	}

	if n <= 0 {
		return false
	}

	current := m.computerList.Index()
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
	if target >= len(items) {
		target = len(items) - 1
	}

	if target == current {
		return false
	}

	m.computerList.Select(target)

	return true
}
