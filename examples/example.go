package main

import (
	"fmt"
	"log"

	logitechprox "github.com/DerPeter77/gohid/devices/logitechProX"
	"github.com/DerPeter77/gohid/internal/linux"
)

func main() {
	fmt.Println("Example")

	ch, err := linux.ReadHidFile("/dev/hidraw5")
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
