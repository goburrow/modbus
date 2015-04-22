// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"bytes"
	"testing"
)

func TestASCIIEncoding(t *testing.T) {
	encoder := asciiPackager{}
	encoder.SlaveId = 17

	pdu := ProtocolDataUnit{}
	pdu.FunctionCode = 3
	pdu.Data = []byte{0, 107, 0, 3}

	adu, err := encoder.Encode(&pdu)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(":1103006B00037E\r\n")
	if !bytes.Equal(expected, adu) {
		t.Fatalf("adu actual: %v, expected %v", adu, expected)
	}
}

func TestASCIIDecoding(t *testing.T) {
	decoder := asciiPackager{}
	decoder.SlaveId = 247
	adu := []byte(":F7031389000A60\r\n")

	pdu, err := decoder.Decode(adu)
	if err != nil {
		t.Fatal(err)
	}

	if 3 != pdu.FunctionCode {
		t.Fatalf("Function code: expected %v, actual %v", 15, pdu.FunctionCode)
	}
	expected := []byte{0x13, 0x89, 0, 0x0A}
	if !bytes.Equal(expected, pdu.Data) {
		t.Fatalf("Data: expected %v, actual %v", expected, pdu.Data)
	}
}

func BenchmarkASCIIEncoder(b *testing.B) {
	encoder := asciiPackager{
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

func BenchmarkASCIIDecoder(b *testing.B) {
	decoder := asciiPackager{
		SlaveId: 10,
	}
	adu := []byte(":F7031389000A60\r\n")
	for i := 0; i < b.N; i++ {
		_, err := decoder.Decode(adu)
		if err != nil {
			b.Fatal(err)
		}
	}
}
