package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/DerPeter77/go-hid/internal/linux"
)

func main() {
	fmt.Println("Example")

	devPath := "/dev/hidraw5"
	fmt.Printf("Starte Live-Stream für %s (Beenden mit Strg+C)...\n\n", devPath)

	err := linux.StreamRawData(devPath, func(data []byte) {
		timestamp := time.Now().Format("15:04:05.000")

		// Hex-String formatieren
		hexDump := hex.EncodeToString(data)
		var formattedHex string
		for i := 0; i < len(hexDump); i += 2 {
			formattedHex += hexDump[i:i+2] + " "
		}

		// Basis-Ausgabe der Hex-Daten
		output := fmt.Sprintf("[%s] %2d Bytes | %s", timestamp, len(data), formattedHex)

		// Automatische Batterie-Erkennung & Berechnung
		if mv, percent, ok := linux.ParseBattery(data); ok {
			output += fmt.Sprintf("  <-- [BATTERIE: %d mV | ca. %d%%]", mv, percent)
		}

		fmt.Println(output)
	})

	if err != nil {
		log.Fatalf("Fehler: %v", err)
	}
}
