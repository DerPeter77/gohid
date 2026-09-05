package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
		case "devices":
			if len(os.Args) > 2 {
				deviceName := os.Args[2]

				// Get the devices list and check if the Name is in there
				devices_list, err := gohid.GetAllUsbDevices()
				if err != nil {
					log.Fatal(err)
				}

				for _, device := range devices_list {
					if strings.Contains(device.Name, deviceName) {

						go func(devPath string, devName string) {
							log.Default().Printf("Using Device Path: %v for Device Name: %v", devPath, devName)

							dev := gohid.NewDevice(devPath)
							ch, err := dev.Read()
							if err != nil {
								log.Fatal(err)
							}

							g6Press := []byte{0x11, 0xff, 0x8, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
							g7Press := []byte{0x11, 0xff, 0x8, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

							for data := range ch {
								fmt.Printf("DEBUG: %#+v\n", data)
								switch string(data) {
								case string(g6Press):
									fmt.Println("Pressed G6")
								case string(g7Press):
									fmt.Println("Pressed G7")
								}
							}
						}(device.Path, device.Name)
					}
				}

				for {
				}
			}
		}
	}
}
