// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"fmt"
	"testing"
)

type encodeTestpair struct {
	Entity   Bencoder
	Expected string
}

var encodeTests = []encodeTestpair{
	{NewBeString("spam"), "4:spam"},
	{NewBeString("boourns"), "7:boourns"},
	{NewBeInteger(0), "i0e"},
	{NewBeInteger(-50), "i-50e"},
	{NewBeInteger(123456), "i123456e"},
	{&BeList{NewBeInteger(123456), NewBeString("spam")}, "li123456e4:spame"},
	{&BeList{&BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}, "lli1ei2ei3eee"},
	{&BeDict{"spam": NewBeInteger(3)}, "d4:spami3ee"},
	{&BeDict{"spam": &BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}, "d4:spamli1ei2ei3eee"},
	{&BeDict{"egg": &BeDict{"spam": &BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}}, "d3:eggd4:spamli1ei2ei3eeee"},
}

func TestEncode(t *testing.T) {
	for _, pair := range encodeTests {
		str := fmt.Sprintf("%s", pair.Entity.Bencode())
		if str != pair.Expected {
			t.Error("For", pair.Entity, "expected", pair.Expected, "got", str)
		}
	}
}

type encodeTestpair2 struct {
	Value    interface{}
	Expected string
}

var encodeTests2 = []encodeTestpair2{
	{"spam", "4:spam"},
	{123, "i123e"},
	{[]interface{}{"spam", 123}, "l4:spami123ee"},
	{map[string]interface{}{"spam": 123}, "d4:spami123ee"},
}

func TestEncodeInterface(t *testing.T) {
	for _, pair := range encodeTests2 {
		b, _ := BeEncode(pair.Value)
		if string(b) != pair.Expected {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", string(b))
		}
	}
}
