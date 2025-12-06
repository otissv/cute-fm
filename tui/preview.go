package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/bimg"

	"cute/console"
	"cute/filesystem"
)

const (
	// imagePreviewDebounce controls how long we wait with the cursor on an image
	// before actually invoking kitty icat. This keeps the UI snappy when moving
	// quickly through many images.
	imagePreviewDebounce = 200 * time.Millisecond

	// maxImagePreviewBytes is a soft limit on file size for which we attempt an
	// image preview. Very large images can make the terminal sluggish or time
	// out the Kitty graphics protocol, so we skip previews above this size.
	maxImagePreviewBytes int64 = 20 * 1024 * 1024 // ~20 MiB

	// maxThumbnailWidth and maxThumbnailHeight control the maximum pixel
	// dimensions for generated thumbnails. Images are resized to fit within
	// this box while preserving aspect ratio. This ensures the preview image
	// is at most ~300px wide on screen.
	maxThumbnailWidth  = 300
	maxThumbnailHeight = 300
)

// UpdatePreview recomputes the right-hand preview panel based on the currently
// selected file. It handles text files (via bat when available), directories,
// and image files (with special handling for Kitty).
func (m *Model) UpdatePreview() {
	// If there are no files, clear the preview.
	if len(m.files) == 0 {
		// If there is a pending image preview timer, cancel it.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}

		// If an image preview was active, clear it from the terminal first so we
		// don't leave a stale image behind.
		if m.imagePreviewActive {
			m.clearImagePreview()
			m.imagePreviewActive = false
		}

		m.rightViewport.SetContent("")
		m.lastPreviewedPath = ""
		return
	}

	idx := m.fileList.Index()
	if idx < 0 || idx >= len(m.files) {
		// If there is a pending image preview timer, cancel it.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}

		// If an image preview was active, clear it from the terminal first.
		if m.imagePreviewActive {
			m.clearImagePreview()
			m.imagePreviewActive = false
		}

		m.rightViewport.SetContent("")
		m.lastPreviewedPath = ""
		return
	}

	fi := m.files[idx]
	path := fi.Path
	if path == "" {
		path = filepath.Join(m.currentDir, fi.Name)
	}

	// When previews are disabled, always show simple file info/properties in
	// the right-hand panel instead of rich text/image previews. This also
	// avoids calling out to external tools like kitty icat or bat.
	if !m.previewEnabled {
		// Cancel any pending image preview timers and clear any active image.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}
		if m.imagePreviewActive {
			m.clearImagePreview()
			m.imagePreviewActive = false
		}

		m.rightViewport.SetContent(renderFileInfoPanel(fi))
		m.lastPreviewedPath = path
		return
	}

	// When the selected file changes, cancel any pending image preview and hide
	// any previously rendered image preview so we don't show a stale image
	// while the new preview is loading.
	if path != m.lastPreviewedPath {
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}
		if m.imagePreviewActive {
			m.clearImagePreview()
			m.imagePreviewActive = false
		}
	}

	switch {
	case fi.IsDir:
		// Cancel any pending image preview when switching to a directory.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}

		m.imagePreviewActive = false
		content := m.previewDirectory(path)
		m.rightViewport.SetContent(content)
	case isImageFile(path):
		// Cancel any previous pending image preview; we'll schedule a new one
		// for this path below.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}

		// Skip previews for very large images to avoid blocking the terminal
		// with slow or timing-out Kitty graphics operations.
		if !canPreviewImage(path) {
			if m.imagePreviewActive {
				m.clearImagePreview()
				m.imagePreviewActive = false
			}
			m.rightViewport.SetContent(
				"Image too large to preview (limit ~20MiB).\nOpen the file directly if you want to view it.",
			)
			break
		}

		// Clear textual content so the preview area appears empty while we wait
		// for the debounced image preview to fire.
		m.rightViewport.SetContent("")

		// Mark the image preview as pending; the actual kitty icat invocation is
		// performed after a short debounce delay, if the selection is still on
		// this image.
		m.imagePreviewActive = false
		console.Log("UpdatePreview: scheduling image preview path=%s term=%s vw=%d vh=%d h=%d", path, m.terminalType, m.viewportWidth, m.viewportHeight, m.height)
		m.scheduleImagePreview(path)
	default:
		// Cancel any pending image preview when switching to a non-image file.
		if m.imagePreviewTimer != nil {
			m.imagePreviewTimer.Stop()
			m.imagePreviewTimer = nil
			m.pendingImagePath = ""
		}

		m.imagePreviewActive = false
		if isTextFile(path) {
			// Use the viewport height as a soft cap for the number of lines.
			maxLines := m.viewportHeight
			if maxLines <= 0 {
				maxLines = 40
			}
			content := renderTextPreview(path, maxLines)
			m.rightViewport.SetContent(content)
		} else {
			m.rightViewport.SetContent("No preview available for this file type.")
		}
	}

	m.lastPreviewedPath = path
}

