// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	tcpProtocolIdentifier uint16 = 0x0000

	// Modbus Application Protocol
	tcpHeaderLength = 7
	tcpMaxLength    = 260
)

type TcpClientHandler struct {
	tcpPackager
	tcpTransporter
}

func TcpClient(connectString string) Client {
	handler := &TcpClientHandler{}
	handler.ConnectString = connectString
	return TcpClientWithHandler(handler)
}

func TcpClientWithHandler(handler *TcpClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

// Implements Encoder and Decoder interface
type tcpPackager struct {
	// For synchronization between messages of server & client
	transactionId uint16
	// Broadcast address is 0
	UnitId byte
}

// Encode adds modbus application protocol header:
//  Transaction identifier: 2 bytes
//  Protocol identifier: 2 bytes
//  Length: 2 bytes
//  Unit identifier: 1 byte
func (mb *tcpPackager) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	var buf bytes.Buffer

	// Transaction identifier
	mb.transactionId++
	if err = binary.Write(&buf, binary.BigEndian, mb.transactionId); err != nil {
		return
	}
	// Protocol identifier
	if err = binary.Write(&buf, binary.BigEndian, tcpProtocolIdentifier); err != nil {
		return
	}
	// Length = sizeof(UnitId) + sizeof(FunctionCode) + Data
	length := uint16(1 + 1 + len(pdu.Data))
	if err = binary.Write(&buf, binary.BigEndian, length); err != nil {
		return
	}
	// Unit identifier
	if err = binary.Write(&buf, binary.BigEndian, mb.UnitId); err != nil {
		return
	}
	// PDU
	var n int
	if err = buf.WriteByte(pdu.FunctionCode); err != nil {
		return
	}
	if n, err = buf.Write(pdu.Data); err != nil {
		return
	}
	if n != len(pdu.Data) {
		err = fmt.Errorf("modbus: encoded pdu size '%v' does not match expected '%v'", len(pdu.Data), n)
		return
	}
	adu = buf.Bytes()
	return
}

// Verify confirm transaction, protocol and unit id
func (mb *tcpPackager) Verify(aduRequest []byte, aduResponse []byte) (err error) {
	// Transaction id
	responseVal := binary.BigEndian.Uint16(aduResponse)
	requestVal := binary.BigEndian.Uint16(aduRequest)
	if responseVal != requestVal {
		err = fmt.Errorf("modbus: adu transaction id '%v' does not match request '%v'", responseVal, requestVal)
		return
	}
	// Protocol id
	responseVal = binary.BigEndian.Uint16(aduResponse[2:])
	requestVal = binary.BigEndian.Uint16(aduRequest[2:])
	if responseVal != requestVal {
		err = fmt.Errorf("modbus: adu protocol id '%v' does not match request '%v'", responseVal, requestVal)
		return
	}
	// Unit id (1 byte)
	if aduResponse[6] != aduRequest[6] {
		err = fmt.Errorf("modbus: adu unit id '%v' does not match request '%v'", aduResponse[6], aduRequest[6])
		return
	}
	return
}

// Decode extracts PDU from TCP frame:
//  Transaction identifier: 2 bytes
//  Protocol identifier: 2 bytes
//  Length: 2 bytes
//  Unit identifier: 1 byte
func (mb *tcpPackager) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	// Read length value in the header
	length := binary.BigEndian.Uint16(adu[4:])
	pduLength := len(adu) - tcpHeaderLength
	if pduLength <= 0 || pduLength != int(length-1) {
		err = fmt.Errorf("modbus: adu length '%v' does not match pdu data '%v'", length-1, pduLength)
		return
	}
	pdu = &ProtocolDataUnit{}
	// The first byte after header is function code
	pdu.FunctionCode = adu[tcpHeaderLength]
	pdu.Data = adu[tcpHeaderLength+1:]
	return
}

