// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package test

import (
	"github.com/goburrow/modbus"
	"log"
	"os"
	"testing"
)

const (
	testTcpServer = "localhost:5020"
)

func TestTcpClientReadCoils(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestReadCoils(t, client)
}

func TestTcpClientReadDiscreteInputs(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	// Read discrete inputs 197-218
	ClientTestDiscreteInputs(t, client)
}

func TestTcpClientReadHoldingRegisters(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestReadHoldingRegisters(t, client)
}

func TestTcpClientReadInputRegisters(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestReadInputRegisters(t, client)
}

func TestTcpClientWriteSingleCoil(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestWriteSingleCoil(t, client)
}

func TestTcpClientWriteSingleRegister(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestWriteSingleRegister(t, client)
}

func TestTcpClientWriteMultipleCoils(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestWriteMultipleCoils(t, client)
}

func TestTcpClientWriteMultipleRegisters(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestWriteMultipleRegisters(t, client)
}

func TestTcpClientMaskWriteRegisters(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestMaskWriteRegisters(t, client)
}

func TestTcpClientReadWriteMultipleRegisters(t *testing.T) {
	client := modbus.TcpClient(testTcpServer)
	ClientTestReadWriteMultipleRegisters(t, client)
}

func TestTcpClientReadFIFOQueue(t *testing.T) {
	handler := &modbus.TcpClientHandler{}
	handler.ConnectString = testTcpServer
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)

	client := modbus.TcpClientWithHandler(handler)
	ClientTestReadFIFOQueue(t, client)
}

func TestTcpClientAdvancedUsage(t *testing.T) {
	var handler modbus.TcpClientHandler
	handler.ConnectString = testTcpServer
	handler.UnitId = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.TcpClientWithHandler(&handler)
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
