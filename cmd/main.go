package main

import (
	"fmt"
	"log"
	"os"

	logitechprox "github.com/DerPeter77/gohid/devices/logitechProX"
	"github.com/DerPeter77/gohid/internal/linux"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1]

		switch args {
		case "list":
			devices_list, err := linux.GetAllUsbDevices()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("USB Devicelist:\n\n")
			for _, device := range devices_list {
				fmt.Printf("Gerät: \"%v\" Pfad: \"%v\"\n", device.Name, device.Path)
			}
			return
		}
	}

	devPath := "/dev/hidraw5"
	readChan := make(chan []byte, 1)

	// 1. Reader im Hintergrund starten (weil Read ewig blockiert)
	go linux.ReadHidFile(devPath, readChan)

	// 3. Batterie-Abfrage direkt und einfach abfeuern!
	request := []byte{
		0x11, 0xFF, 0x06, 0x09, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	fmt.Println("Sende Batterie-Anfrage...")
	if err := linux.WriteHidFile(devPath, request); err != nil {
		log.Fatalf("Fehler beim Senden: %v", err)
	}

	for data := range readChan {
		if mv, percent, ok := logitechprox.ParseBattery(data); ok {
			fmt.Printf("Akkustand empfangen: %d mV (%d%%)\n", mv, percent)
			close(readChan)
		}
	}
}
