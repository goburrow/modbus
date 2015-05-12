// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"github.com/goburrow/serial"
)

const (
	// Default timeout
	serialTimeoutMillis = 5000
)

// serialPort has configuration and I/O controller.
type serialPort struct {
	// Serial port configuration.
	serial.Config
	// port is platform-dependent data structure for serial port.
	port serial.Port
	// Read timeout
	isConnected bool
}

func (mb *serialPort) Connect() (err error) {
	if mb.isConnected {
		return
	}
	if mb.port == nil {
		mb.port, err = serial.Open(&mb.Config)
	} else {
		err = mb.port.Open(&mb.Config)
	}
	if err == nil {
		mb.isConnected = true
	}
	return
}

func (mb *serialPort) Close() (err error) {
	if !mb.isConnected {
		return
	}
	err = mb.port.Close()
	mb.isConnected = false
	return
}
