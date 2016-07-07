// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"testing"
)

type encodeTestpair2 struct {
	Value    interface{}
	Expected string
}

var encodeTests2 = []encodeTestpair2{
	{"spam", "4:spam"},
	{[]byte("spam"), "4:spam"},
	{123, "i123e"},
	{[]interface{}{"spam", 123}, "l4:spami123ee"},
	{map[string]interface{}{"spam": 123}, "d4:spami123ee"},
}

func TestEncodeInterface(t *testing.T) {
	for _, pair := range encodeTests2 {
		b, _ := Marshal(pair.Value)
		if string(b) != pair.Expected {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", string(b))
		}
	}
}
