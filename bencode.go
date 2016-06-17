package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

const (
	Unicoded     byte   = 100
	Unicodee     byte   = 101
	Unicodei     byte   = 105
	Unicodel     byte   = 108
	Unicode0     byte   = 48
	Unicode1     byte   = 49
	Unicode2     byte   = 50
	Unicode3     byte   = 51
	Unicode4     byte   = 52
	Unicode5     byte   = 53
	Unicode6     byte   = 54
	Unicode7     byte   = 55
	Unicode8     byte   = 56
	Unicode9     byte   = 57
	BeIntPattern string = "^(0|-[1-9]\\d*|[1-9]\\d*)$"
)

func BeDecode(b []byte) (entity *BeNode, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = fmt.Errorf("%v", ex)
		}
	}()
	r := bufio.NewReader(bytes.NewReader(b))
	entity = decodeEntity(r)
	return entity, err
}

func decodeEntity(reader *bufio.Reader) *BeNode {
	var bencodeEntity BeNode
	if b, err := reader.Peek(1); err != nil {
		if err == io.EOF {
			return nil
		} else {
			panic(err)
		}
	} else {
		switch b[0] {
		case Unicodei:
			bencodeEntity = BeNode{Integer: decodeInteger(reader), Type: BeIntegerType}
		case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6, Unicode7, Unicode8, Unicode9:
			bencodeEntity = BeNode{String: decodeString(reader), Type: BeStringType}
		case Unicodel:
			bencodeEntity = BeNode{List: decodeList(reader), Type: BeListType}
		case Unicoded:
			bencodeEntity = BeNode{Dictionary: decodeDictionary(reader), Type: BeDictType}
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
		if value == nil {
			break
		}
		list = append(list, *value)
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
		if key == nil || key.String == nil {
			break
		} else if key.Type != BeStringType {
			panic("Dictionary key was not a string.")
		} else {
			// get associated value
			v := decodeEntity(reader)
			if v.IsNil() {
				panic(fmt.Sprintf("Dictionary key '%s' does not have an associated value.", key.String.Val))
			}
			dict[string(key.String.Val)] = *v
		}
	}
	log.Printf("INFO: Decoded dictionary. Returning %v", dict)
	return &dict
}
