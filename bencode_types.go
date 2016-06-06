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
	Key  BeString
	Node BeNode
}

type BeNode struct {
	Val interface{}
}
