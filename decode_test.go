// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"bufio"
	"bytes"
	"testing"
)

type integerTestPair struct {
	Value    string
	Expected int64
}

var integerTests = []integerTestPair{
	{"i123e", 123},
	{"i0e", 0},
	{"i-123e", -123},
}

func TestDecodeInteger(t *testing.T) {
	for _, pair := range integerTests {
		r := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		i := decodeInteger(r)
		if i != pair.Expected {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", i)
		}
	}
}

type stringTestPair struct {
	Value    string
	Expected []byte
}

var stringTests = []stringTestPair{
	{"4:spam", []byte("spam")},
	{"0:", []byte("")},
	{"7:boourns", []byte("boourns")},
}

func TestDecodeString(t *testing.T) {
	for _, pair := range stringTests {
		r := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		s := decodeString(r)
		if !bytes.Equal(s, pair.Expected) {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", s)
		}
	}
}

type listTestPair struct {
	Value    string
	Expected []interface{}
}

var listTests = []listTestPair{
	{"li1ei2ei3ee", []interface{}{1, 2, 3}},
	{"l4:spami222ee", []interface{}{[]byte("spam"), 222}},
	{"li1eli2ei3ei4eee", []interface{}{1, []interface{}{2, 3, 4}}},
}

func TestDecodeList(t *testing.T) {
	for _, pair := range listTests {
		r := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		l := decodeList(r)
		if len(l) != len(pair.Expected) {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", l)
		}
	}
}

func compareSlices(a, b []interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type dictTestPair struct {
	Value    string
	Expected map[string]interface{}
}

var dictTests = []dictTestPair{
	{"d4:spami123ee", map[string]interface{}{"spam": 123}},
}

func TestDecodeDictionary(t *testing.T) {
	for _, pair := range dictTests {
		r := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		d := decodeDictionary(r)
		if len(d) != len(pair.Expected) {
			t.Error("For", pair.Value, "Expected", pair.Expected, "Got", d)
		}
	}
}
