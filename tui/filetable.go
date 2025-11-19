package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lsfm/filesystem"
	"lsfm/theming"
)

// renderFileTable builds a simple eza-like table for the left viewport using
// the provided theme and directory listing.
//
// Columns:
//
//	Permissions Size User  Group Date Modified  Name
func renderFileTable(theme theming.Theme, files []filesystem.FileInfo) string {
	var b strings.Builder

	// Header
	header := fmt.Sprintf(
		"%-11s %-6s %-8s %-8s %-14s %s",
		"Permissions",
		"Size",
		"User",
		"Group",
		"Date Modified",
		"Name",
	)
	headerStyle := lipgloss.NewStyle().Bold(true)
	b.WriteString(headerStyle.Render(header))
	b.WriteRune('\n')

	// Rows
	for _, fi := range files {
		perm := fi.Permissions
		size := fi.Size
		user := fi.User
		group := fi.Group
		date := fi.DateModified
		name := fi.Name

		// Permission color (rough heuristic based on rwx count).
		permStyle := colorForPermissions(theme, perm)
		permText := permStyle.Render(perm)

		// Field colors.
		userStyle := theming.StyleFromSpec(theme.FieldColors["user"])
		groupStyle := theming.StyleFromSpec(theme.FieldColors["group"])
		sizeStyle := theming.StyleFromSpec(theme.FieldColors["size"])
		timeStyle := theming.StyleFromSpec(theme.FieldColors["time"])

		userText := userStyle.Render(user)
		groupText := groupStyle.Render(group)
		sizeText := sizeStyle.Render(size)
		timeText := timeStyle.Render(date)

		// File name color based on file type.
		nameColorSpec := theme.FileTypeColors[fi.Type]
		nameStyle := theming.StyleFromSpec(nameColorSpec)
		nameText := nameStyle.Render(name)

		line := fmt.Sprintf(
			"%-11s %-6s %-8s %-8s %-14s %s",
			permText,
			sizeText,
			userText,
			groupText,
			timeText,
			nameText,
		)

		b.WriteString(line)
		b.WriteRune('\n')
	}

	return b.String()
}

// colorForPermissions selects a permission style based on how permissive the
// rwx bits are, using the "full", "partial", and "none" colors from the theme.
func colorForPermissions(theme theming.Theme, perm string) lipgloss.Style {
	// Drop the leading type character ('d' or '.') if present.
	if len(perm) > 0 {
		perm = perm[1:]
	}

	fullSpec := theme.PermissionColors["full"]
	partialSpec := theme.PermissionColors["partial"]
	noneSpec := theme.PermissionColors["none"]

	// Count non-dash characters – a simple proxy for access level.
	count := 0
	for _, r := range perm {
		if r != '-' {
			count++
		}
	}

	switch {
	case count >= 9:
		return theming.StyleFromSpec(fullSpec)
	case count == 0:
		return theming.StyleFromSpec(noneSpec)
	default:
		return theming.StyleFromSpec(partialSpec)
	}
}
