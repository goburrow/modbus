// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package test

import (
	"github.com/goburrow/modbus"
	"log"
	"os"
	"testing"
	"time"
)

const (
	testTCPServer = "localhost:5020"
)

func TestTCPClientReadCoils(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestReadCoils(t, client)
}

func TestTCPClientReadDiscreteInputs(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	// Read discrete inputs 197-218
	ClientTestDiscreteInputs(t, client)
}

func TestTCPClientReadHoldingRegisters(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestReadHoldingRegisters(t, client)
}

func TestTCPClientReadInputRegisters(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestReadInputRegisters(t, client)
}

func TestTCPClientWriteSingleCoil(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestWriteSingleCoil(t, client)
}

func TestTCPClientWriteSingleRegister(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestWriteSingleRegister(t, client)
}

func TestTCPClientWriteMultipleCoils(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestWriteMultipleCoils(t, client)
}

func TestTCPClientWriteMultipleRegisters(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestWriteMultipleRegisters(t, client)
}

func TestTCPClientMaskWriteRegisters(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestMaskWriteRegisters(t, client)
}

func TestTCPClientReadWriteMultipleRegisters(t *testing.T) {
	client := modbus.TCPClient(testTCPServer)
	ClientTestReadWriteMultipleRegisters(t, client)
}

func TestTCPClientReadFIFOQueue(t *testing.T) {
	handler := &modbus.TCPClientHandler{}
	handler.Address = testTCPServer
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)

	client := modbus.TCPClientWithHandler(handler)
	ClientTestReadFIFOQueue(t, client)
}

func TestTCPClientAdvancedUsage(t *testing.T) {
	var handler modbus.TCPClientHandler
	handler.Address = testTCPServer
	handler.Timeout = 5 * time.Second
	handler.SlaveId = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.TCPClientWithHandler(&handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}
