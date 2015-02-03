// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"bytes"
	"testing"
)

func TestRTUEncoding(t *testing.T) {
	encoder := rtuPackager{}
	encoder.SlaveId = 0x01

	pdu := ProtocolDataUnit{}
	pdu.FunctionCode = 0x03
	pdu.Data = []byte{0x50, 0x00, 0x00, 0x18}

	adu, err := encoder.Encode(&pdu)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte{0x01, 0x03, 0x50, 0x00, 0x00, 0x18, 0x54, 0xC0}
	if !bytes.Equal(expected, adu) {
		t.Fatalf("adu: expected %v, actual %v", expected, adu)
	}
}

func TestRTUDecoding(t *testing.T) {
	decoder := rtuPackager{}
	adu := []byte{0x01, 0x10, 0x8A, 0x00, 0x00, 0x03, 0xAA, 0x10}

	pdu, err := decoder.Decode(adu)
	if err != nil {
		t.Fatal(err)
	}

	if 16 != pdu.FunctionCode {
		t.Fatalf("Function code: expected %v, actual %v", 16, pdu.FunctionCode)
	}
	expected := []byte{0x8A, 0x00, 0x00, 0x03}
	if !bytes.Equal(expected, pdu.Data) {
		t.Fatalf("Data: expected %v, actual %v", expected, pdu.Data)
	}
}
