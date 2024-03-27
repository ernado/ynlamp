package ynlamp

import (
	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

const (
	packetSetColor = 0xa1
	packetSetWhite = 0xaa
	packetPowerOff = 0xa3
)

const (
	colorTypeTemperature = 0x01
	colorTypeRGB         = 0x02
)

func transmitServiceUUID() bluetooth.UUID {
	v := uuid.MustParse("f000aa60-0451-4000-b000-000000000000")
	return bluetooth.NewUUID(v)
}
