package main

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
	Unicoded     byte   = 100
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

func BeDecode(b []byte) (interface{}, error) {
	r := bufio.NewReader(bytes.NewReader(b))
	var entity interface{}
	for {
		entity, err := decodeEntity(r)
		if err != nil {
			return nil, err
		}
		if entity == nil {
			break
		}
	}
	return entity, nil
}

func decodeEntity(reader *bufio.Reader) (interface{}, error) {
	var bencodeEntity interface{}
	if b, err := reader.Peek(1); err != nil {
		if err == io.EOF {
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		var err error
		switch b[0] {
		case Unicodei:
			bencodeEntity, err = decodeInteger(reader)
		case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6,
			Unicode7, Unicode8, Unicode9:
			bencodeEntity, err = decodeString(reader)
		case Unicodel:
			bencodeEntity, err = decodeList(reader)
		case Unicoded:
			bencodeEntity, err = decodeDictionary(reader)
		default:
			bencodeEntity = nil
		}
		if err != nil {
			return nil, err
		}
	}
	return bencodeEntity, nil
}

func decodeInteger(reader *bufio.Reader) (*BeInteger, error) {
	var str string
	if b, err := reader.ReadBytes('e'); err != nil {
		if err == io.EOF {
			return nil, err
		}
	} else {
		str = fmt.Sprintf("%s", string(b[1:len(b)-1]))
		if b, err := regexp.MatchString(BeIntPattern, str); err != nil || b == false {
			return nil, errors.New(fmt.Sprintf("Could not parse integer: %s", str))
		}
	}
	log.Printf("INFO: Decoded integer. Returning %v", BeInteger{str})
	return &BeInteger{str}, nil
}

func decodeString(reader *bufio.Reader) (*BeString, error) {
	var length int
	// get the length of the BeString
	if b, err := reader.ReadBytes(':'); err != nil {
		if err == io.EOF {
			return nil, err
		}
	} else {
		str := fmt.Sprintf("%s", string(b[:len(b)-1]))
		length, err = strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
	}
	// get the value of the BeString
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		if b, err := reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			buf.WriteByte(b)
		}
	}
	log.Printf("INFO: Decoded string. Returning %v in %v", buf.String(), BeString{length, buf.Bytes()})
	return &BeString{length, buf.Bytes()}, nil
}

func decodeList(reader *bufio.Reader) (*BeList, error) {
	var list BeList
	reader.ReadByte()
	for {
		value, err := decodeEntity(reader)
		if err != nil {
			return nil, err
		}
		if value == nil {
			break
		}
		list.Node = append(list.Node, BeNode{value})
	}
	log.Printf("INFO: Decoded list. Returning %v", list)
	return &list, nil
}

func decodeDictionary(reader *bufio.Reader) (*BeDict, error) {
	var dict BeDict
	reader.ReadByte()
	for {
		keyEntity, err := decodeEntity(reader)
		if err != nil {
			return nil, err
		}
		if keyEntity == nil {
			break
		}
		key, ok := keyEntity.(*BeString)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Dictionary key was not a string: %v", key))
		}
		value, err := decodeEntity(reader)
		if err != nil {
			return nil, err
		}
		dict.Entry = append(dict.Entry, BeDictEntry{*key, BeNode{value}})
	}
	log.Printf("INFO: Decoded dictionary. Returning %v", dict)
	return &dict, nil
}
