package tui

import "charm.land/lipgloss/v2"

type FileListComponentArgs struct {
	Width         int
	Height        int
	SplitPaneType ActiveViewportType
}

func FileList(m Model, args FileListComponentArgs) string {
	theme := m.GetTheme()
	fileList := m.GetLeftPaneFileListForViewport(args.SplitPaneType)
	activeViewport := m.GetActiveViewport()
	isSplitPaneOpen := m.GetIsSplitPaneOpen()

	contentWidth := args.Width - 2
	if contentWidth < 1 {
		contentWidth = 1
	}

	header := RenderFileHeaderRow(FileHeaderRowArgs{
		Theme:        theme,
		TotalWidth:   contentWidth,
		Columns:      m.GetColumnVisibilityForViewport(args.SplitPaneType),
		SortColumnBy: m.GetSortFileListColumnBy(),
	})
	body := fileList.View()

	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)

	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Background)).
		BorderBackground(lipgloss.Color(theme.Background)).
		BorderForeground(lipgloss.Color("#1E1E1E")).
		BorderStyle(lipgloss.RoundedBorder()).
		Foreground(lipgloss.Color(theme.Foreground)).
		Height(args.Height).
		Width(args.Width).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	if activeViewport == args.SplitPaneType && isSplitPaneOpen {
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color(theme.ActiveBorderColor)).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true)
	}

	return baseStyle.
		Render(inner)
}
