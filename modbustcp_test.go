package modbus

import (
	"testing"
)

const (
	server = "localhost:5020"
)

var (
	address uint16
	quantity uint16
)

func TestReadCoils(t *testing.T) {
	client := ModbusTcpClient(server)
	address = 100
	quantity = 9
	_, err := client.ReadCoils(address, quantity)
	if err != nil {
		t.Error(err)
	}
}

func TestDiscreteInputs(t *testing.T) {
	client := ModbusTcpClient(server)
	address = 100
	quantity = 9
	_, err := client.ReadDiscreteInputs(address, quantity)
	if err != nil {
		t.Error(err)
	}
}

func TestReadHoldingRegisters(t *testing.T) {
	client := ModbusTcpClient(server)
	address = 100
	quantity = 9
	_, err := client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		t.Error(err)
	}
}

func TestReadInputRegisters(t *testing.T) {
	client := ModbusTcpClient(server)
	address = 100
	quantity = 9
	_, err := client.ReadInputRegisters(address, quantity)
	if err != nil {
		t.Error(err)
	}
}
