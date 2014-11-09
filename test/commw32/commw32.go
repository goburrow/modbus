// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
// +build windows,cgo

// Port of commw32.c
// To generate go types: go tool cgo commw32.go
package main

// #include <windows.h>
import "C"

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

const port = "COM4"

func main() {
	handle, err := syscall.CreateFile(syscall.StringToUTF16Ptr(port),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,   // mode
		nil, // security
		syscall.OPEN_EXISTING, // no creating new
		0,
		0)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("handle created %d\n", handle)
	defer syscall.CloseHandle(handle)

	var dcb C.DCB
	dcb.BaudRate = 9600
	dcb.ByteSize = 8
	dcb.StopBits = C.ONESTOPBIT
	dcb.Parity = C.NOPARITY
	if C.SetCommState(C.HANDLE(handle), &dcb) == 0 {
		fmt.Printf("set comm state error %v\n", syscall.GetLastError())
		return
	}
	fmt.Printf("set comm state succeed\n")

	var timeouts C.COMMTIMEOUTS
	// time-out between charactor for receiving (ms)
	timeouts.ReadIntervalTimeout = 1000
	timeouts.ReadTotalTimeoutMultiplier = 0
	timeouts.ReadTotalTimeoutConstant = 1000
	timeouts.WriteTotalTimeoutMultiplier = 0
	timeouts.WriteTotalTimeoutConstant = 1000
	if C.SetCommTimeouts(C.HANDLE(handle), &timeouts) == 0 {
		fmt.Printf("set comm timeouts error %v\n", syscall.GetLastError())
		return
	}
	fmt.Printf("set comm timeouts succeed\n")

	var n uint32
	data := []byte("abc")
	err = syscall.WriteFile(handle, data, &n, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("write file succeed\n")
	fmt.Printf("Press Enter when ready for reading...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	data = make([]byte, 512)
	err = syscall.ReadFile(handle, data, &n, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("received data %v:\n", n)
	fmt.Printf("%x\n", data[:n])
	fmt.Printf("closed\n")
}
