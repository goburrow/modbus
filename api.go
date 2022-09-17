// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

type Client interface {
	// Bit access

	// ReadCoils reads from 1 to 2000 contiguous status of coils in a
	// remote device and returns coil status.
	ReadCoils(slaveId uint8, address, quantity uint16) (results []byte, err error)
	// ReadDiscreteInputs reads from 1 to 2000 contiguous status of
	// discrete inputs in a remote device and returns input status.
	ReadDiscreteInputs(slaveId uint8, address, quantity uint16) (results []byte, err error)
	// WriteSingleCoil write a single output to either ON or OFF in a
	// remote device and returns output value.
	WriteSingleCoil(slaveId uint8, address, value uint16) (results []byte, err error)
	// WriteMultipleCoils forces each coil in a sequence of coils to either
	// ON or OFF in a remote device and returns quantity of outputs.
	WriteMultipleCoils(slaveId uint8, address, quantity uint16, value []byte) (results []byte, err error)

	// 16-bit access

	// ReadInputRegisters reads from 1 to 125 contiguous input registers in
	// a remote device and returns input registers.
	ReadInputRegisters(slaveId uint8, address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote device and returns register value.
	ReadHoldingRegisters(slaveId uint8, address, quantity uint16) (results []byte, err error)
	// WriteSingleRegister writes a single holding register in a remote
	// device and returns register value.
	WriteSingleRegister(slaveId uint8, address, value uint16) (results []byte, err error)
	// WriteMultipleRegisters writes a block of contiguous registers
	// (1 to 123 registers) in a remote device and returns quantity of
	// registers.
	WriteMultipleRegisters(slaveId uint8, address, quantity uint16, value []byte) (results []byte, err error)
	// ReadWriteMultipleRegisters performs a combination of one read
	// operation and one write operation. It returns read registers value.
	ReadWriteMultipleRegisters(slaveId uint8, readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error)
	// MaskWriteRegister modify the contents of a specified holding
	// register using a combination of an AND mask, an OR mask, and the
	// register's current contents. The function returns
	// AND-mask and OR-mask.
	MaskWriteRegister(slaveId uint8, address, andMask, orMask uint16) (results []byte, err error)
	//ReadFIFOQueue reads the contents of a First-In-First-Out (FIFO) queue
	// of register in a remote device and returns FIFO value register.
	ReadFIFOQueue(slaveId uint8, address uint16) (results []byte, err error)
}
