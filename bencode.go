package main

import (
	"io"
	"bufio"
	"log"
	"fmt"
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
                        fmt.Printf("%q [%d]\n", string(c), sz)
                }
        }
}

