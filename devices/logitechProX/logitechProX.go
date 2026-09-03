package logitechprox

import (
	"github.com/DerPeter77/gohid/internal"
)

// Sends a request for logitech pro x to respond with the battery status.
func RequestBatteryStatus(devPath string) error {
	device := internal.NewDevice(devPath)

	request := []byte{
		0x11, 0xFF, 0x06, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	if err := device.Write(request); err != nil {
		return err
	}
	return nil
}

// ParseBattery attempts to extract voltage and percentage from a HID++ 2.0 byte array
func ParseBattery(data []byte) (mv uint16, percent int, ok bool) {
	// Check if it's a valid HID++ 2.0 Long Report (20 bytes)
	// Byte 0 (0x11): Long Report ID
	// Byte 2 (0x06): Battery feature index for G Pro X
	if len(data) < 7 || data[0] != 0x11 || data[2] != 0x06 {
		return 0, 0, false
	}

	// Combine byte 4 (MSB) and byte 5 (LSB) into millivolts (mV)
	mv = (uint16(data[4]) << 8) | uint16(data[5])

	// Sanity check (headset batteries range between 3000 mV and 4300 mV)
	if mv < 3000 || mv > 4400 {
		return 0, 0, false
	}

	// Convert to percentage (3400 mV = 0%, 4200 mV = 100%)
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
	// Linear scaling for battery voltage
	return int((float64(mv-3400) / float64(4200-3400)) * 100)
}
