package main

import (
	"fmt"
	"log"

	gohid "github.com/DerPeter77/go-hid"
)

func main() {
	fmt.Println("Example")

	devices, err := gohid.GetDevicesList()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		fmt.Printf("%#+v\n", device)
	}
}
