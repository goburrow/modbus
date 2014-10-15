// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

import (
	"bytes"
	"encoding/binary"
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

// Common errors
var (
	ErrBadRequest = errors.New("invalid request")
	ErrNoResponse = errors.New("no response data")
	ErrResponseSize = errors.New("response data size mismatch")
	ErrTransactionId = errors.New("transaction id mismatch")
)

func (e *ModbusError) Error() string {
	// Really need fmt?
	var buf bytes.Buffer
	buf.WriteString("modbus exception: ")
	buf.WriteString(strconv.FormatInt(int64(e.ExceptionCode), 10))
	buf.WriteString(", function: ")
	buf.WriteString(strconv.FormatInt(int64(e.FunctionCode), 10))
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

type Decoder interface {
	Decode(adu []byte) (pdu *ProtocolDataUnit, err error)
}

type Transporter interface {
	Send(aduRequest []byte) (aduResponse []byte, err error)
}

type modbusClient struct {
	encoder     Encoder
	decoder     Decoder
	transporter Transporter
}

// Request:
//  Function code         : 1 byte (0x01)
// 	Starting address      : 2 bytes
// 	Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x01)
// 	Byte count            : 1 byte
// 	Input status          : N* bytes (=N or N+1)
func (mb *modbusClient) ReadCoils(address, quantity uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadCoils,
		Data:         dataRead(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	if count != (len(response.Data)-1) {
		err = ErrResponseSize
		return
	}
	results = response.Data[1:]
	return
}

// Request:
//  Function code         : 1 byte (0x02)
// 	Starting address      : 2 bytes
// 	Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x02)
// 	Byte count            : 1 byte
// 	Input status          : N* bytes (=N or N+1)
func (mb *modbusClient) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadDiscreteInputs,
		Data:         dataRead(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	if count != (len(response.Data)-1) {
		err = ErrResponseSize
		return
	}
	results = response.Data[1:]
	return
}

// Request:
//  Function code         : 1 byte (0x03)
// 	Starting address      : 2 bytes
// 	Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x03)
// 	Byte count            : 1 byte
// 	Register value        : Nx2 bytes
func (mb *modbusClient) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadHoldingRegisters,
		Data:         dataRead(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	if count != (len(response.Data)-1) {
		err = ErrResponseSize
		return
	}
	results = response.Data[1:]
	return
}

// Request:
//  Function code         : 1 byte (0x04)
// 	Starting address      : 2 bytes
// 	Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x04)
// 	Byte count            : 1 byte
// 	Input registers       : Nx2 bytes
func (mb *modbusClient) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadInputRegisters,
		Data:         dataRead(address, quantity),
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	count := int(response.Data[0])
	if count != (len(response.Data)-1) {
		err = ErrResponseSize
		return
	}
	results = response.Data[1:]
	return
}

// Request:
//  Function code         : 1 byte (0x04)
// 	Starting address      : 2 bytes
// 	Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x04)
// 	Byte count            : 1 byte
// 	I value        : Nx2 bytes
func (mb *modbusClient) WriteSingleCoil(address, count int) {
	return
}

func (mb *modbusClient) WriteMultipleCoils(address, count int) {

}

// 16-bit access



func (mb *modbusClient) WriteSingleRegister(address, count int) {

}

func (mb *modbusClient) WriteMultipleRegisters(address, count int) {

}

func (mb *modbusClient) ReadWriteMultipleRegisters(address, count int) {

}

func (mb *modbusClient) MaskWriteRegister(address, count int) {

}

// Send a request and check possible exception in the response
func (mb *modbusClient) send(request *ProtocolDataUnit) (response *ProtocolDataUnit, err error) {
	aduRequest, err := mb.encoder.Encode(request)
	if err != nil {
		return
	}
	aduResponse, err := mb.transporter.Send(aduRequest)
	if err != nil {
		return
	}
	response, err = mb.decoder.Decode(aduResponse)
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
		err = ErrNoResponse
		return
	}
	return
}

// Request data for read functions
func dataRead(address, quantity uint16) []byte {
	data := [4]byte{}
	binary.BigEndian.PutUint16(data[:], address)
	binary.BigEndian.PutUint16(data[2:], quantity)

	return data[:]
}

func responseError(response *ProtocolDataUnit) error {
	mbError := &ModbusError{FunctionCode: response.FunctionCode}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = response.Data[0]
	}
	return mbError
}
