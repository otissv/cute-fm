package tui

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"cute/command"
	"cute/filesystem"
)

func SetQuitMode() {
	if ActiveTuiMode != ModeQuit {
		PreviousTuiMode = ActiveTuiMode
		ActiveTuiMode = ModeQuit
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

		if ActiveTuiMode == ModeAddFile {
			return m.UtilityMode(msg, "touch")
		}

		if ActiveTuiMode == ModeCd {
			return m.UtilityMode(msg, "cd")
		}

		if ActiveTuiMode == ModeColumnVisibiliy {
			return m.ColumnVisibiliyMode(msg)
		}

		if ActiveTuiMode == ModeCommand {
			return m.CommandMode(msg)
		}

		if ActiveTuiMode == ModeCopy {
			return m.UtilityMode(msg, "cp")
		}

		if ActiveTuiMode == ModeFileListSplitPane {
			return m.FileListSplitPaneMode(msg)
		}

		if ActiveTuiMode == ModeFilter {
			return m.FilterMode(msg)
		}

		if ActiveTuiMode == ModeGoto {
			return m.GotoMode(msg)
		}

		if ActiveTuiMode == ModeHelp {
			return m.HelpMode(msg)
		}

		if ActiveTuiMode == ModeMkdir {
			return m.UtilityMode(msg, "mkdir")
		}
		if ActiveTuiMode == ModeMove {
			return m.UtilityMode(msg, "mv")
		}

		if ActiveTuiMode == ModeSelect {
			return m.SelectMode(msg)
		}

		if ActiveTuiMode == ModeNormal {
			return m.NormalMode(msg)
		}

		if ActiveTuiMode == ModeQuit {
			return m.QuitMode(msg)
		}

		if ActiveTuiMode == ModeRename {
			return m.UtilityMode(msg, "rename")
		}

		if ActiveTuiMode == ModeRemove {
			return m.ConfirmMode(msg, "rm -r")
		}

		if ActiveTuiMode == ModeSort {
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
	m.commandHistory = m.LoadCommandHistory()
	env := m.GetCommandEnvironment()

	res, err := command.Execute(env, line)

	return res, err
}

func (m *Model) GetSelectedEntry() *command.SelectedEntry {
	pane := m.GetActivePane()

	selectedIdx := pane.fileList.Index()
	if selectedIdx < 0 || selectedIdx >= len(pane.files) {
		return nil
	}

	fi := pane.files[selectedIdx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(pane.currentDir, fi.Name)
	}

	return &command.SelectedEntry{
		Name:  fi.Name,
		Path:  path,
		IsDir: fi.IsDir,
		Type:  fi.Type,
	}
}

func (m *Model) GetCommandEnvironment() command.Environment {
	pane := m.GetActivePane()

	return command.Environment{
		Cwd:      pane.currentDir,
		Config:   m.runtimeConfig,
		Selected: m.GetSelectedEntry(),
	}
}

func (m *Model) ApplyFilter() {
	pane := m.GetActivePane()
	query := strings.TrimSpace(pane.filterQuery)

	if len(pane.allFiles) == 0 {
		return
	}

	base := filterByViewMode(pane.allFiles)

	if query == "" {
		pane.files = base
	} else {
		lq := strings.ToLower(query)
		var filtered []filesystem.FileInfo
		for _, fi := range base {
			if strings.Contains(strings.ToLower(fi.Name), lq) {
				filtered = append(filtered, fi)
			}
		}
		pane.files = filtered
	}

	m.applySorting(pane)
	items := FileInfosToItems(pane.files, pane.marked)
	pane.fileList.SetItems(items)

	if len(pane.files) > 0 {
		pane.fileList.Select(0)
	}

	m.UpdateFileInfoPane()
}

func (m *Model) changeDirectoryInternal(dir string, trackHistory bool) {
	pane := m.GetActivePane()

	if trackHistory && pane.currentDir != "" && pane.currentDir != dir {
		pane.dirBackStack = append(pane.dirBackStack, pane.currentDir)
		pane.dirForwardStack = nil
	}

	files, err := filesystem.ListDirectory(dir)
	if err != nil {
		m.fileInfoViewport.SetContent("Error reading directory:\n" + err.Error())
		return
	}

	pane.currentDir = dir
	pane.allFiles = files
	pane.files = files
	pane.marked = make(map[string]bool)

	items := FileInfosToItems(files, pane.marked)
	pane.fileList.SetItems(items)

	if len(files) > 0 {
		pane.fileList.Select(0)
	}

	m.ApplyFilter()

	m.UpdateFileInfoPane()
}

func (m *Model) ChangeDirectory(dir string) {
	m.changeDirectoryInternal(dir, true)
}

func (m *Model) ReloadDirectory() {
	pane := m.GetActivePane()

	if pane.currentDir == "" {
		return
	}
	m.changeDirectoryInternal(pane.currentDir, false)
}

func (m *Model) NavigatePreviousDir() {
	pane := m.GetActivePane()

	if len(pane.dirBackStack) == 0 {
		return
	}

	lastIdx := len(pane.dirBackStack) - 1
	prevDir := pane.dirBackStack[lastIdx]
	pane.dirBackStack = pane.dirBackStack[:lastIdx]

	if pane.currentDir != "" && pane.currentDir != prevDir {
		pane.dirForwardStack = append(pane.dirForwardStack, pane.currentDir)
	}

	m.changeDirectoryInternal(prevDir, false)
}

func (m *Model) NavigateNextDir() {
	pane := m.GetActivePane()

	if len(pane.dirForwardStack) == 0 {
		return
	}

	lastIdx := len(pane.dirForwardStack) - 1
	nextDir := pane.dirForwardStack[lastIdx]
	pane.dirForwardStack = pane.dirForwardStack[:lastIdx]

	if pane.currentDir != "" && pane.currentDir != nextDir {
		pane.dirBackStack = append(pane.dirBackStack, pane.currentDir)
	}

	m.changeDirectoryInternal(nextDir, false)
}

func filterByViewMode(files []filesystem.FileInfo) []filesystem.FileInfo {
	if ActiveFileListMode == "" {
		ActiveFileListMode = "ll"
	}

	out := make([]filesystem.FileInfo, 0, len(files))

	switch ActiveFileListMode {
	case "ls":
		// Hide dotfiles (roughly emulating eza/ls without -a)
		for _, fi := range files {
			if strings.HasPrefix(fi.Name, ".") {
				continue
			}
			out = append(out, fi)
		}
	case "ld":
		// Only directories (hidden or not)
		for _, fi := range files {
			if fi.IsDir || fi.Type == "directory" {
				out = append(out, fi)
			}
		}
	case "lf":
		// Only files (hidden or not)
		for _, fi := range files {
			if !fi.IsDir {
				out = append(out, fi)
			}
		}
	default:
		// "ll"  or any unknown mode: show everything
		out = append(out, files...)
	}

	return out
}

// parseHumanSize converts a human-readable size string (e.g. "1.3k", "5.7M")
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

func parseDateModified(s string) time.Time {
	const layout = "02 Jan 15:04"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}
	}

	now := time.Now()
	return time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
}

func (m *Model) applySorting(pane *filePane) {
	if len(pane.files) == 0 {
		return
	}

	sortBy := m.sortColumnBy
	if sortBy.column == "" {
		return
	}

	sort.SliceStable(pane.files, func(i, j int) bool {
		a := pane.files[i]
		b := pane.files[j]

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
