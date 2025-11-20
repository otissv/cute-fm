package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// ViewPrimitive is any Bubble/Bubble Tea primitive that can render itself
// via a View() string method (for example, textinput.Model, viewport.Model,
// list.Model, etc.).
type ViewPrimitive interface {
	View() string
}

// FloatingWindow is a reusable helper for rendering a primitive inside a
// centered "floating" window. It does not manage its own Update logic; you
// embed the primitive in your main model and call View() here when you want
// to draw it as a modal/floating window.
type FloatingWindow struct {
	// Content to render inside the window.
	Content ViewPrimitive

	// Fixed window size in cells. If 0, the content's size is used and only
	// borders/padding are applied.
	Width  int
	Height int

	// Optional title rendered in the window border.
	Title string

	// Base style for the floating window (borders, colors, padding, etc.).
	Style lipgloss.Style
}

// DefaultFloatingStyle returns a simple rounded, bordered style that matches
// the existing theme defaults reasonably well.
func DefaultFloatingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1).
		BorderForeground(lipgloss.Color("#874BFD")).
		Background(lipgloss.Color("#212121"))
}

// View renders the floating window centered within the given outer size.
// outerWidth/outerHeight should typically come from the main window size.
func (fw FloatingWindow) View(outerWidth, outerHeight int) string {
	if fw.Content == nil {
		return ""
	}

	contentView := fw.Content.View()

	style := fw.Style
	// If no custom style was provided, fall back to a default one. We use
	// a heuristic: if neither border nor background is set, assume "empty".
	if style.GetBorderStyle() == (lipgloss.Border{}) {
		style = DefaultFloatingStyle()
	}

	if fw.Width > 0 {
		style = style.Width(fw.Width)
	}
	if fw.Height > 0 {
		style = style.Height(fw.Height)
	}

	box := style.Render(contentView)

	// Center the box in the available space. We intentionally do NOT set
	// WithWhitespaceChars/WithWhitespaceForeground here because this window
	// is later overlaid on top of the base layout: if we filled the
	// whitespace with visible characters or colors, it would cover the
	// entire layout instead of just the dialog area.

	return lipgloss.Place(
		outerWidth,
		outerHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars("///"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#303030")),
	)
}
