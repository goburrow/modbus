// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

// ModbusError implements error interface
type ModbusError struct {
	FunctionCode  byte
	ExceptionCode byte
}

type ModbusClient interface {
	// Bit access
	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
	ReadCoils(address, quantity uint16) (results []byte, err error)
	WriteSingleCoil(address, count int)
	WriteMultipleCoils(address, count int)

	// 16-bit access
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	WriteSingleRegister(address, count int)
	WriteMultipleRegisters(address, count int)
	ReadWriteMultipleRegisters(address, count int)
	MaskWriteRegister(address, count int)
}
