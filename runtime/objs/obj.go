package objs

import (
	"unsafe"
)

type Obj struct {
	Ptr unsafe.Pointer
	Int int

	Type Type
}

func MakeString(s string) Obj {
	strPtr := unsafe.StringData(s)
	strLen := len(s)

	return Obj{Ptr: unsafe.Pointer(strPtr), Int: strLen, Type: String}
}

func MakeInt(i int) Obj {
	return Obj{Int: i, Type: Int}
}

func MakeBool(b bool) Obj {
	var i int
	if b {
		i = 1
	}

	return Obj{Int: i, Type: Bool}
}

func MakeArray(s []Obj) Obj {
	arrPtr := unsafe.SliceData(s)
	arrLen := len(s)

	return Obj{Ptr: unsafe.Pointer(arrPtr), Int: arrLen, Type: Array}
}

func (o Obj) AsString() string {
	return unsafe.String((*byte)(o.Ptr), o.Int)
}

func (o Obj) AsArray() []Obj {
	return unsafe.Slice((*Obj)(o.Ptr), o.Int)
}

func (o Obj) Truthy() bool {
	switch o.Type {
	case Int:
		return o.Int != 0
	case Bool:
		return o.Int != 0
	case String:
		return o.Int > 0
	case Array:
		return o.Int > 0
	case Func:
		return true
	case Builtin:
		return true
	default:
		return false
	}
}
