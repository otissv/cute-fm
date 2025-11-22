package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Environment describes the current execution context for a command.
type Environment struct {
	// Cwd is the current working directory of the file manager.
	Cwd string
	// ConfigCommands holds user-defined commands loaded from the config file.
	ConfigCommands map[string]string
}

// Result captures the outcome of executing a command.
type Result struct {
	// Output is any textual output produced by the command (stdout/stderr).
	Output string
	// Cwd, when non-empty, is the new working directory to use.
	Cwd string
	// Refresh indicates that the file list for the current directory should be
	// refreshed after the command completes.
	Refresh bool
	// ViewMode, when non-empty, indicates a logical file-list view mode
	// (e.g. "ll", "ls", ...).
	ViewMode string
	// OpenHelp indicates that the UI should open the help modal.
	OpenHelp bool
	// Quit indicates that the application should exit.
	Quit bool
}

// Execute parses and executes a single command line within the given
// environment and returns the resulting state changes and output.
func Execute(env Environment, input string) (Result, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{}, nil
	}

	// Handle "sh <command>" specially so we preserve the original text,
	// including quoting and shell operators.
	if strings.HasPrefix(input, "sh ") {
		cmdLine := strings.TrimSpace(strings.TrimPrefix(input, "sh"))
		out, err := runShell(env.Cwd, cmdLine)
		return Result{Output: out}, err
	}

	fields := strings.Fields(input)
	if len(fields) == 0 {
		return Result{}, nil
	}

	name := fields[0]
	args := fields[1:]

	switch name {
	case "cd":
		return cmdCd(env, args)
	case "ll", "ls", "ld", "lf":
		return CmdViewModeDescription(name), nil
	case "help":
		return Result{OpenHelp: true}, nil
	case "touch":
		return cmdTouch(env, args)
	case "mkdir":
		return cmdMkdir(env, args, false)
	case "mkcd":
		return cmdMkcd(env, args)
	case "rm":
		return cmdRm(env, args)
	case "mv":
		return cmdMv(env, args)
	case "cp":
		return cmdCp(env, args)
	case "ln":
		return cmdLn(env, args)
	case "quit", "q":
		return Result{Quit: true}, nil
	default:
		// Fall back to user-defined commands from the config file, if any.
		if env.ConfigCommands != nil {
			if val, ok := env.ConfigCommands[name]; ok {
				return executeConfigCommand(env, val)
			}
		}

		// As a last resort, try to execute the input as a shell command.
		out, err := runShell(env.Cwd, input)
		return Result{Output: out}, err
	}
}

// cmdCd implements "cd <directory>" semantics without changing the process-wide
// working directory. Instead, it validates and returns the new directory path.
func cmdCd(env Environment, args []string) (Result, error) {
	if len(args) < 1 {
		return Result{}, fmt.Errorf("cd: missing directory")
	}

	target := expandPath(args[0], env.Cwd)
	info, err := os.Stat(target)
	if err != nil {
		return Result{}, fmt.Errorf("cd: %w", err)
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("cd: not a directory: %s", target)
	}

	return Result{
		Cwd:    target,
		Output: fmt.Sprintf("changed directory to %s", target),
	}, nil
}

// cmdViewMode maps logical "ls"-style commands to a named view mode.
func CmdViewModeDescription(mode string) Result {
	// These strings document the intended eza-style flags; the TUI can store
	// or display them as needed.
	var desc string

	switch mode {
	case "ll":
		desc = "Lists directory (default)"
	case "ls":
		desc = "Hides dotfiles"
	case "ld":
		desc = "Only directories"
	case "lf":
		desc = "-h -a -g -l --only-files --group-directories-first --git --icons"
	}

	return Result{
		ViewMode: mode,
		Output:   fmt.Sprintf("switched view to %s (%s)", mode, desc),
	}
}

func CmdViewModeStatus(mode string) string {
	var status string

	switch mode {
	case "ll":
		status = "All"
	case "ls":
		status = "List"
	case "ld":
		status = "Dirs:ll"
	case "lf":
		status = "Files"
	}

	return status
}

// cmdTouch implements "touch <file> ..." semantics similar to the shell.
func cmdTouch(env Environment, args []string) (Result, error) {
	if len(args) == 0 {
		return Result{}, fmt.Errorf("touch: missing file operand")
	}

	for _, a := range args {
		target := expandPath(a, env.Cwd)
		f, err := os.OpenFile(target, os.O_RDONLY|os.O_CREATE, 0o666)
		if err != nil {
			return Result{}, fmt.Errorf("touch: %w", err)
		}
		_ = f.Close()
	}

	return Result{Output: "touch: updated files", Refresh: true}, nil
}

