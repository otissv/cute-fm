package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

func Blends(colo1 string, color2 string) []color.Color {
	return gamut.Blends(lipgloss.Color(colo1), lipgloss.Color(color2), 50)
}

func RainbowText(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

func Rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str += base.
			Background(lipgloss.Color(color.Hex())). // use blend as background
			Foreground(lipgloss.Color("#D4D4D4")).   // fixed text color
			Render(string(ss))
	}
	return str
}
