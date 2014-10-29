// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"errors"
)

var errSerialNotImplemented = errors.New("Serial on Windows is not implemented yet")

// serialPort is included in serialTransporter.
type serialPort struct {
}

// Connect opens serial port. Device must be set before calling this method.
func (mb *serialTransporter) Connect() (err error) {
	err = errSerialNotImplemented
	return
}

func (mb *serialTransporter) Close() (err error) {
	err = errSerialNotImplemented
	return
}

// isConnected returns true if serial port has been opened
func (mb *serialTransporter) isConnected() bool {
	return false
}

// read reads from serial port, blocked until data received or timeout after Timeout
func (mb *serialTransporter) read(b []byte) (n int, err error) {
	err = errSerialNotImplemented
	return
}

func (mb *serialTransporter) write(b []byte) (n int, err error) {
	err = errSerialNotImplemented
	return
}
