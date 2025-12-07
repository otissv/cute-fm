package tui

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"cute/filesystem"
	"cute/theming"
)

const (
	colIndex = 5
	colPerms = 11
	colSize  = 6
	colType  = 16
	colUser  = 8
	colGroup = 8
	colDate  = 14
)

// FileItem wraps filesystem.FileInfo to implement the list.Item interface.
type FileItem struct {
	Info filesystem.FileInfo
}

// FilterValue returns the file name for filtering.
func (i FileItem) FilterValue() string {
	return i.Info.Name
}

// FileItemDelegate handles rendering of file items in the list.
type FileItemDelegate struct {
	theme      theming.Theme
	totalWidth int
	// columns defines which FileInfo columns should be rendered for each row,
	// in order, excluding the leading index column which is always shown.
	columns []filesystem.FileInfoColumn
}

// NewFileItemDelegate creates a new delegate for rendering file items.
func NewFileItemDelegate(theme theming.Theme, width int, columns []filesystem.FileInfoColumn) FileItemDelegate {
	if len(columns) == 0 {
		columns = filesystem.ColumnNames
	}
	return FileItemDelegate{
		theme:      theme,
		totalWidth: width,
		columns:    columns,
	}
}

// Height returns the height of each item (1 line per file).
func (d FileItemDelegate) Height() int {
	return 1
}

// Spacing returns the spacing between items.
func (d FileItemDelegate) Spacing() int {
	return 0
}

