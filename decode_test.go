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
