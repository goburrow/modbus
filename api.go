// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

type Client interface {
	// Bit access

	// ReadDiscreteInputs returns input status.
	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
	// ReadCoils returns coils status.
	ReadCoils(address, quantity uint16) (results []byte, err error)
	// WriteSingleCoil returns output value.
	WriteSingleCoil(address, value uint16) (results []byte, err error)
	// WriteMultipleCoils returns quantity of outputs.
	WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error)

	// 16-bit access

	// ReadInputRegisters returns input registers.
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters returns register value.
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	// WriteSingleRegister returns register value.
	WriteSingleRegister(address, value uint16) (results []byte, err error)
	// WriteMultipleRegisters returns quantity of registers.
	WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error)
	// ReadWriteMultipleRegisters returns read registers value.
	ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error)
	// MaskWriteRegister returns AND-mask + OR-mask.
	MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error)
	//ReadFIFOQueue returns FIFO value register.
	ReadFIFOQueue(address uint16) (results []byte, err error)
}
