package tui

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"cute/command"
	"cute/console"
	"cute/filesystem"
)

func SetQuitMode() {
	if ActiveTuiMode != TuiModeQuit {
		PreviousTuiMode = ActiveTuiMode
		ActiveTuiMode = TuiModeQuit
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		m.width = msg.Width
		m.height = msg.Height + 2

		helpWidth := msg.Width / 2
		helpHeight := msg.Height / 2
		if helpWidth < 20 {
			helpWidth = 20
		}
		if helpHeight < 5 {
			helpHeight = 5
		}

		m.CalcLayout()

		return m, nil

	case tea.KeyMsg:

		console.Log("Mode: %s, Key Pressed: %s", ActiveTuiMode, msg.String())
		if ActiveTuiMode == TuiModeAddFile {
			return m.UtilityMode(msg, "touch")
		}

		if ActiveTuiMode == TuiModeCd {
			return m.UtilityMode(msg, "cd")
		}

		if ActiveTuiMode == TuiModeColumnVisibiliy {
			return m.ColumnVisibiliyMode(msg)
		}

		if ActiveTuiMode == TuiModeCommand {
			return m.CommandMode(msg)
		}

		if ActiveTuiMode == TuiModeCopy {
			return m.UtilityMode(msg, "cp")
		}

		if ActiveTuiMode == TuiModeFilter {
			return m.FilterMode(msg)
		}

		if ActiveTuiMode == TuiModeGoto {
			return m.GotoMode(msg)
		}

		if ActiveTuiMode == TuiModeHelp {
			return m.HelpMode(msg)
		}

		if ActiveTuiMode == TuiModeMkdir {
			return m.UtilityMode(msg, "mkdir")
		}
		if ActiveTuiMode == TuiModeMove {
			return m.UtilityMode(msg, "mv")
		}

		if ActiveTuiMode == TuiModeNormal {
			return m.NormalMode(msg)
		}

		if ActiveTuiMode == TuiModeQuit {
			return m.QuitMode(msg)
		}

		if ActiveTuiMode == TuiModeRename {
			return m.UtilityMode(msg, "rename")
		}

		if ActiveTuiMode == TuiModeRemove {
			return m.ConfirmMode(msg, "rm -r")
		}

		if ActiveTuiMode == TuiModeSort {
			return m.SortMode(msg)
		}
	}

	return m, nil
}

func (m *Model) ExecuteCommand(line string) (command.Result, error) {
	if line == "" {
		return command.Result{}, nil
	}

	m.AppendCommandHistory(line)
	// Reload history to include the new command
	m.commandHistory = m.LoadCommandHistory()

	env := m.GetCommandEnvironment()

	res, err := command.Execute(env, line)

	return res, err
}

func (m *Model) GetSelectedEntry() *command.SelectedEntry {
	selectedIdx := m.fileList.Index()
	if selectedIdx < 0 || selectedIdx >= len(m.files) {
		return nil
	}

	fi := m.files[selectedIdx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(m.currentDir, fi.Name)
	}

	return &command.SelectedEntry{
		Name:  fi.Name,
		Path:  path,
		IsDir: fi.IsDir,
		Type:  fi.Type,
	}
}

// GetCommandEnvironment builds the command execution environment using the
// current model state, including the currently selected entry (if any).
func (m *Model) GetCommandEnvironment() command.Environment {
	return command.Environment{
		Cwd:      m.currentDir,
		Config:   m.runtimeConfig,
		Selected: m.GetSelectedEntry(),
	}
}

