package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"cute/filesystem"
)

// UpdateFileInfoPanel recomputes the right-hand preview panel based on the currently
// selected file. It handles text files (via bat when available), directories,
// and image files (with special handling for Kitty).
func (m *Model) UpdateFileInfoPanel() {
	// If there are no files, clear the preview.
	if len(m.files) == 0 {
		m.fileInfoViewport.SetContent("")
		m.lastPreviewedPath = ""
		return
	}

	idx := m.fileList.Index()
	if idx < 0 || idx >= len(m.files) {
		m.fileInfoViewport.SetContent("")
		m.lastPreviewedPath = ""
		return
	}

	fi := m.files[idx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(m.currentDir, fi.Name)
	}

	// Always show simple file info/properties in the right-hand panel instead
	// of rich text/image previews.
	m.fileInfoViewport.SetContent(renderFileInfoPanel(fi))
	m.lastPreviewedPath = path
}

// renderFileInfoPanel renders basic file information and properties for the
// currently selected entry. This is now the only content shown in the
// right-hand panel (rich previews have been removed).
func renderFileInfoPanel(fi filesystem.FileInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, "File info\n\n")
	fmt.Fprintf(&b, "Name: %s\n", fi.Name)
	if fi.Path != "" {
		fmt.Fprintf(&b, "Path: %s\n", fi.Path)
	}
	fmt.Fprintf(&b, "Type: %s\n", fi.Type)

	fmt.Fprintf(&b, "Size: %s\n", fi.Size)
	fmt.Fprintf(&b, "Owner: %s\n", fi.User)
	fmt.Fprintf(&b, "Group: %s\n", fi.Group)
	fmt.Fprintf(&b, "Modified: %s\n", fi.DateModified)

	fmt.Fprintf(&b, "Permissions: %s\n", fi.Permissions)

	// Decode the permission bits (positions 2–10) into Owner / Group / Other
	// sections with explicit Read/Write/Execute or Traverse flags.
	perm := fi.Permissions
	if len(perm) >= 10 {
		owner := perm[1:4]
		group := perm[4:7]
		other := perm[7:10]

		// Helper to print a permission group with a heading.
		printGroup := func(heading string, bits string) {
			if len(bits) != 3 {
				return
			}

			fmt.Fprintf(&b, "\n%s\n", heading)

			// Read
			if bits[0] == 'r' {
				fmt.Fprintf(&b, "Read: True\n")
			} else {
				fmt.Fprintf(&b, "Read: -\n")
			}

			// Write
			if bits[1] == 'w' {
				fmt.Fprintf(&b, "Write: True\n")
			} else {
				fmt.Fprintf(&b, "Write: -\n")
			}

			// Execute / Traverse
			label := "Execute"
			if fi.IsDir {
				label = "Traverse"
			}
			if bits[2] == 'x' {
				fmt.Fprintf(&b, "%s: True\n", label)
			} else {
				fmt.Fprintf(&b, "%s: -\n", label)
			}
		}

		printGroup("Owner", owner)
		printGroup("Group", group)
		printGroup("Other", other)
	}

	return b.String()
}
