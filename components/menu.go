package components

import (
	"fmt"
	"strings"

	"cute/tui"

	"charm.land/lipgloss/v2"
)

type Menu struct {
	choices     []string
	cursor      int
	selected    map[string]bool
	cursorTypes MenuCursor
}

type MenuCursor struct {
	selected   string
	unselected string
	propmt     string
}

func NewMenu(args tui.MenuArgs) Menu {
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
			selected:   selected,
			unselected: unselected,
			propmt:     current,
		},
	}
}

func (menu Menu) View() string {
	var b strings.Builder
	for i, choice := range menu.choices {
		iChoice := choice
		cursor := "[ " + menu.cursorTypes.unselected + " ]"
		if menu.selected != nil && menu.selected[iChoice] {
			cursor = "[ " + menu.cursorTypes.selected + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}

		if i == menu.cursor {
			cursor = "[ " + menu.cursorTypes.propmt + " ]"
			iChoice = lipgloss.NewStyle().
				Bold(true).
				Render(choice)
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, iChoice)
	}
	return strings.TrimRight(b.String(), "\n")
}
