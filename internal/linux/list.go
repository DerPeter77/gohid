package linux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RawDeviceInfo struct {
	Path string
	Name string
}

// GetAllUsbDevices scans /sys/class/hidraw/ for all USB Devices
func GetAllUsbDevices() ([]RawDeviceInfo, error) {
	var devices []RawDeviceInfo

	entries, err := os.ReadDir("/sys/class/hidraw")
	if err != nil {
		if os.IsNotExist(err) {
			return devices, nil
		}
		return nil, fmt.Errorf("fehler beim Lesen von sysfs: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		ueventPath := filepath.Join("/sys/class/hidraw", name, "device", "uevent")

		info := RawDeviceInfo{
			Path: filepath.Join("/dev", name),
			Name: "Unbekanntes Gerät",
		}

		if devName, err := parseUeventName(ueventPath); err == nil && devName != "" {
			info.Name = devName
		}

		devices = append(devices, info)
	}

	return devices, nil
}

func parseUeventName(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "HID_NAME=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", scanner.Err()
}
