package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DerPeter77/gohid"
	"github.com/DerPeter77/gohid/devices"
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
			if len(os.Args) < 3 {
				fmt.Println("Usage: gohid devices <device-name> <command>")
				fmt.Println("Available device handlers:")
				for _, handler := range devices.ListHandlers() {
					fmt.Printf("  %s - commands: %v\n", handler.Name(), handler.Commands())
				}
				return
			}

			deviceName := os.Args[2]
			handler := devices.GetHandler(deviceName)
			if handler == nil {
				fmt.Printf("No handler for device: %s\n", deviceName)
				return
			}

			devicesList, err := gohid.GetAllUsbDevices()
			if err != nil {
				log.Fatal(err)
			}

			var devicePath string
			for _, d := range devicesList {
				if strings.Contains(d.Name, deviceName) || handler.Match(d.Name) {
					devicePath = d.Path
					break
				}
			}

			if devicePath == "" {
				fmt.Printf("Device not found: %s\n", deviceName)
				return
			}

			device := gohid.NewDevice(devicePath)

			if len(os.Args) > 3 {
				if err := handler.Handle(&device, os.Args[3:]); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("Available commands for %s: %v\n", handler.Name(), handler.Commands())
			}
		}
	}
}
