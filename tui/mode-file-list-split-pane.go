package tui

import tea "charm.land/bubbletea/v2"

func (m Model) FileListSplitPaneMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate all other keys to normal mode behaviour so navigation and
	// commands work as expected while staying in select mode unless they
	// explicitly change the TUI mode.
	return m.NormalMode(msg)
}
