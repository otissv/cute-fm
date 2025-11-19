package main

import (
	"fmt"
	"os"

	"lsfm/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the initial model
	m := tui.InitialModel()

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
