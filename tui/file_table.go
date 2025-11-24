package tui

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"

	"cute/filesystem"
	"cute/theming"
)

// renderFileTable builds a simple eza-like table for the left viewport using
// the provided theme and directory listing. selectedIndex is the 0-based row
// index that should be highlighted; pass -1 for "no selection".
// totalWidth is the target width (in terminal cells) of each rendered row;
// when greater than zero, the last column is padded so that the row's
// background color extends all the way to the viewport edge.
//
// Columns:
//
//	Permissions Size User  Group Date Modified  Name
func renderFileTable(theme theming.Theme, files []filesystem.FileInfo, selectedIndex int, totalWidth int) string {
	var b strings.Builder
	const (
		colPerms = 11
		colSize  = 6
		colUser  = 8
		colGroup = 8
		colDate  = 14
	)

	// Header (unstyled for simplicity, then bolded as a whole).
	headerCols := []string{
		padCell("Permissions", colPerms),
		padCell("Size", colSize),
		padCell("User", colUser),
		padCell("Group", colGroup),
		padCell("Date Modified", colDate),
		"Name",
	}

	// When we know the total target width, pad the "Name" header so that its
	// cell lines up with the last column and visually spans the remaining
	// viewport width.
	if totalWidth > 0 {
		const numSeps = 5 // spaces between the 6 header columns
		fixedCols := colPerms + colSize + colUser + colGroup + colDate
		baseWidth := fixedCols + numSeps
		if totalWidth > baseWidth {
			nameColWidth := totalWidth - baseWidth
			if nameColWidth < lipgloss.Width("Name") {
				nameColWidth = lipgloss.Width("Name")
			}
			headerCols[len(headerCols)-1] = padCell("Name", nameColWidth)
		}
	}

	headerLine := strings.Join(headerCols, " ")
	headerStyle := lipgloss.NewStyle().Bold(true)
	b.WriteString(headerStyle.Render(headerLine))
	b.WriteRune('\n')

	// Rows
	for i, fi := range files {
		size := fi.Size
		user := fi.User
		group := fi.Group
		date := fi.DateModified
		name := fi.Name

		// Field colors.
		userStyle := theming.StyleFromSpec(theme.FieldColors["user"])
		groupStyle := theming.StyleFromSpec(theme.FieldColors["group"])
		sizeStyle := theming.StyleFromSpec(theme.FieldColors["size"])
		timeStyle := theming.StyleFromSpec(theme.FieldColors["time"])

		// If this row is selected, apply only the selected background to every
		// column so the entire row is highlighted while preserving the
		// foreground colors coming from file-type and permission-level specs.
		bgColor := ""
		if selectedIndex >= 0 && i == selectedIndex && theme.Selection.Background != "" {
			bgColor = theme.Selection.Background
			bg := lipgloss.Color(bgColor)
			userStyle = userStyle.Background(bg)
			groupStyle = groupStyle.Background(bg)
			sizeStyle = sizeStyle.Background(bg)
			timeStyle = timeStyle.Background(bg)
		} else {
			bgColor = theme.FileList.Background
			bg := lipgloss.Color(bgColor)
			userStyle = userStyle.Background(bg)
			groupStyle = groupStyle.Background(bg)
			sizeStyle = sizeStyle.Background(bg)
			timeStyle = timeStyle.Background(bg)
		}

		// Render permission string with per-character coloring.
		permTextRaw := renderPermissions(theme, fi, bgColor)
		permText := padCellWithBG(permTextRaw, colPerms, bgColor)

		userText := padCellWithBG(userStyle.Render(user), colUser, bgColor)
		groupText := padCellWithBG(groupStyle.Render(group), colGroup, bgColor)
		sizeText := padCellWithBG(sizeStyle.Render(size), colSize, bgColor)
		timeText := padCellWithBG(timeStyle.Render(date), colDate, bgColor)

		// File name color based on file type.
		nameColorSpec := theme.FileTypeColors[fi.Type]
		nameStyle := theming.StyleFromSpec(nameColorSpec)
		if bgColor != "" {
			nameStyle = nameStyle.Background(lipgloss.Color(bgColor))
		}
		nameText := nameStyle.Render(name)

		lineCols := []string{
			permText,
			sizeText,
			userText,
			groupText,
			timeText,
			nameText,
		}

		sep := " "
		if bgColor != "" {
			sep = lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(" ")
		}

		line := strings.Join(lineCols, sep)

		// If we know the viewport width, pad the end of the line so that the
		// row's background extends all the way to the edge. This is especially
		// important for the selected row highlight on the "Name" column.
		if totalWidth > 0 && bgColor != "" {
			lineWidth := lipgloss.Width(line)
			if lineWidth < totalWidth {
				missing := totalWidth - lineWidth
				bg := lipgloss.Color(bgColor)
				spaceStyle := lipgloss.NewStyle().Background(bg)
				pad := spaceStyle.Render(" ")

				var tail strings.Builder
				for i := 0; i < missing; i++ {
					tail.WriteString(pad)
				}
				line += tail.String()
			}
		}

		b.WriteString(line)
		b.WriteRune('\n')
	}

	return b.String()
}

// colorForPermissions selects a permission style based on how permissive the
// rwx bits are, using the "full", "partial", and "none" colors from the theme.
func renderPermissions(theme theming.Theme, fi filesystem.FileInfo, bgColor string) string {
	perm := fi.Permissions
	if perm == "" {
		return ""
	}

	var b strings.Builder

	// Optional background for selected rows.
	var bg color.Color
	hasBG := false
	if bgColor != "" {
		bg = lipgloss.Color(bgColor)
		hasBG = true
	}

	// First character: type indicator ('d', '.', 'l', etc.) colored by file type.
	typeSpec := theme.FileTypeColors[fi.Type]
	typeStyle := theming.StyleFromSpec(typeSpec)
	if hasBG {
		typeStyle = typeStyle.Background(bg)
	}
	b.WriteString(typeStyle.Render(string(perm[0])))

	// Remaining permission bits: color each character separately.
	for _, r := range perm[1:] {
		ch := string(r)

		var spec string
		switch ch {
		case "r":
			spec = theme.Permissions.Read
		case "w":
			spec = theme.Permissions.Write
		case "x":
			spec = theme.Permissions.Exec
		default:
			spec = theme.Permissions.None
		}

		style := theming.StyleFromSpec(spec)
		if hasBG {
			style = style.Background(bg)
		}

		b.WriteString(style.Render(ch))
	}

	return b.String()
}

// padCell right-pads the given content with spaces so that its visible width
// (taking into account ANSI escape sequences used by lipgloss) is at least w.
func padCell(content string, w int) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}
	return content + strings.Repeat(" ", w-width)
}

// padCellWithBG right-pads content like padCell, but if a non-empty bgColor
// is provided, the padding spaces are rendered with that background color so
// that the cell's whitespace is also highlighted.
func padCellWithBG(content string, w int, bgColor string) string {
	width := lipgloss.Width(content)
	if width >= w {
		return content
	}

	// If no background color is specified, fall back to the plain padding.
	if bgColor == "" {
		return padCell(content, w)
	}

	missing := w - width
	bg := lipgloss.Color(bgColor)
	spaceStyle := lipgloss.NewStyle().Background(bg)

	var b strings.Builder
	b.WriteString(content)

	// Render one styled space and reuse it to avoid repeated allocations.
	pad := spaceStyle.Render(" ")
	for i := 0; i < missing; i++ {
		b.WriteString(pad)
	}

	return b.String()
}
