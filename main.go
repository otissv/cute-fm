package main

import (
	"fmt"
	"os"

	"cute/tui"

	tea "github.com/charmbracelet/bubbletea"
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

	// Create a new Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(), // Enable alternate screen buffer
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
