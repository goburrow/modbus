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

// Common errors
var (
	ErrBadRequest    = errors.New("invalid request")
	ErrResponse      = errors.New("invalid response data")
	ErrResponseSize  = errors.New("response data size mismatch")
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
