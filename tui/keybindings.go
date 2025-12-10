package tui

type category string

type Keybinding struct {
	Category    category
	Description string
	On          []string
}

type KeybindingCategoryField struct {
	Name        category
	Description string
}

type KeybindingCategory struct {
	General    KeybindingCategoryField
	Navigation KeybindingCategoryField
	Filter     KeybindingCategoryField
	Editing    KeybindingCategoryField
	Help       KeybindingCategoryField
	Views      KeybindingCategoryField
	Command    KeybindingCategoryField
}

var KeybindingCategories = KeybindingCategory{
	General: KeybindingCategoryField{
		Name:        "General",
		Description: "",
	},
	Navigation: KeybindingCategoryField{
		Name:        "Navigation",
		Description: "",
	},
	Filter: KeybindingCategoryField{
		Name:        "Filter",
		Description: "",
	},
	Editing: KeybindingCategoryField{
		Name:        "Editing",
		Description: "",
	},
	Help: KeybindingCategoryField{
		Name:        "Help",
		Description: "",
	},
	Views: KeybindingCategoryField{
		Name:        "Views",
		Description: "",
	},
	Command: KeybindingCategoryField{
		Name:        "Command",
		Description: "",
	},
}

type Keybindings struct {
	AddFile                Keybinding
	AutoComplete           Keybinding
	Cancel                 Keybinding
	Cd                     Keybinding
	Command                Keybinding
	ColumnVisibiliy        Keybinding
	Copy                   Keybinding
	Directories            Keybinding
	Down                   Keybinding
	Enter                  Keybinding
	FileInfoPane           Keybinding
	Files                  Keybinding
	Filter                 Keybinding
	GoToEnd                Keybinding
	GoToStart              Keybinding
	Goto                   Keybinding
	Help                   Keybinding
	HiddenFiles            Keybinding
	Home                   Keybinding
	List                   Keybinding
	Mkdir                  Keybinding
	Move                   Keybinding
	NextDir                Keybinding
	PageDown               Keybinding
	PageUp                 Keybinding
	Parent                 Keybinding
	Paste                  Keybinding
	PreviewPane            Keybinding
	PreviousDir            Keybinding
	Quit                   Keybinding
	Redo                   Keybinding
	Remove                 Keybinding
	Rename                 Keybinding
	Select                 Keybinding
	SelectAll              Keybinding
	Sort                   Keybinding
	Sudo                   Keybinding
	SwitchBetweenSplitPane Keybinding
	ToggleRightPane        Keybinding
	Undo                   Keybinding
	Up                     Keybinding
}

