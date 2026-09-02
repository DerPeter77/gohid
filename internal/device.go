package internal

import (
	"fmt"
	"os"
)

type Device struct {
	path string
}

// NewDevice creates and returns a new Device with the path given as Argument
func NewDevice(path string) Device {
	return Device{
		path: path,
	}
}

// GetDevicePath returns the Device Path
func (device Device) GetDevicePath() string {
	return device.path
}

// Read reads the HID File from hidpath and sends all of the data to the given channel.
// Max Buffer for reading the file is at [64]byte.
// Returns error or nil if successfull.
func (device Device) Read() (<-chan []byte, error) {
	file, err := os.OpenFile(device.path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte, 10)

	go func() {
		defer file.Close()
		defer close(ch)

		for {
			buf := make([]byte, 64)
			n, err := file.Read(buf)
			if err != nil {
				return
			}

			data := make([]byte, n)
			copy(data, buf[:n])

			ch <- data
		}
	}()
	return ch, nil
}

// Write writes data to the file at hidpath.
// Returns error or nil if successfull.
func (device Device) Write(data []byte) error {
	file, err := os.OpenFile(device.path, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}

	if n < len(data) {
		return fmt.Errorf("Only wrote %d from %d Bytes", n, len(data))
	}

	return nil
}
