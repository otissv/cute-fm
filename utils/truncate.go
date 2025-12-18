package utils

import (
	"charm.land/lipgloss/v2"
)

func TruncateString(content string, maxWidth int) string {
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

func TruncateAndPadCell(content string, w int, bgColor string) string {
	truncated := TruncateString(content, w)
	return PadCellWithBG(truncated, w, bgColor)
}
