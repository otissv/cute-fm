package tui

import (
	"fmt"
	"strings"

	"cute/filesystem"
)

func (m *Model) UpdateFileInfoPane() {
	pane := m.GetActivePane()

	if len(pane.files) == 0 {
		m.fileInfoViewport.SetContent("")
		return
	}

	idx := pane.fileList.Index()
	if idx < 0 || idx >= len(pane.files) {
		m.fileInfoViewport.SetContent("")
		return
	}

	fi := pane.files[idx]

	m.fileInfoViewport.SetContent(renderFileInfoPane(fi))
}

func renderFileInfoPane(fi filesystem.FileInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, "File info\n\n")
	fmt.Fprintf(&b, "Name: %s\n", fi.Name)

	if fi.Path != "" {
		fmt.Fprintf(&b, "Path: %s\n", fi.Path)
	}
	fmt.Fprintf(&b, "Type: %s\n", fi.Type)
	if fi.MimeType != "" {
		fmt.Fprintf(&b, "MIME: %s\n", fi.MimeType)
	}
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
