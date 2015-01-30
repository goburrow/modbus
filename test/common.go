// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

package test

import (
	"runtime"
	"strings"
	"testing"
)

func AssertEquals(t *testing.T, expected, actual interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		// Get file name only
		idx := strings.LastIndex(file, "/")
		if idx >= 0 {
			file = file[idx+1:]
		}
	}

	if expected != actual {
		t.Logf("%s:%d: Expected: %+v (%T), actual: %+v (%T)", file, line,
			expected, expected, actual, actual)
		t.FailNow()
	}
}
