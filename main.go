package main

import (
	"fmt"
	"io/ioutil"
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
	entity, err := BeDecode(b)
	fmt.Printf("Entity is %v\nErr is %v\n", entity, err)
}
