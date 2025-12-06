package tui

// Init initializes the model (required by Bubble Tea)

// CalcLayout recalculates the viewport dimensions based on the current
// window size and whether command mode is active, then updates the file
// list and right viewport dimensions.
func (m *Model) CalcLayout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	// Approximate fixed heights for the non-viewport rows. With the current
	// layout we always render, in order:
	//   - header row
	//   - a combined row containing the search bar (left) and preview tabs (right)
	//   - the main viewports row (file list + preview)
	//   - status bar
	//   - command bar
	//
	// Only the main viewports row should grow/shrink with the terminal height.
	const (
		headerRows  = 2 // Header typically renders on 2 lines (title + padding).
		statusRows  = 2 // Status bar row.
		commandRows = 2 // Command bar row at the bottom.
	)

	// Viewport style height: remaining height after the fixed rows.
	viewportHeight := m.height - (headerRows + statusRows + commandRows)
	if viewportHeight < 3 {
		viewportHeight = 3 // Minimum: 1 content + 2 borders
	}
	// Persist the total viewport box height so Lip Gloss containers (FileList,
	// Preview) can render with a fixed height instead of expanding to fit
	// their content.
	m.viewportHeight = viewportHeight

	// Viewport content height (scrollable area): style height - 2 border lines.
	viewportContentHeight := viewportHeight - 2
	if viewportContentHeight < 1 {
		viewportContentHeight = 1 // Minimum content height
	}

	// Calculate viewport width. When the right panel is hidden, the left
	// viewport should take the full terminal width; otherwise, split evenly.
	if m.showRightPanel {
		m.viewportWidth = m.width / 2
	} else {
		m.viewportWidth = m.width
	}

	// Content width for the list (subtract borders).
	listContentWidth := m.viewportWidth - 2
	if listContentWidth < 1 {
		listContentWidth = 1
	}

	// Update the file list dimensions.
	m.fileList.SetSize(listContentWidth, viewportContentHeight)

	// Update the delegate with the new width for proper row padding.
	m.UpdateFileListDelegate(listContentWidth)

	// Update right viewport dimensions (height is the content height).
	m.fileInfoViewport.SetWidth(m.viewportWidth)
	m.fileInfoViewport.SetHeight(viewportContentHeight)
}
