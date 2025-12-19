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
	"cute/utils"
)

const (
	colMarkerWidth = 5
	colIndexWidth  = 5
	colPermsWidth  = 11
	colSizeWidth   = 6
	colTypeWidth   = 16
	colUserWidth   = 8
	colGroupWidth  = 8
	colDateWidth   = 14
	colNameWidth   = 20
)

type FileItem struct {
	Info   filesystem.FileInfo
	Marked bool
}

func (i FileItem) FilterValue() string {
	return i.Info.Name
}

type FileItemDelegate struct {
	theme      theming.Theme
	totalWidth int
	columns    []filesystem.FileInfoColumn
}

func NewFileItemDelegate(theme theming.Theme, width int, columns []filesystem.FileInfoColumn) FileItemDelegate {
	if len(columns) == 0 {
		columns = filesystem.FileInfoColumnNames
	}
	return FileItemDelegate{
		theme:      theme,
		totalWidth: width,
		columns:    columns,
	}
}

func (d FileItemDelegate) Height() int {
	return 1
}

func (d FileItemDelegate) Spacing() int {
	return 0
}

func (d FileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d FileItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	fi, ok := item.(FileItem)
	if !ok {
		return
	}

	isCursor := index == m.Index()
	isMarked := fi.Marked

	// Compute a Vim-style line number:
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

	bgColor := theme.Background

	if isMarked {
		bgColor = theme.Marked
	}
	if isCursor {
		bgColor = theme.Hovered
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

	permTextRaw := renderPermissions(theme, fi, bgColor)
	permsText := utils.TruncateAndPadCell(permTextRaw, colPermsWidth, bgColor)
	indexText := utils.TruncateAndPadCell(indexStyle.Render(fmt.Sprintf("%d", index)), colIndexWidth, bgColor)
	userText := utils.TruncateAndPadCell(userStyle.Render(user), colUserWidth, bgColor)
	groupText := utils.TruncateAndPadCell(groupStyle.Render(group), colGroupWidth, bgColor)
	sizeText := utils.TruncateAndPadCell(sizeStyle.Render(size), colSizeWidth, bgColor)
	typeText := utils.TruncateAndPadCell(typeStyle.Render(mime), colTypeWidth, bgColor)
	dateText := utils.TruncateAndPadCell(timeStyle.Render(date), colDateWidth, bgColor)

	nameColorSpec := theme.FileTypeColors[fi.Type]
	nameStyle := theming.StyleFromSpec(nameColorSpec)
	if bgColor != "" {
		nameStyle = nameStyle.Background(lipgloss.Color(bgColor))
	}

	nameText := utils.TruncateAndPadCell(nameStyle.Render(name), colNameWidth, bgColor)

	lineCols := []string{}

	if ActiveTuiMode == ModeSelect {
		markerContent := "[   ]"
		if isMarked {
			markerContent = "[ x ]"
		}

		markerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground))

		if bgColor != "" {
			markerStyle = markerStyle.Background(lipgloss.Color(bgColor))
		}

		markerText := utils.TruncateAndPadCell(markerStyle.Render(markerContent), colMarkerWidth, bgColor)
		lineCols = append(lineCols, markerText)
	}

	lineCols = append(lineCols, indexText)

	filteredColumns := getFileInfoFilteredColumns(
		d.columns,
		filesystem.FileInfo{
			Name:         nameText,
			Permissions:  permsText,
			Size:         sizeText,
			MimeType:     typeText,
			User:         userText,
			Group:        groupText,
			DateModified: dateText,
		},
	)

	lineCols = append(lineCols, filteredColumns...)

	sep := " "
	if bgColor != "" {
		sep = lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(" ")
	}

	line := strings.Join(lineCols, sep)

	// Ensure the line never exceeds totalWidth - truncate if necessary
	if d.totalWidth > 0 {
		lineWidth := lipgloss.Width(line)
		if lineWidth > d.totalWidth {
			// Truncate the entire line to fit
			line = utils.TruncateString(line, d.totalWidth)
		}

		// Pad the end of the line so that the row's background extends to the edge.
		if bgColor != "" {
			lineWidth = lipgloss.Width(line)
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
	}

	return line
}

