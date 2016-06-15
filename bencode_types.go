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

func (b BeNode) IsNil() bool {
	return b.String == nil && b.Integer == nil && b.List == nil && b.List == nil
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
