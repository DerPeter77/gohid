//go:build linux

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

// EnumerateSysfs scans /sys/class/hidraw/ for all USB Devices
func EnumerateSysfs() ([]RawDeviceInfo, error) {
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

// StreamRawData reads the file and calls the handler
func StreamRawData(devPath string, handler func([]byte)) error {
	file, err := os.OpenFile(devPath, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("konnte %s nicht öffnen (Rechte geprüft? sudo?): %w", devPath, err)
	}
	defer file.Close()

	buf := make([]byte, 64)

	for {
		n, err := file.Read(buf)
		if err != nil {
			return fmt.Errorf("lesefehler: %w", err)
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		handler(data)
	}
}

// ParseBattery versucht, aus einem HID++ 2.0 Byte-Array die Spannung und % zu extrahieren
func ParseBattery(data []byte) (mv uint16, percent int, ok bool) {
	// Prüfen, ob es ein validier HID++ 2.0 Long Report (20 Bytes) ist
	// Byte 0 (0x11): Long Report ID
	// Byte 2 (0x06): Feature-Index der Batterie beim G Pro X
	if len(data) < 7 || data[0] != 0x11 || data[2] != 0x06 {
		return 0, 0, false
	}

	// Byte 4 (MSB) und Byte 5 (LSB) zu Millivolt (mV) verknüpfen
	mv = (uint16(data[4]) << 8) | uint16(data[5])

	// Plausibilitätsprüfung (Headset-Akkus bewegen sich zwischen 3000 mV und 4300 mV)
	if mv < 3000 || mv > 4400 {
		return 0, 0, false
	}

	// Umrechnung in Prozent (3400 mV = 0%, 4200 mV = 100%)
	percent = calculatePercent(mv)

	return mv, percent, true
}

func calculatePercent(mv uint16) int {
	if mv >= 4200 {
		return 100
	}
	if mv <= 3400 {
		return 0
	}
	// Lineare Skalierung für die Akkuspannung
	return int((float64(mv-3400) / float64(4200-3400)) * 100)
}
