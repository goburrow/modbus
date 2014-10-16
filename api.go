// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

// ModbusError implements error interface
type ModbusError struct {
	FunctionCode  byte
	ExceptionCode byte
}

type Client interface {
	// Bit access

	// Returns input status
	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
	// Returns coils status
	ReadCoils(address, quantity uint16) (results []byte, err error)
	// Returns output value
	WriteSingleCoil(address, value uint16) (results []byte, err error)
	// Returns quantity of outputs
	WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error)

	// 16-bit access

	// Returns input registers
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	// Returns register value
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	// Return register value
	WriteSingleRegister(address, value uint16) (results []byte, err error)
	// Returns quantity of registers
	WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error)
	// Returns read registers value
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error)
	// Returns AND-mask + OR-mask
	MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error)
}
