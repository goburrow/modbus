// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package test

import (
	"testing"
	"github.com/goburrow/modbus"
	"log"
	"os"
)
const (
	testAsciiDevice = "/dev/pts/7"
)

func TestASCIIClientReadCoils(t *testing.T) {
	client := modbus.ASCIIClient(testAsciiDevice)
	ClientTestReadCoils(t, client)
}

func TestASCIIClientAdvancedUsage(t *testing.T) {
	var handler modbus.ASCIIClientHandler
	handler.Device = testAsciiDevice
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.ASCIIClientWithHandler(&handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}
