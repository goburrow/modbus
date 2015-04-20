// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var baudRates = map[int]uint32{
	50:      syscall.B50,
	75:      syscall.B75,
	110:     syscall.B110,
	134:     syscall.B134,
	150:     syscall.B150,
	200:     syscall.B200,
	300:     syscall.B300,
	600:     syscall.B600,
	1200:    syscall.B1200,
	1800:    syscall.B1800,
	2400:    syscall.B2400,
	4800:    syscall.B4800,
	9600:    syscall.B9600,
	19200:   syscall.B19200,
	38400:   syscall.B38400,
	57600:   syscall.B57600,
	115200:  syscall.B115200,
	230400:  syscall.B230400,
	460800:  syscall.B460800,
	500000:  syscall.B500000,
	576000:  syscall.B576000,
	921600:  syscall.B921600,
	1000000: syscall.B1000000,
	1152000: syscall.B1152000,
	1500000: syscall.B1500000,
	2000000: syscall.B2000000,
	2500000: syscall.B2500000,
	3000000: syscall.B3000000,
	3500000: syscall.B3500000,
	4000000: syscall.B4000000,
}

var charSizes = map[int]uint32{
	5: syscall.CS5,
	6: syscall.CS6,
	7: syscall.CS7,
	8: syscall.CS8,
}

// serialPort is included in serialTransporter.
// The reason of using embedded struct rather than interface is to utilise
// Logger and Timeout property without passing their reference.
type serialPort struct {
	// Should use fd directly by using syscall.Open() ?
	file       *os.File
	oldTermios *syscall.Termios
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
	termios, err := newTermios(&mb.serialConfig)
	if err != nil {
		return
	}
	// See man termios(3)
	// O_NOCTTY: no controlling terminal
	// O_NDELAY: no data carrier detect
	mb.file, err = os.OpenFile(mb.Address, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, os.FileMode(0666))
	if err != nil {
		return
	}
	// Backup current termios
	mb.backupTermios()
	if err = tcsetattr(int(mb.file.Fd()), termios); err != nil {
		mb.file.Close()
		mb.file = nil
		mb.oldTermios = nil
		return
	}
	return
}

func (mb *serialTransporter) Close() (err error) {
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: closing '%v'\n", mb.Address)
	}
	if mb.file != nil {
		mb.restoreTermios()
		err = mb.file.Close()
		mb.file = nil
	}
	return
}

// isConnected returns true if serial port has been opened.
func (mb *serialTransporter) isConnected() bool {
	return mb.file != nil
}

// Read reads from serial port, blocked until data received or timeout after Timeout.
func (mb *serialTransporter) read(b []byte) (n int, err error) {
	var rfds syscall.FdSet

	fd := int(mb.file.Fd())
	fdSet(fd, &rfds)

	timeout := syscall.NsecToTimeval(mb.Timeout.Nanoseconds())
	time.Sleep(300 * time.Millisecond)
	if _, err = syscall.Select(fd+1, &rfds, nil, nil, &timeout); err != nil {
		return
	}
	if fdIsSet(fd, &rfds) {
		n, err = mb.file.Read(b)
		return
	}
	// Timeout
	err = fmt.Errorf("modbus: read timeout after %s", mb.Timeout.String())
	return
}

func (mb *serialTransporter) write(b []byte) (n int, err error) {
	n, err = mb.file.Write(b)
	return
}

// backupTermios saves current termios setting.
// Make sure that device file has been opened before calling this function.
func (mb *serialTransporter) backupTermios() {
	oldTermios := &syscall.Termios{}
	if err := tcgetattr(int(mb.file.Fd()), oldTermios); err != nil {
		// Warning only
		if mb.Logger != nil {
			log.Printf("modbus: Could not get current termios setting '%v'\n", err)
		}
	} else {
		// Will be reloaded when closing
		mb.oldTermios = oldTermios
	}
}

