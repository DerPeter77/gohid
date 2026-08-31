package main

import (
	"fmt"
	"log"

	"github.com/DerPeter77/gohid/internal/linux"
)

func main() {
	fmt.Printf("Example loaded!\n\n")

	devices_list, err := linux.GetAllUsbDevices()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("USB Geräteliste:")
	for _, device := range devices_list {
		fmt.Printf("Gerät: %v Pfad: %v\n", device.Name, device.Path)
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
		if mv, percent, ok := linux.ParseBattery(data); ok {
			fmt.Printf("Akkustand empfangen: %d mV (%d%%)\n", mv, percent)
			close(readChan)
		}
	}
}
