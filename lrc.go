// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.
package modbus

// Longitudinal Redundancy Checking
type lRC struct {
	sum uint8
}

func (lrc *lRC) pushByte(b byte) *lRC {
	lrc.sum += b
	return lrc
}

func (lrc *lRC) pushBytes(data []byte) *lRC {
	for _, b := range data {
		lrc.pushByte(b)
	}
	return lrc
}

func (lrc *lRC) value() byte {
	// Return twos complement
	return uint8(-int8(lrc.sum))
}