// previewDirectory renders a directory listing similar to `ls -lh` using the
// same formatting as the main file list.
func (m *Model) previewDirectory(path string) string {
	entries, err := filesystem.ListDirectory(path)
	if err != nil {
		return formatPreviewError("Error reading directory:\n" + err.Error())
	}
	if len(entries) == 0 {
		return "Empty directory."
	}

	// Reuse the file-list delegate so the preview matches list styling.
	delegate := NewFileItemDelegate(m.theme, m.viewportWidth-2)

	var b strings.Builder
	for i, entry := range entries {
		line := delegate.renderFileRow(entry, false, i)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return strings.TrimRight(b.String(), "\n")
}

// isImageFile does a simple extension-based check for common image types.
func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff", ".tif":
		return true
	default:
		return false
	}
}

// isTextFile performs a simple heuristic check by scanning the first few KB
// for NUL bytes. If any are found, we treat the file as binary.
func isTextFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if n == 0 {
		// Empty files are considered text.
		return true
	}
	if err != nil && err != io.EOF {
		// On other errors, be conservative and treat as binary.
		return false
	}

	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return false
		}
	}
	return true
}

// renderTextPreview tries to use `bat` for syntax-highlighted previews, and
// falls back to a simple line-based preview if bat is unavailable.
func renderTextPreview(path string, maxLines int) string {
	if maxLines <= 0 {
		maxLines = 40
	}

	// Prefer bat if available.
	if _, err := exec.LookPath("bat"); err == nil {
		lineRange := fmt.Sprintf("1:%d", maxLines)
		cmd := exec.Command("bat",
			"--color=always",
			"--style=plain",
			"--paging=never",
			"--line-range", lineRange,
			path,
		)
		out, err := cmd.Output()
		if err == nil {
			return string(out)
		}
	}

	// Fallback: read the first maxLines lines directly.
	f, err := os.Open(path)
	if err != nil {
		return formatPreviewError("Error opening file:\n" + err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var b strings.Builder
	lineCount := 0
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		b.WriteByte('\n')
		lineCount++
		if lineCount >= maxLines {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return formatPreviewError("Error reading file:\n" + err.Error())
	}
	return b.String()
}

// formatPreviewError wraps an error message in a simple "modal-like" block
// suitable for rendering inside the right preview viewport.
func formatPreviewError(msg string) string {
	if strings.TrimSpace(msg) == "" {
		return "Error"
	}
	return "Error\n\n" + msg
}

// previewImage renders an image in the right viewport when running under
// Kitty; for other terminals it falls back to a simple message.
func (m *Model) previewImage(path string) {
	if m.terminalType != string(TerminalKitty) {
		m.rightViewport.SetContent("No image preview available in this terminal.\nImage previews are currently supported only in Kitty.")
		return
	}

	// Generate a resized thumbnail of the image using libvips via bimg. This
	// keeps the bytes sent through the Kitty graphics protocol reasonably
	// small, which improves responsiveness when previewing large images.
	thumbPath, err := createThumbnailVips(path, maxThumbnailWidth, maxThumbnailHeight)
	if err != nil {
		m.rightViewport.SetContent(formatPreviewError("Error preparing image for preview:\n" + err.Error()))
		return
	}

	// Mark the image preview as active so the view layer can hide textual
	// content and we know to clear the Kitty graphics layer when changing
	// selection.
	m.imagePreviewActive = true

	// Clear textual content so the image is not obscured by colored cells.
	m.rightViewport.SetContent("")

	// Determine the cell rectangle for the right preview viewport. For Kitty's
	// --place, coordinates are in terminal cells with origin at the top-left
	// of the screen. We approximate the right panel as the right half of the
	// terminal, aligned with the main viewport row.
	widthCells := m.viewportWidth
	if widthCells < 1 {
		widthCells = 1
	}
	heightCells := m.viewportHeight
	if heightCells < 1 {
		heightCells = 1
	}

	// Layout assumptions from calc-layout.go:
	//   - headerRows = 2
	//   - statusRows = 2
	//   - commandRows = 2
	// The combined "viewports" row (file list + preview) starts after:
	//   headerRows + 1 line for the search/tabs row.
	const headerRows = 2
	const searchAndTabsRows = 1

	left := m.viewportWidth
	top := headerRows + searchAndTabsRows

	place := fmt.Sprintf("%dx%d@%dx%d", widthCells, heightCells, left, top)

	// Run kitty icat as a background command, attached to our TTY, placing the
	// image inside the right preview viewport rectangle. If icat fails, show
	// the error text inside the right viewport instead of printing to the
	// main terminal.
	go func(imagePath, placeArg string, model *Model) {
		// Ensure the temporary thumbnail is cleaned up when we're done.
		defer os.Remove(imagePath)

		console.Log("previewImage: starting kitty icat path=%s place=%s", imagePath, placeArg)
		var stderr bytes.Buffer
		cmd := exec.Command("kitty", "+kitten", "icat",
			"--silent",
			"--stdin=no",
			"--place", placeArg,
			imagePath,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = err.Error()
			}
			console.Log("previewImage: kitty icat error: %s", errMsg)

			// If the terminal reports that it does not support the Kitty
			// graphics protocol (or repeatedly times out), disable image
			// previews for the rest of this session so the TUI remains
			// responsive and the layout stable.
			if strings.Contains(errMsg, "does not support the graphics protocol") ||
				strings.Contains(errMsg, "i/o timeout") {
				model.terminalType = string(TerminalUnknown)
				model.imagePreviewActive = false
				model.rightViewport.SetContent(
					"Image preview disabled because this terminal does not support\n" +
						"the Kitty graphics protocol or is too slow to respond.\n\n" +
						"Run cute-fm directly in Kitty/WezTerm/Konsole (without tmux)\n" +
						"to enable graphical image previews.",
				)
				return
			}

			model.rightViewport.SetContent(formatPreviewError("Error rendering image:\n" + errMsg))
		}
	}(thumbPath, place, m)
}

// clearImagePreview clears any previously rendered image from the Kitty
// graphics layer so that moving the file selection or toggling previews off
// hides the old image preview immediately.
func (m *Model) clearImagePreview() {
	go func() {
		console.Log("clearImagePreview: clearing kitty icat (global)")
		var stderr bytes.Buffer
		cmd := exec.Command("kitty", "+kitten", "icat",
			"--silent",
			"--stdin=no",
			"--clear",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = err.Error()
			}
			console.Log("clearImagePreview: kitty icat --clear error: %s", errMsg)
		}
	}()
}

