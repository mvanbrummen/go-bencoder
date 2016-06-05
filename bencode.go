package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
)

func BeDecode(reader io.Reader) {
	r := bufio.NewReader(reader)
	for {
		if c, sz, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			switch c {
			case 'i':
				DecodeInteger(r)
			}
			fmt.Printf("%q [%d]\n", string(c), sz)
		}
	}
}

func DecodeInteger(reader *bufio.Reader) *BeInteger {
	var buf bytes.Buffer
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			switch c {
			case 'i':
				continue
			case 'e':
				break
			default:
				buf.WriteRune(c)
			}
		}
	}
	log.Println("INFO: Decoded integer. Returning", BeInteger{buf.String()})
	return &BeInteger{buf.String()}
}
