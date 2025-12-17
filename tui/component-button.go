package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Button struct {
	label      string
	x, y       int
	width      int
	onClick    func() tea.Cmd
	style      lipgloss.Style
	hoverStyle lipgloss.Style
	isHovered  bool
}

type ButtonArgs struct {
	borderColor string
	foreground  string
	hoverStyle  *lipgloss.Style
	isHovered   bool
	label       string
	margin      []int
	onClick     func() tea.Cmd
	padding     []int
	width       int
	x, y        int
	borders     bool
	align       *lipgloss.Position
}

type ButtonMsg struct {
	ButtonID string
}

func NewButton(m Model, args ButtonArgs) Button {
	theme := m.GetTheme()
	width := 5
	foreground := theme.Foreground
	borderColor := theme.BorderColor
	padding := []int{0, 0}
	margin := []int{0, 0}

	if args.borderColor != "" {
		borderColor = args.borderColor
	}

	if args.foreground != "" {
		foreground = args.foreground
	}

	if args.margin != nil {
		margin = args.margin
	}

	if args.padding != nil {
		padding = args.padding
	}

	if args.width != 0 {
		width = args.width
	}

	style := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Foreground(lipgloss.Color(foreground)).
		Margin(margin...).
		Padding(padding...).
		Width(width)

	hoverStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Foreground(lipgloss.Color(theme.Primary)).
		Margin(margin...).
		Padding(padding...).
		Width(width)

	if args.borders {
		style = style.
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true)

		hoverStyle = hoverStyle.
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true)
	} else {
		style = style.
			BorderTop(false).
			BorderBottom(false).
			BorderLeft(false).
			BorderRight(false)

		hoverStyle = hoverStyle.
			BorderTop(false).
			BorderBottom(false).
			BorderLeft(false).
			BorderRight(false)
	}

	if args.align != nil {
		style = style.Align(*args.align)
		hoverStyle = hoverStyle.Align(*args.align)
	}

	isHovered := false

	if args.hoverStyle != nil {
		hoverStyle = *args.hoverStyle
	}
	if args.isHovered {
		isHovered = args.isHovered
	}

	return Button{
		label:      args.label,
		x:          args.x,
		y:          args.y,
		width:      args.width,
		onClick:    args.onClick,
		style:      style,
		hoverStyle: hoverStyle,
		isHovered:  isHovered,
	}
}

// SetStyles allows customizing button styles
func (b *Button) SetStyles(normal, hover lipgloss.Style) {
	b.style = normal
	b.hoverStyle = hover
}

// View renders the button
func (b Button) View() string {
	style := b.style
	if b.isHovered {
		style = b.hoverStyle
	}
	return style.Render(b.label)
}

// Update handles mouse events for the button
func (b Button) Update(msg tea.Msg) (Button, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseClickMsg:
		// Handle mouse clicks
		if msg.Button == tea.MouseLeft {
			// Check if click is within button bounds
			mouseX := msg.X
			mouseY := msg.Y

			buttonX := b.x
			buttonY := b.y
			buttonWidth := b.width
			buttonHeight := 1 // Height is typically 1 line for buttons

			isOverButton := mouseX >= buttonX && mouseX < buttonX+buttonWidth &&
				mouseY >= buttonY && mouseY < buttonY+buttonHeight

			if isOverButton && b.onClick != nil {
				return b, b.onClick()
			}
		}
	case tea.MouseMotionMsg:
		// Handle mouse movement for hover effects
		mouseX := msg.X
		mouseY := msg.Y

		buttonX := b.x
		buttonY := b.y
		buttonWidth := b.width
		buttonHeight := 1

		b.isHovered = mouseX >= buttonX && mouseX < buttonX+buttonWidth &&
			mouseY >= buttonY && mouseY < buttonY+buttonHeight
	}
	return b, nil
}

// CheckBounds checks if coordinates are within button bounds
func (b Button) CheckBounds(x, y int) bool {
	return x >= b.x && x < b.x+b.width && y >= b.y && y < b.y+1
}

// // Example usage function that creates a button with an action
// func ExampleButtonWithAction(m) Button {
// 	return NewButton("Click Me!", 10, 5, 12, func() tea.Cmd {
// 		return func() tea.Msg {
// 			fmt.Println("Button clicked!")
// 			return ButtonMsg{ButtonID: "example"}
// 		}
// 	})
// }
