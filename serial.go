// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

// SerialConfig is common configuration for serial port.
type serialConfig struct {
	// Device path (/dev/ttyS0)
	Address string
	// Baud rate (default 19200)
	BaudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	CharSize int
	// Stop bits: 1 or 2 (default 1)
	StopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string
}

type serialController interface {
	Connect(config *serialConfig) (err error)
	Close() (err error)

	IsConnected() bool
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}
