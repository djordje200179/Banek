package objs

import (
	"bytes"
	"encoding/binary"
	"slices"
	"strconv"
	"unsafe"
)

type Obj struct {
	PtrData unsafe.Pointer
	IntData uint64

	Tag Tag
}

func MakeInt(value int) Obj {
	return Obj{Tag: TypeInt, IntData: uint64(value)}
}

func MakeBool(value bool) Obj {
	var intData uint64
	if value {
		intData = 1
	}

	return Obj{Tag: TypeBool, IntData: intData}
}

func MakeStr(value string) Obj {
	ptrData := unsafe.Pointer(unsafe.StringData(value))
	intData := uint64(len(value))

	return Obj{Tag: TypeStr, PtrData: ptrData, IntData: intData}
}

func MakeArray(value *Array) Obj {
	ptrData := unsafe.Pointer(value)

	return Obj{Tag: TypeArray, PtrData: ptrData}
}

func (obj Obj) AsInt() int {
	return int(obj.IntData)
}

func (obj Obj) AsBool() bool {
	return obj.IntData != 0
}

func (obj Obj) AsStr() string {
	ptr := (*byte)(obj.PtrData)
	size := int(obj.IntData)

	return unsafe.String(ptr, size)
}

func (obj Obj) AsArray() *Array {
	return (*Array)(obj.PtrData)
}

func (obj Obj) String() string {
	switch obj.Tag {
	case TypeInt:
		integer := obj.AsInt()
		return strconv.Itoa(integer)
	case TypeBool:
		boolean := obj.AsBool()
		return strconv.FormatBool(boolean)
	case TypeStr:
		str := obj.AsStr()
		return str
	case TypeUndefined:
		return "undefined"
	case TypeArray:
		arr := obj.AsArray()
		return arr.String()
	default:
		return Config[obj.Tag].Stringer(obj)
	}
}

func (obj Obj) Equals(other Obj) bool {
	if obj.Tag != other.Tag {
		return false
	}

	switch obj.Tag {
	case TypeInt:
		return obj.AsInt() == other.AsInt()
	case TypeBool:
		return obj.AsBool() == other.AsBool()
	case TypeStr:
		return obj.AsStr() == other.AsStr()
	case TypeArray:
		firstArr := obj.AsArray()
		secondArr := other.AsArray()

		return slices.EqualFunc(firstArr.Slice, secondArr.Slice, Obj.Equals)
	case TypeUndefined:
		return true
	default:
		return Config[obj.Tag].Equaler(obj, other)
	}
}

func (obj Obj) Clone() Obj {
	switch obj.Tag {
	case TypeArray:
		arr := obj.AsArray()

		newArr := new(Array)
		newArr.Slice = slices.Clone(arr.Slice)

		return MakeArray(newArr)
	default:
		return obj
	}
}

func (obj Obj) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	buf.WriteByte(byte(obj.Tag))
	switch obj.Tag {
	case TypeUndefined:
	case TypeBool:
		var intData uint8
		if obj.AsBool() {
			intData = 1
		}

		err = binary.Write(&buf, binary.LittleEndian, intData)
		if err != nil {
			return
		}
	case TypeInt:
		err = binary.Write(&buf, binary.LittleEndian, obj.IntData)
		if err != nil {
			return
		}
	case TypeStr:
		strBytes := []byte(obj.AsStr())
		err = binary.Write(&buf, binary.LittleEndian, strBytes)
		if err != nil {
			return
		}
	case TypeArray:
		arr := obj.AsArray()
		err = binary.Write(&buf, binary.LittleEndian, arr.Slice)
		if err != nil {
			return
		}
	default:
		var objBuf []byte
		objBuf, err = Config[TypeFunc].Marshaller(obj)
		if err != nil {
			return
		}

		buf.Write(objBuf)
	}

	data = buf.Bytes()
	return
}

func (obj *Obj) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	tag, err := buf.ReadByte()
	if err != nil {
		return err
	}
	obj.Tag = Tag(tag)

	switch obj.Tag {
	case TypeUndefined:
	case TypeBool:
		var intData uint8
		err = binary.Read(buf, binary.LittleEndian, &intData)
		if err != nil {
			return err
		}

		if intData == 1 {
			obj.IntData = 1
		}
	case TypeInt:
		err = binary.Read(buf, binary.LittleEndian, &obj.IntData)
		if err != nil {
			return err
		}
	case TypeStr:
		var str string
		err = binary.Read(buf, binary.LittleEndian, &str)
		if err != nil {
			return err
		}

		obj.PtrData = unsafe.Pointer(unsafe.StringData(str))
		obj.IntData = uint64(len(str))
	case TypeArray:
		arr := new(Array)
		err = binary.Read(buf, binary.LittleEndian, &arr.Slice)
		if err != nil {
			return err
		}
	default:
		var err error
		*obj, err = Config[TypeFunc].Unmarshaller(buf)
		if err != nil {
			return err
		}
	}

	return nil
}
