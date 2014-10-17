// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"testing"
	"os"
	"log"
)

const (
	testTcpServer = "localhost:5020"
)

func TestReadCoils(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Read discrete outputs 20-38:
	address := uint16(0x0013)
	quantity := uint16(0x0013)
	results, err := client.ReadCoils(address, quantity)
	if err != nil {
		t.Error(err)
	}
	if 3 != len(results) {
		t.Errorf("expected: %v, actual: %v", 3, len(results))
	}
}

func TestDiscreteInputs(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Read discrete inputs 197-218
	address := uint16(0x00C4)
	quantity := uint16(0x0016)
	results, err := client.ReadDiscreteInputs(address, quantity)
	if err != nil {
		t.Error(err)
	}
	if 3 != len(results) {
		t.Errorf("expected: %v, actual: %v", 3, len(results))
	}
}

func TestReadHoldingRegisters(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Read registers 108-110
	address := uint16(0x006B)
	quantity := uint16(0x0003)
	results, err := client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		t.Error(err)
	}
	if 6 != len(results) {
		t.Errorf("expected: %v, actual: %v", 6, len(results))
	}
}

func TestReadInputRegisters(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Read input register 9
	address := uint16(0x0008)
	quantity := uint16(0x0001)
	results, err := client.ReadInputRegisters(address, quantity)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(results) {
		t.Errorf("expected: %v, actual: %v", 2, len(results))
	}
}

func TestWriteSingleCoil(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Write coil 173 ON
	address := uint16(0x00AC)
	value := uint16(0xFF00)
	results, err := client.WriteSingleCoil(address, value)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(results) {
		t.Errorf("expected: %v, actual: %v", 2, len(results))
	}
}

func TestWriteSingleRegister(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Write register 2 to 00 03 hex
	address := uint16(0x0001)
	value := uint16(0x0003)
	results, err := client.WriteSingleRegister(address, value)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(results) {
		t.Errorf("expected: %v, actual: %v", 2, len(results))
	}
}

func TestWriteMultipleCoils(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Write a series of 10 coils starting at coil 20
	address := uint16(0x0013)
	quantity := uint16(0x000A)
	values := []byte{0xCD, 0x01}
	results, err := client.WriteMultipleCoils(address, quantity, values)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(results) {
		t.Errorf("expected: %v, actual: %v", 2, len(results))
	}
}

func TestWriteMultipleRegisters(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Write two registers starting at 2 to 00 0A and 01 02 hex
	address := uint16(0x0001)
	quantity := uint16(0x0002)
	values := []byte{0x00, 0x0A, 0x01, 0x02}
	results, err := client.WriteMultipleRegisters(address, quantity, values)
	if err != nil {
		t.Error(err)
	}
	if 2 != len(results) {
		t.Errorf("expected: %v, actual: %v", 2, len(results))
	}
}

func TestMaskWriteRegisters(t *testing.T) {
	client := TcpClient(testTcpServer)
	// Mask write to register 5
	address := uint16(0x0004)
	andMask := uint16(0x00F2)
	orMask := uint16(0x0025)
	results, err := client.MaskWriteRegister(address, andMask, orMask)
	if err != nil {
		t.Error(err)
	}
	if 4 != len(results) {
		t.Errorf("expected: %v, actual: %v", 4, len(results))
	}
}

func TestReadWriteMultipleRegisters(t *testing.T) {
	client := TcpClient(testTcpServer)
	// read six registers starting at register 4, and to write three registers starting at register 15
	address := uint16(0x0003)
	quantity := uint16(0x0006)
	writeAddress := uint16(0x000E)
	writeQuantity := uint16(0x0003)
	values := []byte{0x00, 0xFF, 0x00, 0xFF, 0x00, 0xFF}
	results, err := client.ReadWriteMultipleRegisters(address, quantity, writeAddress, writeQuantity, values)
	if err != nil {
		t.Error(err)
	}
	if 12 != len(results) {
		t.Errorf("expected: %v, actual: %v", 12, len(results))
	}
}

func TestReadFIFOQueue(t *testing.T) {
	handler := &TcpClientHandler{}
	handler.Address = testTcpServer
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)

	client := TcpClientWithHandler(handler)
	// Read queue starting at the pointer register 1246
	address := uint16(0x04DE)
	results, err := client.ReadFIFOQueue(address)
	// Server not implemented
	if err != nil {
		if "modbus: exception '1', function '152'" != err.Error() {
			t.Error(err)
		}
	} else {
		if 0 == len(results) {
			t.Errorf("expected: !, actual: %v", 0, len(results))
		}
	}
}
