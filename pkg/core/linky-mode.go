package core

import "go.bug.st/serial"

type LinkyMode struct {
	BaudRate  int
	FrameSize int
	Parity    serial.Parity
	StopBits  serial.StopBits
}

var (
	Standard   = LinkyMode{9600, 7, serial.NoParity, serial.OneStopBit}
	Historical = LinkyMode{1200, 7, serial.NoParity, serial.OneStopBit}
)
