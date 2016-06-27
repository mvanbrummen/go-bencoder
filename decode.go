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
	"log"
	"regexp"
	"strconv"
)

const (
	Unicoded     byte   = 0x64
	Unicodee     byte   = 0x65
	Unicodei     byte   = 0x69
	Unicodel     byte   = 0x6c
	Unicode0     byte   = 0x30
	Unicode1     byte   = 0x31
	Unicode2     byte   = 0x32
	Unicode3     byte   = 0x33
	Unicode4     byte   = 0x34
	Unicode5     byte   = 0x35
	Unicode6     byte   = 0x36
	Unicode7     byte   = 0x37
	Unicode8     byte   = 0x38
	Unicode9     byte   = 0x39
	BeIntPattern string = "^(0|-[1-9]\\d*|[1-9]\\d*)$"
)

func BeDecode(b []byte) (dict *BeDict, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = fmt.Errorf("%v", ex)
		}
	}()
	r := bufio.NewReader(bytes.NewReader(b))
	entity := *decodeEntity(r)
	if entity == nil {
		err = errors.New("Root element was not a dictionary.")
	}
	d, ok := entity.(*BeDict)
	if !ok {
		err = errors.New("Root element was not a dictionary.")
	}
	dict = d
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
		case Unicodei:
			bencodeEntity = decodeInteger(reader)
		case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6, Unicode7, Unicode8, Unicode9:
			bencodeEntity = decodeString(reader)
		case Unicodel:
			bencodeEntity = decodeList(reader)
		case Unicoded:
			bencodeEntity = decodeDictionary(reader)
		case Unicodee:
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
	log.Printf("INFO: Decoded integer. Returning %v", BeInteger{str})
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
	log.Printf("INFO: Decoded string. Returning %v in %v", buf.String(), BeString{length, buf.Bytes()})
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
	log.Printf("INFO: Decoded list. Returning %v", list)
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
	log.Printf("INFO: Decoded dictionary. Returning %v", dict)
	return &dict
}
