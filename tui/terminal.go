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

// detectTerminalType inspects environment variables to determine the
// current terminal. For now we only care about Kitty, but this can be
// extended in the future.
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
