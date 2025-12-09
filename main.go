package main

import (
	"fmt"
	"os"

	"cute/components"
	"cute/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	// Determine starting directory: if a positional argument is provided,
	// use it as the initial directory; otherwise fall back to the current
	// working directory (handled inside InitialModel).
	startDir := ""
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	// Create the initial model
	m := tui.InitialModel(startDir)

	// Inject UI components .

	m.SearchBar = components.SearchBar
	m.CurrentDir = components.CurrentDir
	m.Header = components.Header
	m.StatusBar = components.StatusBar
	m.ViewModeText = components.ViewModeText
	m.FileInfo = components.FileInfo
	m.FileListView = components.FileList
	m.TuiMode = components.TuiMode
	m.SudoMode = components.SudoMode

	// Inject UI modals.
	m.HelpModal = components.HelpModal
	m.CommandModal = components.CommandModal
	m.DialogModal = components.DialogModal
	m.ColumnModal = components.ColumnModal

	// Create a new Bubble Tea program
	p := tea.NewProgram(m)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
