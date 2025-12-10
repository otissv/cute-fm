package tui

import (
	"os"
	"strings"
)

type TerminalType string

const (
	TerminalUnknown TerminalType = ""
	TerminalKitty   TerminalType = "kitty"
)

func detectTerminalType() TerminalType {
	// Kitty sets KITTY_WINDOW_ID and usually TERM=xterm-kitty.
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return TerminalKitty
	}
	if strings.Contains(strings.ToLower(os.Getenv("TERM")), "kitty") {
		return TerminalKitty
	}
	return TerminalUnknown
}
