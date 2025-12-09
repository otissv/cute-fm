package components

import (
	"charm.land/lipgloss/v2"

	"cute/tui"
)

func FileList(m tui.Model, args tui.FileListComponentArgs) string {
	theme := m.GetTheme()
	fileList := m.GetLeftPaneFileListForViewport(args.SplitPaneType)
	activeViewport := m.GetActiveViewport()
	isSplitPaneOpen := m.GetIsSplitPaneOpen()

	// Content width is viewport width minus left/right borders.
	contentWidth := args.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	header := tui.RenderFileHeaderRow(tui.FileHeaderRowArgs{
		Theme:        theme,
		TotalWidth:   contentWidth,
		Columns:      m.GetColumnVisibilityForViewport(args.SplitPaneType),
		SortColumnBy: m.GetSortColumnBy(),
	})
	body := fileList.View()

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.FileList.Background)).
		BorderBackground(lipgloss.Color(theme.FileList.Background)).
		BorderForeground(lipgloss.Color("#1E1E1E")).
		BorderStyle(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color(theme.FileList.Foreground)).
		Height(args.Height).
		Width(args.Width).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	if activeViewport == args.SplitPaneType && isSplitPaneOpen {
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color(theme.BorderColor)).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true)
	}

	return baseStyle.
		Render(inner)
}
