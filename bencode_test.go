package main

import (
	"bufio"
	"bytes"
	"testing"
)

type integerTestpair struct {
	Value    []byte
	Error    error
	Expected *BeInteger
}

type stringTestpair struct {
	Value    []byte
	Error    error
	Expected *BeString
}

var integerTests = []integerTestpair{
	{[]byte("i1e"), nil, &BeInteger{"1"}},
	{[]byte("i1345435e"), nil, &BeInteger{"1345435"}},
	{[]byte("i-1e"), nil, &BeInteger{"-1"}},
	{[]byte("i0e"), nil, &BeInteger{"0"}},
}

var stringTests = []stringTestpair{
	{[]byte("4:spam"), nil, &BeString{4, []byte("spam")}},
	{[]byte("10:kickboxing"), nil, &BeString{10, []byte("kickboxing")}},
	{[]byte("0:"), nil, &BeString{0, []byte("")}},
}

func TestDecodeInteger(t *testing.T) {
	for _, pair := range integerTests {
		reader := bufio.NewReader(bytes.NewReader(pair.Value))
		result, err := DecodeInteger(reader)
		if err != nil {
			if err != pair.Error {
				t.Error("For", pair.Value, "Expected", pair.Error, "Got", err)
			}
		}
		if result.Val != pair.Expected.Val {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", result)
		}
	}
}

func TestDecodeString(t *testing.T) {
	for _, pair := range stringTests {
		reader := bufio.NewReader(bytes.NewReader(pair.Value))
		result, err := DecodeString(reader)
		if err != nil {
			if err != pair.Error {
				t.Error("For", pair.Value, "Expected", pair.Error, "Got", err)
			}
		}
		if string(result.Val) != string(pair.Expected.Val) {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", string(result.Val))
		}
		if result.Len != pair.Expected.Len {
			t.Error("For", pair.Value, "Expected", pair.Expected.Len, "Got", result.Len)
		}
	}
}
