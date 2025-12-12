package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"cute/filesystem"
)

// SettingsTOML represents the TOML structure for settings
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
		case "info", "file_info":
			merged.SplitPane = FileInfoSplitPaneType
		case "list", "file_list":
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
		for _, colName := range tomlSettings.Columns {
			col := parseColumnName(colName)
			if col != "" {
				columns = append(columns, col)
			}
		}
		if len(columns) > 0 {
			merged.ColumnVisibility = columns
		}
	}

	if tomlSettings.Sorting.Column != "" {
		col := parseColumnName(tomlSettings.Sorting.Column)
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

func parseColumnName(name string) filesystem.FileInfoColumn {
	switch strings.ToLower(name) {
	case "permissions", "Permissions":
		return filesystem.ColumnPermissions
	case "size", "Size":
		return filesystem.ColumnSize
	case "type", "mimetype", "Type", "Mimetype":
		return filesystem.ColumnMimeType
	case "user", "User":
		return filesystem.ColumnUser
	case "group", "Group":
		return filesystem.ColumnGroup
	case "datemodified", "date_modified", "date", "dateModified", "Date":
		return filesystem.ColumnDateModified
	case "name", "Name":
		return filesystem.ColumnName
	default:
		return ""
	}
}
