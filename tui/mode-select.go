package tui

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

// toggleCurrentSelection toggles the marked state of the currently focused row
// in the active pane and refreshes the list items so the marker column updates.
func (m *Model) toggleCurrentSelection() {
	pane := m.activePane()
	if len(pane.files) == 0 {
		return
	}

	idx := pane.fileList.Index()
	if idx < 0 || idx >= len(pane.files) {
		return
	}

	if pane.marked == nil {
		pane.marked = make(map[string]bool)
	}

	fi := pane.files[idx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(pane.currentDir, fi.Name)
	}

	if pane.marked[path] {
		delete(pane.marked, path)
	} else {
		pane.marked[path] = true
	}

	pane.fileList.SetItems(FileInfosToItems(pane.files, pane.marked))
}

// toggleSelectAll toggles between marking all visible rows and clearing all
// marks in the active pane.
func (m *Model) toggleSelectAll() {
	pane := m.activePane()
	if len(pane.files) == 0 {
		return
	}

	if pane.marked == nil {
		pane.marked = make(map[string]bool)
	}

	allMarked := true
	for _, fi := range pane.files {
		path := fi.Path
		if path == "" {
			path = filepath.Join(pane.currentDir, fi.Name)
		}
		if !pane.marked[path] {
			allMarked = false
			break
		}
	}

	if allMarked {
		// Clear all marks.
		pane.marked = make(map[string]bool)
	} else {
		// Mark all visible files.
		for _, fi := range pane.files {
			path := fi.Path
			if path == "" {
				path = filepath.Join(pane.currentDir, fi.Name)
			}
			pane.marked[path] = true
		}
	}

	pane.fileList.SetItems(FileInfosToItems(pane.files, pane.marked))
}

func (m Model) SelectMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	bindings := GetKeyBindings()

	// Only handle key messages here; ignore everything else.
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	key := keyMsg.String()

	switch {
	// Quit application
	case bindings.Quit.Matches(key):
		SetQuitMode()
		return m, nil

	// Leave select mode on cancel (esc / ctrl+q).
	case bindings.Cancel.Matches(key):
		ActiveTuiMode = TuiModeNormal
		return m, nil

	// Toggle select all rows in the current pane.
	case bindings.SelectAll.Matches(key):
		m.toggleSelectAll()
		return m, nil

	// Toggle the current row's marked state.
	case bindings.Select.Matches(key):
		m.toggleCurrentSelection()
		return m, nil
	}

	// Delegate all other keys to normal mode behaviour so navigation and
	// commands work as expected while staying in select mode unless they
	// explicitly change the TUI mode.
	return m.NormalMode(msg)
}
