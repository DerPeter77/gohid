package linux

import (
	"fmt"
	"os"
)

// ReadHidFile reads the HID File from hidpath and sends all of the data to the given channel.
// Max Buffer for reading the file is at [64]byte.
// Returns error or nil if successfull.
func ReadHidFile(hidpath string) (<-chan []byte, error) {
	file, err := os.OpenFile(hidpath, os.O_RDONLY, 0666)
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
				close(ch)
				return
			}

			data := make([]byte, n)
			copy(data, buf[:n])

			ch <- data
		}
	}()
	return ch, nil
}

// Same like ReadHidFile but takes the channel given in Arguments. Use this function only if necessary
func ReadHidFileChannel(hidpath string, ch chan<- []byte) error {
	file, err := os.OpenFile(hidpath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		buf := make([]byte, 64)
		n, err := file.Read(buf)
		if err != nil {
			close(ch)
			return err
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		ch <- data
	}
}

// WriteHidFile writes data to the file at hidpath.
// Returns error or nil if successfull.
func WriteHidFile(hidpath string, data []byte) error {
	file, err := os.OpenFile(hidpath, os.O_WRONLY, 0666)
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
