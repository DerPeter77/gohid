package devices

import (
	"fmt"
	"strings"

	"github.com/DerPeter77/gohid"
)

type GProXHandler struct{}

func (handler GProXHandler) Name() string {
	return "Logitech PRO X Wireless Gaming Headset"
}

func (handler GProXHandler) Match(deviceName string) bool {
	return strings.Contains(handler.Name(), deviceName)
}

func (handler GProXHandler) Handle(device *gohid.Device, args []string) error {
	switch args[0] {
	case "battery":
		ch, err := device.Read()
		if err != nil {
			return err
		}

		request := []byte{
			0x11, 0xFF, 0x06, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

		if err := device.Write(request); err != nil {
			return err
		}

		for data := range ch {
			if mv, percent, ok := ParseBattery(data); ok {
				fmt.Printf("Battery: %v%% - %vmV", percent, mv)
				return nil
			}
		}
	}
	return nil
}

func (handler GProXHandler) Commands() []string {
	return []string{
		"battery",
	}
}

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

func init() {
	Register(GProXHandler{})
}