// cmdMkdir implements "mkdir [-p] <dirs>...".
func cmdMkdir(env Environment, args []string, cdInto bool) (Result, error) {
	if len(args) == 0 {
		return Result{}, fmt.Errorf("mkdir: missing operand")
	}

	var lastPath string
	for _, a := range args {
		target := expandPath(a, env.Cwd)
		if err := os.MkdirAll(target, 0o777); err != nil {
			return Result{}, fmt.Errorf("mkdir: %w", err)
		}
		lastPath = target
	}

	res := Result{Output: "mkdir: created directories", Refresh: true}
	if cdInto && lastPath != "" {
		res.Cwd = lastPath
	}
	return res, nil
}

// cmdMkcd is equivalent to "mkdir -p <dirs>..." followed by "cd" into the last.
func cmdMkcd(env Environment, args []string) (Result, error) {
	return cmdMkdir(env, args, true)
}

// cmdRm implements a simplified "rm" that removes files or directories
// recursively (similar to `rm -r`).
func cmdRm(env Environment, args []string) (Result, error) {
	if len(args) == 0 {
		return Result{}, fmt.Errorf("rm: missing operand")
	}

	for _, a := range args {
		target := expandPath(a, env.Cwd)
		if err := os.RemoveAll(target); err != nil {
			return Result{}, fmt.Errorf("rm: %w", err)
		}
	}

	return Result{Output: "rm: removed", Refresh: true}, nil
}

// cmdMv implements "mv <source> <destination>" semantics.
func cmdMv(env Environment, args []string) (Result, error) {
	if len(args) < 2 {
		return Result{}, fmt.Errorf("mv: missing operand")
	}

	src := expandPath(args[0], env.Cwd)
	dst := expandPath(args[1], env.Cwd)

	// If destination is an existing directory, move into it.
	if info, err := os.Stat(dst); err == nil && info.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err := os.Rename(src, dst); err != nil {
		return Result{}, fmt.Errorf("mv: %w", err)
	}

	return Result{Output: fmt.Sprintf("mv: %s -> %s", src, dst), Refresh: true}, nil
}

// cmdCp implements a basic "cp <source> <destination>" for regular files.
func cmdCp(env Environment, args []string) (Result, error) {
	if len(args) < 2 {
		return Result{}, fmt.Errorf("cp: missing operand")
	}

	src := expandPath(args[0], env.Cwd)
	dst := expandPath(args[1], env.Cwd)

	info, err := os.Stat(src)
	if err != nil {
		return Result{}, fmt.Errorf("cp: %w", err)
	}
	if info.IsDir() {
		return Result{}, fmt.Errorf("cp: copying directories is not supported yet")
	}

	// If destination is an existing directory, copy into it.
	if dInfo, err := os.Stat(dst); err == nil && dInfo.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	in, err := os.Open(src)
	if err != nil {
		return Result{}, fmt.Errorf("cp: %w", err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return Result{}, fmt.Errorf("cp: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return Result{}, fmt.Errorf("cp: %w", err)
	}

	return Result{Output: fmt.Sprintf("cp: %s -> %s", src, dst), Refresh: true}, nil
}

// cmdLn implements a simple "ln <source> <destination>" using hard links.
func cmdLn(env Environment, args []string) (Result, error) {
	if len(args) < 2 {
		return Result{}, fmt.Errorf("ln: missing operand")
	}

	src := expandPath(args[0], env.Cwd)
	dst := expandPath(args[1], env.Cwd)

	// If destination is an existing directory, create the link inside it.
	if info, err := os.Stat(dst); err == nil && info.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err := os.Link(src, dst); err != nil {
		return Result{}, fmt.Errorf("ln: %w", err)
	}

	return Result{Output: fmt.Sprintf("ln: %s -> %s", src, dst), Refresh: true}, nil
}

// executeConfigCommand runs a command defined in the config file. If the value
// resolves to an existing directory, it behaves like "cd" to that directory;
// otherwise it is executed as a shell command.
func executeConfigCommand(env Environment, value string) (Result, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Result{}, nil
	}

	target := expandPath(value, env.Cwd)
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		return Result{
			Cwd:    target,
			Output: fmt.Sprintf("changed directory to %s", target),
		}, nil
	}

	out, err := runShell(env.Cwd, value)
	return Result{Output: out}, err
}

// expandPath resolves "~" and relative paths against the provided cwd.
func expandPath(path, cwd string) string {
	if path == "" {
		return cwd
	}

	// Handle leading "~" using the current user's home directory.
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			if path == "~" {
				return home
			}
			if strings.HasPrefix(path, "~/") {
				return filepath.Join(home, path[2:])
			}
		}
	}

	if filepath.IsAbs(path) {
		return path
	}
	if cwd == "" {
		return path
	}
	return filepath.Join(cwd, path)
}

// runShell executes the given shell command using "bash -lc" in the specified
// directory and returns its combined standard output and standard error.
func runShell(dir, cmdLine string) (string, error) {
	if cmdLine == "" {
		return "", nil
	}

	cmd := exec.Command("bash", "-lc", cmdLine)
	if dir != "" {
		cmd.Dir = dir
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	return buf.String(), err
}
