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
}

func TestEncode(t *testing.T) {
	for _, pair := range encodeTests {
		str := fmt.Sprintf("%s", pair.Entity.Bencode())
		if str != pair.Expected {
			t.Error("For", pair.Entity, "expected", pair.Expected, "got", str)
		}
	}
}
