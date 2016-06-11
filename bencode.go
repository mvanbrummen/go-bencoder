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
	Unicoded		 byte   = 100
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

func BeDecode(reader io.Reader) {
	r := bufio.NewReader(reader)
outer:
	for {
		if b, err := r.Peek(1); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			switch b[0] {
			case Unicodei:
				DecodeInteger(r)
			case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6, Unicode7, Unicode8, Unicode9:
				DecodeString(r)
			case Unicodel:
				DecodeList(r)
			case Unicoded:
				DecodeDictionary(r)
			default:
				break outer
			}
		}
	}
}

func DecodeInteger(reader *bufio.Reader) *BeInteger {
	var str string
	if b, err := reader.ReadBytes('e'); err != nil {
		if err == io.EOF {
			panic(err)
		} else {
			log.Fatal(err)
		}
	} else {
		str = fmt.Sprintf("%s", string(b[1:len(b)-1]))
		if b, err := regexp.MatchString(BeIntPattern, str); err != nil || b == false {
			panic(fmt.Sprintf("Could not parse integer: %s", str))
		}
	}
	log.Printf("INFO: Decoded integer. Returning %v", BeInteger{str})
	return &BeInteger{str}
}

func DecodeString(reader *bufio.Reader) *BeString {
	var length int
	// get the length of the BeString
	if b, err := reader.ReadBytes(':'); err != nil {
		if err == io.EOF {
			panic(err)
		} else {
			log.Fatal(err)
		}
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
				log.Fatal(err)
			}
		} else {
			buf.WriteByte(b)
		}
	}
	log.Printf("INFO: Decoded string. Returning %v in %v", buf.String(), BeString{length, buf.Bytes()})
	return &BeString{length, buf.Bytes()}
}

func DecodeList(reader *bufio.Reader) *BeList {
	var list BeList
	reader.ReadByte()
outer:
	for {
		if b, err := reader.Peek(1); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			switch b[0] {
			case Unicodei:
				list.Node = append(list.Node, BeNode{DecodeInteger(reader)})
			case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6, Unicode7, Unicode8, Unicode9:
				list.Node = append(list.Node, BeNode{DecodeString(reader)})
			case Unicodel:
				list.Node = append(list.Node, BeNode{DecodeList(reader)})
			case Unicoded:
				list.Node = append(list.Node, BeNode{DecodeDictionary(reader)})
			default:
				break outer
			}
		}
	}
	fmt.Printf("List is: %v\n", list)
	return &list
}

func DecodeDictionary(reader *bufio.Reader) *BeDict {
	var dict BeDict
	reader.ReadByte()
	for {
		key, ok := getBencodeEntity(reader).(*BeString)
		if key == nil {
			break
		}
		if !ok {
			log.Fatal("Dict key was not a string.")
		}
		value := BeNode{getBencodeEntity(reader)}
		dict.Entry = append(dict.Entry, BeDictEntry{*key, value})
	}
	fmt.Printf("Dictionary is: %v\n", dict)
	return &dict
}

func getBencodeEntity(reader *bufio.Reader) interface{} {
	var bencodeEntity interface{}
	if b, err := reader.Peek(1); err != nil {
		if err == io.EOF {
			return nil
		} else {
			log.Fatal(err)
		}
	} else {
		switch b[0] {
		case Unicodei:
			bencodeEntity = DecodeInteger(reader)
		case Unicode0, Unicode1, Unicode2, Unicode3, Unicode4, Unicode5, Unicode6, Unicode7, Unicode8, Unicode9:
			bencodeEntity = DecodeString(reader)
		case Unicodel:
			bencodeEntity = DecodeList(reader)
		case Unicoded:
			bencodeEntity = DecodeDictionary(reader)
		default:
			bencodeEntity = nil
		}
	}
	return bencodeEntity
}
