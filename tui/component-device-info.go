package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cute/filesystem"
)

func (m *Model) UpdateDeviceInfoPane() {
	if len(m.lastDevices) == 0 {
		m.deviceInfoViewport.SetContent("")
		return
	}

	idx := m.deviceList.Index()
	if idx < 0 || idx >= len(m.lastDevices) {
		m.deviceInfoViewport.SetContent("")
		return
	}

	device := m.lastDevices[idx]

	m.deviceInfoViewport.SetContent(renderDeviceInfoPane(device))
}

func renderDeviceInfoPane(di filesystem.DeviceInfo) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Device info\n\n")
	fmt.Fprintf(&b, "Name: %s\n", di.Name)
	fmt.Fprintf(&b, "Device: %s\n", di.Device)
	fmt.Fprintf(&b, "Mount Point: %s\n", di.MountPoint)
	fmt.Fprintf(&b, "FS Type: %s\n", di.FSType)
	fmt.Fprintf(&b, "Size: %s\n", di.Size)
	fmt.Fprintf(&b, "Used: %s\n", di.Used)
	fmt.Fprintf(&b, "Available: %s\n", di.Avail)
	fmt.Fprintf(&b, "Use Percent: %s\n", di.UsePercent)
	fmt.Fprintf(&b, "Free Percent: %s\n", di.FreePercent)

	// Get and display UUID
	uuid := getDeviceUUID(di.Device)
	if uuid != "" {
		fmt.Fprintf(&b, "UUID: %s\n", uuid)
	}

	if di.Options != "" {
		fmt.Fprintf(&b, "\nOptions:\n")
		options := strings.Split(di.Options, ",")
		for _, opt := range options {
			opt = strings.TrimSpace(opt)
			if opt == "" {
				continue
			}

			// Check if option has a key=value format
			if strings.Contains(opt, "=") {
				parts := strings.SplitN(opt, "=", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Capitalize first letter of key
				label := capitalizeFirst(key)
				fmt.Fprintf(&b, "%s: %s\n", label, value)
			} else {
				// Simple flag option
				label := capitalizeFirst(opt)
				fmt.Fprintf(&b, "%s: true\n", label)
			}
		}
	}

	return b.String()
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getDeviceUUID(device string) string {
	byUUIDDir := "/dev/disk/by-uuid"
	entries, err := os.ReadDir(byUUIDDir)
	if err != nil {
		return ""
	}

	deviceBase := filepath.Base(device)
	for _, entry := range entries {
		linkPath := filepath.Join(byUUIDDir, entry.Name())
		target, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}

		// Resolve relative symlinks
		if !filepath.IsAbs(target) {
			target = filepath.Join(byUUIDDir, target)
		}
		targetBase := filepath.Base(target)

		// Check if this UUID points to our device
		if targetBase == deviceBase {
			return entry.Name()
		}
	}

	return ""
}
