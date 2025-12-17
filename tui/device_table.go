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

	indexText := truncateAndPadCell(indexStyle.Render(fmt.Sprintf("%d", index)), colDeviceIndexWidth, bgColor)
	nameText := truncateAndPadCell(nameStyle.Render(di.Name), colDeviceNameWidth, bgColor)
	devicePathText := truncateAndPadCell(pathStyle.Render(di.Device), colDevicePathWidth, bgColor)
	mountText := truncateAndPadCell(mountStyle.Render(di.MountPoint), colDeviceMountWidth, bgColor)
	typeText := truncateAndPadCell(typeStyle.Render(di.FSType), colDeviceFSTypeWidth, bgColor)
	sizeText := truncateAndPadCell(sizeStyle.Render(di.Size), colDeviceSizeWidth, bgColor)
	usedText := truncateAndPadCell(usedStyle.Render(di.Used), colDeviceUsedWidth, bgColor)
	availText := truncateAndPadCell(availStyle.Render(di.Avail), colDeviceAvailWidth, bgColor)
	usePercentText := truncateAndPadCell(usePercentStyle.Render(di.UsePercent), colDeviceUsePercentWidth, bgColor)
	freePercentText := truncateAndPadCell(freePercentStyle.Render(di.FreePercent), colDeviceFreePercentWidth, bgColor)

	lineCols := []string{
		indexText,
		nameText,
	}

	for _, col := range d.columns {
		// Skip name column since it's already added above
		if col == filesystem.DeviceInfoColumns.Name {
			continue
		}

		switch col {
		case filesystem.DeviceInfoColumns.Device:
			lineCols = append(lineCols, devicePathText)

		case filesystem.DeviceInfoColumns.MountPoint:
			lineCols = append(lineCols, mountText)

		case filesystem.DeviceInfoColumns.FsType:
			lineCols = append(lineCols, typeText)

		case filesystem.DeviceInfoColumns.Size:
			lineCols = append(lineCols, sizeText)

		case filesystem.DeviceInfoColumns.Used:
			lineCols = append(lineCols, usedText)

		case filesystem.DeviceInfoColumns.Avail:
			lineCols = append(lineCols, availText)

		case filesystem.DeviceInfoColumns.UsePercent:
			lineCols = append(lineCols, usePercentText)

		case filesystem.DeviceInfoColumns.FreePercent:
			lineCols = append(lineCols, freePercentText)
		}

	}

	sep := " "
	if bgColor != "" {
		sep = lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(" ")
	}

	line := strings.Join(lineCols, sep)

	// Ensure the line never exceeds totalWidth - truncate if necessary
	if d.totalWidth > 0 {
		lineWidth := lipgloss.Width(line)
		if lineWidth > d.totalWidth {
			line = truncateString(line, d.totalWidth)
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
	Theme      theming.Theme
	TotalWidth int
}

func RenderDeviceHeaderRow(args DeviceHeaderRowArgs) string {
	bgColor := args.Theme.FileList.Background
	bg := lipgloss.Color(bgColor)

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theming.GetTheme().Foreground))

	indexText := truncateAndPadCell(baseStyle.Render(" "), colDeviceIndexWidth, bgColor)
	nameText := truncateAndPadCell(baseStyle.Render("Name"), colDeviceNameWidth, bgColor)
	pathText := truncateAndPadCell(baseStyle.Render("Device"), colDevicePathWidth, bgColor)
	mountText := truncateAndPadCell(baseStyle.Render("Mount Point"), colDeviceMountWidth, bgColor)
	typeText := truncateAndPadCell(baseStyle.Render("FS Type"), colDeviceFSTypeWidth, bgColor)
	sizeText := truncateAndPadCell(baseStyle.Render("Size"), colDeviceSizeWidth, bgColor)
	usedText := truncateAndPadCell(baseStyle.Render("Used"), colDeviceUsedWidth, bgColor)
	availText := truncateAndPadCell(baseStyle.Render("Avail"), colDeviceAvailWidth, bgColor)
	usePercentText := truncateAndPadCell(baseStyle.Render("Use%"), colDeviceUsePercentWidth, bgColor)
	freePercentText := truncateAndPadCell(baseStyle.Render("Free%"), colDeviceFreePercentWidth, bgColor)

	lineCols := []string{
		indexText,
		nameText,
		pathText,
		mountText,
		typeText,
		sizeText,
		usedText,
		availText,
		usePercentText,
		freePercentText,
	}

	sep := lipgloss.NewStyle().Background(bg).Render(" ")
	line := strings.Join(lineCols, sep)

	// Ensure header line never exceeds totalWidth
	if args.TotalWidth > 0 {
		lineWidth := lipgloss.Width(line)
		if lineWidth > args.TotalWidth {
			line = truncateString(line, args.TotalWidth)
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
