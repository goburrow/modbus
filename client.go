// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"encoding/binary"
	"fmt"
)

const (
	ExceptionCodePDUMaxRequestedQuantity  = 1
	ExceptionCodePDUWrongResponseDataSize = 2
	ExceptionCodePDUWrongResponseAddress  = 3
	ExceptionCodePDUWrongResponseValue    = 4
	ExceptionCodePDUWrongANDMask          = 5
	ExceptionCodePDUWrongORMask           = 6
	ExceptionCodePDUFifoGreater           = 7
	ExceptionCodePDUEmptyResponse         = 8
	ExceptionCodePDUWrongDevIDResponse    = 9
	ExceptionCodePDUWrongValueSet         = 10
)

// ModbusPDUError implements error interface for PDU data
type ModbusPDUError struct {
	ExceptionCode byte
	Request       interface{}
	Response      interface{}
}

// Error converts known modbus TCP exception code to error message.
func (e *ModbusPDUError) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodePDUMaxRequestedQuantity:
		name = "quantity '%v' must be between 1 and '%v'"
	case ExceptionCodePDUWrongResponseDataSize:
		name = "response data size '%v' does not match count '%v'"
	case ExceptionCodePDUWrongResponseAddress:
		name = "response address '%v' does not match request '%v'"
	case ExceptionCodePDUWrongResponseValue:
		name = "response value '%v' does not match request '%v'"
	case ExceptionCodePDUWrongANDMask:
		name = "response AND-mask '%v' does not match request '%v'"
	case ExceptionCodePDUWrongORMask:
		name = "response OR-mask '%v' does not match request '%v'"
	case ExceptionCodePDUFifoGreater:
		name = "fifo count '%v' is greater than expected '%v'"
	case ExceptionCodePDUEmptyResponse:
		return "response data is empty"
	case ExceptionCodePDUWrongDevIDResponse:
		return fmt.Sprintf("Read Devce ID should response minimum 14 bytes: %v received", e.Request)
	case ExceptionCodePDUWrongValueSet:
		name = "state '%v' must be either 0xFF00 (ON) or 0x0000 (OFF)"
	default:
		name = "unknown"
	}
	return fmt.Sprintf(name, e.Request, e.Response)
}

// ClientHandler is the interface that groups the Packager and Transporter methods.
type ClientHandler interface {
	Packager
	Transporter
}

type client struct {
	packager    Packager
	transporter Transporter
}

// NewClient creates a new modbus client with given backend handler.
func NewClient(handler ClientHandler) Client {
	return &client{packager: handler, transporter: handler}
}

// NewClient2 creates a new modbus client with given backend packager and transporter.
func NewClient2(packager Packager, transporter Transporter) Client {
	return &client{packager: packager, transporter: transporter}
}

// Request:
//
//	Function code         : 1 byte (0x01)
//	Starting address      : 2 bytes
//	Quantity of coils     : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x01)
//	Byte count            : 1 byte
//	Coil status           : N* bytes (=N or N+1)
func (mb *client) ReadCoils(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      2000,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadCoils,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       length,
			Response:      count,
		}
	}
	results = response.Data[1:]
	return
}

// Request:
//
//	Function code         : 1 byte (0x02)
//	Starting address      : 2 bytes
//	Quantity of inputs    : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x02)
//	Byte count            : 1 byte
//	Input status          : N* bytes (=N or N+1)
func (mb *client) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 2000 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      2000,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadDiscreteInputs,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       length,
			Response:      count,
		}
	}
	results = response.Data[1:]
	return
}

// Request:
//
//	Function code         : 1 byte (0x03)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x03)
//	Byte count            : 1 byte
//	Register value        : Nx2 bytes
func (mb *client) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 125 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      125,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadHoldingRegisters,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		return response.Data[1:], &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       count,
			Response:      length,
		}
	}
	results = response.Data[1:]
	return
}

// Request:
//
//	Function code         : 1 byte (0x04)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x04)
//	Byte count            : 1 byte
//	Input registers       : N bytes
func (mb *client) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	if quantity < 1 || quantity > 125 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      125,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadInputRegisters,
		Data:         dataBlock(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	length := len(response.Data) - 1
	if count != length {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       length,
			Response:      count,
		}
	}
	results = response.Data[1:]
	return
}

