// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	udpTimeout = 10 * time.Second
)

// ASCIIClientHandler implements Packager and Transporter interface.
type ASCIIUDPClientHandler struct {
	asciiPackager
	asciiUdpTransporter
}

// NewASCIIClientHandler allocates and initializes a ASCIIClientHandler.
// The address format is ip:port
func NewASCIIUDPClientHandler(address string, slaveId int) *ASCIIUDPClientHandler {
	handler := &ASCIIUDPClientHandler{}
	handler.Address = address
	handler.SlaveId = byte(slaveId)
	handler.Timeout = udpTimeout
	return handler
}

// ASCIIClient creates ASCII client with default handler and given connect string.
func ASCIIUDPClient(address string, slaveId int) Client {
	handler := NewASCIIUDPClientHandler(address, slaveId)
	return NewClient(handler)
}


func (mb *asciiUdpTransporter) logf(format string, v ...interface{}) {
	if mb.Logger != nil {
		mb.Logger.Printf(format, v...)
	}
}


// asciiUdpTransporter implements Transporter interface.
type asciiUdpTransporter struct {
	// Connect string
	Address string
	// Connect & Read timeout
	Timeout time.Duration
	// Transmission logger
	Logger *log.Logger

	// UDP "connection"
	mu           sync.Mutex
	conn         *net.UDPConn
}

func (mb *asciiUdpTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// Make sure port is connected
	if err = mb.connect(); err != nil {
		return
	}

	// Send the request
	mb.logf("modbus: sending %q\n", aduRequest)
	if _, err = mb.conn.Write(aduRequest); err != nil {
		return
	}
	// Get the response
	var length int
	data := make([]byte, asciiMaxSize)
	mb.conn.SetDeadline(time.Now().Add(mb.Timeout))
	length, _, err = mb.conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	aduResponse = data[0:length]
	mb.logf("modbus: received %q\n", aduResponse)
	return
}


func (mb *asciiUdpTransporter) connect() error {
	if mb.conn == nil {
		s, err := net.ResolveUDPAddr("udp4", mb.Address)
		conn, err := net.DialUDP("udp4", nil, s)
		if err != nil {
			return err
		}
		mb.conn = conn
	}
	return nil
}
