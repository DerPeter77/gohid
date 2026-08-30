package gohid

import "github.com/DerPeter77/go-hid/internal/linux"

type DeviceInfo struct {
	Path string
	Name string
}

func GetDevicesList() ([]DeviceInfo, error) {
	rawDevs, err := linux.EnumerateSysfs()
	if err != nil {
		return nil, err
	}

	var result []DeviceInfo
	for _, rd := range rawDevs {
		result = append(result, DeviceInfo{
			Path: rd.Path,
			Name: rd.Name,
		})
	}

	return result, nil
}
