package tui

import (
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"cute/filesystem"
	"cute/theming"
)

type (
	WindowKind         string
	SplitPaneType      string
	TUIMode            string
	ActiveViewportType string
	ActiveSetting      string
)

type TUIModes struct {
	ModeAddFile           TUIMode
	ModeCd                TUIMode
	ModeColumnVisibility  TUIMode
	ModeCommand           TUIMode
	ModeCopy              TUIMode
	ModeFilter            TUIMode
	ModeGoto              TUIMode
	ModeHelp              TUIMode
	ModeMkdir             TUIMode
	ModeMove              TUIMode
	ModeNormal            TUIMode
	ModeQuit              TUIMode
	ModeRemove            TUIMode
	ModeRename            TUIMode
	ModeSelect            TUIMode
	ModeSort              TUIMode
	ModeFileListSplitPane TUIMode
	ModeSettings          TUIMode
}

type FileListMode string

type ViewPrimitiveView interface {
	View() string
}

type ViewPrimitive string

func (t ViewPrimitive) View() string {
	return string(t)
}

type filePane struct {
	currentDir      string
	allFiles        []filesystem.FileInfo
	files           []filesystem.FileInfo
	fileList        list.Model
	filterQuery     string
	dirBackStack    []string
	dirForwardStack []string
	columns         []filesystem.FileInfoColumn

	marked map[string]bool
}

type SelectedEntry struct {
	Name  string
	Path  string
	IsDir bool
	Type  string
}

type Settings struct {
	ColumnVisibilityFileList    []filesystem.FileInfoColumn
	ColumnVisibilityDevice      []filesystem.DeviceInfoColumn
	FileListMode                FileListMode
	SortDeviceColumnBy          filesystem.FileInfoColumn
	SortDeviceColumnDirection   SortColumnByDirection
	SortFileListColumnBy        filesystem.FileInfoColumn
	SortFileListColumnDirection SortColumnByDirection
	SplitPane                   SplitPaneType
	StartDir                    string
}

type SortColumnByDirection string

type SortDeviceColumnBy struct {
	column    filesystem.DeviceInfoColumn
	direction SortColumnByDirection
}

type SortFileListColumnBy struct {
	column    filesystem.FileInfoColumn
	direction SortColumnByDirection
}

const (
	ModeAddFile           TUIMode = "ADD_FILE"
	ModeCd                TUIMode = "CD"
	ModeColumnVisibility  TUIMode = "COLUMN VISIBILITY"
	ModeCommand           TUIMode = "COMMAND"
	ModeDevice            TUIMode = "COMPUTER"
	ModeCopy              TUIMode = "COPY"
	ModeFilter            TUIMode = "FILTER"
	ModeGoto              TUIMode = "GOTO"
	ModeHelp              TUIMode = "HELP"
	ModeMkdir             TUIMode = "MKDIR"
	ModeMove              TUIMode = "MOVE"
	ModeNormal            TUIMode = "NORMAL"
	ModeParent            TUIMode = "PARENT"
	ModeQuit              TUIMode = "QUIT"
	ModeRemove            TUIMode = "REMOVE"
	ModeRename            TUIMode = "RENAME"
	ModeSelect            TUIMode = "SELECT"
	ModeSort              TUIMode = "SORT"
	ModeFileListSplitPane TUIMode = "SPLIT PANE"
	ModeSettings          TUIMode = "SETTINGS"

	SortingAsc  SortColumnByDirection = "ASC"
	SortingDesc SortColumnByDirection = "DESC"

	DeviceInfoSplitPaneType SplitPaneType = "DEVICE_INFO"
	FileInfoSplitPaneType   SplitPaneType = "FILE_INFO"
	FileListSplitPaneType   SplitPaneType = "FILE_LIST"
	PreviewPaneType         SplitPaneType = "PREVIEW"

	LeftViewportType  ActiveViewportType = "LEFT"
	RightViewportType ActiveViewportType = "RIGHT"

	FileListModeList   FileListMode = "ll"
	FileListModeFile   FileListMode = "lf"
	FileListModeDir    FileListMode = "ld"
	FileListModeDevice FileListMode = "device"

	SETTING_START                 ActiveSetting = "SETTING_START"
	SETTING_SPLIT_PANE            ActiveSetting = "SETTING_SPLIT_PANE"
	SETTING_FILE_LIST_MODE        ActiveSetting = "SETTING_FILE_LIST_MODE"
	SETTING_COLUMN_VISIBILITY     ActiveSetting = "SETTING_COLUMN_VISIBILITY"
	SETTING_SORT_BY_COLUMN        ActiveSetting = "SETTING_SORT_BY_COLUMN"
	SETTING_SORT_COLUMN_DIRECTION ActiveSetting = "SETTING_SORT_COLUMN_DIRECTION"
)

var (
	ActiveTuiMode   TUIMode = ModeNormal
	PreviousTuiMode TUIMode = ModeNormal

	ActiveFileListMode   FileListMode = FileListModeList
	PreviousFileListMode FileListMode = FileListModeList

	TuiModes = TUIModes{
		ModeAddFile:          ModeAddFile,
		ModeCd:               ModeCd,
		ModeColumnVisibility: ModeColumnVisibility,
		ModeCommand:          ModeCommand,
		ModeCopy:             ModeCopy,
		ModeFilter:           ModeFilter,
		ModeGoto:             ModeGoto,
		ModeHelp:             ModeHelp,
		ModeMkdir:            ModeMkdir,
		ModeNormal:           ModeNormal,
		ModeRemove:           ModeRemove,
		ModeRename:           ModeRename,
		ModeSelect:           ModeSelect,
		ModeSort:             ModeSort,
	}
)

