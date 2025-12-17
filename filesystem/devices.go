package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
)

type DeviceInfo struct {
	Device      string
	Name        string
	MountPoint  string
	FSType      string
	Options     string
	Size        string
	Used        string
	Avail       string
	UsePercent  string
	FreePercent string
}

type DeviceInfoColumn string

type DeviceInfoColumnHeadings struct {
	Device      DeviceInfoColumn
	Name        DeviceInfoColumn
	MountPoint  DeviceInfoColumn
	FsType      DeviceInfoColumn
	Size        DeviceInfoColumn
	Used        DeviceInfoColumn
	Avail       DeviceInfoColumn
	UsePercent  DeviceInfoColumn
	FreePercent DeviceInfoColumn
}

var (
	DeviceInfoColumns = DeviceInfoColumnHeadings{
		Device:      "Device",
		Name:        "Name",
		MountPoint:  "Mount",
		Size:        "Size",
		Used:        "Used",
		Avail:       "Avail",
		FsType:      "FS Type",
		UsePercent:  "UsePercent",
		FreePercent: "FreePercent",
	}

	DeviceInfoColumnNames = []DeviceInfoColumn{
		DeviceInfoColumns.Device,
		DeviceInfoColumns.Name,
		DeviceInfoColumns.MountPoint,
		DeviceInfoColumns.Size,
		DeviceInfoColumns.Used,
		DeviceInfoColumns.Avail,
		DeviceInfoColumns.UsePercent,
		DeviceInfoColumns.FreePercent,
	}
)

func ListDevices() ([]DeviceInfo, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/mounts: %w", err)
	}
	defer file.Close()

	var devices []DeviceInfo
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]
		options := fields[3]

		// Skip virtual filesystems and special mounts
		if shouldSkipDevice(device, mountPoint, fsType) {
			continue
		}

		// Get device name/label
		deviceName := getDeviceName(device, mountPoint)

		// Get disk usage information
		size, used, avail, usePercent, freePercent := getDiskUsage(mountPoint)

		devices = append(devices, DeviceInfo{
			Device:      device,
			Name:        deviceName,
			MountPoint:  mountPoint,
			FSType:      fsType,
			Options:     options,
			Size:        size,
			Used:        used,
			Avail:       avail,
			UsePercent:  usePercent,
			FreePercent: freePercent,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/mounts: %w", err)
	}

	return devices, nil
}

// getDiskUsage retrieves disk usage information for a mount point using syscall.Statfs
func getDiskUsage(mountPoint string) (size, used, avail, usePercent, freePercent string) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(mountPoint, &stat)
	if err != nil {
		return "-", "-", "-", "-", "-"
	}

	// Calculate sizes in bytes
	// Blocks are in units of stat.Bsize bytes
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	availBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - (stat.Bfree * uint64(stat.Bsize))

	// Format sizes
	size = formatDiskSize(totalBytes)
	used = formatDiskSize(usedBytes)
	avail = formatDiskSize(availBytes)

	// Calculate usage and free percentages
	var usePercentVal, freePercentVal float64
	if totalBytes > 0 {
		usePercentVal = float64(usedBytes) / float64(totalBytes) * 100
		freePercentVal = float64(availBytes) / float64(totalBytes) * 100
	}
	usePercent = fmt.Sprintf("%.0f%%", usePercentVal)
	freePercent = fmt.Sprintf("%.0f%%", freePercentVal)

	return size, used, avail, usePercent, freePercent
}

// formatDiskSize formats disk size in human-readable format (base 10)
func formatDiskSize(bytes uint64) string {
	const (
		KB = 1000
		MB = KB * 1000
		GB = MB * 1000
		TB = GB * 1000
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1fT", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.1fG", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1fM", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1fk", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d", bytes)
	}
}

// CleanDeviceNames replaces escape sequences like \\x20 with spaces in device names
func CleanDeviceNames(devices []DeviceInfo) {
	for i := range devices {
		devices[i].Name = strings.ReplaceAll(devices[i].Name, "\\x20", " ")
		devices[i].Name = strings.ReplaceAll(devices[i].Name, "\\040", " ")
	}
}

