package utils

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func PadCell(content string, w int) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}
	return content + strings.Repeat(" ", w-width)
}

func PadCellWithBG(content string, w int, bgColor string) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}

	if bgColor == "" {
		return PadCell(content, w)
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
