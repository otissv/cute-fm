package components

import (
	"fmt"
	"reflect"
	"strings"

	"cute/tui"

	"charm.land/lipgloss/v2"
)

func HelpModal(m tui.Model) *lipgloss.Layer {
	width, height := m.GetSize()
	theme := m.GetTheme()

	helpContent := strings.TrimSpace(buildHelpContent())
	scrollOffset := m.GetHelpScrollOffset()

	modalWidth := 90
	modalHeight := height / 2

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

	v := reflect.ValueOf(bindings)

	allBindings := make([]tui.Keybinding, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		kb, ok := v.Field(i).Interface().(tui.Keybinding)
		if ok {
			allBindings[i] = kb
		}
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
