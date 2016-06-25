package bencode

import (
	"bufio"
	"bytes"
	"testing"
)

type integerTestpair struct {
	Value    []byte
	Expected *BeInteger
}

var integerTests = []integerTestpair{
	{[]byte("i1e"), &BeInteger{"1"}},
	{[]byte("i1345435e"), &BeInteger{"1345435"}},
	{[]byte("i-1e"), &BeInteger{"-1"}},
	{[]byte("i0e"), &BeInteger{"0"}},
}

type stringTestpair struct {
	Value    []byte
	Expected *BeString
}

var stringTests = []stringTestpair{
	{[]byte("4:spam"), &BeString{4, []byte("spam")}},
	{[]byte("10:kickboxing"), &BeString{10, []byte("kickboxing")}},
	{[]byte("0:"), &BeString{0, []byte("")}},
}

func TestDecodeInteger(t *testing.T) {
	for _, pair := range integerTests {
		reader := bufio.NewReader(bytes.NewReader(pair.Value))
		result := decodeInteger(reader)
		if result.Val != pair.Expected.Val {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", result)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, pair := range stringTests {
		reader := bufio.NewReader(bytes.NewReader(pair.Value))
		result := decodeString(reader)
		if string(result.Val) != string(pair.Expected.Val) {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", string(result.Val))
		}
		if result.Len != pair.Expected.Len {
			t.Error("For", pair.Value, "Expected", pair.Expected.Len, "Got", result.Len)
		}
	}
}
