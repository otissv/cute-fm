package tui

type Keybinding struct {
	Description string
	On          []string
}

type Keybindings struct {
	Add         Keybinding
	Cancel      Keybinding
	Cd          Keybinding
	CdParent    Keybinding
	Command     Keybinding
	Copy        Keybinding
	Directories Keybinding
	Down        Keybinding
	Enter       Keybinding
	Files       Keybinding
	Filter      Keybinding
	GoToStart   Keybinding
	GoToEnd     Keybinding
	Help        Keybinding
	HiddenFiles Keybinding
	List        Keybinding
	Mkdir       Keybinding
	Move        Keybinding
	Paste       Keybinding
	Preview     Keybinding
	Quit        Keybinding
	Redo        Keybinding
	Rename      Keybinding
	Select      Keybinding
	Tab         Keybinding
	Undo        Keybinding
	Up          Keybinding
}

func GetKeyBindings() Keybindings {
	bindings := Keybindings{
		Add: Keybinding{
			On:          []string{"n"},
			Description: "Create new file.",
		},
		Cancel: Keybinding{
			On:          []string{"esc", "ctrl+q"},
			Description: "Close window.",
		},
		Cd: Keybinding{
			On:          []string{"c"},
			Description: "Change directory.",
		},
		CdParent: Keybinding{
			On:          []string{"backspace", "backspace2"},
			Description: "Change directory.",
		},
		Command: Keybinding{
			On:          []string{":"},
			Description: "Enter Commands.",
		},
		Copy: Keybinding{
			On:          []string{"y"},
			Description: "Copy file or directory.",
		},
		Directories: Keybinding{
			On:          []string{"ctrl+d"},
			Description: "List directories only.",
		},
		Down: Keybinding{
			On:          []string{"down"},
			Description: "Move selection down.",
		},
		Enter: Keybinding{
			On:          []string{"enter"},
			Description: "Execute a command.",
		},
		Files: Keybinding{
			On:          []string{"ctrl+f"},
			Description: "List files only.",
		},
		Filter: Keybinding{
			On:          []string{"f"},
			Description: "Filter directory content.",
		},

		GoToStart: Keybinding{
			On:          []string{"g"},
			Description: "Goto start",
		},

		GoToEnd: Keybinding{
			On:          []string{"G"},
			Description: "Goto end",
		},

		Help: Keybinding{
			On:          []string{"?"},
			Description: "View help (This window).",
		},
		HiddenFiles: Keybinding{
			On:          []string{"h"},
			Description: "Toggle hidden files.",
		},
		List: Keybinding{
			On:          []string{"ctrl+l"},
			Description: "List directory contents.",
		},
		Move: Keybinding{
			On:          []string{"m"},
			Description: "Move file or directory.",
		},
		Mkdir: Keybinding{
			On:          []string{"k"},
			Description: "Create a new directory.",
		},
		Paste: Keybinding{
			On:          []string{"v"},
			Description: "Paste file or directory.",
		},
		Preview: Keybinding{
			On:          []string{"w"},
			Description: "Preview file or folder.",
		},
		Quit: Keybinding{
			On:          []string{"q"},
			Description: "Quit the application.",
		},
		Rename: Keybinding{
			On:          []string{"r"},
			Description: "Rename file or folder.",
		},
		Redo: Keybinding{
			On:          []string{"ctrl+z"},
			Description: "Redo.",
		},
		Select: Keybinding{
			On:          []string{"s"},
			Description: "Select files or directories.",
		},
		Tab: Keybinding{
			On:          []string{"tab"},
			Description: "Auto complete.",
		},
		Undo: Keybinding{
			On:          []string{"z"},
			Description: "Undo.",
		},
		Up: Keybinding{
			On:          []string{"up", "h"},
			Description: "Move selection up.",
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
