package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func QuitWindow(m Model) *lipgloss.Layer {
	cancel := NewButton(m, ButtonArgs{
		label:       "(esc) Cancel",
		borderColor: m.theme.Foreground,
		borders:     true,
		margin:      []int{0, 1},
		width:       20,
		onClick: func() tea.Cmd {
			return func() tea.Msg {
				return ButtonMsg{ButtonID: "cancel"}
			}
		},
	})
	quit := NewButton(m, ButtonArgs{
		label:       "(q) Quit",
		borderColor: m.theme.Foreground,
		borders:     true,
		margin:      []int{0, 1},
		width:       20,
		onClick: func() tea.Cmd {
			return func() tea.Msg {
				return ButtonMsg{ButtonID: "quit"}
			}
		},
	})

	content := fmt.Sprintf(`Are you sure you want to quit?

%s
`, lipgloss.JoinHorizontal(lipgloss.Center, cancel.View(), quit.View()))

	return DialogWindow(m, DialogWindowArgs{
		Title:   "Quit",
		Content: content,
		Width:   50,
	})
}
