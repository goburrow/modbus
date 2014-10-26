// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

const (
	asciiStart     = ":"
	asciiEnd       = "\r\n"
	asciiMinLength = 3
	asciiMaxLength = 513

	asciiTimeoutMillis = 5000
	asciiSleepMillis   = 1000

	hexTable = "0123456789ABCDEF"
)

type ASCIIClientHandler struct {
	asciiPackager
	asciiSerialTransporter
}

func ASCIIClient(device string) Client {
	handler := &ASCIIClientHandler{}
	handler.Device = device
	return ASCIIClientWithHandler(handler)
}

func ASCIIClientWithHandler(handler *ASCIIClientHandler) Client {
	return NewClient(handler, handler)
}

// Implements Encoder and Decoder interface
type asciiPackager struct {
	SlaveId byte
}

// Encode encodes PDU in a ASCII frame:
//  Start           : 1 char
//  Address         : 2 chars
//  Function        : 2 chars
//  Data            : 0 up to 2x252 chars
//  LRC             : 2 chars
//  End             : 2 chars
func (mb *asciiPackager) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	var buf bytes.Buffer

	if _, err = buf.WriteString(asciiStart); err != nil {
		return
	}
	if err = writeHex(&buf, []byte{mb.SlaveId, pdu.FunctionCode}); err != nil {
		return
	}
	if err = writeHex(&buf, pdu.Data); err != nil {
		return
	}
	// Exclude the beginning colon and terminating CRLF pair characters
	var lrc lRC
	lrc.pushByte(mb.SlaveId).pushByte(pdu.FunctionCode).pushBytes(pdu.Data)
	if err = writeHex(&buf, []byte{lrc.value()}); err != nil {
		return
	}
	if _, err = buf.WriteString(asciiEnd); err != nil {
		return
	}
	adu = buf.Bytes()
	return
}

// Verify verifies response length, frame boundary and slave id
func (mb *asciiPackager) Verify(aduRequest []byte, aduResponse []byte) (err error) {
	length := len(aduResponse)
	// Minimum size (including address, function and LRC)
	if length < asciiMinLength+6 {
		err = fmt.Errorf("modbus: response length '%v' does not meet minimum '%v'", length, 9)
		return
	}
	// Length excluding colon must be an even number
	if length%2 != 1 {
		err = fmt.Errorf("modbus: response length '%v' is not an even number", length-1)
		return
	}
	// First char must be a colon
	str := string(aduResponse[0:len(asciiStart)])
	if str != asciiStart {
		err = fmt.Errorf("modbus: response frame '%v'... is not started with '%v'", str, asciiStart)
		return
	}
	// 2 last chars must be \r\n
	str = string(aduResponse[len(aduResponse)-len(asciiEnd):])
	if str != asciiEnd {
		err = fmt.Errorf("modbus: response frame ...'%v' is not ended with '%v'", str, asciiEnd)
		return
	}
	// Slave id
	responseVal, err := readHex(aduResponse[1:])
	if err != nil {
		return
	}
	requestVal, err := readHex(aduRequest[1:])
	if err != nil {
		return
	}
	if responseVal != requestVal {
		err = fmt.Errorf("modbus: response slave id '%v' does not match request '%v'", responseVal, requestVal)
		return
	}
	return
}

// Decode extracts PDU from ASCII frame and verify LRC
func (mb *asciiPackager) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	pdu = &ProtocolDataUnit{}
	// Slave address
	address, err := readHex(adu[1:])
	if err != nil {
		return
	}
	// Function code
	if pdu.FunctionCode, err = readHex(adu[3:]); err != nil {
		return
	}
	// Data
	dataEnd := len(adu) - 4
	data := adu[5:dataEnd]
	pdu.Data = make([]byte, hex.DecodedLen(len(data)))
	if _, err = hex.Decode(pdu.Data, data); err != nil {
		return
	}
	// LRC
	lrcVal, err := readHex(adu[dataEnd:])
	if err != nil {
		return
	}
	// Calculate checksum
	var lrc lRC
	lrc.pushByte(address).pushByte(pdu.FunctionCode).pushBytes(pdu.Data)
	if lrcVal != lrc.value() {
		err = fmt.Errorf("modbus: response lrc '%v' does not match expected '%v'", lrcVal, lrc.value())
		return
	}
	return
}

// asciiSerialTransporter implements Transporter interface
type asciiSerialTransporter struct {
	// Serial port configuration
	serialConfig
	// Read timeout
	Timeout time.Duration
	Logger  *log.Logger

	// Serial controller
	serial serial
}

func (mb *asciiSerialTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	if mb.serial.IsConnected() {
		// flush current data pending in serial port
	} else {
		if err = mb.Connect(); err != nil {
			return
		}
	}
	if mb.Logger != nil {
		mb.Logger.Printf("modbus: sending %v\n", aduRequest)
	}
	var n int
	if n, err = mb.serial.Write(aduRequest); err != nil {
		return
	}
	// TODO: use channel to handle timeout
	// TODO: use syscall.select to wait for data arrived
	var data [asciiMaxLength]byte
	length := 0
	deadline := time.Now().Add(mb.Timeout)

	for {
		if n, err = mb.serial.Read(data[length:]); err != nil {
			return
		}
		length += n
		if length >= asciiMaxLength {
			break
		}
		// Expect end of frame in the data received
		if length > asciiMinLength {
			if string(data[length-len(asciiEnd):]) == asciiEnd {
				break
			}
		}
		if time.Now().After(deadline) {
			err = fmt.Errorf("modbus: Read timeout after %s", mb.Timeout.String())
			return
		}
		time.Sleep(asciiSleepMillis * time.Millisecond)
	}
	aduResponse = data[:length]
	return
}

func (mb *asciiSerialTransporter) Connect() (err error) {
	// timeout is required
	if mb.Timeout <= 0 {
		mb.Timeout = asciiTimeoutMillis * time.Millisecond
	}
	err = mb.serial.Connect(&mb.serialConfig)
	return
}

func (mb *asciiSerialTransporter) Close() (err error) {
	err = mb.serial.Close()
	return
}

// writeHex encodes byte to string in hexadecimal, e.g. 0xA5 => "A5"
// (encoding/hex only supports lowercase string)
func writeHex(buf *bytes.Buffer, value []byte) (err error) {
	var str [2]byte
	for _, v := range value {
		str[0] = hexTable[v>>4]
		str[1] = hexTable[v&0x0F]

		if _, err = buf.Write(str[:]); err != nil {
			return
		}
	}
	return
}

// readHex decodes hexa string to byte, e.g. "8C" => 0x8C
func readHex(data []byte) (value byte, err error) {
	var dst [1]byte
	if _, err = hex.Decode(dst[:], data[0:2]); err != nil {
		return
	}
	value = dst[0]
	return
}
