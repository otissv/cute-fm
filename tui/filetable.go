package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lsfm/filesystem"
	"lsfm/theming"
)

// renderFileTable builds a simple eza-like table for the left viewport using
// the provided theme and directory listing. selectedIndex is the 0-based row
// index that should be highlighted; pass -1 for "no selection".
//
// Columns:
//
//	Permissions Size User  Group Date Modified  Name
func renderFileTable(theme theming.Theme, files []filesystem.FileInfo, selectedIndex int) string {
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
		if selectedIndex >= 0 && i == selectedIndex && theme.SelectedBackground != "" {
			bgColor = theme.SelectedBackground
			bg := lipgloss.Color(bgColor)
			userStyle = userStyle.Background(bg)
			groupStyle = groupStyle.Background(bg)
			sizeStyle = sizeStyle.Background(bg)
			timeStyle = timeStyle.Background(bg)
		}

		// Render permission string with per-character coloring.
		permTextRaw := renderPermissions(theme, fi, bgColor)
		permText := padCell(permTextRaw, colPerms)

		userText := padCell(userStyle.Render(user), colUser)
		groupText := padCell(groupStyle.Render(group), colGroup)
		sizeText := padCell(sizeStyle.Render(size), colSize)
		timeText := padCell(timeStyle.Render(date), colDate)

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

		line := strings.Join(lineCols, " ")
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
	var bg lipgloss.Color
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
			spec = theme.PermRead
		case "w":
			spec = theme.PermWrite
		case "x":
			spec = theme.PermExec
		default:
			spec = theme.PermissionColors["none"]
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
