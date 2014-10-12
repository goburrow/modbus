// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"errors"
	"strconv"
)

const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
)

const (
	ExceptionCodeIllegalFunction                    = 1
	ExceptionCodeIllegalDataAddress                 = 2
	ExceptionCodeIllegalDataValue                   = 3
	ExceptionCodeServerDeviceFailure                = 4
	ExceptionCodeAcknowledge                        = 5
	ExceptionCodeServerDeviceBusy                   = 6
	ExceptionCodeMemoryParityError                  = 8
	ExceptionCodeGatewayPathUnavailable             = 10
	ExceptionCodeGatewayTargetDeviceFailedToRespond = 11
)

var (
	ErrNoData   = errors.New("no data")
	ErrDataSize = errors.New("data size mismatch")
)

func (e *ModbusError) Error() string {
	// Really need fmt?
	var buf bytes.Buffer
	buf.WriteString("modbus function: ")
	buf.WriteString(strconv.Itoa(e.FunctionCode))
	buf.WriteString(", exception: ")
	buf.WriteString(strconv.Itoa(e.ExceptionCode))
	return buf.String()
}

// PDU
type ProtocolDataUnit struct {
	FunctionCode byte
	Data         []byte
}

type Encoder interface {
	Encode(pdu *ProtocolDataUnit) (adu []byte, err error)
}

type Transporter interface {
	Send(adu []byte, pdu *ProtocolDataUnit) (err error)
}

type modbus struct {
	encoder     Encoder
	transporter Transporter
}

// Request:
//  Function code: 1 byte (0x01)
// 	Starting address: 2 bytes
// 	Quantities of coils: 2 bytes
// Response:
//  Function code: 1 byte (0x01)
// 	Byte count: 1 byte
// 	Coil status: n bytes
// Error:
//  Error code: 1 byte (Function code + 0x80)
//  Exception code: 1 byte
func (mb *modbus) ReadCoils(address, count int) (results []byte, err error) {
	data := [4]byte{
		// Starting address high & low
		byte(address >> 16),
		byte(address),
		// Quantity of coils
		byte(count >> 16),
		byte(count),
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadCoils,
		Data:         data,
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count = int(response.Data[0])
	if count != (len(response.Data) - 1) {
		err = ErrDataSize
		return
	}
	results = response.Data[1:]
	return
}

// Send a request and check possible exception in the response
func (mb *modbus) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	adu, err := mb.encoder.Encode(request)
	if err != nil {
		return
	}
	// Use existing pdu for the response
	err = mb.transporter.Send(adu, response)
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
		err = ErrNoData
		return
	}
	return
}

func responseError(response *ProtocolDataUnit) error {
	mbError := &ModbusError{FunctionCode: response.FunctionCode}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = response.Data[0]
	}
	return mbError
}
