// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	TcpProtocolIdentifier uint16 = 0x0000
	TcpUnitIdentifier     byte   = 0xFF

	// Modbus Application Protocol
	TcpHeaderLength = 7
	TcpMaxADULength = 260
)

func TcpClient(address string) Client {
	encodeDecoder := &TcpEncodeDecoder{}
	transporter := &TcpTransporter{address: address}

	return &client{encoder: encodeDecoder, decoder: encodeDecoder, transporter: transporter}
}

// Implements Encoder and Decoder interface
type TcpEncodeDecoder struct {
	// For synchronization between messages of server & client
	// TODO put in a context for the sake of thread-safe
	transactionId uint16
	unitId        byte
}

// Adds modbus application protocol header:
//  Transaction identifier: 2 bytes
//  Protocol identifier: 2 bytes
//  Length: 2 bytes
//  Unit identifier: 1 byte
func (mb *TcpEncodeDecoder) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	var buf bytes.Buffer

	// Transaction identifier
	mb.transactionId++
	if err = binary.Write(&buf, binary.BigEndian, mb.transactionId); err != nil {
		return
	}
	// Protocol identifier
	if err = binary.Write(&buf, binary.BigEndian, TcpProtocolIdentifier); err != nil {
		return
	}
	// Length = sizeof(UnitId) + sizeof(FunctionCode) + Data
	length := uint16(1 + 1 + len(pdu.Data))
	if err = binary.Write(&buf, binary.BigEndian, length); err != nil {
		return
	}
	// Unit identifier
	if err = binary.Write(&buf, binary.BigEndian, mb.unitId); err != nil {
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
		err = ErrResponseSize
		return
	}
	adu = buf.Bytes()
	return
}

func (mb *TcpEncodeDecoder) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	var (
		transactionId uint16
		protocolId    uint16
		length        uint16
		unitId        uint8
	)

	buf := bytes.NewReader(adu)
	if err = binary.Read(buf, binary.BigEndian, &transactionId); err != nil {
		return
	}
	// Not thread safe yet
	if transactionId != mb.transactionId {
		err = ErrTransactionId
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &protocolId); err != nil {
		return
	}
	if protocolId != TcpProtocolIdentifier {
		err = ErrTransactionId
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &length); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &unitId); err != nil {
		return
	}
	if unitId != mb.unitId {
		err = ErrTransactionId
		return
	}
	pduLength := buf.Len()
	if pduLength == 0 || pduLength != int(length-1) {
		err = ErrResponseSize
		return
	}
	pdu = &ProtocolDataUnit{}
	if err = binary.Read(buf, binary.BigEndian, &pdu.FunctionCode); err != nil {
		return
	}
	pdu.Data = make([]byte, pduLength-1)
	var n int
	if n, err = buf.Read(pdu.Data); err != nil {
		return
	}
	if n != pduLength-1 {
		err = ErrResponseSize
		return
	}
	return
}

// Implements Transporter interface
type TcpTransporter struct {
	address string
	timeout time.Duration
}

func (mb *TcpTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	dialer := net.Dialer{Timeout: mb.timeout}
	conn, err := dialer.Dial("tcp", mb.address)
	if err != nil {
		return
	}
	defer conn.Close()

	var n int
	if n, err = conn.Write(aduRequest); err != nil {
		return
	}
	if n != len(aduRequest) {
		err = ErrResponseSize
		// TODO: flush
		return
	}
	fmt.Printf("%v\n", aduRequest)
	// Read header first
	data := [TcpMaxADULength]byte{}
	if n, err = conn.Read(data[:TcpHeaderLength]); err != nil {
		fmt.Printf("%v %v %v\n", data, n, err)
		return
	}
	if n != TcpHeaderLength {
		err = ErrResponseSize
		return
	}
	// Read length, ignore transaction & protocol id (4 bytes)
	length := int(binary.BigEndian.Uint16(data[4:]))
	if length <= 0 {
		err = ErrResponse
		return
	}
	// Skip unit id
	length = TcpHeaderLength - 1 + length
	idx := TcpHeaderLength
	for idx < length {
		if n, err = conn.Read(data[idx:length]); err != nil {
			return
		}
		idx += n
	}
	aduResponse = data[:idx]
	return
}