// shouldSkipDevice determines if a device should be skipped from the list
// Skips virtual filesystems, special mounts, and system directories
func shouldSkipDevice(device, mountPoint, fsType string) bool {
	// Skip virtual filesystems
	virtualFS := []string{
		"proc", "sysfs", "tmpfs", "devtmpfs", "devpts",
		"cgroup", "cgroup2", "pstore", "bpf", "tracefs",
		"debugfs", "securityfs", "hugetlbfs", "mqueue",
		"overlay", "autofs", "binfmt_misc", "configfs",
		"efivarfs", "fusectl", "hugetlbfs", "ramfs",
		"rpc_pipefs", "systemd-1", "none",
	}

	for _, vfs := range virtualFS {
		if fsType == vfs {
			return true
		}
	}

	// Skip if device doesn't start with /dev/ (not a block device)
	if !strings.HasPrefix(device, "/dev/") {
		return true
	}

	// Skip loop devices unless they're explicitly wanted
	if strings.HasPrefix(device, "/dev/loop") {
		return true
	}

	// Skip system directories that are typically not user-accessible storage
	systemMounts := []string{
		"/proc", "/sys", "/dev", "/run", "/tmp",
		"/var/run", "/var/lock", "/boot/efi",
	}

	for _, sysMount := range systemMounts {
		if mountPoint == sysMount {
			return true
		}
	}

	return false
}

func getDeviceName(device, mountPoint string) string {
	// Method 1: Try to get filesystem label from /dev/disk/by-label/
	label := getDeviceLabel(device)
	if label != "" {
		return label
	}

	// Method 2: Try to get device model from /sys/block
	model := getDeviceModel(device)
	if model != "" {
		return model
	}

	// Method 3: Use mount point name (last component)
	if mountPoint != "" && mountPoint != "/" {
		name := filepath.Base(mountPoint)
		if name != "" && name != "." {
			return name
		}
	}

	// Method 4: Fall back to device name (e.g., "sda1")
	return filepath.Base(device)
}

func getDeviceLabel(device string) string {
	byLabelDir := "/dev/disk/by-label"
	entries, err := os.ReadDir(byLabelDir)
	if err != nil {
		return ""
	}

	deviceBase := filepath.Base(device)
	for _, entry := range entries {
		linkPath := filepath.Join(byLabelDir, entry.Name())
		target, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}

		// Resolve relative symlinks
		if !filepath.IsAbs(target) {
			target = filepath.Join(byLabelDir, target)
		}
		targetBase := filepath.Base(target)

		// Check if this label points to our device
		if targetBase == deviceBase {
			return entry.Name()
		}
	}

	return ""
}

func getDeviceModel(device string) string {
	deviceBase := filepath.Base(device)

	baseDevice := deviceBase

	// Remove partition suffix patterns: p1, p2, etc. (for NVMe and eMMC)
	if strings.Contains(baseDevice, "p") {
		parts := strings.Split(baseDevice, "p")
		if len(parts) > 1 {
			// Check if last part is numeric
			lastPart := parts[len(parts)-1]
			if len(lastPart) > 0 && strings.Trim(lastPart, "0123456789") == "" {
				baseDevice = strings.Join(parts[:len(parts)-1], "p")
			}
		}
	} else {
		baseDevice = strings.TrimRight(deviceBase, "0123456789")
	}

	modelPath := filepath.Join("/sys/block", baseDevice, "device", "model")
	modelData, err := os.ReadFile(modelPath)
	if err != nil {
		return ""
	}

	model := strings.TrimSpace(string(modelData))
	if model != "" {
		return model
	}

	return ""
}

type DeviceChangeMsg struct {
	Devices []DeviceInfo
	Added   []DeviceInfo
	Removed []DeviceInfo
}

// Monitors devices and compares with previous state
func WatchDevicesWithState(interval time.Duration, lastDevices []DeviceInfo) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		currentDevices, err := ListDevices()
		if err != nil {
			return DeviceChangeMsg{
				Devices: lastDevices,
				Added:   []DeviceInfo{},
				Removed: []DeviceInfo{},
			}
		}

		added, removed := diffDevices(lastDevices, currentDevices)

		if len(added) > 0 || len(removed) > 0 {
			return DeviceChangeMsg{
				Devices: currentDevices,
				Added:   added,
				Removed: removed,
			}
		}

		return nil
	})
}

func diffDevices(old, new []DeviceInfo) (added, removed []DeviceInfo) {
	// Create maps for quick lookup
	oldMap := make(map[string]DeviceInfo)
	for _, dev := range old {
		oldMap[dev.MountPoint] = dev
	}

	newMap := make(map[string]DeviceInfo)
	for _, dev := range new {
		newMap[dev.MountPoint] = dev
	}

	// Find added devices (in new but not in old)
	for mountPoint, dev := range newMap {
		if _, exists := oldMap[mountPoint]; !exists {
			added = append(added, dev)
		}
	}

	// Find removed devices (in old but not in new)
	for mountPoint, dev := range oldMap {
		if _, exists := newMap[mountPoint]; !exists {
			removed = append(removed, dev)
		}
	}

	return added, removed
}
