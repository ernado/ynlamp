package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-faster/errors"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func packet(command byte, data [3]byte) []byte {
	v := make([]byte, 6)
	v[0] = 0xAE
	v[1] = command
	v[2] = data[0]
	v[3] = data[1]
	v[4] = data[2]
	v[5] = 0x56
	return v
}

func run() error {
	if err := adapter.Enable(); err != nil {
		return errors.Wrap(err, "enable BLE stack")
	}
	var foundDevice bluetooth.ScanResult
	if err := adapter.Scan(func(a *bluetooth.Adapter, d bluetooth.ScanResult) {
		name := d.LocalName()
		if name == "" {
			return
		}
		d.LocalName()
		if name == "YN360II" {
			foundDevice = d
			_ = a.StopScan()
		}
	}); err != nil {
		return errors.Wrap(err, "start scan")
	}

	if foundDevice.LocalName() == "" {
		return errors.New("device not found")
	}

	fmt.Println("Found device:", foundDevice.LocalName())
	fmt.Println("Manufacturer data:", foundDevice.ManufacturerData())

	fmt.Println("Connecting")
	dev, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return errors.Wrap(err, "connect to device")
	}
	defer func() {
		_ = dev.Disconnect()
		fmt.Println("Disconnected")
	}()

	fmt.Println("Connected")

	targetServiceUUID := transmitServiceUUID()
	fmt.Println("target uuid:", targetServiceUUID)
	services, err := dev.DiscoverServices(nil)
	if err != nil {
		return errors.Wrap(err, "discover services")
	}
	if len(services) == 0 {
		return errors.New("no services found")
	}
	var dc bluetooth.DeviceCharacteristic
Services:
	for _, svc := range services {
		characteristics, err := svc.DiscoverCharacteristics(nil)
		if err != nil {
			return errors.Wrap(err, "discover characteristics")
		}
		for _, ch := range characteristics {
			if ch.UUID().String() != "f000aa61-0451-4000-b000-000000000000" {
				continue
			}
			dc = ch
			break Services
		}
	}

	if dc.UUID() == (bluetooth.UUID{}) {
		return errors.New("no characteristics found")
	}

	sendCommand := func(command byte, data [3]byte) error {
		v := packet(command, data)
		if _, err := dc.WriteWithoutResponse(v); err != nil {
			return errors.Wrap(err, "write packet")
		}
		return nil
	}

	defer func() {
		fmt.Println("Powering off")
		_ = sendCommand(packetPowerOff, [3]byte{})
	}()

	fmt.Println("Sending light")

	for i := 0; i >= 0; i++ {
		var args [3]byte
		if i%3 == 0 {
			args[1] = 255
		} else if i%3 == 1 {
			args[2] = 255
		} else {
			args[0] = 255
		}
		if err := sendCommand(packetSetColor, args); err != nil {
			return errors.Wrap(err, "send command")
		}
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}
