package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type Menu struct {
	choices     []string
	cursor      int
	selected    map[string]bool
	cursorTypes MenuCursor
}

type MenuCursor struct {
	Selected   string
	Unselected string
	Prompt     string
}

type MenuArgs struct {
	Choices     []string
	Cursor      int
	Selected    map[string]string
	CursorTypes MenuCursor
}

func NewMenu(args MenuArgs) Menu {
	selectedSet := make(map[string]bool, len(args.Selected))
	for _, col := range args.Selected {
		selectedSet[col] = true
	}

	selected := "x"
	unselected := " "
	propmt := ">"

	if args.CursorTypes.Selected != "" {
		selected = args.CursorTypes.Selected
	}

	if args.CursorTypes.Unselected != "" {
		unselected = args.CursorTypes.Unselected
	}

	if args.CursorTypes.Prompt != "" {
		propmt = args.CursorTypes.Prompt
	}

	current := lipgloss.NewStyle().
		Bold(true).
		Render(propmt)

	return Menu{
		choices:  args.Choices,
		cursor:   args.Cursor,
		selected: selectedSet,
		cursorTypes: MenuCursor{
			Selected:   selected,
			Unselected: unselected,
			Prompt:     current,
		},
	}
}

func (menu Menu) View() string {
	var b strings.Builder
	for i, choice := range menu.choices {
		iChoice := choice
		cursor := "[ " + menu.cursorTypes.Unselected + " ]"
		if menu.selected != nil && menu.selected[iChoice] {
			cursor = "[ " + menu.cursorTypes.Selected + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}

		if i == menu.cursor {
			cursor = "[ " + menu.cursorTypes.Prompt + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, iChoice)
	}
	return strings.TrimRight(b.String(), "\n")
}