// Update handles item-level updates.
func (d FileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render renders a single file item.
func (d FileItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	fi, ok := item.(FileItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	// Compute a Vim-style line number:
	//   - The currently selected row shows 0.
	//   - All other rows show the absolute distance from the selection.
	displayIndex := index + 1
	current := m.Index()
	if current >= 0 {
		if index == current {
			displayIndex = 0
		} else {
			diff := index - current
			if diff < 0 {
				diff = -diff
			}
			displayIndex = diff
		}
	}

	line := d.renderFileRow(fi.Info, isSelected, displayIndex)
	_, _ = io.WriteString(w, line)
}

// renderFileRow renders a single file row with all columns styled.
// index is the precomputed display index (already relative/absolute as desired).
func (d FileItemDelegate) renderFileRow(fi filesystem.FileInfo, isSelected bool, index int) string {
	theme := d.theme

	size := fi.Size
	mime := strings.Split(fi.MimeType, "/")[1]
	user := fi.User
	group := fi.Group
	date := fi.DateModified
	name := fi.Name

	// Field colors.
	indexStyle := theming.StyleFromSpec(theme.FieldColors["nlink"])
	userStyle := theming.StyleFromSpec(theme.FieldColors["user"])
	groupStyle := theming.StyleFromSpec(theme.FieldColors["group"])
	sizeStyle := theming.StyleFromSpec(theme.FieldColors["size"])
	typeStyle := theming.StyleFromSpec(theme.FieldColors["type"])
	timeStyle := theming.StyleFromSpec(theme.FieldColors["time"])

	// Background color for the row.
	bgColor := ""
	if isSelected && theme.Selection.Background != "" {
		bgColor = theme.Selection.Background
		bg := lipgloss.Color(bgColor)
		indexStyle = indexStyle.Background(bg)
		userStyle = userStyle.Background(bg)
		groupStyle = groupStyle.Background(bg)
		sizeStyle = sizeStyle.Background(bg)
		typeStyle = typeStyle.Background(bg)
		timeStyle = timeStyle.Background(bg)
	} else {
		bgColor = theme.FileList.Background
		bg := lipgloss.Color(bgColor)
		indexStyle = indexStyle.Background(bg)
		userStyle = userStyle.Background(bg)
		groupStyle = groupStyle.Background(bg)
		sizeStyle = sizeStyle.Background(bg)
		typeStyle = typeStyle.Background(bg)
		timeStyle = timeStyle.Background(bg)
	}

	// Render permission string with per-character coloring.
	permTextRaw := renderPermissions(theme, fi, bgColor)
	permText := padCellWithBG(permTextRaw, colPerms, bgColor)

	// Index column: value is already prepared by the caller (relative/absolute).
	indexText := padCellWithBG(indexStyle.Render(fmt.Sprintf("%d", index)), colIndex, bgColor)
	userText := padCellWithBG(userStyle.Render(user), colUser, bgColor)
	groupText := padCellWithBG(groupStyle.Render(group), colGroup, bgColor)
	sizeText := padCellWithBG(sizeStyle.Render(size), colSize, bgColor)
	typeText := padCellWithBG(typeStyle.Render(mime), colType, bgColor)
	timeText := padCellWithBG(timeStyle.Render(date), colDate, bgColor)

	// File name color based on file type.
	nameColorSpec := theme.FileTypeColors[fi.Type]
	nameStyle := theming.StyleFromSpec(nameColorSpec)
	if bgColor != "" {
		nameStyle = nameStyle.Background(lipgloss.Color(bgColor))
	}
	nameText := nameStyle.Render(name)

	// Build the list of columns to render based on the delegate configuration.
	lineCols := []string{
		indexText, // always show index
	}

	for _, col := range d.columns {
		switch col {
		case filesystem.ColumnPermissions:
			lineCols = append(lineCols, permText)
		case filesystem.ColumnSize:
			lineCols = append(lineCols, sizeText)
		case filesystem.ColumnMimeType:
			lineCols = append(lineCols, typeText)
		case filesystem.ColumnUser:
			lineCols = append(lineCols, userText)
		case filesystem.ColumnGroup:
			lineCols = append(lineCols, groupText)
		case filesystem.ColumnDateModified:
			lineCols = append(lineCols, timeText)
		case filesystem.ColumnName:
			lineCols = append(lineCols, nameText)
		}
	}

	sep := " "
	if bgColor != "" {
		sep = lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(" ")
	}

	line := strings.Join(lineCols, sep)

	// Pad the end of the line so that the row's background extends to the edge.
	if d.totalWidth > 0 && bgColor != "" {
		lineWidth := lipgloss.Width(line)
		if lineWidth < d.totalWidth {
			missing := d.totalWidth - lineWidth
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

	return line
}

// renderPermissions renders the permission string with per-character coloring.
func renderPermissions(theme theming.Theme, fi filesystem.FileInfo, bgColor string) string {
	perm := fi.Permissions
	if perm == "" {
		return ""
	}

	var b strings.Builder

	// Optional background for selected rows.
	hasBG := bgColor != ""

	// First character: type indicator ('d', '.', 'l', etc.) colored by file type.
	typeSpec := theme.FileTypeColors[fi.Type]
	typeStyle := theming.StyleFromSpec(typeSpec)
	if hasBG {
		typeStyle = typeStyle.Background(lipgloss.Color(bgColor))
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
			style = style.Background(lipgloss.Color(bgColor))
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

// FileInfosToItems converts a slice of FileInfo to a slice of list.Item.
func FileInfosToItems(files []filesystem.FileInfo) []list.Item {
	items := make([]list.Item, len(files))
	for i, f := range files {
		items[i] = FileItem{Info: f}
	}
	return items
}

// RenderFileHeaderRow renders a single header row for the file list, aligned
// with the same columns and widths used for file rows.
func RenderFileHeaderRow(theme theming.Theme, totalWidth int, columns []filesystem.FileInfoColumn) string {
	const (
		colIndex = 5
		colPerms = 11
		colSize  = 6
		colType  = 16
		colUser  = 8
		colGroup = 8
		colDate  = 14
	)

	bgColor := theme.FileList.Background
	bg := lipgloss.Color(bgColor)

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theming.DefaultTheme().Foreground))

	indexText := padCellWithBG(baseStyle.Render(" "), colIndex, bgColor)
	permsText := padCellWithBG(baseStyle.Render("Permissions"), colPerms, bgColor)
	sizeText := padCellWithBG(baseStyle.Render("Size"), colSize, bgColor)
	typeText := padCellWithBG(baseStyle.Render("Type"), colType, bgColor)
	userText := padCellWithBG(baseStyle.Render("User"), colUser, bgColor)
	groupText := padCellWithBG(baseStyle.Render("Group"), colGroup, bgColor)
	dateText := padCellWithBG(baseStyle.Render("Last Modified"), colDate, bgColor)
	nameText := baseStyle.Render("Name") // last column can flow to the right

	lineCols := []string{indexText} // index column always present

	for _, col := range columns {
		switch col {
		case filesystem.ColumnPermissions:
			lineCols = append(lineCols, permsText)
		case filesystem.ColumnSize:
			lineCols = append(lineCols, sizeText)
		case filesystem.ColumnMimeType:
			lineCols = append(lineCols, typeText)
		case filesystem.ColumnUser:
			lineCols = append(lineCols, userText)
		case filesystem.ColumnGroup:
			lineCols = append(lineCols, groupText)
		case filesystem.ColumnDateModified:
			lineCols = append(lineCols, dateText)
		case filesystem.ColumnName:
			lineCols = append(lineCols, nameText)
		}
	}

	sep := lipgloss.NewStyle().Background(bg).Render(" ")
	line := strings.Join(lineCols, sep)

	// Pad out to totalWidth so the background fills the entire content area.
	if totalWidth > 0 && bgColor != "" {
		lineWidth := lipgloss.Width(line)
		if lineWidth < totalWidth {
			missing := totalWidth - lineWidth
			spaceStyle := lipgloss.NewStyle().Background(bg)
			pad := spaceStyle.Render(" ")

			var tail strings.Builder
			for i := 0; i < missing; i++ {
				tail.WriteString(pad)
			}
			line += tail.String()
		}
	}

	return line
}
