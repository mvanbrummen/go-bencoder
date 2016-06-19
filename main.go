package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Need to provide a .torrent file!\n")
		os.Exit(1)
	}
	filePath := os.Args[1]
	fmt.Println(filePath)

	// check that file exists
	if _, err := os.Stat(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "%s doesn't exist\n", filePath)
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file\n")
		os.Exit(1)
	}
	dict, err := BeDecode(b)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	fmt.Printf("dict is %v\nErr is %v\n", dict, err)
	for k, v := range *dict {
		fmt.Printf("%s is %+v\n", k, v)
	}
}
