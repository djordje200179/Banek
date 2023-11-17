package bytecode

import (
	"banek/runtime/objs"
	"bytes"
	"encoding/binary"
	"slices"
	"strconv"
	"unsafe"
)

type Func struct {
	TemplateIndex int

	Captures []*objs.Obj
}

func GetFunc(obj objs.Obj) *Func {
	return (*Func)(obj.PtrData)
}

func (function *Func) MakeObj() objs.Obj {
	ptrData := unsafe.Pointer(function)

	return objs.Obj{Tag: objs.TypeFunc, PtrData: ptrData}
}

func funcString(obj objs.Obj) string {
	function := GetFunc(obj)

	return "func#" + strconv.Itoa(function.TemplateIndex)
}

func funcEquals(first, second objs.Obj) bool {
	firstFunc := GetFunc(first)
	secondFunc := GetFunc(second)

	if firstFunc.TemplateIndex != secondFunc.TemplateIndex {
		return false
	}

	if !slices.Equal(firstFunc.Captures, secondFunc.Captures) {
		return false
	}

	return true
}

func funcMarsal(obj objs.Obj) ([]byte, error) {
	function := GetFunc(obj)

	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, uint64(function.TemplateIndex))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func funcUnmarshal(buf *bytes.Buffer) (objs.Obj, error) {
	var templateIndex uint64
	err := binary.Read(buf, binary.LittleEndian, &templateIndex)
	if err != nil {
		return objs.Obj{}, err
	}

	function := Func{
		TemplateIndex: int(templateIndex),
	}

	return function.MakeObj(), nil
}

func init() {
	objs.Config[objs.TypeFunc] = objs.TypeConfig{
		Stringer: funcString,
		Equaler:  funcEquals,

		Marshaller:   funcMarsal,
		Unmarshaller: funcUnmarshal,
	}
}
