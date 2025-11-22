package config

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// LoadCommands parses the [Commands] section from the given configuration file
// path and returns a map of command name -> value.
//
// The parser understands a very small TOML-like subset:
//
//   - Section headers in the form [SectionName]
//   - Key/value lines: key = "value"
//   - Lines starting with '#' and blank lines are ignored
//
// Only entries that appear under the [Commands] section are returned.
func LoadCommands(path string) map[string]string {
	data, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}
	}

	result := map[string]string{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Section header, e.g. [Commands]
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		// We only care about key/value pairs in the [Commands] section.
		if currentSection != "Commands" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Strip surrounding quotes if present.
		val = strings.Trim(val, `"`)

		if key == "" {
			continue
		}

		result[key] = val
	}

	return result
}


