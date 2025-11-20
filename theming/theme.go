package theming

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	backgroundColor = "#212121"
	foregroundColor = "#F0EDED"
	borderColor     = "#888888"
)

type Style struct {
	Background    string
	Foreground    string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
}

type StyleColor struct {
	Background string
	Foreground string
}

type DefaultDialogStyle struct {
	Background    string
	Foreground    string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	BorderColor   string
}

type BarStyle struct {
	Background    string
	Foreground    string
	Placeholder   string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	BorderColor   string
}

type PermissionsStyle struct {
	Exec  string
	Read  string
	Write string
	None  string
}

type Theme struct {
	FileTypeColors map[string]string
	FieldColors    map[string]string
	Permissions    PermissionsStyle
	BorderColor    string
	Selection      StyleColor
	Foreground     string
	Background     string
	FileList       Style
	Preview        Style
	StatusBar      Style
	DefaultDialog  DefaultDialogStyle
	SearchBar      BarStyle
	CommandBar     BarStyle
}

// DefaultTheme returns a sane fallback theme used when the config
// file cannot be read or parsed.
func DefaultTheme() Theme {
	return Theme{
		Background:  backgroundColor,
		Foreground:  foregroundColor,
		BorderColor: borderColor,
		FileTypeColors: map[string]string{
			"directory":  "#0000FF+bold",
			"symlink":    "#00FFFF",
			"socket":     "#FFFF00",
			"pipe":       "#FF00FF",
			"device":     "#FFFF00+bold",
			"executable": "#00FF00",
			"regular":    foregroundColor,
		},
		FieldColors: map[string]string{
			"user":  "#00FFFF",
			"group": "#00FFFF",
			"size":  "#FFFF00",
			"time":  foregroundColor,
			"nlink": foregroundColor,
		},
		Permissions: PermissionsStyle{
			Exec:  "#FF5F00",
			Read:  "#00FF00",
			Write: "#FFFF00",
			None:  foregroundColor,
		},

		SearchBar: BarStyle{
			Background:    backgroundColor,
			Foreground:    foregroundColor,
			Placeholder:   "#3B3B3B",
			PaddingTop:    0,
			PaddingBottom: 0,
			PaddingLeft:   1,
			PaddingRight:  1,
			BorderColor:   borderColor,
		},

		CommandBar: BarStyle{
			Background:    backgroundColor,
			Foreground:    foregroundColor,
			Placeholder:   "#3B3B3B",
			PaddingTop:    0,
			PaddingBottom: 0,
			PaddingLeft:   1,
			PaddingRight:  1,
			BorderColor:   borderColor,
		},

		Selection: StyleColor{
			Background: "#3B3B3B",
			Foreground: backgroundColor,
		},

		FileList: Style{
			Background:    "",
			Foreground:    "",
			PaddingTop:    0,
			PaddingBottom: 1,
			PaddingLeft:   1,
			PaddingRight:  1,
		},

		Preview: Style{
			Background:    "",
			Foreground:    "",
			PaddingTop:    0,
			PaddingBottom: 1,
			PaddingLeft:   1,
			PaddingRight:  1,
		},

		StatusBar: Style{
			Background:    foregroundColor,
			Foreground:    backgroundColor,
			PaddingTop:    0,
			PaddingBottom: 0,
			PaddingLeft:   1,
			PaddingRight:  1,
		},

		DefaultDialog: DefaultDialogStyle{
			Background:    backgroundColor,
			Foreground:    foregroundColor,
			PaddingTop:    1,
			PaddingBottom: 1,
			PaddingLeft:   1,
			PaddingRight:  1,
			BorderColor:   "#f00f00",
		},
	}
}

// LoadTheme loads theme colors from the given path. The format is a very small
// subset of TOML: "key = \"value\"" lines, comments starting with '#', and
// blank lines are ignored. This is intentionally lenient and does not require
// a full TOML parser.
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
