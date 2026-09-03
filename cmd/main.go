package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DerPeter77/gohid/internal"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1]

		switch args {
		case "list":
			devices_list, err := internal.GetAllUsbDevices()
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
}
