package tui

import (
	"fmt"
	"strings"

	"cute/theming"

	"charm.land/lipgloss/v2"
)

type MenuChoiceType string

const (
	HEADING_CHOICE_TYPE MenuChoiceType = "HEADING_CHOICE_TYPE"
	CHOICE_TYPE         MenuChoiceType = "CHOICE_TYPE"
)

type MenuChoice struct {
	Type  MenuChoiceType
	Label string
}

type Menu struct {
	choices     []MenuChoice
	cursor      int
	selected    map[string]bool
	cursorTypes MenuCursor
	theme       theming.Theme
}

type MenuCursor struct {
	Selected   string
	Unselected string
	Prompt     string
	Numbered   bool
}

type MenuArgs struct {
	Choices     []MenuChoice
	CursorIndex int
	CursorTypes MenuCursor
	Selected    map[string]string
	Theme       theming.Theme
}

func NewMenu(args MenuArgs) Menu {
	selectedSet := make(map[string]bool, len(args.Selected))
	for _, col := range args.Selected {
		selectedSet[col] = true
	}

	selected := "x"
	unselected := " "
	prompt := ">"
	numbered := false

	if args.CursorTypes.Selected != "" {
		selected = args.CursorTypes.Selected
	}

	if args.CursorTypes.Unselected != "" {
		unselected = args.CursorTypes.Unselected
	}

	if args.CursorTypes.Prompt != "" {
		prompt = args.CursorTypes.Prompt
	}

	if args.CursorTypes.Numbered {
		numbered = args.CursorTypes.Numbered
	}

	current := lipgloss.NewStyle().
		Bold(true).
		Render(prompt)

	return Menu{
		choices: args.Choices,
		cursor:  args.CursorIndex,
		cursorTypes: MenuCursor{
			Selected:   selected,
			Unselected: unselected,
			Prompt:     current,
			Numbered:   numbered,
		},
		selected: selectedSet,
		theme:    args.Theme,
	}
}

func (menu Menu) View() string {
	var b strings.Builder

	number := 1

	for i, choice := range menu.choices {
		label := choice.Label
		cursor := "[ " + menu.cursorTypes.Unselected + " ]"
		prefix := ""

		if menu.cursorTypes.Numbered {
			prefix = lipgloss.NewStyle().
				Foreground(lipgloss.Color(menu.theme.Muted)).
				Render(fmt.Sprintf("%d ", number))
		}

		if menu.selected != nil && menu.selected[label] {
			cursor = "[ " + menu.cursorTypes.Selected + " ]"
			label = lipgloss.NewStyle().
				Bold(true).
				Render(label)
		}

		if i == menu.cursor {
			cursor = "[ " + menu.cursorTypes.Prompt + " ]"
			label = lipgloss.NewStyle().
				Bold(true).
				Render(label)
		}

		if choice.Type == HEADING_CHOICE_TYPE {
			label = " " + lipgloss.NewStyle().
				Bold(true).
				Render(label)

			fmt.Fprintf(&b, "\n%s%s\n", prefix, label)

		} else {
			if menu.cursorTypes.Numbered {
				cursor = fmt.Sprintf("%s %s", prefix, cursor)
			}
			fmt.Fprintf(&b, "%s %s\n", cursor, label)
		}

		number += 1
	}
	return strings.TrimRight(b.String(), "\n")
}
