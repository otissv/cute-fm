package components

// viewPrimitive is a minimal implementation of tui.ViewPrimitive that just renders
// the given string. This lets us render the command input directly without an
// extra viewport layer.
type viewPrimitive string

func (t viewPrimitive) View() string {
	return string(t)
}
