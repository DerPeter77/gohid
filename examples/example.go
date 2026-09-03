package main

import (
	"fmt"
	"log"

	logitechprox "github.com/DerPeter77/gohid/devices/logitechProX"
	"github.com/DerPeter77/gohid/internal"
)

func main() {
	exampleHeadsetBattery()

}

func exampleHeadsetBattery() {
	fmt.Println("Example Headset Battery")

	device := internal.NewDevice("/dev/hidraw5")

	ch, err := device.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = logitechprox.RequestBatteryStatus("/dev/hidraw5")
	if err != nil {
		log.Fatal(err)
	}

	for data := range ch {
		if mv, percent, ok := logitechprox.ParseBattery(data); ok {
			fmt.Printf("Battery: %v%% - %vmV", percent, mv)
			return
		}
	}
}
