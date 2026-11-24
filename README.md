# gohid

Go Lib to send and receive from HID Devices in Linux

## Installation

```bash
go install github.com/DerPeter77/gohid/cmd@latest
```

## Verwendung

Das Kommandozeilen-Tool `gohid` bietet folgende Befehle:

### `list`

Listet alle angeschlossenen USB-HID-Geräte mit Name und Gerätepfad auf.

```bash
gohid list
```

Beispielausgabe:

```text
USB Devicelist:

Gerät: "Logitech Gaming Keyboard G910"
  Pfad:    "/dev/hidraw0"
  Phys:    "usb-0000:00:14.0-2/input0"
  Uniq:    ""
  Modalias:"hid:b0003g0001v0000046dp0000c073"

Gerät: "Logitech Gaming Keyboard G910"
  Pfad:    "/dev/hidraw1"
  Phys:    "usb-0000:00:14.0-2/input1"
  Uniq:    ""
  Modalias:"hid:b0003g0001v0000046dp0000c073"

Gerät: "Logitech PRO X Wireless Gaming Headset"
  Pfad:    "/dev/hidraw2"
  Phys:    "usb-0000:00:14.0-3/input0"
  Uniq:    "00:11:22:33:44:55"
  Modalias:"hid:b0003g0001v0000046dp0000c073"
```

### `read <pfad>`

Liest kontinuierlich Rohdaten von einem HID-Gerät am angegebenen Pfad und gibt sie auf der Konsole aus.

```bash
gohid read /dev/hidraw0
```

### `devices <gerätename> [pfad] <befehl>`

Steuert ein spezifisches Gerät über einen Geräte-Handler an.

- `<gerätename>` – Name (oder Teilname) des Geräts, z. B. `Logitech Gaming Keyboard G910`
- `[pfad]` – optional: direkter Gerätepfad (z. B. `/dev/hidraw0`). Wird er weggelassen, sucht `gohid` automatisch nach einem passenden Gerät.
- `<befehl>` – der vom Handler unterstützte Befehl

Ohne Befehl werden die verfügbaren Handler und deren Befehle angezeigt:

```bash
gohid devices
```

Beispielausgabe:

```text
Usage: gohid devices <device-name> <command>
Available device handlers:
  Logitech Gaming Keyboard G910 - commands: [service]
  Logitech PRO X Wireless Gaming Headset - commands: [battery]
```

#### Verfügbare Geräte-Handler und Befehle

| Gerät | Handler-Name | Befehle | Beschreibung |
|-------|--------------|---------|--------------|
| Logitech Gaming Keyboard G910 | `Logitech Gaming Keyboard G910` | `service` | Startet den Tasten-Rebind-Dienst (siehe `service/service.go`) |
| Logitech PRO X Wireless Gaming Headset | `Logitech PRO X Wireless Gaming Headset` | `battery` | Liest den Akkustand des Headsets aus |

Beispiele:

```bash
# Automatische Gerätesuche nach Name und Befehl ausführen
gohid devices "Logitech Gaming Keyboard G910" service

# Mit explizitem Gerätepfad
gohid devices "Logitech PRO X Wireless Gaming Headset" /dev/hidraw1 battery

# Verfügbare Befehle für einen Handler anzeigen
gohid devices "Logitech Gaming Keyboard G910"
```

### Unterschiedliche Interfaces desselben Geräts erkennen

Ein physisches HID-Gerät (z. B. eine Logitech-Tastatur) legt oft mehrere `/dev/hidraw*` Interfaces an – eines für normale Tasten, eines für die Logitech-G-Tasten und eines für Maussensor-Eingaben. Alle Interfaces tragen denselben Gerätenamen (`HID_NAME`), unterscheiden sich aber in:

| Feld | Bedeutung |
|------|-----------|
| `Pfad` | Der Gerätepfad unter `/dev/hidraw*` (einzigartig pro Interface) |
| `Phys` | Die physische USB-Position (identisch für alle Interfaces desselben Geräts) |
| `Uniq` | Eindeutige Geräte-ID / Seriennummer (identisch für alle Interfaces desselben Geräts) |
| `Modalias` | HID-Modul-Alias (identisch für alle Interfaces desselben Geräts) |

**Beispiel:** Zwei Einträge mit demselben Namen `"Logitech Gaming Keyboard G910"` und demselben `Phys`-Wert, aber unterschiedlichen Pfaden (`/dev/hidraw0` und `/dev/hidraw1`), sind zwei Interfaces desselben physischen Tastatur-Geräts. Über `Phys` oder `Uniq` lässt sich zuordnen, welcher Pfad zu welchem Gerät gehört.

## Bibliothek nutzen

Neben dem CLI-Tool kann das Paket auch als Go-Bibliothek verwendet werden:

```go
package main

import (
    "fmt"
    "github.com/DerPeter77/gohid"
)

func main() {
    // Alle USB-HID-Geräte auflisten
    devices, err := gohid.GetAllUsbDevices()
    if err != nil {
        panic(err)
    }

    for _, d := range devices {
        fmt.Printf("Gerät: %s (Pfad: %s)\n", d.Name, d.Path)
    }

    // Ein Gerät öffnen und lesen
    device := gohid.NewDevice("/dev/hidraw0")
    ch, err := device.Read()
    if err != nil {
        panic(err)
    }

    for data := range ch {
        fmt.Printf("%#+v\n", data)
    }

    // Daten schreiben
    err = device.Write([]byte{0x00, 0x01, 0x02})
    if err != nil {
        panic(err)
    }
}
```

## Voraussetzungen

- Linux (nutzt `/sys/class/hidraw/` und `/dev/hidraw*`)
- Go 1.27 oder höher
- Lese-/Schreibrechte auf die Geräte unter `/dev/hidraw*` (ggf. `sudo` oder udev-Regeln erforderlich)

## Lizenz

Siehe [LICENSE](LICENSE).
