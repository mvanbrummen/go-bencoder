// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"bytes"
	"fmt"
	"strconv"
)

type Bencoder interface {
	Bencode() []byte
}

type BeString struct {
	Len int
	Val []byte
}

func NewBeString(str string) *BeString {
	return &BeString{len(str), []byte(str)}
}

func (str *BeString) Bencode() []byte {
	s := fmt.Sprintf("%d:%s", str.Len, str.Val)
	return []byte(s)
}

type BeInteger struct {
	Val string
}

func NewBeInteger(i int) *BeInteger {
	return &BeInteger{strconv.Itoa(i)}
}

func (integer *BeInteger) Bencode() []byte {
	s := fmt.Sprintf("i%se", integer.Val)
	return []byte(s)
}

type BeList []Bencoder

func (list *BeList) Bencode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('l')
	for _, item := range *list {
		buf.Write(item.Bencode())
	}
	buf.WriteRune('e')
	return buf.Bytes()
}

type BeDict map[string]Bencoder

func (dict *BeDict) Bencode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('d')
	for k, v := range *dict {
		str := fmt.Sprintf("%d:%s", len(k), k)
		buf.Write([]byte(str))
		buf.Write(v.Bencode())
	}
	buf.WriteRune('e')
	return buf.Bytes()
}
