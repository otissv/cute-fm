package tui

import (
	"strings"
)

// filterHistoryMatches filters the command history based on the current input
// and returns matching commands. Matches are case-insensitive prefix matches.
func (m *Model) filterHistoryMatches(input string) []string {
	if input == "" {
		return []string{}
	}

	inputLower := strings.ToLower(input)
	var matches []string

	// Search through history in reverse order (most recent first)
	for i := len(m.commandHistory) - 1; i >= 0; i-- {
		cmd := m.commandHistory[i]
		if strings.HasPrefix(strings.ToLower(cmd), inputLower) {
			// Avoid duplicates
			isDuplicate := false
			for _, existing := range matches {
				if existing == cmd {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				matches = append(matches, cmd)
			}
		}
	}

	return matches
}

// updateHistoryMatches updates the history matches based on current input
// and resets the history index.
func (m *Model) updateHistoryMatches() {
	currentInput := strings.TrimSpace(m.commandInput.Value())
	m.historyMatches = m.filterHistoryMatches(currentInput)
	m.historyIndex = -1
}

// completeCommand completes the command input with the next match from history.
// If historyIndex is -1, it uses the first match. Otherwise, it cycles through matches.
func (m *Model) completeCommand() {
	if len(m.historyMatches) == 0 {
		return
	}

	// Cycle through matches
	m.historyIndex++
	if m.historyIndex >= len(m.historyMatches) {
		m.historyIndex = 0
	}

	// Set the input to the matched command
	m.commandInput.SetValue(m.historyMatches[m.historyIndex])
	// Move cursor to end
	m.commandInput.CursorEnd()
}

// navigateHistory moves through history matches using Up/Down arrows.
// delta is -1 for up (older), +1 for down (newer).
// If there are no matches based on current input, it navigates through all history.
func (m *Model) navigateHistory(delta int) {
	currentInput := strings.TrimSpace(m.commandInput.Value())

	// If there's no input or no matches, navigate through all history
	if currentInput == "" || len(m.historyMatches) == 0 {
		if len(m.commandHistory) == 0 {
			return
		}

		// Use all history as matches
		if m.historyIndex < 0 {
			// Start from the most recent (last item)
			m.historyIndex = len(m.commandHistory) - 1
		} else {
			m.historyIndex += delta
		}

		// Wrap around
		if m.historyIndex < 0 {
			m.historyIndex = len(m.commandHistory) - 1
		} else if m.historyIndex >= len(m.commandHistory) {
			m.historyIndex = 0
		}

		// Set the input to the history command
		m.commandInput.SetValue(m.commandHistory[m.historyIndex])
		m.commandInput.CursorEnd()
		return
	}

	// Navigate through filtered matches
	m.historyIndex += delta

	// Wrap around
	if m.historyIndex < 0 {
		m.historyIndex = len(m.historyMatches) - 1
	} else if m.historyIndex >= len(m.historyMatches) {
		m.historyIndex = 0
	}

	// Set the input to the matched command
	m.commandInput.SetValue(m.historyMatches[m.historyIndex])
	// Move cursor to end
	m.commandInput.CursorEnd()
}
