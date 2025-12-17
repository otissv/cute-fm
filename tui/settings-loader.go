package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"cute/filesystem"
)

type SettingsTOML struct {
	StartDir  string   `toml:"startDir"`
	SplitPane string   `toml:"splitPane"`
	FileMode  string   `toml:"fileMode"`
	Sorting   Sorting  `toml:"sorting"`
	Columns   []string `toml:"ColumnVisibility"`
}

type Sorting struct {
	Column    string `toml:"column"`
	Direction string `toml:"direction"`
}

func LoadSettings(configDir string) (*SettingsTOML, error) {
	settingsPath := filepath.Join(configDir, "settings.toml")

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings.toml: %w", err)
	}

	var settings SettingsTOML
	if err := toml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings.toml: %w", err)
	}

	return &settings, nil
}

func SaveDefaultSettings(configDir string) error {
	settingsPath := filepath.Join(configDir, "settings.toml")

	if _, err := os.Stat(settingsPath); err == nil {
		return nil
	}

	defaultSettings := SettingsTOML{
		StartDir:  "home",
		SplitPane: "info",
		FileMode:  "all",
		Sorting: Sorting{
			Column:    "name",
			Direction: "asc",
		},
		Columns: []string{
			"Permissions",
			"User",
			"Group",
			"DateModified",
			"Name",
		},
	}

	data, err := toml.Marshal(defaultSettings)
	if err != nil {
		return fmt.Errorf("failed to marshal default settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write settings.toml: %w", err)
	}

	return nil
}

func MergeSettings(defaultSettings Settings, tomlSettings *SettingsTOML) Settings {
	if tomlSettings == nil {
		return defaultSettings
	}

	merged := defaultSettings

	if tomlSettings.StartDir != "" {
		if tomlSettings.StartDir == "Home" || tomlSettings.StartDir == "home" || tomlSettings.StartDir == "$HOME" || tomlSettings.StartDir == "~/" {
			if homeDir, err := os.UserHomeDir(); err == nil {
				merged.StartDir = homeDir
			}
		} else {
			merged.StartDir = tomlSettings.StartDir
		}
	}

	if tomlSettings.SplitPane != "" {
		switch strings.ToLower(tomlSettings.SplitPane) {
		case "info":
			merged.SplitPane = FileInfoSplitPaneType
		case "list":
			merged.SplitPane = FileListSplitPaneType
		case "preview":
			merged.SplitPane = PreviewPaneType
		default:
			merged.SplitPane = FileInfoSplitPaneType
		}
	}

	if tomlSettings.FileMode != "" {
		switch strings.ToLower(tomlSettings.FileMode) {
		case "all", "ll":
			merged.FileListMode = FileListModeList
		case "files", "lf":
			merged.FileListMode = FileListModeFile
		case "dirs", "directories", "ld":
			merged.FileListMode = FileListModeDir
		default:
			merged.FileListMode = FileListModeList
		}
	}

	if len(tomlSettings.Columns) > 0 {
		columns := make([]filesystem.FileInfoColumn, 0, len(tomlSettings.Columns))
		for _, colNameWidth := range tomlSettings.Columns {
			col := parseName(colNameWidth)
			if col != "" {
				columns = append(columns, col)
			}
		}
		if len(columns) > 0 {
			merged.ColumnVisibility = columns
		}
	}

	if tomlSettings.Sorting.Column != "" {
		col := parseName(tomlSettings.Sorting.Column)
		if col != "" {
			merged.SortColumnBy = col
		}
	}

	if tomlSettings.Sorting.Direction != "" {
		switch strings.ToLower(tomlSettings.Sorting.Direction) {
		case "asc", "ascending":
			merged.SortColumnDirection = SortingAsc
		case "desc", "descending":
			merged.SortColumnDirection = SortingDesc
		default:
			merged.SortColumnDirection = SortingAsc
		}
	}

	return merged
}

func parseName(name string) filesystem.FileInfoColumn {
	switch strings.ToLower(name) {
	case "permissions", "Permissions":
		return filesystem.FileInfoColumns.Permissions
	case "size", "Size":
		return filesystem.FileInfoColumns.Size
	case "type", "mimetype", "Type", "Mimetype":
		return filesystem.FileInfoColumns.MimeType
	case "user", "User":
		return filesystem.FileInfoColumns.User
	case "group", "Group":
		return filesystem.FileInfoColumns.Group
	case "datemodified", "date_modified", "date", "dateModified", "Date":
		return filesystem.FileInfoColumns.Modified
	case "name", "Name":
		return filesystem.FileInfoColumns.Name
	default:
		return ""
	}
}

func SaveSettings(m Model) error {
	configDir := m.GetConfigDir()
	settingsPath := filepath.Join(configDir, "settings.toml")
	tomlSettings := SettingsToTOML(m)

	data, err := toml.Marshal(tomlSettings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write settings.toml: %w", err)
	}

	return nil
}

func SettingsToTOML(m Model) SettingsTOML {
	settings := m.settings
	sortBy := m.GetSortColumnBy()
	startDir := settings.StartDir

	if startDir != "" {
		if homeDir, err := os.UserHomeDir(); err == nil && startDir == homeDir {
			startDir = "Home"
		}
	}

	splitPane := ""
	switch settings.SplitPane {
	case FileInfoSplitPaneType:
		splitPane = "info"
	case FileListSplitPaneType:
		splitPane = "list"
	case PreviewPaneType:
		splitPane = "preview"
	default:
		splitPane = "info"
	}

	fileMode := ""
	switch settings.FileListMode {
	case FileListModeList:
		fileMode = "all"
	case FileListModeFile:
		fileMode = "files"
	case FileListModeDir:
		fileMode = "dirs"
	default:
		fileMode = "all"
	}

	columns := make([]string, 0, len(settings.ColumnVisibility))
	for _, col := range settings.ColumnVisibility {
		columns = append(columns, string(col))
	}

	sortColumn := ""
	if sortBy.Column() != "" {
		sortColumn = string(sortBy.Column())
	} else if settings.SortColumnBy != "" {
		sortColumn = string(settings.SortColumnBy)
	} else {
		sortColumn = "Name"
	}

	sortDirection := ""
	if sortBy.Direction() != "" {
		switch sortBy.Direction() {
		case SortingAsc:
			sortDirection = "asc"
		case SortingDesc:
			sortDirection = "desc"
		default:
			sortDirection = "asc"
		}
	} else if settings.SortColumnDirection != "" {
		switch settings.SortColumnDirection {
		case SortingAsc:
			sortDirection = "asc"
		case SortingDesc:
			sortDirection = "desc"
		default:
			sortDirection = "asc"
		}
	} else {
		sortDirection = "asc"
	}

	return SettingsTOML{
		StartDir:  startDir,
		SplitPane: splitPane,
		FileMode:  fileMode,
		Sorting: Sorting{
			Column:    sortColumn,
			Direction: sortDirection,
		},
		Columns: columns,
	}
}
