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
