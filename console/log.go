package console

import (
	"fmt"
	"os"
	"time"
)

// FilePath is the path to the debug log file. Override this if you want
// the log to go somewhere else.
var FilePath = "debug.log"

// Log appends a timestamped debug message to the debug log file.
// The message supports fmt-style formatting.
func Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), msg)

	f, err := os.OpenFile(FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// If we can't open the debug log, fail silently so we don't
		// break normal application behavior.
		return
	}
	defer f.Close()

	_, _ = f.WriteString(line)
}