// ApplyFilter recomputes the visible file list based on the current value of
// the text input. The filter is a case-insensitive substring match on the file
// name. When the filter changes, the list is updated with the new items.
func (m *Model) ApplyFilter() {
	query := strings.TrimSpace(m.searchInput.Value())

	// If there is no backing data yet, nothing to do.
	if len(m.allFiles) == 0 {
		return
	}

	base := filterByViewMode(m.allFiles)

	// Then apply the search query on top.
	if query == "" {
		m.files = base
	} else {
		lq := strings.ToLower(query)
		var filtered []filesystem.FileInfo
		for _, fi := range base {
			if strings.Contains(strings.ToLower(fi.Name), lq) {
				filtered = append(filtered, fi)
			}
		}
		m.files = filtered
	}

	// Apply the active column-based sorting to the filtered list.
	m.applySorting()

	// Update the list with new items.
	items := FileInfosToItems(m.files)
	m.fileList.SetItems(items)

	// Reset selection to first item if we have items.
	if len(m.files) > 0 {
		m.fileList.Select(0)
	}

	// Update preview for the new selection after filtering.
	m.UpdateFileInfoPanel()
}

// changeDirectoryInternal updates the model to point at a new current
// directory and reloads the file list. When trackHistory is true, the
// previous directory is pushed onto the "back" stack and the "forward"
// stack is cleared, mirroring typical browser navigation behaviour.
func (m *Model) changeDirectoryInternal(dir string, trackHistory bool) {
	// Record history before we actually change directories.
	if trackHistory && m.currentDir != "" && m.currentDir != dir {
		m.dirBackStack = append(m.dirBackStack, m.currentDir)
		// Any new navigation invalidates the "forward" history.
		m.dirForwardStack = nil
	}

	files, err := filesystem.ListDirectory(dir)
	if err != nil {
		m.fileInfoViewport.SetContent("Error reading directory:\n" + err.Error())
		return
	}

	m.currentDir = dir
	m.allFiles = files
	m.files = files

	// Update the list with new items.
	items := FileInfosToItems(files)
	m.fileList.SetItems(items)

	// Select the first item if we have items.
	if len(files) > 0 {
		m.fileList.Select(0)
	}

	// Re-apply search/view filters for the new directory.
	m.ApplyFilter()

	// And recompute the preview for the new directory/selection.
	m.UpdateFileInfoPanel()
}

// ChangeDirectory is the public helper used throughout the TUI when the user
// explicitly navigates to a new directory (e.g. via Enter, :cd, parent, etc.).
// It records the change in the navigation history so that PreviousDir/NextDir
// can traverse it.
func (m *Model) ChangeDirectory(dir string) {
	m.changeDirectoryInternal(dir, true)
}

// ReloadDirectory reloads the current directory without adding a new history
// entry. This is used when commands request a simple refresh of the listing.
func (m *Model) ReloadDirectory() {
	if m.currentDir == "" {
		return
	}
	m.changeDirectoryInternal(m.currentDir, false)
}

// NavigatePreviousDir moves to the previously visited directory, if any.
// It updates both the back and forward stacks so that repeated invocations
// allow walking backward through the navigation history.
func (m *Model) NavigatePreviousDir() {
	if len(m.dirBackStack) == 0 {
		return
	}

	// Pop the last entry from the back stack.
	lastIdx := len(m.dirBackStack) - 1
	prevDir := m.dirBackStack[lastIdx]
	m.dirBackStack = m.dirBackStack[:lastIdx]

	// Current directory becomes a "forward" target.
	if m.currentDir != "" && m.currentDir != prevDir {
		m.dirForwardStack = append(m.dirForwardStack, m.currentDir)
	}

	// Do not record this as a new history entry; we're traversing history.
	m.changeDirectoryInternal(prevDir, false)
}

// NavigateNextDir moves forward in the directory history, if possible.
func (m *Model) NavigateNextDir() {
	if len(m.dirForwardStack) == 0 {
		return
	}

	// Pop the last entry from the forward stack.
	lastIdx := len(m.dirForwardStack) - 1
	nextDir := m.dirForwardStack[lastIdx]
	m.dirForwardStack = m.dirForwardStack[:lastIdx]

	// Current directory becomes part of the "back" history.
	if m.currentDir != "" && m.currentDir != nextDir {
		m.dirBackStack = append(m.dirBackStack, m.currentDir)
	}

	// Do not record this as a new history entry; we're traversing history.
	m.changeDirectoryInternal(nextDir, false)
}

