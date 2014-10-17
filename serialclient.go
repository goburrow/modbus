// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const (
	ASCIIStart = ":"
	ASCIIEnd   = "\r\n"

	hexTable = "0123456789ABCDEF"
)

// Implements Encoder and Decoder interface
type ASCIIEncodeDecoder struct {
	SlaveId byte
}

// ASCII framing:
//  Start           : 1 char
//  Address         : 2 chars
//  Function        : 2 chars
//  Data            : 0 up to 2x252 chars
//  LRC             : 2 chars
//  End             : 2 chars
func (mb *ASCIIEncodeDecoder) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	var buf bytes.Buffer

	if _, err = buf.WriteString(ASCIIStart); err != nil {
		return
	}
	if err = writeHex(&buf, []byte{mb.SlaveId, pdu.FunctionCode}); err != nil {
		return
	}
	if err = writeHex(&buf, pdu.Data); err != nil {
		return
	}
	// Exclude the beginning colon and terminating CRLF pair characters
	var lrc LRC
	lrc.pushByte(mb.SlaveId).pushByte(pdu.FunctionCode).pushBytes(pdu.Data)
	if err = writeHex(&buf, []byte{lrc.value()}); err != nil {
		return
	}
	if _, err = buf.WriteString(ASCIIEnd); err != nil {
		return
	}
	adu = buf.Bytes()
	return
}

func (mb *ASCIIEncodeDecoder) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	var (
		value byte
	)
	length := len(adu)
	// Minimum size
	if length < 9 {
		err = fmt.Errorf("modbus: adu size '%v' does not meet minimum '%v'", length, 9)
		return
	}
	// Length including colon
	if length%2 != 1 {
		err = fmt.Errorf("modbus: adu size '%v' is not an even number", length-1)
		return
	}
	// Skip first colon, read address
	if value, err = readHex(adu[1:]); err != nil {
		return
	}
	if value != mb.SlaveId {
		err = fmt.Errorf("modbus: adu slave id '%v' does not match request '%v'", value, mb.SlaveId)
		return
	}
	// Function code
	pdu = &ProtocolDataUnit{}
	if pdu.FunctionCode, err = readHex(adu[3:]); err != nil {
		return
	}
	// Data
	data := adu[5 : length-4]
	pdu.Data = make([]byte, hex.DecodedLen(len(data)))
	if _, err = hex.Decode(pdu.Data, data); err != nil {
		return
	}
	fmt.Printf("%v", data)

	// LRC
	if value, err = readHex(adu[length-4:]); err != nil {
		return
	}
	var lrc LRC
	lrc.pushByte(mb.SlaveId).pushByte(pdu.FunctionCode).pushBytes(pdu.Data)
	if value != lrc.value() {
		err = fmt.Errorf("modbus: adu lrc '%v' does not match expected '%v'", value, lrc.value())
		return
	}
	return
}

// encoding/hex only supports lowercase string
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

func readHex(data []byte) (value byte, err error) {
	var dst [1]byte
	if _, err = hex.Decode(dst[:], data[0:2]); err != nil {
		return
	}
	value = dst[0]
	return
}
