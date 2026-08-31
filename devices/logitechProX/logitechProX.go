package logitechprox

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
