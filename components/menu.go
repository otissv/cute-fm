package components

import (
	"fmt"
	"strings"

	"cute/filesystem"

	"charm.land/lipgloss/v2"
)

type MenuModel struct {
	choices  []string
	cursor   int
	selected map[string]bool
}

func Menu(choices []string, cursor int, selected []filesystem.FileInfoColumn) MenuModel {
	selectedSet := make(map[string]bool, len(selected))
	for _, col := range selected {
		selectedSet[string(col)] = true
	}

	return MenuModel{
		choices:  choices,
		cursor:   cursor,
		selected: selectedSet,
	}
}

func (m MenuModel) View() string {
	var b strings.Builder
	for i, choice := range m.choices {
		iChoice := choice
		cursor := "[   ]"
		if m.selected != nil && m.selected[iChoice] {
			cursor = "[ x ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}

		if i == m.cursor {
			prompt := lipgloss.NewStyle().
				Bold(true).
				Render(">")

			cursor = "[ " + prompt + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, iChoice)
	}
	return strings.TrimRight(b.String(), "\n")
}
