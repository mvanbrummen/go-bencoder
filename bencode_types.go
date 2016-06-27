package bencode

import (
	"bytes"
	"fmt"
	"strconv"
)

type BeType int

const (
	BeStringType BeType = iota
	BeIntegerType
	BeListType
	BeDictType
)

type Bencoder interface {
	Bencode() []byte
}

type BeNode struct {
	String     *BeString
	Integer    *BeInteger
	List       *BeList
	Dictionary *BeDict
	Type       BeType
}

func (b *BeNode) IsNil() bool {
	return b.String == nil && b.Integer == nil && b.List == nil && b.Dictionary == nil
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

type BeList []BeNode

func (list *BeList) Bencode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('l')
	for _, item := range *list {
		switch item.Type {
		case BeIntegerType:
			buf.Write(item.Integer.Bencode())
		case BeStringType:
			buf.Write(item.String.Bencode())
		case BeListType:
			buf.Write(item.List.Bencode())
		case BeDictType:
			buf.Write(item.Dictionary.Bencode())
		}
	}
	buf.WriteRune('e')
	return buf.Bytes()
}

type BeDict map[string]BeNode

func (dict *BeDict) Bencode() []byte {
	var buf bytes.Buffer
	buf.WriteRune('d')
	for k, v := range *dict {
		buf.Write([]byte(k))
		switch v.Type {
		case BeIntegerType:
			buf.Write(v.Integer.Bencode())
		case BeStringType:
			buf.Write(v.String.Bencode())
		case BeListType:
			buf.Write(v.List.Bencode())
		case BeDictType:
			buf.Write(v.Dictionary.Bencode())
		}
	}
	buf.WriteRune('e')
	return buf.Bytes()
}