// filterByViewMode filters the given file list according to the current view
// mode. It does not modify the original slice.
func filterByViewMode(files []filesystem.FileInfo) []filesystem.FileInfo {
	if ActiveFileListMode == "" {
		ActiveFileListMode = "ll"
	}

	out := make([]filesystem.FileInfo, 0, len(files))

	switch ActiveFileListMode {
	case "ls":
		// Hide dotfiles (roughly emulating eza/ls without -a).
		for _, fi := range files {
			if strings.HasPrefix(fi.Name, ".") {
				continue
			}
			out = append(out, fi)
		}
	case "ld":
		// Only directories (hidden or not). Prefer the IsDir flag, but also
		// fall back to the classified file type to be defensive.
		for _, fi := range files {
			if fi.IsDir || fi.Type == "directory" {
				out = append(out, fi)
			}
		}
	case "lf":
		// Only files (hidden or not).
		for _, fi := range files {
			if !fi.IsDir {
				out = append(out, fi)
			}
		}
	default:
		// "ll"  or any unknown mode: show everything.
		out = append(out, files...)
	}

	return out
}

// parseHumanSize converts a human-readable size string (e.g. "1.3k", "5.7M")
// into an approximate byte count suitable for numeric comparisons. It mirrors
// the format produced by filesystem.formatSize, which uses base-10 units.
func parseHumanSize(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// Detect optional unit suffix.
	n := len(s)
	last := s[n-1]

	mult := float64(1)
	switch last {
	case 'k', 'K':
		mult = 1_000
		s = s[:n-1]
	case 'M':
		mult = 1_000_000
		s = s[:n-1]
	case 'G':
		mult = 1_000_000_000
		s = s[:n-1]
	case 'T':
		mult = 1_000_000_000_000
		s = s[:n-1]
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(val * mult)
}

// parseDateModified parses the formatted DateModified string ("02 Jan 15:04")
// back into a time.Time value for chronological comparisons. Because the
// original format omits the year, we assume the current year when constructing
// the timestamp.
func parseDateModified(s string) time.Time {
	const layout = "02 Jan 15:04"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}
	}

	now := time.Now()
	return time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
}

// applySorting sorts m.files in-place according to the current SortColumnBy
// configuration on the model. When no sort column is set, the slice is left
// unchanged.
func (m *Model) applySorting() {
	if len(m.files) == 0 {
		return
	}

	sortBy := m.sortColumnBy
	if sortBy.column == "" {
		return
	}

	sort.SliceStable(m.files, func(i, j int) bool {
		a := m.files[i]
		b := m.files[j]

		var less bool

		switch sortBy.column {
		case filesystem.ColumnPermissions:
			less = a.Permissions < b.Permissions
		case filesystem.ColumnSize:
			less = parseHumanSize(a.Size) < parseHumanSize(b.Size)
		case filesystem.ColumnMimeType:
			less = a.MimeType < b.MimeType
		case filesystem.ColumnUser:
			less = strings.ToLower(a.User) < strings.ToLower(b.User)
		case filesystem.ColumnGroup:
			less = strings.ToLower(a.Group) < strings.ToLower(b.Group)
		case filesystem.ColumnDateModified:
			less = parseDateModified(a.DateModified).Before(parseDateModified(b.DateModified))
		case filesystem.ColumnName:
			// Preserve the existing "directories first" behaviour when sorting
			// by name so the default view remains intuitive.
			if a.IsDir && !b.IsDir {
				return true
			}
			if !a.IsDir && b.IsDir {
				return false
			}
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		default:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}

		if sortBy.direction == SortingDesc {
			return !less
		}
		return less
	})
}
