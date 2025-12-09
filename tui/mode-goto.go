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

		// Delegate the actual movement to a shared helper so that command
		// mode and goto mode both support the same relative-jump syntax:
		//
		//   "10"   -> move 10 rows down
		//   "10-"  -> move 10 rows up
		//   "-10"  -> move 10 rows up
		m.applyRelativeGoto(inputValue)

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

// applyRelativeGoto moves the current selection in the file list according to
// a Vim-style relative offset encoded in inputValue.
//
// Supported forms (after trimming whitespace):
//
//	"10"   -> move 10 rows down
//	"10-"  -> move 10 rows up
//	"-10"  -> move 10 rows up
//
// It returns true if a valid movement was performed, or false if the input was
// not a valid relative offset or if there are no files to move between.
func (m *Model) applyRelativeGoto(inputValue string) bool {
	inputValue = strings.TrimSpace(inputValue)
	pane := m.GetActivePane()

	if inputValue == "" || len(pane.files) == 0 {
		return false
	}

	moveBackward := false

	// A trailing "-" means "move up".
	if strings.HasSuffix(inputValue, "-") {
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
	m.UpdateFileInfoPanel()

	return true
}
