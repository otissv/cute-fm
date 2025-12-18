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

// Device-specific column width constants
const (
	colDeviceIndexWidth       = 5
	colDeviceNameWidth        = 20
	colDevicePathWidth        = 15
	colDeviceMountWidth       = 20
	colDeviceFSTypeWidth      = 10
	colDeviceSizeWidth        = 10
	colDeviceUsedWidth        = 10
	colDeviceAvailWidth       = 10
	colDeviceUsePercentWidth  = 6
	colDeviceFreePercentWidth = 6
)

type DeviceItem struct {
	Info filesystem.DeviceInfo
}

// FilterValue returns the device name for filtering.
func (i DeviceItem) FilterValue() string {
	return i.Info.Name
}

// DeviceItemDelegate handles rendering of device items in the list.
type DeviceItemDelegate struct {
	theme      theming.Theme
	totalWidth int
	columns    []filesystem.DeviceInfoColumn
}

// NewDeviceItemDelegate creates a new delegate for rendering device items.
func NewDeviceItemDelegate(theme theming.Theme, width int, columns []filesystem.DeviceInfoColumn) DeviceItemDelegate {
	if len(columns) == 0 {
		columns = filesystem.DeviceInfoColumnNames
	}

	return DeviceItemDelegate{
		theme:      theme,
		totalWidth: width,
		columns:    columns,
	}
}

func (d DeviceItemDelegate) Height() int {
	return 1
}

func (d DeviceItemDelegate) Spacing() int {
	return 0
}

func (d DeviceItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d DeviceItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	di, ok := item.(DeviceItem)
	if !ok {
		return
	}

	isCursor := index == m.Index()

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

	line := d.renderDeviceRow(di.Info, isCursor, displayIndex)
	_, _ = io.WriteString(w, line)
}

