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
	AddFile      Keybinding
	AutoComplete Keybinding
	Cancel       Keybinding
	Cd           Keybinding
	Command      Keybinding
	Copy         Keybinding
	Directories  Keybinding
	Down         Keybinding
	Enter        Keybinding
	Files        Keybinding
	Filter       Keybinding
	GoToEnd      Keybinding
	GoToStart    Keybinding
	Help         Keybinding
	HiddenFiles  Keybinding
	List         Keybinding
	Mkdir        Keybinding
	Move         Keybinding
	Parent       Keybinding
	Paste        Keybinding
	Preview      Keybinding
	Quit         Keybinding
	Remove       Keybinding
	Redo         Keybinding
	Rename       Keybinding
	Select       Keybinding
	Undo         Keybinding
	Up           Keybinding
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
		Parent: Keybinding{
			On:          []string{"backspace", "backspace2"},
			Description: "Change directory to parent directory.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Command: Keybinding{
			On:          []string{":"},
			Description: "Enter Commands.",
			Category:    KeybindingCategories.Command.Name,
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
			On:          []string{"down"},
			Description: "Move selection down.",
			Category:    KeybindingCategories.Navigation.Name,
		},
		Enter: Keybinding{
			On:          []string{"enter"},
			Description: "Execute a command.",
			Category:    KeybindingCategories.General.Name,
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
			On:          []string{"h"},
			Description: "Toggle hidden files.",
			Category:    KeybindingCategories.Views.Name,
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
			On:          []string{"k"},
			Description: "Create a new directory.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Paste: Keybinding{
			On:          []string{"v"},
			Description: "Paste file or directory.",
			Category:    KeybindingCategories.Editing.Name,
		},
		Preview: Keybinding{
			On:          []string{"w"},
			Description: "Preview file or folder.",
			Category:    KeybindingCategories.Views.Name,
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
			Description: "Redo.",
			Category:    KeybindingCategories.General.Name,
		},
		Select: Keybinding{
			On:          []string{"s"},
			Description: "Select files or directories.",
			Category:    KeybindingCategories.Editing.Name,
		},

		Undo: Keybinding{
			On:          []string{"z"},
			Description: "Undo.",
			Category:    KeybindingCategories.General.Name,
		},
		Up: Keybinding{
			On:          []string{"up", "h"},
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
