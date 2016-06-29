// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const BeIntPattern string = "^(0|-[1-9]\\d*|[1-9]\\d*)$"

func Unmarshal(b []byte) (dict *BeDict, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = fmt.Errorf("%v", ex)
		}
	}()
	r := bufio.NewReader(bytes.NewReader(b))
	entity := *decodeEntity(r)
	if entity == nil {
		err = errors.New("Failed to parse bencoded data.")
	}
	dict, ok := entity.(*BeDict)
	if !ok {
		err = errors.New("Root element was not a dictionary.")
	}
	return dict, err
}

func decodeEntity(reader *bufio.Reader) *Bencoder {
	var bencodeEntity Bencoder
	if b, err := reader.Peek(1); err != nil {
		if err == io.EOF {
			return nil
		} else {
			panic(err)
		}
	} else {
		switch b[0] {
		case 'i':
			bencodeEntity = decodeInteger(reader)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			bencodeEntity = decodeString(reader)
		case 'l':
			bencodeEntity = decodeList(reader)
		case 'd':
			bencodeEntity = decodeDictionary(reader)
		case 'e':
			reader.ReadByte()
		default:
			return nil
		}
	}
	return &bencodeEntity
}

func decodeInteger(reader *bufio.Reader) *BeInteger {
	var str string
	if b, err := reader.ReadBytes('e'); err != nil {
		panic(err)
	} else {
		str = fmt.Sprintf("%s", string(b[1:len(b)-1]))
		if b, err := regexp.MatchString(BeIntPattern, str); err != nil || b == false {
			panic(fmt.Sprintf("Could not parse integer: %s", str))
		}
	}
	return &BeInteger{str}
}

func decodeString(reader *bufio.Reader) *BeString {
	var length int
	// get the length of the BeString
	if b, err := reader.ReadBytes(':'); err != nil {
		panic(err)
	} else {
		str := fmt.Sprintf("%s", string(b[:len(b)-1]))
		length, err = strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
	}
	// get the value of the BeString
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		if b, err := reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			buf.WriteByte(b)
		}
	}
	return &BeString{length, buf.Bytes()}
}

func decodeList(reader *bufio.Reader) *BeList {
	var list BeList
	reader.ReadByte()
	for {
		value := decodeEntity(reader)
		if value == nil || *value == nil {
			break
		} else {
			list = append(list, *value)
		}
	}
	return &list
}

func decodeDictionary(reader *bufio.Reader) *BeDict {
	var dict BeDict = make(BeDict)
	reader.ReadByte()
	for {
		// get dictionary key
		key := decodeEntity(reader)
		if key == nil || *key == nil {
			break
		}
		k, ok := (*key).(*BeString)
		if !ok {
			panic("Dictionary key was not a string.")
		} else {
			// get associated value
			v := decodeEntity(reader)
			if v == nil {
				panic(fmt.Sprintf("Dictionary key '%s' does not have an associated value.", k.Val))
			}
			dict[string(k.Val)] = *v
		}
	}
	return &dict
}