func (d DeviceItemDelegate) renderDeviceRow(di filesystem.DeviceInfo, isCursor bool, index int) string {
	theme := d.theme

	// Field colors
	indexStyle := theming.StyleFromSpec(theme.FieldColors["nlink"])
	nameStyle := theming.StyleFromSpec(theme.FileList.Foreground)
	pathStyle := theming.StyleFromSpec(theme.FieldColors["user"])
	mountStyle := theming.StyleFromSpec(theme.FieldColors["group"])
	typeStyle := theming.StyleFromSpec(theme.FieldColors["type"])
	sizeStyle := theming.StyleFromSpec(theme.FieldColors["size"])
	usedStyle := theming.StyleFromSpec(theme.FieldColors["size"])
	availStyle := theming.StyleFromSpec(theme.FieldColors["size"])
	usePercentStyle := theming.StyleFromSpec(theme.FieldColors["size"])
	freePercentStyle := theming.StyleFromSpec(theme.FieldColors["size"])

	bgColor := theme.FileList.Background
	if isCursor {
		bgColor = theme.Selection.Background
	}

	if bgColor != "" {
		bg := lipgloss.Color(bgColor)
		indexStyle = indexStyle.Background(bg)
		nameStyle = nameStyle.Background(bg)
		pathStyle = pathStyle.Background(bg)
		mountStyle = mountStyle.Background(bg)
		typeStyle = typeStyle.Background(bg)
		sizeStyle = sizeStyle.Background(bg)
		usedStyle = usedStyle.Background(bg)
		availStyle = availStyle.Background(bg)
		usePercentStyle = usePercentStyle.Background(bg)
		freePercentStyle = freePercentStyle.Background(bg)
	}

	indexText := utils.TruncateAndPadCell(indexStyle.Render(fmt.Sprintf("%d", index)), colDeviceIndexWidth, bgColor)
	nameText := utils.TruncateAndPadCell(nameStyle.Render(di.Name), colDeviceNameWidth, bgColor)
	devicePathText := utils.TruncateAndPadCell(pathStyle.Render(di.Device), colDevicePathWidth, bgColor)
	mountText := utils.TruncateAndPadCell(mountStyle.Render(di.MountPoint), colDeviceMountWidth, bgColor)
	typeText := utils.TruncateAndPadCell(typeStyle.Render(di.FSType), colDeviceFSTypeWidth, bgColor)
	sizeText := utils.TruncateAndPadCell(sizeStyle.Render(di.Size), colDeviceSizeWidth, bgColor)
	usedText := utils.TruncateAndPadCell(usedStyle.Render(di.Used), colDeviceUsedWidth, bgColor)
	availText := utils.TruncateAndPadCell(availStyle.Render(di.Avail), colDeviceAvailWidth, bgColor)
	usePercentText := utils.TruncateAndPadCell(usePercentStyle.Render(di.UsePercent), colDeviceUsePercentWidth, bgColor)
	freePercentText := utils.TruncateAndPadCell(freePercentStyle.Render(di.FreePercent), colDeviceFreePercentWidth, bgColor)

	lineCols := []string{indexText}

	filteredColumns := getDeviceInfoFilteredColumns(
		d.columns,
		filesystem.DeviceInfo{
			Device:      devicePathText,
			Name:        nameText,
			MountPoint:  mountText,
			FSType:      typeText,
			Size:        sizeText,
			Used:        usedText,
			Avail:       availText,
			UsePercent:  usePercentText,
			FreePercent: freePercentText,
		})

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

func DeviceInfosToItems(devices []filesystem.DeviceInfo) []list.Item {
	items := make([]list.Item, len(devices))
	for i, d := range devices {
		items[i] = DeviceItem{Info: d}
	}
	return items
}

type DeviceHeaderRowArgs struct {
	Theme        theming.Theme
	TotalWidth   int
	columns      []filesystem.DeviceInfoColumn
	SortColumnBy SortDeviceColumnBy
}

func RenderDeviceHeaderRow(args DeviceHeaderRowArgs) string {
	bgColor := args.Theme.FileList.Background
	bg := lipgloss.Color(bgColor)

	nameHeading := "Name"
	devicePathHeading := "Device"
	mountHeading := "Mount Point"
	typeHeading := "FS Type"
	sizeHeading := "Size"
	usedHeading := "Used"
	availHeading := "Avail"
	usePercentHeading := "Use%"
	freePercentHeading := "Free%"

	sortByDirection := "↓ "

	if args.SortColumnBy.direction == "desc" {
		sortByDirection = "↑ "
	}

	switch args.SortColumnBy.column {
	case filesystem.DeviceInfoColumns.Name:
		nameHeading = sortByDirection + nameHeading
	case filesystem.DeviceInfoColumns.Device:
		devicePathHeading = sortByDirection + devicePathHeading
	case filesystem.DeviceInfoColumns.MountPoint:
		mountHeading = sortByDirection + mountHeading
	case filesystem.DeviceInfoColumns.FsType:
		typeHeading = sortByDirection + typeHeading
	case filesystem.DeviceInfoColumns.Size:
		sizeHeading = sortByDirection + sizeHeading
	case filesystem.DeviceInfoColumns.Used:
		usedHeading = sortByDirection + usedHeading
	case filesystem.DeviceInfoColumns.Avail:
		availHeading = sortByDirection + availHeading
	case filesystem.DeviceInfoColumns.UsePercent:
		usePercentHeading = sortByDirection + usePercentHeading
	case filesystem.DeviceInfoColumns.FreePercent:
		freePercentHeading = sortByDirection + freePercentHeading
	}

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theming.GetTheme().Foreground))
	indexText := utils.TruncateAndPadCell(baseStyle.Render(" "), colDeviceIndexWidth, bgColor)
	nameText := utils.TruncateAndPadCell(baseStyle.Render(nameHeading), colDeviceNameWidth, bgColor)
	devicePathText := utils.TruncateAndPadCell(baseStyle.Render(devicePathHeading), colDevicePathWidth, bgColor)
	mountText := utils.TruncateAndPadCell(baseStyle.Render(mountHeading), colDeviceMountWidth, bgColor)
	typeText := utils.TruncateAndPadCell(baseStyle.Render(typeHeading), colDeviceFSTypeWidth, bgColor)
	sizeText := utils.TruncateAndPadCell(baseStyle.Render(sizeHeading), colDeviceSizeWidth, bgColor)
	usedText := utils.TruncateAndPadCell(baseStyle.Render(usedHeading), colDeviceUsedWidth, bgColor)
	availText := utils.TruncateAndPadCell(baseStyle.Render(availHeading), colDeviceAvailWidth, bgColor)
	usePercentText := utils.TruncateAndPadCell(baseStyle.Render(usePercentHeading), colDeviceUsePercentWidth, bgColor)
	freePercentText := utils.TruncateAndPadCell(baseStyle.Render(freePercentHeading), colDeviceFreePercentWidth, bgColor)

	lineCols := []string{indexText}

	filteredColumns := getDeviceInfoFilteredColumns(
		args.columns,
		filesystem.DeviceInfo{
			Device:      devicePathText,
			Name:        nameText,
			MountPoint:  mountText,
			FSType:      typeText,
			Size:        sizeText,
			Used:        usedText,
			Avail:       availText,
			UsePercent:  usePercentText,
			FreePercent: freePercentText,
		})

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

func getDeviceInfoFilteredColumns(columns []filesystem.DeviceInfoColumn, info filesystem.DeviceInfo) []string {
	lineCols := []string{}

	for _, col := range columns {
		switch col {
		case filesystem.DeviceInfoColumns.Name:
			lineCols = append(lineCols, info.Name)

		case filesystem.DeviceInfoColumns.Device:
			lineCols = append(lineCols, info.Device)

		case filesystem.DeviceInfoColumns.MountPoint:
			lineCols = append(lineCols, info.MountPoint)

		case filesystem.DeviceInfoColumns.FsType:
			lineCols = append(lineCols, info.FSType)

		case filesystem.DeviceInfoColumns.Size:
			lineCols = append(lineCols, info.Size)

		case filesystem.DeviceInfoColumns.Used:
			lineCols = append(lineCols, info.Used)

		case filesystem.DeviceInfoColumns.Avail:
			lineCols = append(lineCols, info.Avail)

		case filesystem.DeviceInfoColumns.UsePercent:
			lineCols = append(lineCols, info.UsePercent)

		case filesystem.DeviceInfoColumns.FreePercent:
			lineCols = append(lineCols, info.FreePercent)
		}
	}

	return lineCols
}
