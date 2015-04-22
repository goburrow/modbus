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

var responseLengthTests = []struct {
	adu    []byte
	length int
}{
	{[]byte{4, 1, 0, 0xA, 0, 0xD, 0xDD, 0x98}, 7},
	{[]byte{4, 2, 0, 0xA, 0, 0xD, 0x99, 0x98}, 7},
	{[]byte{1, 3, 0, 0, 0, 2, 0xC4, 0xB}, 9},
	{[]byte{0x11, 5, 0, 0xAC, 0xFF, 0, 0x4E, 0x8B}, 8},
	{[]byte{0x11, 6, 0, 1, 0, 3, 0x9A, 0x9B}, 8},
	{[]byte{0x11, 0xF, 0, 0x13, 0, 0xA, 2, 0xCD, 1, 0xBF, 0xB}, 8},
	{[]byte{0x11, 0x10, 0, 1, 0, 2, 4, 0, 0xA, 1, 2, 0xC6, 0xF0}, 8},
}

func TestCalculateResponseLength(t *testing.T) {
	for _, input := range responseLengthTests {
		output := calculateResponseLength(input.adu)
		if output != input.length {
			t.Errorf("Response length of %x: expected %v, actual: %v",
				input.adu, input.length, output)
		}
	}
}

func BenchmarkRTUEncoder(b *testing.B) {
	encoder := rtuPackager{
		SlaveId: 10,
	}
	pdu := ProtocolDataUnit{
		FunctionCode: 1,
		Data:         []byte{2, 3, 4, 5, 6, 7, 8, 9},
	}
	for i := 0; i < b.N; i++ {
		_, err := encoder.Encode(&pdu)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRTUDecoder(b *testing.B) {
	decoder := rtuPackager{
		SlaveId: 10,
	}
	adu := []byte{0x01, 0x10, 0x8A, 0x00, 0x00, 0x03, 0xAA, 0x10}
	for i := 0; i < b.N; i++ {
		_, err := decoder.Decode(adu)
		if err != nil {
			b.Fatal(err)
		}
	}
}
