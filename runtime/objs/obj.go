package objs

import (
	"unsafe"
)

type Obj struct {
	Ptr unsafe.Pointer

	typeAndNum int
}

func (o Obj) Type() Type {
	return Type(o.typeAndNum & TypeMask)
}

func (o Obj) getNum() int {
	return o.typeAndNum >> TypeBits
}

func Make(t Type, ptr unsafe.Pointer, i int) Obj {
	return Obj{Ptr: ptr, typeAndNum: int(t) | (i << TypeBits)}
}

func MakeString(s string) Obj {
	strPtr := unsafe.StringData(s)
	strLen := len(s)

	return Make(String, unsafe.Pointer(strPtr), strLen)
}

func MakeInt(i int) Obj {
	return Make(Int, nil, i)
}

func MakeBool(b bool) Obj {
	var i int
	if b {
		i = 1
	}

	return Make(Bool, nil, i)
}

func MakeArray(s []Obj) Obj {
	arrPtr := unsafe.SliceData(s)
	arrLen := len(s)

	return Make(Array, unsafe.Pointer(arrPtr), arrLen)
}

func (o Obj) AsString() string {
	return unsafe.String((*byte)(o.Ptr), o.getNum())
}

func (o Obj) AsArray() []Obj {
	return unsafe.Slice((*Obj)(o.Ptr), o.getNum())
}

func (o Obj) AsInt() int {
	return o.getNum()
}

func (o Obj) AsBool() bool {
	return o.getNum() != 0
}

func (o Obj) Truthy() bool {
	switch o.Type() {
	case Int:
		return o.AsInt() != 0
	case Bool:
		return o.AsBool()
	case String:
		return o.getNum() > 0
	case Array:
		return o.getNum() > 0
	case Func:
		return true
	case Builtin:
		return true
	default:
		return false
	}
}
