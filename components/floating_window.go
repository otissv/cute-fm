package components

import (
	"strings"

	"cute/theming"
	"cute/tui"

	"charm.land/lipgloss/v2"
)

type FloatingWindow struct {
	Content      tui.ViewPrimitive
	Width        int
	Height       int
	Title        string
	Style        lipgloss.Style
	ScrollOffset int
}

func DefaultFloatingStyle(theme theming.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Dialog.Background)).
		Foreground(lipgloss.Color(theme.Dialog.Foreground)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Dialog.Border)).
		BorderBackground(lipgloss.Color(theme.Dialog.Background)).
		PaddingTop(theme.Dialog.PaddingTop).
		PaddingBottom(theme.Dialog.PaddingBottom).
		PaddingLeft(theme.Dialog.PaddingLeft).
		PaddingRight(theme.Dialog.PaddingRight)
}

func (fw FloatingWindow) View(outerWidth, outerHeight int) string {
	if fw.Content == nil {
		return ""
	}

	style := fw.Style

	if fw.Width > 0 {
		style = style.Width(fw.Width)
	}
	if fw.Height > 0 {
		style = style.Height(fw.Height)
	}

	contentView := fw.Content.View()

	// Apply vertical scrolling based on ScrollOffset. We slice the content
	// lines before rendering the box so the border and padding remain intact.
	if fw.Height > 0 {
		lines := strings.Split(contentView, "\n")
		if len(lines) > 0 {
			totalLines := len(lines)

			start := fw.ScrollOffset
			if start < 0 {
				start = 0
			}
			if start >= totalLines {
				start = totalLines - 1
			}

			end := start + fw.Height
			if end > totalLines {
				end = totalLines
			}

			canScrollUp := start > 0
			canScrollDown := end < totalLines

			visible := lines[start:end]

			// Add simple scroll indicators to the first/last visible lines
			// without changing the total number of lines, so the box height
			// remains stable.
			if len(visible) > 0 {
				if canScrollUp {
					visible[0] = "↑\n" + visible[0]
				}
				if canScrollDown {
					lastIdx := len(visible) - 1
					visible[lastIdx] = visible[lastIdx] + "\n↓"
				}
			}

			contentView = strings.Join(visible, "\n")
		}
	}

	box := style.Render(contentView)

	// If there's no title, just return the styled box.
	if fw.Title == "" {
		return box
	}

	// of the rendered box. This avoids needing an extra canvas or layer.
	lines := strings.Split(box, "\n")
	if len(lines) == 0 {
		return box
	}

	row := 0 // 0 = draw directly on the top border; use 1 to move it one row down.
	if row >= len(lines) {
		return box
	}

	line := lines[row]

	// The border line contains ANSI escape codes added by lipgloss. We must not
	// overwrite those, so we only modify the substring between the visible
	// corner runes (for the default rounded border: '╭' ... '╮').
	leftCorner := '╭'
	rightCorner := '╮'

	leftIdx := strings.IndexRune(line, leftCorner)
	rightIdx := strings.LastIndex(line, string(rightCorner))
	if leftIdx == -1 || rightIdx == -1 || rightIdx <= leftIdx {
		// Fallback: unexpected border characters; keep the box as-is.
		return box
	}

	prefix := line[:leftIdx]
	borderSegment := line[leftIdx : rightIdx+len("╮")]
	suffix := line[rightIdx+len("╮"):]

	segmentRunes := []rune(borderSegment)

	displayTitle := "┤ " + fw.Title + " ├"
	titleRunes := []rune(displayTitle)

	boxWidth := len(segmentRunes)
	titleWidth := len(titleRunes)

	// Keep both corner runes intact, so only draw within the inner span.
	if boxWidth < 2 {
		return box
	}

	innerWidth := boxWidth - 2
	if innerWidth <= 0 {
		return box
	}

	if titleWidth > innerWidth {
		titleRunes = titleRunes[:innerWidth]
		titleWidth = innerWidth
	}

	// Left-align within the inner area, starting just after the left corner.
	start := 2
	for i := 0; i < titleWidth && start+i < boxWidth-1; i++ {
		segmentRunes[start+i] = titleRunes[i]
	}

	lines[row] = prefix + string(segmentRunes) + suffix
	return strings.Join(lines, "\n")
}
