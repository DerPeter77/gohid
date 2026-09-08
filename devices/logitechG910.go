package devices

import (
	"log"
	"os/exec"
	"slices"
	"strings"

	"github.com/DerPeter77/gohid"
	"github.com/DerPeter77/gohid/service"
)

type G910Handler struct{}

func (handler G910Handler) Name() string {
	return "Logitech Gaming Keyboard G910"
}

func (handler G910Handler) Match(deviceName string) bool {
	return strings.Contains(handler.Name(), deviceName)
}

func (handler G910Handler) Handle(device *gohid.Device, args []string) error {
	switch args[0] {
	case "service":
		// Get Key Rebinds
		for _, config := range service.CurrentConfig.ServiceDeviceConfig {
			if config.Name == handler.Name() {
				handler.StartService(device, config.KeyRebinds)
			}
		}
	}
	return nil
}

func (handler G910Handler) Commands() []string {
	return []string{
		"service",
	}
}

func (handler G910Handler) StartService(device *gohid.Device, keyRebinds []service.KeyRebind) {
	ch, err := device.Read()
	if err != nil {
		log.Fatal(err)
	}

	for data := range ch {
		for _, keyBind := range keyRebinds {
			if slices.Compare(keyBind.Event, data) == 0 {
				exec.Command(keyBind.Command, keyBind.Args...).Run()
				// Only for Debug Reasons
				// fmt.Printf("Data: %#+v\n", keyBind.KeyName)
			}
		}
	}
}

func init() {
	Register(G910Handler{})
}
