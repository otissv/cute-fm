package filesystem

import (
	"fmt"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// FileInfo represents file or directory information
type FileInfo struct {
	Permissions  string // File permissions (e.g., "drwxr-xr-x", ".rw-r--r--")
	Size         string // File size (e.g., "1.3k", "5.7M"); for directories this is a byte total of direct children
	User         string // Owner username
	Group        string // Group name
	DateModified string // Date modified (e.g., "19 Nov 18:41")
	Name         string // File or directory name
	IsDir        bool   // Whether this is a directory
	Path         string // Full path to the file/directory
	// Type is a high-level classification used for coloring, e.g.:
	// "directory", "symlink", "socket", "pipe", "device", "executable", "regular".
	Type string
	// MimeType is the best-effort MIME type for the entry, such as
	// "text/plain" or "image/png". Directories are reported as
	// "inode/directory".
	MimeType string
}

// FileInfoColumn is an identifier for a column that can be shown for a FileInfo.
// Using a dedicated type avoids sprinkling raw strings like "Permissions" or
// "Size" throughout the codebase.
type FileInfoColumn string

const (
	ColumnPermissions  FileInfoColumn = "Permissions"
	ColumnSize         FileInfoColumn = "Size"
	ColumnMimeType     FileInfoColumn = "Type"
	ColumnUser         FileInfoColumn = "User"
	ColumnGroup        FileInfoColumn = "Group"
	ColumnDateModified FileInfoColumn = "DateModified"
	ColumnName         FileInfoColumn = "Name"
)

var ColumnNames = []FileInfoColumn{
	ColumnPermissions,
	ColumnSize,
	ColumnMimeType,
	ColumnUser,
	ColumnGroup,
	ColumnDateModified,
	ColumnName,
}

// ListDirectory lists the contents of a directory and returns file information
// Returns a slice of FileInfo structs sorted by name (directories first)
func ListDirectory(dirPath string) ([]FileInfo, error) {
	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var fileInfos []FileInfo

	// Process each entry
	for _, entry := range entries {
		// Get full path
		fullPath := filepath.Join(dirPath, entry.Name())

		// Get file info
		info, err := entry.Info()
		if err != nil {
			// Skip files we can't stat
			continue
		}

		// Get system-specific file info for user/group
		sysInfo := info.Sys()
		var uid, gid uint32
		if sysInfo != nil {
			if stat, ok := sysInfo.(*syscall.Stat_t); ok {
				uid = stat.Uid
				gid = stat.Gid
			}
		}

		// Get username
		username := "unknown"
		if u, err := user.LookupId(fmt.Sprintf("%d", uid)); err == nil {
			username = u.Username
		}

		// Get group name
		groupname := "unknown"
		if g, err := user.LookupGroupId(fmt.Sprintf("%d", gid)); err == nil {
			groupname = g.Name
		}

		// Format permissions
		permissions := formatPermissions(info.Mode(), entry.IsDir())

		// Format size. For regular files we use the file size directly. For
		// directories, we compute a shallow size by summing the sizes of
		// non-directory entries in that directory.
		var size string
		if entry.IsDir() {
			dirSize := calculateDirectorySize(fullPath)
			size = formatSize(dirSize, false)
		} else {
			size = formatSize(info.Size(), false)
		}

		// Format date modified
		dateModified := formatDateModified(info.ModTime())

		// Determine file type for colorization.
		fileType := classifyFileType(info, entry.IsDir())

		// Best-effort MIME type detection for the new "Type" column.
		mimeType := detectMimeType(entry.Name(), entry.IsDir())

		fileInfo := FileInfo{
			Permissions:  permissions,
			Size:         size,
			User:         username,
			Group:        groupname,
			DateModified: dateModified,
			Name:         entry.Name(),
			IsDir:        entry.IsDir(),
			Path:         fullPath,
			Type:         fileType,
			MimeType:     mimeType,
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(fileInfos, func(i, j int) bool {
		// If one is a directory and the other isn't, directory comes first
		if fileInfos[i].IsDir && !fileInfos[j].IsDir {
			return true
		}
		if !fileInfos[i].IsDir && fileInfos[j].IsDir {
			return false
		}
		// Both are same type, sort alphabetically
		return fileInfos[i].Name < fileInfos[j].Name
	})

	return fileInfos, nil
}

// detectMimeType returns a best-effort MIME type for the given entry name.
// Directories are always reported as "inode/directory". For regular files we
// consult the standard library's extension-based lookup and fall back to
// "application/octet-stream" when no mapping is known.
func detectMimeType(name string, isDir bool) string {
	if isDir {
		return "inode/directory"
	}

	ext := strings.ToLower(filepath.Ext(name))
	if ext != "" {
		if mt := mime.TypeByExtension(ext); mt != "" {
			return mt
		}
	}

	return "application/octet-stream"
}

// classifyFileType classifies a file into a high-level type used for styling.
func classifyFileType(info os.FileInfo, isDir bool) string {
	if isDir {
		return "Directory"
	}

	mode := info.Mode()

	switch {
	case mode&os.ModeSymlink != 0:
		return "Symlink"
	case mode&os.ModeSocket != 0:
		return "Socket"
	case mode&os.ModeNamedPipe != 0:
		return "pipe"
	case mode&os.ModeDevice != 0:
		return "Device"
	}

	// Treat any non-directory file with an execute bit set as "executable".
	if mode&0o111 != 0 {
		return "Executable"
	}

	return "File"
}

// formatPermissions formats file permissions in Unix-style format
// Directories: "drwxr-xr-x", Files: ".rw-r--r--"
func formatPermissions(mode os.FileMode, isDir bool) string {
	perm := mode.Perm()
	var result string

	if isDir {
		result = "d"
	} else {
		result = "."
	}

	// Owner permissions
	if perm&0o400 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perm&0o200 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perm&0o100 != 0 {
		result += "x"
	} else {
		result += "-"
	}

	// Group permissions
	if perm&0o040 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perm&0o020 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perm&0o010 != 0 {
		result += "x"
	} else {
		result += "-"
	}

	// Other permissions
	if perm&0o004 != 0 {
		result += "r"
	} else {
		result += "-"
	}
	if perm&0o002 != 0 {
		result += "w"
	} else {
		result += "-"
	}
	if perm&0o001 != 0 {
		result += "x"
	} else {
		result += "-"
	}

	return result
}

// formatSize formats file size in human-readable format (base 10).
// Format: "1.3k", "5.7M", etc.
func formatSize(size int64, isDir bool) string {
	const (
		KB = 1000 // Base 10 (decimal)
		MB = KB * 1000
		GB = MB * 1000
		TB = GB * 1000
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.1fT", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fk", float64(size)/float64(KB)) // lowercase 'k' to match format
	default:
		return fmt.Sprintf("%d", size) // No "B" suffix for bytes
	}
}

// calculateDirectorySize returns a shallow size for the given directory path by
// summing the sizes of non-directory entries directly inside it. Errors are
// treated as zero so that directory listings remain robust.
func calculateDirectorySize(dirPath string) int64 {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0
	}

	var total int64
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}
		total += info.Size()
	}
	return total
}

// formatDateModified formats the modification time
// Format: "DD MMM HH:MM" (e.g., "19 Nov 18:41")
func formatDateModified(modTime time.Time) string {
	// Format: "02 Jan 15:04"
	return modTime.Format("02 Jan 15:04")
}
