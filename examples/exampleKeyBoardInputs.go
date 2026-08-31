package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/DerPeter77/gohid/internal/linux"
)

func exampleKeyBoardInputs() {
	fmt.Println("Example Keyboard Inputs")

	ch, err := linux.ReadHidFile("/dev/input/event5")
	if err != nil {
		log.Fatal(err)
	}

	for data := range ch {
		for i := 0; i < len(data); i += 24 {
			// Verhindern, dass wir außerhalb des Arrays lesen, falls Reste übrig sind
			if i+24 > len(data) {
				break
			}

			var event InputEvent
			buf := bytes.NewReader(data[i : i+24])

			// Linux-Input-Daten liegen immer im Little-Endian-Format vor
			err := binary.Read(buf, binary.LittleEndian, &event)
			if err != nil {
				fmt.Println("Fehler beim Dekodieren:", err)
				continue
			}

			// Nur EV_KEY (Type = 1) ausgeben, um Rauschen (wie EV_SYN oder EV_MSC) zu filtern
			if event.Type == 1 {
				status := "Losgelassen"
				if event.Value == 1 {
					status = "Gedrückt"
				} else if event.Value == 2 {
					status = "Gehalten"
				}

				fmt.Printf("Taste Code: %d | Aktion: %s\n", event.Code, status)
			}
		}
	}
}

type TimeVal struct {
	Sec  int64 // 8 Byte: Sekunden
	Usec int64 // 8 Byte: Mikrosekunden
}

// InputEvent entspricht der Linux 'struct input_event' (Gesamt: 24 Byte)
type InputEvent struct {
	Time  TimeVal // 16 Byte
	Type  uint16  // 2 Byte (z.B. EV_KEY = 1)
	Code  uint16  // 2 Byte (Der Key-Code)
	Value int32   // 4 Byte (1 = Gedrückt, 0 = Losgelassen, 2 = Gehalten)
}