func GetKeyBindings() Keybindings {
	bindings := Keybindings{
		AddFile: Keybinding{
			On:          []string{"a"},
			Description: "Create new file.",
			Category:    KeybindingCategories.Editing.Name,
		},
		AutoComplete: Keybinding{
			On:          []string{"tab"},
			Description: "Auto complete.",
			Category:    KeybindingCategories.Command.Name,
		},
		Cancel: Keybinding{
			On:          []string{"esc", "ctrl+q"},
			Description: "Close window.",
			Category:    KeybindingCategories.General.Name,
		},
		Cd: Keybinding{
			On:          []string{"c"},
			Description: "Change directory.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Command: Keybinding{
			On:          []string{":"},
			Description: "Enter Commands.",
			Category:    KeybindingCategories.Command.Name,
		},
		ColumnVisibiliy: Keybinding{
			On:          []string{"["},
			Description: "Show and hide columns",
			Category:    KeybindingCategories.General.Name,
		},
		Copy: Keybinding{
			On:          []string{"y"},
			Description: "Copy file or directory.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Directories: Keybinding{
			On:          []string{"ctrl+d"},
			Description: "List directories only.",
			Category:    KeybindingCategories.Views.Name,
		},
		Down: Keybinding{
			On:          []string{"down", "j"},
			Description: "Move selection down.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Enter: Keybinding{
			On:          []string{"enter"},
			Description: "Execute a command.",
			Category:    KeybindingCategories.General.Name,
		},
		FileInfoPane: Keybinding{
			On:          []string{"ctrl+i"},
			Description: "Execute a command.",
			Category:    KeybindingCategories.Views.Name,
		},
		Files: Keybinding{
			On:          []string{"ctrl+f"},
			Description: "List files only.",
			Category:    KeybindingCategories.Views.Name,
		},
		Filter: Keybinding{
			On:          []string{"f"},
			Description: "Filter directory content.",
			Category:    KeybindingCategories.Filter.Name,
		},
		Goto: Keybinding{
			On:          []string{"-", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"},
			Description: "Enter goto mode to jump to a selection",
			Category:    KeybindingCategories.Navigation.Name,
		},
		GoToStart: Keybinding{
			On:          []string{"g", "home"},
			Description: "Goto start",
			Category:    KeybindingCategories.Navigation.Name,
		},
		GoToEnd: Keybinding{
			On:          []string{"G", "end"},
			Description: "Goto end",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Help: Keybinding{
			On:          []string{"?"},
			Description: "View KeybindingCategories.Help (This window).",
			Category:    KeybindingCategories.Help.Name,
		},
		HiddenFiles: Keybinding{
			On:          []string{"i"},
			Description: "Toggle hidden files.",
			Category:    KeybindingCategories.Views.Name,
		},
		Home: Keybinding{
			On:          []string{"~"},
			Description: "Goto home directory",
			Category:    KeybindingCategories.General.Name,
		},
		List: Keybinding{
			On:          []string{"ctrl+l"},
			Description: "List directory contents.",
			Category:    KeybindingCategories.Views.Name,
		},
		Move: Keybinding{
			On:          []string{"m"},
			Description: "Move file or directory.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Mkdir: Keybinding{
			On:          []string{"A"},
			Description: "Create a new directory.",
			Category:    KeybindingCategories.Editing.Name,
		},
		NextDir: Keybinding{
			On:          []string{"right", "l"},
			Description: "Go to next directory in history.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		PageDown: Keybinding{
			On:          []string{"pgdown"},
			Description: "Move selection one page down.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		PageUp: Keybinding{
			On:          []string{"pgup"},
			Description: "Move selection one page up.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Parent: Keybinding{
			On:          []string{"backspace"},
			Description: "Change directory to parent directory.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		// leave for slipt view
		// Paste: Keybinding{
		// 	On:          []string{"v"},
		// 	Description: "Paste file or directory.",
		// 	Category:    KeybindingCategories.Editing.Name,
		// },
		PreviewPane: Keybinding{
			On:          []string{"ctrl+w"},
			Description: "Preview file or folder.",
			Category:    KeybindingCategories.Views.Name,
		},
		PreviousDir: Keybinding{
			On:          []string{"left", "h"},
			Description: "Go to previous directory in history.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Quit: Keybinding{
			On:          []string{"q"},
			Description: "Quit the application.",
			Category:    KeybindingCategories.General.Name,
		},
		Remove: Keybinding{
			On:          []string{"r"},
			Description: "Remove files or folders.",
			Category:    KeybindingCategories.Command.Name,
		},
		Rename: Keybinding{
			On:          []string{"n"},
			Description: "Rename file or folder.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Redo: Keybinding{
			On:          []string{"ctrl+z"},
			Description: "Redo previuos action.",
			Category:    KeybindingCategories.General.Name,
		},
		Select: Keybinding{
			On:          []string{"space"},
			Description: "Select item.",
			Category:    KeybindingCategories.Editing.Name,
		},
		SelectAll: Keybinding{
			On:          []string{"ctrl+a"},
			Description: "Toggle select all items.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Sort: Keybinding{
			On:          []string{"]"},
			Description: "Enter sudo mode",
			Category:    KeybindingCategories.General.Name,
		},
		Sudo: Keybinding{
			On:          []string{"u"},
			Description: "Sort file list colomns",
			Category:    KeybindingCategories.General.Name,
		},
		SwitchBetweenSplitPane: Keybinding{
			On:          []string{"tab"},
			Description: "Toggle beright pane.",
			Category:    KeybindingCategories.Views.Name,
		},
		ToggleRightPane: Keybinding{
			On:          []string{"t"},
			Description: "Toggle right pane.",
			Category:    KeybindingCategories.Views.Name,
		},
		Undo: Keybinding{
			On:          []string{"z"},
			Description: "Undo last action.",
			Category:    KeybindingCategories.General.Name,
		},
		Up: Keybinding{
			On:          []string{"up", "k"},
			Description: "Move selection up.",
			Category:    KeybindingCategories.Navigation.Name,
		},
	}

	return bindings
}

func (k Keybinding) Matches(key string) bool {
	for _, v := range k.On {
		if v == key {
			return true
		}
	}
	return false
}
