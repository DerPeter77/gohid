package devices

var RegisteredDeviceHandlers []DeviceHandler

func Register(handler DeviceHandler) {
	RegisteredDeviceHandlers = append(RegisteredDeviceHandlers, handler)
}

func ListHandlers() []DeviceHandler {
	return RegisteredDeviceHandlers
}

func GetHandler(deviceName string) DeviceHandler {
	for _, handler := range RegisteredDeviceHandlers {
		if handler.Match(deviceName) {
			return handler
		}
	}
	return nil
}