func renderPermissions(theme theming.Theme, fi filesystem.FileInfo, bgColor string) string {
	perm := fi.Permissions
	if perm == "" {
		return ""
	}

	var b strings.Builder

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
	SortColumnBy SortFileListColumnBy
}

func RenderFileHeaderRow(args FileHeaderRowArgs) string {
	bgColor := args.Theme.Background
	bg := lipgloss.Color(bgColor)

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theming.GetTheme().Foreground))

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
	case filesystem.FileInfoColumns.Permissions:
		permsHeading = sortByDirection + permsHeading
	case filesystem.FileInfoColumns.Size:
		sizeHeading = sortByDirection + sizeHeading
	case filesystem.FileInfoColumns.MimeType:
		typeHeading = sortByDirection + typeHeading
	case filesystem.FileInfoColumns.User:
		userHeading = sortByDirection + userHeading
	case filesystem.FileInfoColumns.Group:
		groupHeading = sortByDirection + groupHeading
	case filesystem.FileInfoColumns.Modified:
		dateHeading = sortByDirection + dateHeading
	case filesystem.FileInfoColumns.Name:
		nameHeading = sortByDirection + nameHeading
	}

	indexText := utils.TruncateAndPadCell(baseStyle.Render(" "), colIndexWidth, bgColor)
	permsText := utils.TruncateAndPadCell(baseStyle.Render(permsHeading), colPermsWidth, bgColor)
	sizeText := utils.TruncateAndPadCell(baseStyle.Render(sizeHeading), colSizeWidth, bgColor)
	typeText := utils.TruncateAndPadCell(baseStyle.Render(typeHeading), colTypeWidth, bgColor)
	userText := utils.TruncateAndPadCell(baseStyle.Render(userHeading), colUserWidth, bgColor)
	groupText := utils.TruncateAndPadCell(baseStyle.Render(groupHeading), colGroupWidth, bgColor)
	dateText := utils.TruncateAndPadCell(baseStyle.Render(dateHeading), colDateWidth, bgColor)
	nameText := utils.TruncateAndPadCell(baseStyle.Render(nameHeading), colNameWidth, bgColor) // Now also truncated and padded

	lineCols := []string{}

	// Optional selection marker header column in select mode.
	if ActiveTuiMode == ModeSelect {
		markerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(args.Theme.Foreground))
		markerText := utils.TruncateAndPadCell(markerStyle.Render("[   ]"), colMarkerWidth, bgColor)
		lineCols = append(lineCols, markerText)
	}

	lineCols = append(lineCols, indexText)

	filteredColumns := getFileInfoFilteredColumns(
		args.Columns,
		filesystem.FileInfo{
			Name:         nameText,
			Permissions:  permsText,
			Size:         sizeText,
			MimeType:     typeText,
			User:         userText,
			Group:        groupText,
			DateModified: dateText,
		},
	)

	lineCols = append(lineCols, filteredColumns...)

	sep := lipgloss.NewStyle().Background(bg).Render(" ")
	line := strings.Join(lineCols, sep)

	// Ensure header line never exceeds totalWidth
	if args.TotalWidth > 0 {
		lineWidth := lipgloss.Width(line)
		if lineWidth > args.TotalWidth {
			line = utils.TruncateString(line, args.TotalWidth)
		}

		// Pad out to totalWidth so the background fills the entire content area.
		if bgColor != "" {
			lineWidth = lipgloss.Width(line)
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
	}

	return line
}

func getFileInfoFilteredColumns(columns []filesystem.FileInfoColumn, info filesystem.FileInfo) []string {
	lineCols := []string{}
	for _, col := range columns {
		switch col {
		case filesystem.FileInfoColumns.Name:
			lineCols = append(lineCols, info.Name)
		case filesystem.FileInfoColumns.Permissions:
			lineCols = append(lineCols, info.Permissions)
		case filesystem.FileInfoColumns.Size:
			lineCols = append(lineCols, info.Size)
		case filesystem.FileInfoColumns.MimeType:
			lineCols = append(lineCols, info.MimeType)
		case filesystem.FileInfoColumns.User:
			lineCols = append(lineCols, info.User)
		case filesystem.FileInfoColumns.Group:
			lineCols = append(lineCols, info.Group)
		case filesystem.FileInfoColumns.Modified:
			lineCols = append(lineCols, info.DateModified)
		}
	}

	return lineCols
}
