package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DerPeter77/gohid"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1]

		switch args {
		case "list":
			devices_list, err := gohid.GetAllUsbDevices()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("USB Devicelist:\n\n")
			for _, device := range devices_list {
				fmt.Printf("Gerät: \"%v\" Pfad: \"%v\"\n", device.Name, device.Path)
			}
			return
		case "read":
			if len(os.Args) > 2 {
				path := os.Args[2]

				device := gohid.NewDevice(path)
				ch, err := device.Read()
				if err != nil {
					log.Fatal(err)
				}

				for data := range ch {
					fmt.Printf("%#+v\n", data)
				}
			}
		}
	}
}
