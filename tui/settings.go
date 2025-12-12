package tui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

func GetSettings() []MenuChoice {
	customStartDirInput := textinput.New()
	customStartDirInput.Placeholder = "Custom..."
	customStartDirInputStyle := lipgloss.NewStyle().
		Width(20).
		BorderBottom(true).
		Render(customStartDirInput.View())

	settings := []MenuChoice{
		{
			Label: "Start Directory",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "Home directory",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Current directory",
		},
		{
			Label: customStartDirInputStyle,
			Type:  CHOICE_TYPE,
		},

		{
			Label: "Split Pane",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "None",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "File Info",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "File List",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Preview",
		},

		{
			Label: "File List Mode",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "List all files",
		},
		{
			Label: "Directories only",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Files only",
			Type:  CHOICE_TYPE,
		},

		{
			Label: "Column Visibility",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "Permissions",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Size",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Type",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "User",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Group",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "DateModified",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Name",
			Type:  CHOICE_TYPE,
		},

		{
			Label: "Sort Columns By",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "Permissions",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Size",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Type",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "User",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Group",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "DateModified",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Name",
			Type:  CHOICE_TYPE,
		},

		{
			Label: "Sort Columns Direction",
			Type:  HEADING_CHOICE_TYPE,
		},
		{
			Label: "Ascending",
			Type:  CHOICE_TYPE,
		},
		{
			Label: "Descending",
			Type:  CHOICE_TYPE,
		},
	}

	return settings
}
