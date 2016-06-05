package main

import (
	_"fmt"
)

type BeString struct {
	Len int
	Val []byte
}

type BeInteger struct {
	Val string
}

type BeList struct {
	Node []BeNode
}

type BeDict struct {
	Nodes []BeDictEntry
}

type BeDictEntry struct {
	Key  string
	Node BeNode
}

type BeNode struct {
	value interface{}
}
