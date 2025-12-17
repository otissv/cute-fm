package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func truncateString(content string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	currentWidth := lipgloss.Width(content)
	if currentWidth <= maxWidth {
		return content
	}

	ellipsis := "..."
	ellipsisWidth := lipgloss.Width(ellipsis)
	if maxWidth < ellipsisWidth {
		if maxWidth >= 3 {
			return ellipsis[:maxWidth]
		}
		return ""
	}

	truncateWidth := maxWidth - ellipsisWidth
	runes := []rune(content)
	truncated := ""

	for _, r := range runes {
		test := truncated + string(r)
		if lipgloss.Width(test) > truncateWidth {
			break
		}
		truncated = test
	}

	return truncated + ellipsis
}

func truncateAndPadCell(content string, w int, bgColor string) string {
	truncated := truncateString(content, w)
	return padCellWithBG(truncated, w, bgColor)
}

func padCell(content string, w int) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}
	return content + strings.Repeat(" ", w-width)
}

func padCellWithBG(content string, w int, bgColor string) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}

	if bgColor == "" {
		return padCell(content, w)
	}

	missing := w - width
	bg := lipgloss.Color(bgColor)
	spaceStyle := lipgloss.NewStyle().Background(bg)

	var b strings.Builder
	b.WriteString(content)

	pad := spaceStyle.Render(" ")
	for i := 0; i < missing; i++ {
		b.WriteString(pad)
	}

	return b.String()
}
