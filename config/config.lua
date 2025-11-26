-- Lua configuration for the Cute file manager.
--
-- This file replaces the old cute.toml configuration. It defines both:
--   - the UI theme (via the global `theme` table), and
--   - user-defined commands (via the global `commands` table).
--
-- The Go side loads this file once at startup. Commands are invoked from the
-- command bar with access to the currently selected file or directory.

-- Theme -----------------------------------------------------------------------
--
-- Keys here correspond to the simple key/value overrides understood by the
-- Go theming layer. Any key you omit falls back to the built-in defaults.
--
-- Common keys:
--   foreground, background, border
--   selected_foreground, selected_background
--   directory, regular, symlink, socket, pipe, device, executable
--   nlink, user, group, size, time

theme = {
  foreground = "#F0EDED",
  background = "#1E1E1E",
  border = "#F25D94",

  -- Selection colors
  selected_foreground = "#1E1E1E",
  selected_background = "#3B3B3B",

  -- File type colors
  directory  = "#A8D2FF",
  regular    = "#F0EDED",
  executable = "#FF9BC0",
}

-- Commands --------------------------------------------------------------------
--
-- Each command is a function of the form:
--
--   function(ctx, args)
--     -- ctx describes the current selection:
--     --   ctx.cwd    : current working directory
--     --   ctx.path   : full path of the selected item (if any)
--     --   ctx.name   : base name of the selected item (if any)
--     --   ctx.is_dir : boolean
--     --   ctx.type   : file type string ("directory", "regular", ...)
--     --
--     -- args is an array-like table of extra CLI args (1-based).
--     --
--     -- The function should return either:
--     --   * a table:
--     --       { output = "...", cwd = "...", refresh = true,
--     --         view_mode = "ll", open_help = false, quit = false }
--     --   * or a string, which is treated as `output`.
--   end
--
-- Example usage from the command bar:
--   :open
--   :edit some-file.txt

commands = {}

-- Open the selected file or directory using the system default handler.
function commands.open(ctx, args)
  local target = ctx.path or ctx.cwd
  if not target or target == "" then
    return "no selection to open"
  end

  -- On Linux this uses xdg-open; adjust if needed for other platforms.
  os.execute(string.format("xdg-open %q &", target))
  return { refresh = false }
end

-- Edit the selected file (or a provided path) with $EDITOR (falling back to nvim).
function commands.edit(ctx, args)
  local target = ctx.path
  if #args >= 1 then
    -- Allow overriding the target via an explicit argument.
    target = args[1]
  end

  if not target or target == "" then
    return "no file selected to edit"
  end

  local editor = os.getenv("EDITOR") or "nvim"
  os.execute(string.format("%s %q", editor, target))
  return { refresh = false }
end


