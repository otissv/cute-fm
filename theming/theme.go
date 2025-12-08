package theming

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	color0  = "#1E1E1E"
	color1  = "#F0EDED"
	color2  = "#F25D94"
	color3  = "#FFF1A8"
	color4  = "#FF9BC0"
	color5  = "#FAD2E1"
	color6  = "#7CFFD2"
	color7  = "#E37CFF"
	color8  = "#A8D2FF"
	color9  = "#2072D5"
	color10 = "#F00F00"

	background               = ""
	foreground               = color1
	borderColor              = color2
	primary                  = color2
	secondary                = color3
	placeholder              = "#33282E"
	viewModBackground        = background
	viewModForeground        = foreground
	commandBarBackground     = background
	commandBarBorder         = borderColor
	commandBarForeground     = foreground
	commandBarPlaceholder    = "#A8A7A7"
	leftCurrentDirBackground = background
	leftCurrentDirForeground = foreground
	fieldGroup               = color3
	fieldNlink               = foreground
	fieldSize                = color4
	fieldTime                = foreground
	fieldUser                = color3
	fileListBackground       = background
	fileListForeGround       = foreground
	fileListBorder           = background
	fileTypeDevice           = color5
	fileTypeDirectory        = color8
	fileTypeExecutable       = color4
	fileTypePipe             = color6
	fileTypeRegular          = foreground
	fileTypeSocket           = color3
	fileTypeSymlink          = color5
	fileListMarked           = color9
	permExec                 = color4
	permNone                 = foreground
	permRead                 = color6
	permWrite                = color5
	searchBackground         = background
	searchBorder             = background
	searchForeground         = foreground
	headerBackground         = background
	previewBackground        = background
	previewBorderBackground  = ""
	previewForeground        = foreground
	previewBorder            = borderColor
	normalModeBackground     = color4
	normalModeForeground     = color0
	commandModeBackground    = color7
	commandModeForeground    = color0
	filterModeBackground     = color4
	filterModeForeground     = color0
	helpModeBackground       = color6
	helpModeForeground       = color0
	quitModeBackground       = "#000000"
	quitModeForeground       = "#F0EDED"
	dialogTitle              = color9
	sudoBackground           = color10
	sudoForeground           = color1
)

type Style struct {
	Background       string
	BorderBackground string
	Foreground       string
	PaddingTop       int
	PaddingBottom    int
	PaddingLeft      int
	PaddingRight     int
	Border           string
}

type StyleColor struct {
	Background string
	Foreground string
}

type FileListStyle struct {
	Background       string
	BorderBackground string
	Foreground       string
	PaddingTop       int
	PaddingBottom    int
	PaddingLeft      int
	PaddingRight     int
	Border           string
	Marked           string
}

type DialogStyle struct {
	Background    string
	Foreground    string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	Border        string
	Title         string
}

type BarStyle struct {
	Background    string
	Foreground    string
	Placeholder   string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	Border        string
}

type PermissionsStyle struct {
	Exec  string
	Read  string
	Write string
	None  string
}

type TuiMode struct {
	NormalModeBackground  string
	NormalModeForeground  string
	CommandModeBackground string
	CommandModeForeground string
	FilterModeBackground  string
	FilterModeForeground  string
	HelpModeBackground    string
	HelpModeForeground    string
	QuitModeBackground    string
	QuitModeForeground    string
}

type FilelistMode struct {
	ListModeBackground     string
	ListModeModeForeground string
	FileModeBackground     string
	DirModeForeground      string
	DirerModeBackground    string
}

type Theme struct {
	Foreground     string
	Background     string
	Primary        string
	Secondary      string
	BorderColor    string
	CommandBar     BarStyle
	CurrentDir     StyleColor
	Dialog         DialogStyle
	FieldColors    map[string]string
	FileList       FileListStyle
	FileTypeColors map[string]string
	Header         StyleColor
	Permissions    PermissionsStyle
	Preview        Style
	SearchBar      BarStyle
	Selection      StyleColor
	StatusBar      Style
	SudoMode       StyleColor
	TuiMode        TuiMode
	ViewMode       StyleColor
}

