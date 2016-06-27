// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"bufio"
	"bytes"
	"testing"
)

type integerTestpair struct {
	Value    string
	Expected *BeInteger
}

var integerTests = []integerTestpair{
	{"i1e", NewBeInteger(1)},
	{"i1345435e", NewBeInteger(1345435)},
	{"i-1e", NewBeInteger(-1)},
	{"i0e", NewBeInteger(0)},
}

type stringTestpair struct {
	Value    string
	Expected *BeString
}

var stringTests = []stringTestpair{
	{"4:spam", NewBeString("spam")},
	{"10:kickboxing", NewBeString("kickboxing")},
	{"0:", NewBeString("")},
}

type listTestpair struct {
	Value    string
	Expected *BeList
}

var listTests = []listTestpair{
	{"l4:spami1ee", &BeList{NewBeString("spam"), NewBeInteger(1)}},
	{"lli1ei2ei3eee", &BeList{&BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}},
	{"le", &BeList{}},
}

type dictTestpair struct {
	Value    string
	Expected *BeDict
}

var dictTests = []dictTestpair{
	{"d4:spami3ee", &BeDict{"spam": NewBeInteger(3)}},
	{"d4:spamli1ei2ei3eee", &BeDict{"spam": &BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}},
	{"d3:eggd4:spamli1ei2ei3eeee", &BeDict{"egg": &BeDict{"spam": &BeList{NewBeInteger(1), NewBeInteger(2), NewBeInteger(3)}}}},
}

func TestDecodeInteger(t *testing.T) {
	for _, pair := range integerTests {
		reader := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		result := decodeInteger(reader)
		if result.Val != pair.Expected.Val {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", result)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, pair := range stringTests {
		reader := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		result := decodeString(reader)
		if string(result.Val) != string(pair.Expected.Val) {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", string(result.Val))
		}
		if result.Len != pair.Expected.Len {
			t.Error("For", pair.Value, "Expected", pair.Expected.Len, "Got", result.Len)
		}
	}
}

func TestDecodeList(t *testing.T) {
	for _, pair := range listTests {
		reader := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		result := decodeList(reader)
		if len(*result) != len(*pair.Expected) {
			t.Error("For", pair.Value, "Expected", len(*pair.Expected), "got", len(*result))
		}
	}
}

func TestDecodeDictionary(t *testing.T) {
	for _, pair := range dictTests {
		reader := bufio.NewReader(bytes.NewReader([]byte(pair.Value)))
		result := decodeDictionary(reader)
		if len(*result) != len(*pair.Expected) {
			t.Error("For", pair.Value, "Expected", len(*pair.Expected), "got", len(*result))
		}
	}
}
