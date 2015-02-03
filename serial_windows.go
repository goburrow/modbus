// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

// #include <windows.h>
import "C"

import (
	"fmt"
	"syscall"
	"time"
)

// serialPort is included in serialTransporter.
type serialPort struct {
	handle syscall.Handle
}

// Connect opens serial port. Device must be set before calling this method.
func (mb *serialTransporter) Connect() (err error) {
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: connecting '%v'\n", mb.Address)
	}
	// Timeout is required
	if mb.Timeout <= 0 {
		mb.Timeout = serialTimeoutMillis * time.Millisecond
	}
	handle, err := newHandle(&mb.serialConfig)
	if err != nil {
		return
	}
	// Read and write timeout
	timeoutMillis := C.DWORD(mb.Timeout.Nanoseconds() / 1E6)
	if timeoutMillis < 1 {
		timeoutMillis = 1
	}
	var timeouts C.COMMTIMEOUTS
	// wait until a byte arrived or time out
	timeouts.ReadIntervalTimeout = C.MAXDWORD
	timeouts.ReadTotalTimeoutMultiplier = C.MAXDWORD
	timeouts.ReadTotalTimeoutConstant = timeoutMillis
	timeouts.WriteTotalTimeoutMultiplier = 0
	timeouts.WriteTotalTimeoutConstant = timeoutMillis
	if C.SetCommTimeouts(C.HANDLE(handle), &timeouts) == 0 {
		err = fmt.Errorf("modbus: could not set device timeouts: %v", syscall.GetLastError())
		syscall.CloseHandle(handle)
		return
	}
	mb.handle = handle
	return
}

func (mb *serialTransporter) Close() (err error) {
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: closing '%v'\n", mb.Address)
	}
	if mb.handle != syscall.InvalidHandle {
		err = syscall.CloseHandle(mb.handle)
		mb.handle = syscall.InvalidHandle
	}
	return
}

// isConnected returns true if serial port has been opened.
func (mb *serialTransporter) isConnected() bool {
	return (mb.handle != 0) && (mb.handle != syscall.InvalidHandle)
}

// read reads from serial port, blocks until data received or times out.
func (mb *serialTransporter) read(b []byte) (n int, err error) {
	var done uint32
	if err = syscall.ReadFile(mb.handle, b, &done, nil); err != nil {
		err = fmt.Errorf("modbus: could not read from device: %v", err)
		return
	}
	if done == 0 {
		err = fmt.Errorf("modbus: read timed out after %s", mb.Timeout.String())
	}
	n = int(done)
	return
}

// write sends data to serial port.
func (mb *serialTransporter) write(b []byte) (n int, err error) {
	var done uint32
	if err = syscall.WriteFile(mb.handle, b, &done, nil); err != nil {
		err = fmt.Errorf("modbus: could not write to device: %v", err)
		return
	}
	n = int(done)
	return
}

func newHandle(config *serialConfig) (handle syscall.Handle, err error) {
	handle, err = syscall.CreateFile(
		syscall.StringToUTF16Ptr(config.Address),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,   // mode
		nil, // security
		syscall.OPEN_EXISTING, // create mode
		0, // attributes
		0) // templates
	if err != nil {
		err = fmt.Errorf("modbus: could not create device handle: %v", err)
		return
	}
	defer func() {
		if err != nil {
			syscall.CloseHandle(handle)
		}
	}()
	var dcb C.DCB
	dcb.BaudRate = C.DWORD(config.BaudRate)
	// Data bits
	if config.DataBits == 0 {
		dcb.ByteSize = 8
	} else {
		dcb.ByteSize = C.BYTE(config.DataBits)
	}
	// Stop bits
	switch config.StopBits {
	case 0:
		// Default is one stop bit
		fallthrough
	case 1:
		dcb.StopBits = C.ONESTOPBIT
	case 2:
		dcb.StopBits = C.TWOSTOPBITS
	default:
		err = fmt.Errorf("modbus: stop bits '%v' is not supported", config.StopBits)
		return
	}
	// Parity
	switch config.Parity {
	case "":
		// Default parity mode must be Even parity
		fallthrough
	case "E":
		dcb.Parity = C.EVENPARITY
	case "O":
		dcb.Parity = C.ODDPARITY
	case "N":
		dcb.Parity = C.NOPARITY
	default:
		err = fmt.Errorf("modbus: parity '%v' is not supported", config.Parity)
		return
	}
	if C.SetCommState(C.HANDLE(handle), &dcb) == 0 {
		err = fmt.Errorf("modbus: could not set device state: %v", syscall.GetLastError())
		return
	}
	return
}
