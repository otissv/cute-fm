package filesystem

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"syscall"
	"time"
)

// FileInfo represents file or directory information
type FileInfo struct {
	Permissions  string // File permissions (e.g., "drwxr-xr-x", ".rw-r--r--")
	Size         string // File size (e.g., "1.3k", "5.7M", "-" for directories)
	User         string // Owner username
	Group        string // Group name
	DateModified string // Date modified (e.g., "19 Nov 18:41")
	Name         string // File or directory name
	IsDir        bool   // Whether this is a directory
	Path         string // Full path to the file/directory
	// Type is a high-level classification used for coloring, e.g.:
	// "directory", "symlink", "socket", "pipe", "device", "executable", "regular".
	Type string
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

		// Format size
		size := formatSize(info.Size(), entry.IsDir())

		// Format date modified
		dateModified := formatDateModified(info.ModTime())

		// Determine file type for colorization.
		fileType := classifyFileType(info, entry.IsDir())

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

// classifyFileType classifies a file into a high-level type used for styling.
func classifyFileType(info os.FileInfo, isDir bool) string {
	if isDir {
		return "directory"
	}

	mode := info.Mode()

	switch {
	case mode&os.ModeSymlink != 0:
		return "symlink"
	case mode&os.ModeSocket != 0:
		return "socket"
	case mode&os.ModeNamedPipe != 0:
		return "pipe"
	case mode&os.ModeDevice != 0:
		return "device"
	}

	// Treat any non-directory file with an execute bit set as "executable".
	if mode&0o111 != 0 {
		return "executable"
	}

	return "file"
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

// formatSize formats file size in human-readable format (base 10)
// Returns "-" for directories
// Format: "1.3k", "5.7M", etc.
func formatSize(size int64, isDir bool) string {
	if isDir {
		return "-"
	}

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

// formatDateModified formats the modification time
// Format: "DD MMM HH:MM" (e.g., "19 Nov 18:41")
func formatDateModified(modTime time.Time) string {
	// Format: "02 Jan 15:04"
	return modTime.Format("02 Jan 15:04")
}
