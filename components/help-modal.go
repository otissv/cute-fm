package components

import (
	"fmt"
	"strings"

	"cute/tui"

	"charm.land/lipgloss/v2"
)

func HelpModal(m tui.Model) *lipgloss.Layer {
	width, height := m.GetSize()
	theme := m.GetTheme()

	helpContent := strings.TrimSpace(buildHelpContent())
	scrollOffset := m.GetHelpScrollOffset()

	// Choose a dialog-sized window, not full-screen.
	modalWidth := width / 2
	if modalWidth > 60 {
		modalWidth = 60
	}
	if modalWidth < 30 {
		modalWidth = 30
	}

	modalHeight := height / 2
	if modalHeight > 16 {
		modalHeight = 16
	}
	if modalHeight < 6 {
		modalHeight = 6
	}

	fw := FloatingWindow{
		Content:      textView(helpContent),
		Width:        modalWidth,
		Height:       modalHeight,
		Style:        DefaultFloatingStyle(theme),
		Title:        "Help",
		ScrollOffset: scrollOffset,
	}

	modalContent := fw.View(width, height)
	return CenterModal(modalContent, width, height)
}

func buildHelpContent() string {
	bindings := tui.GetKeyBindings()

	// Flatten all keybindings into a slice so we can group them dynamically
	allBindings := []tui.Keybinding{
		bindings.AddFile,
		bindings.Cancel,
		bindings.Cd,
		bindings.Parent,
		bindings.Command,
		bindings.Copy,
		bindings.Directories,
		bindings.Down,
		bindings.Enter,
		bindings.Files,
		bindings.Filter,
		bindings.GoToStart,
		bindings.GoToEnd,
		bindings.Help,
		bindings.HiddenFiles,
		bindings.List,
		bindings.Mkdir,
		bindings.Move,
		bindings.Paste,
		bindings.Preview,
		bindings.Quit,
		bindings.Redo,
		bindings.Rename,
		bindings.Select,
		bindings.AutoComplete,
		bindings.Undo,
		bindings.Up,
	}

	// Group bindings by category name
	groups := map[string][]tui.Keybinding{}
	for _, kb := range allBindings {
		categoryName := string(kb.Category)
		groups[categoryName] = append(groups[categoryName], kb)
	}

	var b strings.Builder

	// Category order
	categories := []tui.KeybindingCategoryField{
		tui.KeybindingCategories.General,
		tui.KeybindingCategories.Navigation,
		tui.KeybindingCategories.Filter,
		tui.KeybindingCategories.Editing,
		tui.KeybindingCategories.Help,
		tui.KeybindingCategories.Views,
		tui.KeybindingCategories.Command,
	}

	for _, cat := range categories {
		name := string(cat.Name)
		kbs := groups[name]
		if len(kbs) == 0 {
			continue
		}

		// Category title
		b.WriteString(name + "\n")

		// Keybindings within the category
		for _, kb := range kbs {
			keys := strings.Join(kb.On, " / ")
			line := fmt.Sprintf("  %-18s %s\n", keys, kb.Description)
			b.WriteString(line)
		}

		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}
