package theming

import (
	"os"
	"path/filepath"
	"strings"

	"cute/config"

	"charm.land/lipgloss/v2"
	"github.com/pelletier/go-toml/v2"
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
	color10 = "#3B3B3B"
	color11 = "#33282E"

	background               = ""
	foreground               = color1
	borderColor              = color11
	activeBorderColor        = color2
	primary                  = color2
	muted                    = color11
	secondary                = color3
	placeholder              = color11
	viewModBackground        = background
	viewModForeground        = foreground
	leftCurrentDirBackground = background
	leftCurrentDirForeground = foreground
	fieldName                = foreground
	fieldGroup               = color3
	fieldNlink               = foreground
	fieldSize                = color4
	fieldTime                = foreground
	fieldUser                = color3
	fileTypeDevice           = color5
	fileTypeDirectory        = color8
	fileTypeExecutable       = color4
	fileTypePipe             = color6
	fileTypeRegular          = foreground
	fileTypeSocket           = color3
	fileTypeSymlink          = color5
	marked                   = color9
	hovered                  = color10
	permExec                 = color4
	permNone                 = foreground
	permRead                 = color6
	permWrite                = color5
	normalModeBackground     = ""
	normalModeForeground     = color4
	commandModeBackground    = ""
	commandModeForeground    = color7
	filterModeBackground     = ""
	filterModeForeground     = color4
	helpModeBackground       = ""
	helpModeForeground       = color6
	quitModeBackground       = "#000000"
	quitModeForeground       = "#F0EDED"
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

type FileListMode struct {
	ListModeBackground     string
	ListModeModeForeground string
	FileModeBackground     string
	DirModeForeground      string
	DirerModeBackground    string
}

type Theme struct {
	Foreground        string
	Background        string
	Primary           string
	Secondary         string
	Muted             string
	Placeholder       string
	BorderColor       string
	ActiveBorderColor string
	Hovered           string
	Marked            string
	CurrentDir        StyleColor
	FieldColors       map[string]string
	FileTypeColors    map[string]string
	Permissions       PermissionsStyle
	TuiMode           TuiMode
	ViewMode          StyleColor
}

func GetTheme() Theme {
	theme := Theme{
		Background:        background,
		Foreground:        foreground,
		BorderColor:       borderColor,
		Primary:           primary,
		Secondary:         secondary,
		Muted:             muted,
		ActiveBorderColor: activeBorderColor,
		Placeholder:       placeholder,
		Marked:            marked,
		Hovered:           hovered,

		CurrentDir: StyleColor{
			Background: leftCurrentDirBackground,
			Foreground: leftCurrentDirForeground,
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
			"name":  fieldName,
			"group": fieldGroup,
			"nlink": fieldNlink,
			"size":  fieldSize,
			"time":  fieldTime,
			"user":  fieldUser,
		},

		Permissions: PermissionsStyle{
			Exec:  permExec,
			Read:  permRead,
			Write: permWrite,
			None:  permNone,
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

	customTheme := loadThemeConfig()
	return mergeTheme(theme, customTheme)
}

func loadThemeConfig() Theme {
	customTheme := Theme{}

	configDir := config.GetConfigDir()
	themePath := filepath.Join(configDir, "theme.toml")

	// Check if theme.toml exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return customTheme
	}

	// Read the theme file
	data, err := os.ReadFile(themePath)
	if err != nil {
		// If we can't read it, return default theme
		return customTheme
	}

	// Parse the TOML file
	if err := toml.Unmarshal(data, &customTheme); err != nil {
		// If parsing fails, return default theme
		return customTheme
	}

	return customTheme
}

// mergeTheme merges a custom theme into the default theme, only overriding non-empty values
func mergeTheme(defaultTheme, customTheme Theme) Theme {
	merged := defaultTheme

	// Merge top-level string fields
	if customTheme.Foreground != "" {
		merged.Foreground = customTheme.Foreground
	}
	if customTheme.Background != "" {
		merged.Background = customTheme.Background
	}
	if customTheme.Primary != "" {
		merged.Primary = customTheme.Primary
	}
	if customTheme.Secondary != "" {
		merged.Secondary = customTheme.Secondary
	}
	if customTheme.Muted != "" {
		merged.Muted = customTheme.Muted
	}
	if customTheme.BorderColor != "" {
		merged.BorderColor = customTheme.BorderColor
	}

	// Merge nested structs
	merged.Permissions = mergePermissionsStyle(defaultTheme.Permissions, customTheme.Permissions)
	merged.TuiMode = mergeTuiMode(defaultTheme.TuiMode, customTheme.TuiMode)
	merged.ViewMode = mergeStyleColor(defaultTheme.ViewMode, customTheme.ViewMode)

	// Merge maps (only override keys that exist in custom theme)
	if customTheme.FieldColors != nil {
		if merged.FieldColors == nil {
			merged.FieldColors = make(map[string]string)
		}
		for k, v := range customTheme.FieldColors {
			if v != "" {
				merged.FieldColors[k] = v
			}
		}
	}

	if customTheme.FileTypeColors != nil {
		if merged.FileTypeColors == nil {
			merged.FileTypeColors = make(map[string]string)
		}
		for k, v := range customTheme.FileTypeColors {
			if v != "" {
				merged.FileTypeColors[k] = v
			}
		}
	}

	return merged
}

// Helper functions to merge nested structs
func mergeStyleColor(defaultStyle, customStyle StyleColor) StyleColor {
	merged := defaultStyle
	if customStyle.Background != "" {
		merged.Background = customStyle.Background
	}
	if customStyle.Foreground != "" {
		merged.Foreground = customStyle.Foreground
	}
	return merged
}

func mergePermissionsStyle(defaultStyle, customStyle PermissionsStyle) PermissionsStyle {
	merged := defaultStyle
	if customStyle.Exec != "" {
		merged.Exec = customStyle.Exec
	}
	if customStyle.Read != "" {
		merged.Read = customStyle.Read
	}
	if customStyle.Write != "" {
		merged.Write = customStyle.Write
	}
	if customStyle.None != "" {
		merged.None = customStyle.None
	}
	return merged
}

func mergeTuiMode(defaultMode, customMode TuiMode) TuiMode {
	merged := defaultMode
	if customMode.NormalModeBackground != "" {
		merged.NormalModeBackground = customMode.NormalModeBackground
	}
	if customMode.NormalModeForeground != "" {
		merged.NormalModeForeground = customMode.NormalModeForeground
	}
	if customMode.CommandModeBackground != "" {
		merged.CommandModeBackground = customMode.CommandModeBackground
	}
	if customMode.CommandModeForeground != "" {
		merged.CommandModeForeground = customMode.CommandModeForeground
	}
	if customMode.FilterModeBackground != "" {
		merged.FilterModeBackground = customMode.FilterModeBackground
	}
	if customMode.FilterModeForeground != "" {
		merged.FilterModeForeground = customMode.FilterModeForeground
	}
	if customMode.HelpModeBackground != "" {
		merged.HelpModeBackground = customMode.HelpModeBackground
	}
	if customMode.HelpModeForeground != "" {
		merged.HelpModeForeground = customMode.HelpModeForeground
	}
	if customMode.QuitModeBackground != "" {
		merged.QuitModeBackground = customMode.QuitModeBackground
	}
	if customMode.QuitModeForeground != "" {
		merged.QuitModeForeground = customMode.QuitModeForeground
	}
	return merged
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