// DefaultTheme returns a sane fallback theme used when the config
// file cannot be read or parsed.
func DefaultTheme() Theme {
	return Theme{
		Background:  background,
		Foreground:  foreground,
		BorderColor: borderColor,
		Primary:     primary,
		Secondary:   secondary,

		CommandBar: BarStyle{
			Background:    commandBarBackground,
			Border:        commandBarBorder,
			Foreground:    commandBarForeground,
			PaddingBottom: 0,
			PaddingLeft:   1,
			PaddingRight:  1,
			PaddingTop:    0,
			Placeholder:   commandBarPlaceholder,
		},

		Dialog: DialogStyle{
			Background:    background,
			Border:        borderColor,
			Foreground:    foreground,
			PaddingBottom: 1,
			PaddingLeft:   1,
			PaddingRight:  1,
			PaddingTop:    1,
			Title:         dialogTitle,
		},

		CurrentDir: StyleColor{
			Background: leftCurrentDirBackground,
			Foreground: leftCurrentDirForeground,
		},

		FileList: FileListStyle{
			Background:    fileListBackground,
			Foreground:    fileListForeGround,
			PaddingBottom: 1,
			Border:        fileListBorder,
			PaddingLeft:   1,
			PaddingRight:  1,
			PaddingTop:    0,
			Marked:        fileListMarked,
		},

		FileTypeColors: map[string]string{
			"directory":  fileTypeDirectory,
			"symlink":    fileTypeSymlink,
			"socket":     fileTypeSocket,
			"pipe":       fileTypePipe,
			"device":     fileTypeDevice,
			"executable": fileTypeExecutable,
			"regular":    fileTypeRegular,
		},

		FieldColors: map[string]string{
			"group": fieldGroup,
			"nlink": fieldNlink,
			"size":  fieldSize,
			"time":  fieldTime,
			"user":  fieldUser,
		},

		Header: StyleColor{
			Background: headerBackground,
		},

		Permissions: PermissionsStyle{
			Exec:  permExec,
			Read:  permRead,
			Write: permWrite,
			None:  permNone,
		},

		Preview: Style{
			Background:       previewBackground,
			BorderBackground: previewBorderBackground,
			Foreground:       previewForeground,
			Border:           previewBorder,
			PaddingBottom:    1,
			PaddingLeft:      1,
			PaddingRight:     1,
			PaddingTop:       0,
		},

		SearchBar: BarStyle{
			Background:    background,
			Border:        borderColor,
			Foreground:    foreground,
			PaddingBottom: 0,
			PaddingLeft:   1,
			PaddingRight:  1,
			PaddingTop:    0,
			Placeholder:   placeholder,
		},

		StatusBar: Style{
			Background:    searchBackground,
			Foreground:    searchForeground,
			Border:        searchBorder,
			PaddingBottom: 0,
			PaddingLeft:   0,
			PaddingRight:  0,
			PaddingTop:    0,
		},

		Selection: StyleColor{
			Background: "#3B3B3B",
			Foreground: background,
		},

		SudoMode: StyleColor{
			Background: sudoBackground,
			Foreground: sudoForeground,
		},

		TuiMode: TuiMode{
			CommandModeBackground: commandModeBackground,
			CommandModeForeground: commandModeForeground,
			FilterModeBackground:  filterModeBackground,
			FilterModeForeground:  filterModeForeground,
			HelpModeBackground:    helpModeBackground,
			HelpModeForeground:    helpModeForeground,
			NormalModeBackground:  normalModeBackground,
			NormalModeForeground:  normalModeForeground,
			QuitModeBackground:    quitModeBackground,
			QuitModeForeground:    quitModeForeground,
		},

		ViewMode: StyleColor{
			Background: viewModBackground,
			Foreground: viewModForeground,
		},
	}
}

// LoadThemeFromMap constructs a theme from a simple key/value map. The keys are
// the same as the ones previously used in the TOML-style configuration:
//
//   - File type colors: "directory", "symlink", "socket", "pipe", "device",
//     "executable", "regular"
//   - Field colors: "nlink", "user", "group", "size", "time"
//   - Interface colors: "border", "selected_foreground", "selected_background",
//     "foreground", "background"
//
// The map is applied as overrides on top of DefaultTheme.
func LoadThemeFromMap(raw map[string]string) Theme {
	theme := DefaultTheme()

	for k, v := range raw {
		switch k {
		// File type colors
		case "directory", "symlink", "socket", "pipe", "device", "executable", "regular":
			if theme.FileTypeColors == nil {
				theme.FileTypeColors = map[string]string{}
			}
			theme.FileTypeColors[k] = v

		// Field colors
		case "nlink", "user", "group", "size", "time":
			if theme.FieldColors == nil {
				theme.FieldColors = map[string]string{}
			}
			theme.FieldColors[k] = v

		// Interface colors
		case "border":
			theme.BorderColor = v
		case "selected_foreground":
			theme.Selection.Foreground = v
		case "selected_background":
			theme.Selection.Background = v
		case "foreground":
			theme.Foreground = v
		case "background":
			theme.Background = v
		}
	}

	return theme
}

// LoadTheme loads theme colors from the given path. The format is a very small
// subset of TOML: "key = \"value\"" lines, comments starting with '#', and
// blank lines are ignored. This is intentionally lenient and does not require
// a full TOML parser.
//
// This function remains for compatibility, but the main configuration path now
// uses Lua (see config.LoadRuntimeConfig).
func LoadTheme(path string) Theme {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultTheme()
	}

	raw := map[string]string{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Strip surrounding quotes if present.
		val = strings.Trim(val, `"`)

		if key != "" {
			raw[key] = val
		}
	}

	return LoadThemeFromMap(raw)
}

// StyleFromSpec builds a lipgloss style from a specification string, such as:
//
//	"#0000FF+bold"
//	"blue+bold"
//	"dim"
//
// Attributes supported: bold, dim, underline, italic.
func StyleFromSpec(spec string) lipgloss.Style {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return lipgloss.NewStyle()
	}

	style := lipgloss.NewStyle()
	parts := strings.Split(spec, "+")

	for _, p := range parts {
		token := strings.TrimSpace(p)
		if token == "" {
			continue
		}

		switch strings.ToLower(token) {
		case "bold":
			style = style.Bold(true)
		case "dim":
			style = style.Faint(true)
		case "underline":
			style = style.Underline(true)
		case "italic":
			style = style.Italic(true)
		default:
			// Treat as a color (hex, ANSI color name, or number)
			style = style.Foreground(lipgloss.Color(token))
		}
	}

	return style
}