// Request:
//
//	Device ID Code		  : 1 byte
//	Object ID		      : 1 byte
//
// Response:
//
//	Objects number        : 1 byte
//	Objects data		  : N bytes
func (mb *client) ReadDeviceIdentification(devIdCode byte, objectId byte) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeDevId,
		Data:         []byte{MEITypeDevId, devIdCode, objectId},
	}

	aduRequest, err := mb.packager.Encode(&request)
	if err != nil {
		return
	}

	aduResponse, err := mb.transporter.Send(aduRequest)
	if err != nil {
		return
	}

	length := len(aduResponse)
	if length < 14 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongDevIDResponse,
			Request:       length,
		}
	}

	rxLen := binary.BigEndian.Uint16(aduResponse[4:6])
	length -= 6

	if int(rxLen) != length {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       length,
			Response:      rxLen,
		}
	}

	results = aduResponse[13 : 14+rxLen]
	return
}

func (mb *client) GetDeviceIdentification(devIdCode byte, objectId byte) (results DeviceIdentification, err error) {
	raw, err := mb.ReadDeviceIdentification(devIdCode, objectId)

	if err != nil {
		return
	}

	results = ParseDeviceIdentification(raw)
	return
}

// Request:
//
//	Function code         : 1 byte (0x05)
//	Output address        : 2 bytes
//	Output value          : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x05)
//	Output address        : 2 bytes
//	Output value          : 2 bytes
func (mb *client) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	// The requested ON/OFF state can only be 0xFF00 and 0x0000
	if value != 0xFF00 && value != 0x0000 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongValueSet,
			Request:       value,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteSingleCoil,
		Data:         dataBlock(address, value),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      4,
		}
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseAddress,
			Request:       respValue,
			Response:      address,
		}
	}
	results = response.Data[2:]
	respValue = binary.BigEndian.Uint16(results)
	if value != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseValue,
			Request:       respValue,
			Response:      value,
		}
	}
	return
}

// Request:
//
//	Function code         : 1 byte (0x06)
//	Register address      : 2 bytes
//	Register value        : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x06)
//	Register address      : 2 bytes
//	Register value        : 2 bytes
func (mb *client) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteSingleRegister,
		Data:         dataBlock(address, value),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      4,
		}
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseAddress,
			Request:       respValue,
			Response:      address,
		}
	}
	results = response.Data[2:]
	respValue = binary.BigEndian.Uint16(results)
	if value != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseValue,
			Request:       respValue,
			Response:      value,
		}
	}
	return
}

// Request:
//
//	Function code         : 1 byte (0x0F)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
//	Byte count            : 1 byte
//	Outputs value         : N* bytes
//
// Response:
//
//	Function code         : 1 byte (0x0F)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
func (mb *client) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	if quantity < 1 || quantity > 1968 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      1968,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteMultipleCoils,
		Data:         dataBlockSuffix(value, address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      4,
		}
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseAddress,
			Request:       respValue,
			Response:      address,
		}
	}
	results = response.Data[2:]
	respValue = binary.BigEndian.Uint16(results)
	if quantity != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseValue,
			Request:       respValue,
			Response:      quantity,
		}
	}
	return
}

// Request:
//
//	Function code         : 1 byte (0x10)
//	Starting address      : 2 bytes
//	Quantity of outputs   : 2 bytes
//	Byte count            : 1 byte
//	Registers value       : N* bytes
//
// Response:
//
//	Function code         : 1 byte (0x10)
//	Starting address      : 2 bytes
//	Quantity of registers : 2 bytes
func (mb *client) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	if quantity < 1 || quantity > 123 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       quantity,
			Response:      123,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteMultipleRegisters,
		Data:         dataBlockSuffix(value, address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 4 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      4,
		}
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseAddress,
			Request:       respValue,
			Response:      address,
		}
	}
	results = response.Data[2:]
	respValue = binary.BigEndian.Uint16(results)
	if quantity != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseValue,
			Request:       respValue,
			Response:      quantity,
		}
	}
	return
}

