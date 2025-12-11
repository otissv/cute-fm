package tui

import (
	"cute/filesystem"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

type ACTIVE_SECTION string

const (
	SECTION_START                 ACTIVE_SECTION = "SECTION_START"
	SECTION_SPLITPANE             ACTIVE_SECTION = "SECTION_SPLITPANE"
	SECTION_FILE_LIST_MODE        ACTIVE_SECTION = "SECTION_FILE_LIST_MODE"
	SECTION_COLUMN_VISIBILITY     ACTIVE_SECTION = "SECTION_COLUMN_VISIBILITY"
	SECTION_SORT_BY_COLUMN        ACTIVE_SECTION = "SECTION_SORT_BY_COLUMN"
	SECTIOn_SORT_COLUMN_DIRECTION ACTIVE_SECTION = "SECTION_SORT_COLUMN_DIRECTION"
)

func SettingsWindow(m Model) *lipgloss.Layer {
	theme := m.GetTheme()
	width, height := m.GetSize()

	// activeSection := SECTION_START

	startDirSection := startDir()
	splitpaneSection := splitpane()
	fileListModeSection := fileListMode()
	columnVisibiliySection := columnVisibiliy()
	sortByColumnSection := sortByColumn()
	sortByColumnDirectionSection := sortByColumnDirection()

	contentItems := []string{}
	contentItems = append(contentItems, startDirSection...)
	contentItems = append(contentItems, splitpaneSection...)
	contentItems = append(contentItems, fileListModeSection...)
	contentItems = append(contentItems, columnVisibiliySection...)
	contentItems = append(contentItems, sortByColumnSection...)
	contentItems = append(contentItems, sortByColumnDirectionSection...)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		contentItems...,
	)

	fw := FloatingWindow{
		Content: ViewPrimitive(content),
		Width:   50,
		Height:  50,
		Style:   DefaultFloatingStyle(theme),
		Title:   "Settings",
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}

func menuHeading() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Left)
}

func menuChoices(choices []string) []string {
	items := []string{}
	style := lipgloss.NewStyle().Align(lipgloss.Left)

	for _, choice := range choices {
		items = append(items, style.Render(choice))
	}

	return items
}

func startDir() []string {
	startDirInput := textinput.New()
	startDirInput.Placeholder = "Custom..."
	startDirInputStyle := lipgloss.NewStyle().Render(startDirInput.View())

	Choices := menuChoices([]string{
		"Home directory",
		"Current directory",
		startDirInputStyle,
	})

	Cursor := 0
	selected := map[string]string{
		"Home directory": "Home directory",
	}

	menu := NewMenu(MenuArgs{
		Choices:  Choices,
		Cursor:   Cursor,
		Selected: selected,
	})

	return []string{
		menuHeading().Render("Start Directory"),
		menu.View(),
	}
}

func splitpane() []string {
	choices := []string{
		"None",
		"File Info",
		"File List",
		"Preview",
	}

	cursor := 0
	selected := map[string]string{
		"File Info": "File Info",
	}

	menu := NewMenu(MenuArgs{
		Choices:  choices,
		Cursor:   cursor,
		Selected: selected,
	})

	return []string{
		menuHeading().
			PaddingTop(1).
			Render("Split pane"),
		menu.View(),
	}
}

func fileListMode() []string {
	choices := []string{
		"List all files",
		"Directories only",
		"Files only",
	}

	cursor := 0
	selected := map[string]string{
		"List all files": "List all files",
	}

	menu := NewMenu(MenuArgs{
		Choices:  choices,
		Cursor:   cursor,
		Selected: selected,
	})

	return []string{
		menuHeading().
			PaddingTop(1).
			Render("File List Modes"),
		menu.View(),
	}
}

func columnVisibiliy() []string {
	choices := make([]string, len(filesystem.ColumnNames))
	for i, c := range filesystem.ColumnNames {
		choices[i] = string(c)
	}

	cursor := 0
	selected := map[string]string{
		string(filesystem.ColumnPermissions):  string(filesystem.ColumnPermissions),
		string(filesystem.ColumnSize):         string(filesystem.ColumnSize),
		string(filesystem.ColumnUser):         string(filesystem.ColumnUser),
		string(filesystem.ColumnGroup):        string(filesystem.ColumnGroup),
		string(filesystem.ColumnDateModified): string(filesystem.ColumnDateModified),
		string(filesystem.ColumnName):         string(filesystem.ColumnName),
	}

	menu := NewMenu(MenuArgs{
		Choices:  choices,
		Cursor:   cursor,
		Selected: selected,
	})

	return []string{
		menuHeading().
			PaddingTop(1).
			Render("Column Visibility"),
		menu.View(),
	}
}

func sortByColumn() []string {
	choices := make([]string, len(filesystem.ColumnNames))
	for i, c := range filesystem.ColumnNames {
		choices[i] = string(c)
	}

	cursor := 0
	selected := map[string]string{
		"Name": "Name",
	}

	menu := NewMenu(MenuArgs{
		Choices:  choices,
		Cursor:   cursor,
		Selected: selected,
	})

	return []string{
		menuHeading().
			PaddingTop(1).
			Render("Sory By Column"),
		menu.View(),
	}
}

func sortByColumnDirection() []string {
	choices := []string{
		string(SortingAsc),
		string(SortingDesc),
	}

	cursor := 0
	StartDirSelected := make(map[string]string, 1)

	menu := NewMenu(MenuArgs{
		Choices:  choices,
		Cursor:   cursor,
		Selected: StartDirSelected,
	})

	return []string{
		menuHeading().
			PaddingTop(1).
			Render("Sort Column Direction"),
		menu.View(),
	}
}
