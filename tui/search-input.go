package tui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

func (m Model) SearchInput(prompt string, placeholder string) textinput.Model {
	searchInput := textinput.New()
	searchInput.Prompt = prompt
	searchInput.Placeholder = placeholder
	searchInput.CharLimit = 256
	searchInput.SetWidth(50)
	searchInput.Blur()

	searchBaseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.CommandBar.Background)).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground))

	searchPlaceholderStyle := searchBaseStyle.
		Foreground(lipgloss.Color(m.theme.CommandBar.Placeholder))

	searchStyles := searchInput.Styles()
	searchStyles.Focused.Text = searchBaseStyle
	searchStyles.Focused.Placeholder = searchPlaceholderStyle
	searchStyles.Focused.Prompt = searchBaseStyle

	searchStyles.Blurred.Text = searchBaseStyle
	searchStyles.Blurred.Placeholder = searchPlaceholderStyle
	searchStyles.Blurred.Prompt = searchBaseStyle

	searchStyles.Cursor.Color = lipgloss.Color(m.theme.CommandBar.Foreground)

	searchInput.SetStyles(searchStyles)

	return searchInput
}