// Request:
//
//	Function code         : 1 byte (0x16)
//	Reference address     : 2 bytes
//	AND-mask              : 2 bytes
//	OR-mask               : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x16)
//	Reference address     : 2 bytes
//	AND-mask              : 2 bytes
//	OR-mask               : 2 bytes
func (mb *client) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeMaskWriteRegister,
		Data:         dataBlock(address, andMask, orMask),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Fixed response length
	if len(response.Data) != 6 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      6,
		}
	}
	respValue := binary.BigEndian.Uint16(response.Data)
	if address != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseAddress,
			Request:       respValue,
			Response:      address,
		}
	}
	respValue = binary.BigEndian.Uint16(response.Data[2:])
	if andMask != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongANDMask,
			Request:       respValue,
			Response:      andMask,
		}
	}
	respValue = binary.BigEndian.Uint16(response.Data[4:])
	if orMask != respValue {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongORMask,
			Request:       respValue,
			Response:      orMask,
		}
	}
	results = response.Data[2:]
	return
}

// Request:
//
//	Function code         : 1 byte (0x17)
//	Read starting address : 2 bytes
//	Quantity to read      : 2 bytes
//	Write starting address: 2 bytes
//	Quantity to write     : 2 bytes
//	Write byte count      : 1 byte
//	Write registers value : N* bytes
//
// Response:
//
//	Function code         : 1 byte (0x17)
//	Byte count            : 1 byte
//	Read registers value  : Nx2 bytes
func (mb *client) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	if readQuantity < 1 || readQuantity > 125 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       readQuantity,
			Response:      125,
		}
	}
	if writeQuantity < 1 || writeQuantity > 121 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUMaxRequestedQuantity,
			Request:       writeQuantity,
			Response:      121,
		}
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadWriteMultipleRegisters,
		Data:         dataBlockSuffix(value, readAddress, readQuantity, writeAddress, writeQuantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	if count != (len(response.Data) - 1) {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data) - 1,
			Response:      count,
		}
	}
	results = response.Data[1:]
	return
}

// Request:
//
//	Function code         : 1 byte (0x18)
//	FIFO pointer address  : 2 bytes
//
// Response:
//
//	Function code         : 1 byte (0x18)
//	Byte count            : 2 bytes
//	FIFO count            : 2 bytes
//	FIFO count            : 2 bytes (<=31)
//	FIFO value register   : Nx2 bytes
func (mb *client) ReadFIFOQueue(address uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadFIFOQueue,
		Data:         dataBlock(address),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	if len(response.Data) < 4 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data),
			Response:      4,
		}
	}
	count := int(binary.BigEndian.Uint16(response.Data))
	if count != (len(response.Data) - 1) {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUWrongResponseDataSize,
			Request:       len(response.Data) - 1,
			Response:      count,
		}
	}
	count = int(binary.BigEndian.Uint16(response.Data[2:]))
	if count > 31 {
		return []byte{}, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUFifoGreater,
			Request:       count,
			Response:      31,
		}
	}
	results = response.Data[4:]
	return
}

// Helpers

// send sends request and checks possible exception in the response.
func (mb *client) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	aduRequest, err := mb.packager.Encode(request)
	if err != nil {
		return
	}
	aduResponse, err := mb.transporter.Send(aduRequest)
	if err != nil {
		return
	}
	if err = mb.packager.Verify(aduRequest, aduResponse); err != nil {
		return
	}
	response, err = mb.packager.Decode(aduResponse)
	if err != nil {
		return
	}
	// Check correct function code returned (exception)
	if response.FunctionCode != request.FunctionCode {
		err = responseError(response)
		return
	}
	if response.Data == nil || len(response.Data) == 0 {
		// Empty response
		return response, &ModbusPDUError{
			ExceptionCode: ExceptionCodePDUEmptyResponse,
		}
	}
	return
}

// dataBlock creates a sequence of uint16 data.
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}

// dataBlockSuffix creates a sequence of uint16 data and append the suffix plus its length.
func dataBlockSuffix(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+1+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	data[length] = uint8(len(suffix))
	copy(data[length+1:], suffix)
	return data
}

func responseError(response *ProtocolDataUnit) error {
	mbError := &ModbusError{FunctionCode: response.FunctionCode}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = response.Data[0]
	}
	return mbError
}
