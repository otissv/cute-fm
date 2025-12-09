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
	colMarker = 5
	colIndex  = 5
	colPerms  = 11
	colSize   = 6
	colType   = 16
	colUser   = 8
	colGroup  = 8
	colDate   = 14
)

// FileItem wraps filesystem.FileInfo to implement the list.Item interface.
type FileItem struct {
	Info   filesystem.FileInfo
	Marked bool
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

	isCursor := index == m.Index()
	isMarked := fi.Marked
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

	line := d.renderFileRow(fi.Info, isCursor, isMarked, displayIndex)
	_, _ = io.WriteString(w, line)
}

// renderFileRow renders a single file row with all columns styled.
// index is the precomputed display index (already relative/absolute as desired).
func (d FileItemDelegate) renderFileRow(fi filesystem.FileInfo, isCursor bool, isMarked bool, index int) string {
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
	// In select mode, the cursor row uses FileList.Marked as its background.
	// Otherwise, fall back to the generic Selection background when available,
	// or the file-list background.
	bgColor := theme.FileList.Background

	if isMarked {
		bgColor = theme.FileList.Marked
	}
	if isCursor {
		bgColor = theme.Selection.Background
	}

	if bgColor != "" {
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
	lineCols := []string{}

	// Optional selection marker column in select mode.
	if ActiveTuiMode == ModeSelect {
		markerContent := "[   ]"
		if isMarked {
			markerContent = "[ x ]"
		}

		markerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FileList.Foreground))

		if bgColor != "" {
			markerStyle = markerStyle.Background(lipgloss.Color(bgColor))
		}

		markerText := padCellWithBG(markerStyle.Render(markerContent), colMarker, bgColor)
		lineCols = append(lineCols, markerText)
	}

	// Index column: always present.
	lineCols = append(lineCols, indexText)

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
func FileInfosToItems(files []filesystem.FileInfo, marked map[string]bool) []list.Item {
	items := make([]list.Item, len(files))
	for i, f := range files {
		if marked != nil && marked[f.Path] {
			items[i] = FileItem{Info: f, Marked: true}
		} else {
			items[i] = FileItem{Info: f, Marked: false}
		}
	}
	return items
}

type FileHeaderRowArgs struct {
	Theme        theming.Theme
	TotalWidth   int
	Columns      []filesystem.FileInfoColumn
	SortColumnBy SortColumnBy
}

// RenderFileHeaderRow renders a single header row for the file list, aligned
// with the same columns and widths used for file rows.
func RenderFileHeaderRow(args FileHeaderRowArgs) string {
	bgColor := args.Theme.FileList.Background
	bg := lipgloss.Color(bgColor)

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theming.DefaultTheme().Foreground))

	permsHeading := "Permissions"
	sizeHeading := "Size"
	typeHeading := "Type"
	userHeading := "User"
	groupHeading := "Group"
	dateHeading := "Date Modified"
	nameHeading := "Name"

	sortByDirection := "↓ "

	if args.SortColumnBy.direction == "desc" {
		sortByDirection = "↑ "
	}

	switch args.SortColumnBy.column {
	case filesystem.ColumnPermissions:
		permsHeading = sortByDirection + permsHeading
	case filesystem.ColumnSize:
		sizeHeading = sortByDirection + sizeHeading
	case filesystem.ColumnMimeType:
		typeHeading = sortByDirection + typeHeading
	case filesystem.ColumnUser:
		userHeading = sortByDirection + userHeading
	case filesystem.ColumnGroup:
		groupHeading = sortByDirection + groupHeading
	case filesystem.ColumnDateModified:
		dateHeading = sortByDirection + dateHeading
	case filesystem.ColumnName:
		nameHeading = sortByDirection + nameHeading
	}

	indexText := padCellWithBG(baseStyle.Render(" "), colIndex, bgColor)
	permsText := padCellWithBG(baseStyle.Render(permsHeading), colPerms, bgColor)
	sizeText := padCellWithBG(baseStyle.Render(sizeHeading), colSize, bgColor)
	typeText := padCellWithBG(baseStyle.Render(typeHeading), colType, bgColor)
	userText := padCellWithBG(baseStyle.Render(userHeading), colUser, bgColor)
	groupText := padCellWithBG(baseStyle.Render(groupHeading), colGroup, bgColor)
	dateText := padCellWithBG(baseStyle.Render(dateHeading), colDate, bgColor)
	nameText := baseStyle.Render(nameHeading) // last column can flow to the right

	lineCols := []string{}

	// Optional selection marker header column in select mode.
	if ActiveTuiMode == ModeSelect {
		markerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(args.Theme.FileList.Foreground))
		markerText := padCellWithBG(markerStyle.Render("[   ]"), colMarker, bgColor)
		lineCols = append(lineCols, markerText)
	}

	// Index column always present.
	lineCols = append(lineCols, indexText)

	for _, col := range args.Columns {
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
	if args.TotalWidth > 0 && bgColor != "" {
		lineWidth := lipgloss.Width(line)
		if lineWidth < args.TotalWidth {
			missing := args.TotalWidth - lineWidth
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
