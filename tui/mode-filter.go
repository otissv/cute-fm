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

	pane := m.GetActivePane()
	before := pane.filterQuery

	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	pane.filterQuery = m.searchInput.Value()

	if pane.filterQuery != before {
		m.ApplyFilter()
	}

	switch {
	// Quit application
	case bindings.Quit.Matches(keyMsg.String()):
		SetQuitMode()
		return m, nil

	// Enter normal mode
	case bindings.Cancel.Matches(keyMsg.String()):
		ActiveTuiMode = ModeNormal
		return m, nil

	// Switch panes in file list slipt mode
	case bindings.Tab.Matches(keyMsg.String()):
		if m.isSplitPaneOpen {
			if m.activeViewport == LeftViewportType {
				m.activeViewport = RightViewportType
			} else {
				m.activeViewport = LeftViewportType
			}

			newPane := m.GetActivePane()
			m.searchInput.SetValue(newPane.filterQuery)
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
