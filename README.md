## cute-fm

`cute-fm` is a terminal file manager built with Bubble Tea and Lip Gloss, featuring:

- A two‑pane layout (file list + preview)
- Text previews (with `bat` when available)
- Image previews using Kitty graphics (with debounced, libvips‑powered thumbnails)
- Lua‑based configuration for themes and commands

---

## Prerequisites

- **Go toolchain**
  - **Version**: Go **1.24+** (matches the `go` directive in `go.mod`)
  - Install from the official Go downloads or your distro’s packages.

- **libvips (for high‑performance image thumbnails)**
  - `cute-fm` uses [`github.com/h2non/bimg`](https://github.com/h2non/bimg) which requires **libvips** and `pkg-config`.
  - Install using your package manager:
    - Debian/Ubuntu: `sudo apt install libvips-dev pkg-config`
    - Fedora: `sudo dnf install vips-devel pkg-config`
    - Arch/Manjaro: `sudo pacman -S libvips pkgconf`
    - macOS (Homebrew): `brew install vips pkg-config`

- **Terminal with Kitty graphics protocol support (for image preview)**
  - Recommended: **Kitty** (`TERM=xterm-kitty` or `KITTY_WINDOW_ID` set).
  - Also supported: terminals that implement the Kitty graphics protocol (e.g. WezTerm, Konsole via compatibility layers), but this project is primarily tested with Kitty.
  - If your terminal does **not** support the graphics protocol, image previews will be skipped or fall back to text messages.

- **Optional but recommended**
  - [`bat`](https://github.com/sharkdp/bat) for syntax‑highlighted text previews.
  - A Nerd Font or other powerline‑friendly font for nicer glyphs.

---

## Building

From the project root:

```bash
go build -o cute-fm ./...
```

This will produce a `cute-fm` (or `lsfm`, depending on your build target) binary in the current directory.

You can also install it into your `$GOBIN`:

```bash
go install ./...
```

Ensure `$GOBIN` (or `$GOPATH/bin`) is on your `PATH` so you can run the binary directly.

---

## Running

Run the TUI from any directory:

```bash
cute-fm            # start in current directory
cute-fm /path/to/dir
```

Inside the TUI:

- Use **`j`/`k`** or **arrow keys** to move the cursor in the file list.
- Use `:` to open the command bar, `?` for help, `f` to filter, etc. (see the built‑in help for full keybindings).

Image previews will appear on the right when:

- You are running inside a Kitty‑compatible terminal, and
- The selected file is an image within the configured size limit.

---

## Configuration (Lua)

`cute-fm` loads a Lua configuration file to customize **themes** and **commands**.

Search order for `config.lua`:

1. `$XDG_CONFIG_HOME/cute/config.lua` (or `~/.config/cute/config.lua` on most systems)
2. `<binary-dir>/config/config.lua`
3. `./config/config.lua` (useful during development)

The Lua file can define:

- A `theme` table with simple color overrides.
- A `commands` table mapping command names to Lua functions that receive context about the selected file and can return output, new working directory, view mode, etc.

See `config/config.lua` in the repo as a starting point.

---

## Notes on performance

- Image previews are **debounced**: the image is only rendered after the cursor rests on a file briefly, which keeps navigation smooth.
- Thumbnails are generated using **libvips** via `bimg`, downscaling large images before sending them to the terminal, which significantly reduces lag and timeouts in the Kitty graphics protocol.