// tcpTransporter implements Transporter interface
type tcpTransporter struct {
	ConnectString string
	// Connect & Read timeout
	Timeout time.Duration
	// Transmission logger
	Logger *log.Logger

	// TCP connection
	conn net.Conn
}

// Send sends data to server and ensures response length is greater than header length
func (mb *tcpTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	var data [tcpMaxLength]byte

	if mb.conn != nil {
		// Flush current data and check if the connection is alive
		if err = mb.flush(data[:]); err != nil {
			return
		}
	} else {
		// Establish a new connection and close it when complete
		if err = mb.Connect(); err != nil {
			return
		}
		defer mb.Close()
	}
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: sending %v\n", aduRequest)
	}
	if err = mb.write(aduRequest); err != nil {
		return
	}
	// Read header first
	var n int
	if n, err = mb.read(data[:tcpHeaderLength]); err != nil {
		return
	}
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: received header %v\n", data[:tcpHeaderLength])
	}
	if n != tcpHeaderLength {
		err = fmt.Errorf("modbus: response header size '%v' does not match expected '%v'", n, tcpHeaderLength)
		return
	}
	// Read length, ignore transaction & protocol id (4 bytes)
	length := int(binary.BigEndian.Uint16(data[4:]))
	if length <= 0 {
		err = fmt.Errorf("modbus: length in response header '%v' must not be zero", length)
		return
	}
	if length > (tcpMaxLength - tcpHeaderLength + 1) {
		err = fmt.Errorf("modbus: length in response header '%v' must not greater than '%v'", length, tcpMaxLength-tcpHeaderLength+1)
		return
	}
	// Skip unit id
	length = tcpHeaderLength - 1 + length
	idx := tcpHeaderLength
	for idx < length {
		if n, err = mb.read(data[idx:length]); err != nil {
			return
		}
		idx += n
	}
	aduResponse = data[:idx]
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: received %v\n", aduResponse)
	}
	return
}

// Connects to the address in ConnectString
// Connect and Close are exported so that multiple requests can be done with one session
func (mb *tcpTransporter) Connect() (err error) {
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: connecting '%v'\n", mb.ConnectString)
	}
	dialer := net.Dialer{Timeout: mb.Timeout}
	mb.conn, err = dialer.Dial("tcp", mb.ConnectString)
	return
}

// Closes current connection
func (mb *tcpTransporter) Close() (err error) {
	if mb.conn != nil {
		err = mb.conn.Close()
		mb.conn = nil
		if mb.Logger != nil {
			mb.Logger.Printf("modbus: closed connection '%v'\n", mb.ConnectString)
		}
	}
	return
}

// These methods must only be called after Connect()
func (mb *tcpTransporter) write(b []byte) (err error) {
	var n int
	if mb.Timeout > 0 {
		if err = mb.conn.SetWriteDeadline(time.Now().Add(mb.Timeout)); err != nil {
			return
		}
	}
	if n, err = mb.conn.Write(b); err != nil {
		return
	}
	// Is this checking necessary?
	if n != len(b) {
		err = fmt.Errorf("modbus: sent length '%v' does not match expected '%v'", n, len(b))
		return
	}
	return
}

func (mb *tcpTransporter) read(b []byte) (n int, err error) {
	if mb.Timeout > 0 {
		if err = mb.conn.SetReadDeadline(time.Now().Add(mb.Timeout)); err != nil {
			return
		}
	}
	n, err = mb.conn.Read(b)
	return
}

// Flushes pending data in the connection,
// returns io.EOF if connection is closed
func (mb *tcpTransporter) flush(b []byte) (err error) {
	if err = mb.conn.SetReadDeadline(time.Now()); err != nil {
		return
	}
	// Reset timeout setting
	if mb.Timeout <= 0 {
		defer mb.conn.SetReadDeadline(time.Time{})
	}
	if _, err = mb.conn.Read(b); err != nil {
		// Ignore timeout error
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			err = nil
		}
	}
	return
}
