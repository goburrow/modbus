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
	buf.WriteString(strconv.Itoa(e.ExceptionCode))
	buf.WriteString(", function: ")
	buf.WriteString(strconv.Itoa(e.FunctionCode))
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
//  Function code: 1 byte (0x01)
// 	Starting address: 2 bytes
// 	Quantity of coils: 2 bytes
// Response:
//  Function code: 1 byte (0x01)
// 	Byte count: 1 byte
// 	Coil status: n bytes (=N or N + 1)
func (mb *modbusClient) ReadCoils(address, quantity int) (results []byte, err error) {
	data := [4]byte{}
	binary.BigEndian.PutUint16(data[:], uint16(address))
	binary.BigEndian.PutUint16(data[2:], uint16(quantity))

	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadCoils,
		Data:         data[:],
	}
	response, err := mb.send(&request)
	if err != nil {
		return
	}
	// Expect count / 8 bytes
	count := int(response.Data[0])
	if count != quantity {
		err = ErrResponseSize
		return
	}
	results = response.Data[1:]
	return
}

// Bit access
func (mb *modbusClient) ReadDiscreteInputs(address, quantity int) (results []byte, err error) {
	return
}

func (mb *modbusClient) WriteSingleCoil(address, count int) {
	return
}

func (mb *modbusClient) WriteMultipleCoils(address, count int) {

}

// 16-bit access
func (mb *modbusClient) ReadInputRegisters(address, count int) {

}

func (mb *modbusClient) ReadHoldingRegisters(address, count int) {

}

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
	if len(response.Data) == 0 {
		// Empty response
		err = ErrNoResponse
		return
	}
	return
}

func responseError(response *ProtocolDataUnit) error {
	mbError := &ModbusError{FunctionCode: int(response.FunctionCode)}
	if response.Data != nil && len(response.Data) > 0 {
		mbError.ExceptionCode = int(response.Data[0])
	}
	return mbError
}
