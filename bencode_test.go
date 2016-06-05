package main

import (
	"bufio"
	"bytes"
	"testing"
)

type testpair struct {
	Value    []byte
	Expected *BeInteger
}

var tests = []testpair{
	{[]byte("i1e"), &BeInteger{"1"}},
  {[]byte("i1345435e"), &BeInteger{"1345435"}},
  {[]byte("i-1e"), &BeInteger{"-1"}},
  {[]byte("i0e"), &BeInteger{"0"}},
}

func TestDecodeInteger(t *testing.T) {
	for _, pair := range tests {
		reader := bufio.NewReader(bytes.NewReader(pair.Value))
		result := DecodeInteger(reader)
		if result.Val != pair.Expected.Val {
			t.Error("For", pair.Value, "Expected", pair.Expected.Val, "Got", result)
		}
	}
}
