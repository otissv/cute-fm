package main

import (
	"fmt"
	"os"

	"cute/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	startDir := ""
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	m := tui.InitialModel(startDir)

	tui.InjectIntoModel(&m)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
