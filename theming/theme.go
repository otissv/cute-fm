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
)

type Style struct {
	Background    string
	Foreground    string
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
}

// Theme holds color configuration loaded from lsfm.toml.
// Keys are fairly free-form and mapped from the flat key/value config file.
type Theme struct {
	// FileTypeColors maps keys like "directory", "symlink", "socket",
	// "pipe", "device", "executable", "regular" to style specs.
	FileTypeColors map[string]string

	// FieldColors maps keys like "user", "group", "size", "time", "nlink".
	FieldColors map[string]string

	// PermissionColors maps keys like "full", "partial", "none".
	PermissionColors map[string]string

	// Permission bit colors for individual permission characters.
	// These are used to color the "rwx" bits in the permissions column.
	PermRead  string
	PermWrite string
	PermExec  string

	// Interface colors.
	BorderColor        string
	SelectedForeground string
	SelectedBackground string
	Foreground         string
	Background         string
	SearchPlaceholder  string
	HelpText           string
	TitleForeground    string
	TitleBackground    string
	FileList           Style
	Preview            Style
	StatusBar          Style
}

// DefaultTheme returns a sane fallback theme used when the config
// file cannot be read or parsed.
func DefaultTheme() Theme {
	return Theme{
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
		PermissionColors: map[string]string{
			"full":    "#00FF00",
			"partial": "#FFFF00",
			"none":    foregroundColor,
		},
		Background:         backgroundColor,
		BorderColor:        "#888888",
		Foreground:         foregroundColor,
		HelpText:           foregroundColor,
		PermExec:           "#FF5F00",
		PermRead:           "#00FF00",
		PermWrite:          "#FFFF00",
		SearchPlaceholder:  "#3B3B3B",
		SelectedBackground: "#3B3B3B",
		SelectedForeground: backgroundColor,
		TitleBackground:    "#25A065",
		TitleForeground:    foregroundColor,

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

		// Permission colors
		case "full", "partial", "none":
			if theme.PermissionColors == nil {
				theme.PermissionColors = map[string]string{}
			}
			theme.PermissionColors[k] = v

		// Permission bit colors
		case "perm_read":
			theme.PermRead = v
		case "perm_write":
			theme.PermWrite = v
		case "perm_exec":
			theme.PermExec = v

		// Interface colors
		case "border":
			theme.BorderColor = v
		case "selected_foreground":
			theme.SelectedForeground = v
		case "selected_background":
			theme.SelectedBackground = v
		case "foreground":
			theme.Foreground = v
		case "background":
			theme.Background = v
		case "search_placeholder":
			theme.SearchPlaceholder = v
		case "help_text":
			theme.HelpText = v
		case "title_foreground":
			theme.TitleForeground = v
		case "title_background":
			theme.TitleBackground = v
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
