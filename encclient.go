package modbus

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	encMaxSize     = 256
	encHeaderSize  = 3
	encCrcSize     = 2
	encTimeout     = 10 * time.Second
	encIdleTimeout = 60 * time.Second
)

// region - Enc Client

// EncClient creates TCP-RTU client with default handler and given connect string.
func EncClient(address string) Client {
	handler := NewEncClientHandler(address)
	return NewClient(handler)
}

// endregion

// region - Enc Client Handler

type EncClientHandler struct {
	rtuPackager
	encTransporter
}

func NewEncClientHandler(address string) *EncClientHandler {
	h := &EncClientHandler{}
	h.Address = address
	h.Timeout = encTimeout
	h.IdleTimeout = encIdleTimeout
	return h
}

// endregion

// region - Enc Transporter

type encTransporter struct {

	// Connect string
	Address string

	// Connect & Read timeout
	Timeout time.Duration

	// Idle timeout to close the connection
	IdleTimeout time.Duration

	// Transmission logger
	Logger *log.Logger

	// TCP connection
	mu           sync.Mutex
	conn         net.Conn
	closeTimer   *time.Timer
	lastActivity time.Time
}

// region = public API

func (mb *encTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {

	mb.mu.Lock()
	defer mb.mu.Unlock()

	// Establish a new connection if not connected
	if err = mb.connect(); err != nil {
		return
	}

	// Set timer to close when idle
	mb.lastActivity = time.Now()
	mb.startCloseTimer()

	// Set write and read timeout
	var timeout time.Time
	if mb.Timeout > 0 {
		timeout = mb.lastActivity.Add(mb.Timeout)
	}
	if err = mb.conn.SetDeadline(timeout); err != nil {
		return
	}

	// Send data
	mb.logf("modbus: sending % x", aduRequest)
	if _, err = mb.conn.Write(aduRequest); err != nil {
		mb.logf("modbus: could not send request: %s", err)
		return
	}

	//time.Sleep(mb.calculateDelay(len(aduRequest) + bytesToRead))
	//time.Sleep(1000)

	var data [encMaxSize]byte

	// read header
	if _, err = io.ReadFull(mb.conn, data[:encHeaderSize]); err != nil {
		return
	}
	//mb.logf("modbus: header % x", data[:encHeaderSize])

	length := data[2]
	if length <= 0 {
		mb.flush(data[:])
		err = fmt.Errorf("modbus: length in response header '%v' must not be zero", length)
		return
	}
	if length > (encMaxSize - encHeaderSize) {
		mb.flush(data[:])
		err = fmt.Errorf("modbus: length in response header '%v' must not greater than '%v'", length, encMaxSize-encHeaderSize)
		return
	}

	// read data
	length += encHeaderSize
	if _, err = io.ReadFull(mb.conn, data[encHeaderSize:length]); err != nil {
		return
	}
	//mb.logf("modbus: data % x", data[encHeaderSize:length])

	//read CRC
	if _, err = io.ReadFull(mb.conn, data[length:length+encCrcSize]); err != nil {
		return
	}
	//mb.logf("modbus: crc % x", data[length:length+encCrcSize])

	aduResponse = data[:length+encCrcSize]
	mb.logf("modbus: received % x\n", aduResponse)
	return

}

// Connect establishes a new connection to the address in Address.
// Connect and Close are exported so that multiple requests can be done with one session
func (mb *encTransporter) Connect() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return mb.connect()
}

// Close closes current connection.
func (mb *encTransporter) Close() error {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return mb.close()
}

// endregion

// region = helper methods

func (mb *encTransporter) connect() error {
	if mb.conn == nil {
		dialer := net.Dialer{Timeout: mb.Timeout}
		conn, err := dialer.Dial("tcp", mb.Address)
		if err != nil {
			return err
		}
		mb.conn = conn
	}
	return nil
}
func (mb *encTransporter) startCloseTimer() {
	if mb.IdleTimeout <= 0 {
		return
	}
	if mb.closeTimer == nil {
		mb.closeTimer = time.AfterFunc(mb.IdleTimeout, mb.closeIdle)
	} else {
		mb.closeTimer.Reset(mb.IdleTimeout)
	}
}

// flush flushes pending data in the connection,
// returns io.EOF if connection is closed.
func (mb *encTransporter) flush(b []byte) (err error) {
	if err = mb.conn.SetReadDeadline(time.Now()); err != nil {
		return
	}
	// Timeout setting will be reset when reading
	if _, err = mb.conn.Read(b); err != nil {
		// Ignore timeout error
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			err = nil
		}
	}
	return
}

func (mb *encTransporter) logf(format string, v ...interface{}) {
	if mb.Logger != nil {
		mb.Logger.Printf(format, v...)
	}
}

// closeLocked closes current connection. Caller must hold the mutex before calling this method.
func (mb *encTransporter) close() (err error) {
	if mb.conn != nil {
		err = mb.conn.Close()
		mb.conn = nil
	}
	return
}

// closeIdle closes the connection if last activity is passed behind IdleTimeout.
func (mb *encTransporter) closeIdle() {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	if mb.IdleTimeout <= 0 {
		return
	}
	idle := time.Now().Sub(mb.lastActivity)
	if idle >= mb.IdleTimeout {
		mb.logf("modbus: closing connection due to idle timeout: %v", idle)
		mb.close()
	}
}

// endregion

// endregion