type Model struct {
	activeWindow         WindowKind
	activeSplitPane      SplitPaneType
	previousSplitPane    SplitPaneType
	activeViewport       ActiveViewportType
	commandHistory       []string // Command history for auto-complete
	commandInput         textinput.Model
	deviceList           list.Model // Device/device list for navigation
	configDir            string
	countPrefix          int            // countPrefix stores a pending numeric prefix for Vim-style navigation (e.g. "10j" / "3↓" in the file list). A value of 0 means "no active prefix".
	fileInfoViewport     viewport.Model // Independent state for each file-list pane.
	deviceInfoViewport   viewport.Model
	height               int
	helpScrollOffset     int      // Help window scroll state
	historyIndex         int      // Current index in historyMatches for navigation
	historyMatches       []string // Filtered matches based on current input
	isActionInProgress   bool
	isSidePanelOpen      bool
	isSplitPaneOpen      bool
	isSudo               bool
	jumpTo               string
	lastDevices          []filesystem.DeviceInfo // Track last known devices for change detection
	layout               string
	layoutRows           []string
	leftPane             filePane
	menuCursorIndex      int
	rightPane            filePane
	searchInput          textinput.Model
	settings             Settings
	showRightPane        bool
	sortDeviceColumnBy   SortDeviceColumnBy
	sortFileListColumnBy SortFileListColumnBy
	theme                theming.Theme
	titleText            string
	viewportHeight       int
	viewportWidth        int
	width                int
	deviceColumns        []filesystem.DeviceInfoColumn
}

func (m Model) Init() tea.Cmd {
	// Start device monitoring and text input blinking
	// Note: lastDevices should already be initialized in InitialModel()
	return tea.Batch(
		textinput.Blink,
		filesystem.WatchDevicesWithState(2*time.Second, m.lastDevices),
	)
}

func (m Model) GetActiveWindow() WindowKind {
	return m.activeWindow
}

func (m *Model) GetActivePane() *filePane {
	if m.activeViewport == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return &m.rightPane
	}
	return &m.leftPane
}

func (m Model) GetActiveViewport() ActiveViewportType {
	return m.activeViewport
}

func (m Model) GetFileListColumnVisibility() []filesystem.FileInfoColumn {
	return m.leftPane.columns
}

func (m Model) GetDeviceColumnVisibility() []filesystem.DeviceInfoColumn {
	return m.deviceColumns
}

func (m Model) GetColumnVisibilityForViewport(view ActiveViewportType) []filesystem.FileInfoColumn {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.columns
	}
	return m.leftPane.columns
}

func (m Model) GetCommandInput() textinput.Model {
	return m.commandInput
}

func (m Model) GetCommandInputView() string {
	return m.commandInput.View()
}

func (m Model) GetConfigDir() string {
	return m.configDir
}

func (m Model) GetIsActionInProgress() bool {
	return m.isSplitPaneOpen
}

func (m Model) GetIsSplitPaneOpen() bool {
	return m.isSplitPaneOpen
}

func (m Model) GetLeftPaneCurrentDir() string {
	return m.leftPane.currentDir
}

func (m Model) GetLeftPaneFileList() list.Model {
	return m.leftPane.fileList
}

func (m Model) GetLeftPaneFileListForViewport(view ActiveViewportType) list.Model {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.fileList
	}
	return m.leftPane.fileList
}

func (m Model) GetMenuCursorIndex() int {
	return m.menuCursorIndex
}

func (m Model) GetFileInfoViewport() viewport.Model {
	return m.fileInfoViewport
}

func (m Model) GetDeviceInfoViewport() viewport.Model {
	return m.deviceInfoViewport
}

func (m Model) GetRightCurrentDir() string {
	return m.rightPane.currentDir
}

func (m Model) GetRightPaneCurrentDir() string {
	return m.rightPane.currentDir
}

func (m Model) GetSearchInputText() string {
	// Default to the active viewport's filter text.
	if m.activeViewport == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.filterQuery
	}
	return m.leftPane.filterQuery
}

func (m Model) GetSearchInputView() string {
	return m.searchInput.View()
}

func (m Model) GetSearchInputTextForViewport(view ActiveViewportType) string {
	if view == RightViewportType && m.isSplitPaneOpen && m.activeSplitPane == FileListSplitPaneType {
		return m.rightPane.filterQuery
	}
	return m.leftPane.filterQuery
}

func (m Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m Model) GetSortDeviceColumnBy() SortDeviceColumnBy {
	return m.sortDeviceColumnBy
}

func (m Model) GetSortFileListColumnBy() SortFileListColumnBy {
	return m.sortFileListColumnBy
}

func (m Model) GetTheme() theming.Theme {
	return m.theme
}

func (m Model) GetTitleText() string {
	return m.titleText
}

func (m Model) GetHelpScrollOffset() int {
	return m.helpScrollOffset
}
