package tui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

func (m Model) CommandInput(prompt string, placeholder string) textinput.Model {
	commandInput := textinput.New()
	commandInput.Prompt = prompt
	commandInput.Placeholder = placeholder
	commandInput.CharLimit = 256
	commandInput.SetWidth(50)
	commandInput.Blur()

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.theme.CommandBar.Background)).
		Foreground(lipgloss.Color(m.theme.CommandBar.Foreground))

	placeholderStyle := baseStyle.
		Foreground(lipgloss.Color(m.theme.CommandBar.Placeholder))

	styles := commandInput.Styles()
	styles.Focused.Text = baseStyle
	styles.Focused.Placeholder = placeholderStyle
	styles.Focused.Prompt = baseStyle

	styles.Blurred.Text = baseStyle
	styles.Blurred.Placeholder = placeholderStyle
	styles.Blurred.Prompt = baseStyle

	styles.Cursor.Color = lipgloss.Color(m.theme.CommandBar.Foreground)

	commandInput.SetStyles(styles)

	return commandInput
}
