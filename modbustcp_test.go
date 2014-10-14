package modbus

import (
	"testing"
)

func TestReadCoils(t *testing.T) {
	address := "localhost:5020"
	client := ModbusTcpClient(address)
	_, err := client.ReadCoils(100, 1)
	if err != nil {
		t.Error(err)
	}
}
