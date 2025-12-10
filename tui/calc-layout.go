package tui

// CalcLayout recalculates the viewport dimensions based on the current
// window size and whether command mode is active, then updates the file
// list and right viewport dimensions.
func (m *Model) CalcLayout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	const (
		headerRow = 4
		statusRow = 3
	)

	// Viewport style height: remaining height after the fixed rows.
	viewportHeight := m.height - (headerRow + statusRow)
	if viewportHeight < 3 {
		viewportHeight = 3 // Minimum: 1 content + 2 borders
	}

	m.viewportHeight = viewportHeight

	// Viewport content height (scrollable area): style height - 2 border lines.
	viewportContentHeight := viewportHeight - 2
	if viewportContentHeight < 1 {
		viewportContentHeight = 1 // Minimum content height
	}

	// Calculate viewport width. When the right pane is hidden, the left
	// viewport should take the full terminal width; otherwise, split evenly.
	if m.showRightPane {
		m.viewportWidth = m.width / 2
	} else {
		m.viewportWidth = m.width
	}

	// Content width for the list (subtract borders).
	listContentWidth := m.viewportWidth - 2
	if listContentWidth < 1 {
		listContentWidth = 1
	}

	// Update the file list dimensions for both panes.
	m.leftPane.fileList.SetSize(listContentWidth, viewportContentHeight)
	m.rightPane.fileList.SetSize(listContentWidth, viewportContentHeight)

	// Update the delegates with the new width for proper row padding.
	m.UpdateFileListDelegate(listContentWidth)

	// Update right viewport dimensions (height is the content height).
	m.fileInfoViewport.SetWidth(m.viewportWidth)
	m.fileInfoViewport.SetHeight(viewportContentHeight)
}
