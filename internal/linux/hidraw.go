package linux

import (
	"fmt"
	"os"
)

// ReadHidFile reads the HID File from hidpath and sends all of the data to the given channel.
// Max Buffer for reading the file is at [64]byte
func ReadHidFile(hidpath string, ch chan<- []byte) error {
	file, err := os.OpenFile(hidpath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		buf := make([]byte, 64)
		n, err := file.Read(buf)
		if err != nil {
			close(ch)
			return err
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		ch <- data
	}
}

func WriteHidFile(hidpath string, data []byte) error {
	file, err := os.OpenFile(hidpath, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}

	if n < len(data) {
		return fmt.Errorf("Only wrote %d from %d Bytes", n, len(data))
	}

	return nil
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
