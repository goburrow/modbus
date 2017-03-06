// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func TestTCPEncoding(t *testing.T) {
	packager := tcpPackager{}
	pdu := ProtocolDataUnit{}
	pdu.FunctionCode = 3
	pdu.Data = []byte{0, 4, 0, 3}

	adu, err := packager.Encode(&pdu)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{0, 1, 0, 0, 0, 6, 0, 3, 0, 4, 0, 3}
	if !bytes.Equal(expected, adu) {
		t.Fatalf("Expected %v, actual %v", expected, adu)
	}
}

func TestTCPDecoding(t *testing.T) {
	packager := tcpPackager{}
	packager.transactionId = 1
	packager.SlaveId = 17
	adu := []byte{0, 1, 0, 0, 0, 6, 17, 3, 0, 120, 0, 3}

	pdu, err := packager.Decode(adu)
	if err != nil {
		t.Fatal(err)
	}

	if 3 != pdu.FunctionCode {
		t.Fatalf("Function code: expected %v, actual %v", 3, pdu.FunctionCode)
	}
	expected := []byte{0, 120, 0, 3}
	if !bytes.Equal(expected, pdu.Data) {
		t.Fatalf("Data: expected %v, actual %v", expected, adu)
	}
}

func TestTCPTransporter(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		_, err = io.Copy(conn, conn)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	client := &tcpTransporter{
		Address:     ln.Addr().String(),
		Timeout:     1 * time.Second,
		IdleTimeout: 100 * time.Millisecond,
	}
	req := []byte{0, 1, 0, 2, 0, 2, 1, 2}
	rsp, err := client.Send(req)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(req, rsp) {
		t.Fatalf("unexpected response: %x", rsp)
	}
	time.Sleep(150 * time.Millisecond)
	if client.conn != nil {
		t.Fatalf("connection is not closed: %+v", client.conn)
	}
}

func BenchmarkTCPEncoder(b *testing.B) {
	encoder := tcpPackager{
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

func BenchmarkTCPDecoder(b *testing.B) {
	decoder := tcpPackager{
		SlaveId: 10,
	}
	adu := []byte{0, 1, 0, 0, 0, 6, 17, 3, 0, 120, 0, 3}
	for i := 0; i < b.N; i++ {
		_, err := decoder.Decode(adu)
		if err != nil {
			b.Fatal(err)
		}
	}
}