// scheduleImagePreview sets up a debounced image preview for the given path.
// If the cursor moves away from this image before the debounce completes, the
// pending preview is cancelled.
func (m *Model) scheduleImagePreview(path string) {
	m.pendingImagePath = path

	m.imagePreviewTimer = time.AfterFunc(imagePreviewDebounce, func() {
		m.runDebouncedImagePreview(path)
	})
}

// runDebouncedImagePreview is invoked after the debounce delay. It verifies
// that the selection is still on the same image before actually calling
// previewImage.
func (m *Model) runDebouncedImagePreview(path string) {
	// Clear the timer reference so future calls can schedule new previews.
	m.imagePreviewTimer = nil

	// If the pending path has changed, the user moved the selection; skip.
	if m.pendingImagePath != "" && m.pendingImagePath != path {
		return
	}

	// Ensure we still have files and a valid selection.
	if len(m.files) == 0 {
		return
	}

	idx := m.fileList.Index()
	if idx < 0 || idx >= len(m.files) {
		return
	}

	fi := m.files[idx]
	currentPath := fi.Path
	if currentPath == "" {
		currentPath = filepath.Join(m.currentDir, fi.Name)
	}

	// If the selection moved to a different file, don't render this image.
	if currentPath != path {
		return
	}

	// Finally, render the image preview.
	m.previewImage(path)
}

// canPreviewImage returns true if the given file is small enough to be safely
// previewed as an image without likely causing terminal slowdowns.
func canPreviewImage(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return info.Size() <= maxImagePreviewBytes
}

// renderFileInfoPanel renders basic file information and properties for the
// currently selected entry. This is used when previews are disabled so the
// right-hand panel still shows something useful without invoking external
// tools.
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

// createThumbnailVips uses libvips via the bimg library to generate a resized
// thumbnail for the given image. The resulting thumbnail is written to a
// temporary file whose path is returned.
func createThumbnailVips(srcPath string, maxW, maxH int) (string, error) {
	if maxW <= 0 {
		maxW = maxThumbnailWidth
	}
	if maxH <= 0 {
		maxH = maxThumbnailHeight
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", err
	}

	image := bimg.NewImage(data)

	options := bimg.Options{
		Width:   maxW,
		Height:  maxH,
		Enlarge: false,
		Quality: 80,
		Type:    bimg.JPEG,
	}

	thumb, err := image.Process(options)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "cute-fm-thumb-*.jpg")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(thumb); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
