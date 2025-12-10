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
	pane := m.GetActivePane()
	before := pane.filterQuery

	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	// Keep the pane-specific filter query in sync with the text input so
	// that each pane remembers its own filter while sharing a single
	// searchInput UI.
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
	case bindings.SwitchBetweenSplitPane.Matches(keyMsg.String()):
		if m.isSplitPaneOpen {
			if m.activeViewport == LeftViewportType {
				m.activeViewport = RightViewportType
			} else {
				m.activeViewport = LeftViewportType
			}

			// After switching panes, load that pane's existing filter text
			// into the shared search input so each pane can be edited
			// independently.
			newPane := m.GetActivePane()
			m.searchInput.SetValue(newPane.filterQuery)
		}
		return m, nil
	}

	return m, tea.Batch(cmds...)
}
