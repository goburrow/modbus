// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

// ModbusError implements error interface
type ModbusError struct {
	FunctionCode  int
	ExceptionCode int
}

type Modbus interface {
	// Bit access
	ReadDiscreteInputs(address, count int) (results []byte, err error)
	ReadCoils(address, count int) ([]byte, error)
	WriteSingleCoil(address, count int)
	WriteMultipleCoils(address, count int)

	// 16-bit access
	ReadInputRegisters(address, count int)
	ReadHoldingRegisters(address, count int)
	WriteSingleRegister(address, count int)
	WriteMultipleRegisters(address, count int)
	ReadWriteMultipleRegisters(address, count int)
	MaskWriteRegister(address, count int)
}
