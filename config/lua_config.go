package config

import (
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"

	"cute/theming"
)

// RuntimeConfig holds the theme and user-defined commands loaded from a Lua
// configuration file.
//
// The Lua file is expected to define:
//
//	theme = {
//	  -- keys matching the simple key/value theme overrides understood by
//	  -- theming.LoadThemeFromMap, e.g.:
//	  --   foreground = "#F0EDED",
//	  --   background = "#1E1E1E",
//	  --   border     = "#F25D94",
//	  --   directory  = "#A8D2FF",
//	  --   regular    = "#F0EDED",
//	}
//
//	commands = {
//	  mycmd = function(ctx, args)
//	    -- ctx describes the selected file/dir:
//	    --   ctx.cwd      : current working directory
//	    --   ctx.path     : full path to the selected entry (if any)
//	    --   ctx.name     : base name of the selected entry (if any)
//	    --   ctx.is_dir   : boolean
//	    --   ctx.type     : file type string (\"directory\", \"regular\", ...)
//	    --
//	    -- args is an array-like table of additional CLI arguments.
//	    --
//	    -- The function must return either:
//	    --   * a table describing the command result:
//	    --       { output = \"...\", cwd = \"/new/dir\", refresh = true,
//	    --         view_mode = \"ll\", open_help = false, quit = false }
//	    --     (all fields are optional), or
//	    --   * a string, which is treated as Output.
//	  end,
//	}
//
// Command invocation and result decoding are handled by the command package;
// RuntimeConfig only stores the Lua state and registered functions.
type RuntimeConfig struct {
	// L is the Lua state backing this configuration.
	L *lua.LState

	// Theme is the fully-resolved TUI theme, produced from the Lua "theme"
	// table (when present) layered over theming.DefaultTheme().
	Theme theming.Theme

	// commands maps command names (as typed in the command bar) to the
	// corresponding Lua function objects.
	commands map[string]*lua.LFunction
}

// Command looks up a user-defined command function by name.
func (rc *RuntimeConfig) Command(name string) *lua.LFunction {
	if rc == nil {
		return nil
	}
	if rc.commands == nil {
		return nil
	}
	return rc.commands[name]
}

// LoadRuntimeConfig discovers and loads the Lua configuration for the given
// config directory. It never falls back to TOML; if no Lua file can be found
// or loaded, a RuntimeConfig with the default theme and no commands is
// returned.
//
// The search order for the Lua file is:
//
//  1. <configDir>/config.lua
//  2. <binaryDir>/config/config.lua
//  3. ./config/config.lua  (useful during development)
func LoadRuntimeConfig(configDir string) *RuntimeConfig {
	theme := theming.DefaultTheme()
	rc := &RuntimeConfig{
		Theme:    theme,
		commands: map[string]*lua.LFunction{},
	}

	path := findLuaConfigPath(configDir)
	if path == "" {
		return rc
	}

	L := lua.NewState()

	if err := L.DoFile(path); err != nil {
		// If the Lua file fails to load, fall back to the default theme and
		// no commands, but still close the Lua state.
		L.Close()
		return rc
	}

	// Extract theme overrides from global "theme" table, if present.
	overrides := map[string]string{}
	if v := L.GetGlobal("theme"); v.Type() == lua.LTTable {
		tbl := v.(*lua.LTable)
		tbl.ForEach(func(k, v lua.LValue) {
			ks, ok1 := k.(lua.LString)
			vs, ok2 := v.(lua.LString)
			if !ok1 || !ok2 {
				return
			}
			overrides[string(ks)] = string(vs)
		})
	}

	if len(overrides) > 0 {
		rc.Theme = theming.LoadThemeFromMap(overrides)
	}

	// Extract user-defined commands from global "commands" table, if present.
	if v := L.GetGlobal("commands"); v.Type() == lua.LTTable {
		tbl := v.(*lua.LTable)
		commands := map[string]*lua.LFunction{}
		tbl.ForEach(func(k, v lua.LValue) {
			name, ok := k.(lua.LString)
			if !ok {
				return
			}
			fn, ok := v.(*lua.LFunction)
			if !ok {
				return
			}
			commands[string(name)] = fn
		})
		rc.commands = commands
	}

	rc.L = L
	return rc
}

// findLuaConfigPath returns the first existing Lua config file path following
// the search order described in LoadRuntimeConfig.
func findLuaConfigPath(configDir string) string {
	candidates := []string{}

	if configDir != "" {
		candidates = append(candidates, filepath.Join(configDir, "config.lua"))
	}

	if exe, err := os.Executable(); err == nil && exe != "" {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates, filepath.Join(exeDir, "config", "config.lua"))
	}

	// Development-friendly relative path.
	candidates = append(candidates, filepath.Join("config", "config.lua"))

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}
