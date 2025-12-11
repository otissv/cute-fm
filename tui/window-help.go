package tui

import (
	"fmt"
	"reflect"
	"strings"

	"charm.land/lipgloss/v2"
)

func HelpWindow(m Model) *lipgloss.Layer {
	width, height := m.GetSize()
	theme := m.GetTheme()

	helpContent := strings.TrimSpace(buildHelpContent())
	scrollOffset := m.GetHelpScrollOffset()

	windowWidth := 90
	windowHeight := height / 2

	fw := FloatingWindow{
		Content:      ViewPrimitive(helpContent),
		Width:        windowWidth,
		Height:       windowHeight,
		Style:        DefaultFloatingStyle(theme),
		Title:        "Help",
		ScrollOffset: scrollOffset,
	}

	windowContent := fw.View(width, height)
	return CenterWindow(windowContent, width, height)
}

func buildHelpContent() string {
	bindings := GetKeyBindings()

	v := reflect.ValueOf(bindings)

	allBindings := make([]Keybinding, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		kb, ok := v.Field(i).Interface().(Keybinding)
		if ok {
			allBindings[i] = kb
		}
	}

	groups := map[string][]Keybinding{}
	for _, kb := range allBindings {
		categoryName := string(kb.Category)
		groups[categoryName] = append(groups[categoryName], kb)
	}

	var b strings.Builder

	categories := []KeybindingCategoryField{
		KeybindingCategories.General,
		KeybindingCategories.Navigation,
		KeybindingCategories.Filter,
		KeybindingCategories.Editing,
		KeybindingCategories.Help,
		KeybindingCategories.Views,
		KeybindingCategories.Command,
	}

	for _, cat := range categories {
		name := string(cat.Name)
		kbs := groups[name]
		if len(kbs) == 0 {
			continue
		}

		b.WriteString(name + "\n")

		for _, kb := range kbs {
			keys := strings.Join(kb.On, " / ")
			line := fmt.Sprintf("  %-18s %s\n", keys, kb.Description)
			b.WriteString(line)
		}

		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}
