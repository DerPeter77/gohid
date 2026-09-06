package devices

import "github.com/DerPeter77/gohid"

type DeviceHandler interface {
	Name() string
	Match(name string) bool
	Handle(device *gohid.Device, args []string) error
	Commands() []string
}