// restoreTermios restores backed up termios setting.
// Make sure that device file has been opened before calling this function.
func (mb *serialTransporter) restoreTermios() {
	if mb.oldTermios == nil {
		return
	}
	if err := tcsetattr(int(mb.file.Fd()), mb.oldTermios); err != nil {
		// Warning only
		if mb.Logger != nil {
			mb.Logger.Printf("modbus: Could not restore termios setting '%v'\n", err)
		}
	} else {
		mb.oldTermios = nil
	}
}

// Helpers for termios

func newTermios(config *serialConfig) (termios *syscall.Termios, err error) {
	termios = &syscall.Termios{}
	var flag uint32
	// Baud rate
	if config.BaudRate == 0 {
		// 19200 is the required default
		flag = syscall.B19200
	} else {
		flag = baudRates[config.BaudRate]
		if flag == 0 {
			err = fmt.Errorf("modbus: Baud rate '%v' is not supported", config.BaudRate)
			return
		}
	}
	// Input baud
	termios.Ispeed = flag
	// Output baud
	termios.Ospeed = flag
	// Character size
	if config.DataBits == 0 {
		flag = syscall.CS8
	} else {
		flag = charSizes[config.DataBits]
		if flag == 0 {
			err = fmt.Errorf("modbus: Character size '%v' is not supported", config.DataBits)
			return
		}
	}
	termios.Cflag |= flag
	// Stop bits
	switch config.StopBits {
	case 0:
		// Default is one stop bit
		fallthrough
	case 1:
		// noop
	case 2:
		// CSTOPB: Set two stop bits
		termios.Cflag |= syscall.CSTOPB
	default:
		err = fmt.Errorf("modbus: Stop bits '%v' is not supported", config.StopBits)
		return
	}
	switch config.Parity {
	case "N":
		// noop
	case "O":
		// PARODD: Parity is odd
		termios.Cflag |= syscall.PARODD
		fallthrough
	case "":
		// As mentioned in the spec, the default parity mode must be Even parity
		fallthrough
	case "E":
		// PARENB: Enable parity generation on output
		termios.Cflag |= syscall.PARENB
		// INPCK: Enable input parity checking
		termios.Iflag |= syscall.INPCK
	default:
		err = fmt.Errorf("modbus: Parity '%v' is not supported", config.Parity)
		return
	}
	// Control modes
	// CREAD: Enable receiver
	// CLOCAL: Ignore control lines
	termios.Cflag |= syscall.CREAD | syscall.CLOCAL
	// Special characters
	// VMIN: Minimum number of characters for noncanonical read
	// VTIME: Time in deciseconds for noncanonical read
	// Both unused as NDELAY is we utilized to open device
	return
}

// tcsetattr sets terminal file descriptor parameters.
// See man tcsetattr(3).
func tcsetattr(fd int, termios *syscall.Termios) (err error) {
	r, _, errno := syscall.Syscall(uintptr(syscall.SYS_IOCTL),
		uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		err = errno
		return
	}
	if r != 0 {
		err = fmt.Errorf("modbus: tcsetattr failed '%v'", r)
	}
	return
}

// tcgetattr gets terminal file descriptor parameters.
// See man tcgetattr(3).
func tcgetattr(fd int, termios *syscall.Termios) (err error) {
	r, _, errno := syscall.Syscall(uintptr(syscall.SYS_IOCTL),
		uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		err = errno
		return
	}
	if r != 0 {
		err = fmt.Errorf("modbus: tcgetattr failed '%v'", r)
	}
	return
}

// fdGet returns index and offset of fd in fds
func fdGet(fd int, fds *syscall.FdSet) (index, offset int) {
	index = fd / (syscall.FD_SETSIZE / len(fds.Bits)) % len(fds.Bits)
	offset = fd % (syscall.FD_SETSIZE / len(fds.Bits))
	return
}

// fdSet implements FD_SET macro
func fdSet(fd int, fds *syscall.FdSet) {
	idx, pos := fdGet(fd, fds)
	fds.Bits[idx] = 1 << uint(pos)
}

// fdIsSet implements FD_ISSET macro
func fdIsSet(fd int, fds *syscall.FdSet) bool {
	idx, pos := fdGet(fd, fds)
	return fds.Bits[idx]&(1<<uint(pos)) != 0
}
