package main

type BeType int

const (
	BeStringType BeType = iota
	BeIntegerType
	BeListType
	BeDictType
)

type BeNode struct {
	String     *BeString
	Integer    *BeInteger
	List       *BeList
	Dictionary *BeDict
	Type       BeType
}

type BeString struct {
	Len int
	Val []byte
}

type BeInteger struct {
	Val string
}

type BeList []BeNode

type BeDict map[string]BeNode
