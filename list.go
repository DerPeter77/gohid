package gohid

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RawDeviceInfo struct {
	Path     string
	Name     string
	Phys     string
	Uniq     string
	Modalias string
}

// GetAllUsbDevices scans /sys/class/hidraw/ for all USB Devices and returns them
func GetAllUsbDevices() ([]RawDeviceInfo, error) {
	var devices []RawDeviceInfo

	entries, err := os.ReadDir("/sys/class/hidraw")
	if err != nil {
		if os.IsNotExist(err) {
			return devices, nil
		}
		return nil, fmt.Errorf("Error while reading sysfs: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		deviceDir := filepath.Join("/sys/class/hidraw", name, "device")

		info := RawDeviceInfo{
			Path: filepath.Join("/dev", name),
			Name: "Unknown Device",
		}

		uevent := parseUevent(filepath.Join(deviceDir, "uevent"))
		if name, ok := uevent["HID_NAME"]; ok && name != "" {
			info.Name = name
		}
		info.Phys = uevent["HID_PHYS"]
		info.Uniq = uevent["HID_UNIQ"]

		if modalias, err := os.ReadFile(filepath.Join(deviceDir, "modalias")); err == nil {
			info.Modalias = strings.TrimSpace(string(modalias))
		}

		devices = append(devices, info)
	}

	return devices, nil
}

// parseUevent reads a uevent file and returns its key/value pairs.
func parseUevent(path string) map[string]string {
	result := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return result
}
