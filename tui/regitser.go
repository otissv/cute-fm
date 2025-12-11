package tui

import "charm.land/lipgloss/v2"

type Windows struct {
	Column   func(m Model, args ColumnWindowArgs) *lipgloss.Layer
	Command  func(m Model, args CommandWindowArgs) *lipgloss.Layer
	Dialog   func(m Model, args DialogWindowArgs) *lipgloss.Layer
	Help     func(m Model) *lipgloss.Layer
	Settings func(m Model) *lipgloss.Layer
}

type Components struct {
	CurrentDir   func(m Model, args CurrentDirComponentArgs) string
	FileListView func(m Model, args FileListComponentArgs) string
	FileInfo     func(m Model, args ComponentArgs) string
	Header       func(m Model, args ComponentArgs) string
	PreviewTabs  func(m Model, args ComponentArgs) string
	SearchBar    func(m Model, args ComponentArgs) string
	SearchText   func(m Model, view ActiveViewportType) string
	StatusBar    func(m Model, args ComponentArgs, items ...string) string
	SudoMode     func(m Model, args ComponentArgs) string
	TuiMode      func(m Model, args ComponentArgs) string
	ViewModeText func(m Model, args ComponentArgs) string
}

func InjectIntoModel(m *Model) {
	// Components
	m.Components.SearchBar = SearchBar
	m.Components.CurrentDir = CurrentDir
	m.Components.Header = Header
	m.Components.StatusBar = StatusBar
	m.Components.ViewModeText = ViewModeText
	m.Components.FileInfo = FileInfo
	m.Components.FileListView = FileList
	m.Components.TuiMode = TuiMode
	m.Components.SudoMode = SudoMode
	m.Components.SearchText = SearchText

	// Windows
	m.Windows.Help = HelpWindow
	m.Windows.Command = CommandWindow
	m.Windows.Dialog = DialogWindow
	m.Windows.Column = ColumnWindow
	m.Windows.Settings = SettingsWindow
}
