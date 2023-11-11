package objs

import (
	"banek/runtime/types"
	"fmt"
	"slices"
	"strings"
)

type Array struct {
	Slice []types.Obj
}

func (array *Array) Type() types.Type { return types.TypeArray }

func (array *Array) Clone() types.Obj {
	return &Array{Slice: slices.Clone(array.Slice)}
}

func (array *Array) Equals(other types.Obj) bool {
	otherArray, ok := other.(*Array)
	if !ok {
		return false
	}

	return slices.Equal(array.Slice, otherArray.Slice)
}

func (array *Array) String() string {
	var sb strings.Builder

	elements := make([]string, len(array.Slice))
	for i, element := range array.Slice {
		elements[i] = element.String()
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteByte(']')

	return sb.String()
}

func (array *Array) Size() int {
	return len(array.Slice)
}

func (array *Array) CanIndex(key types.Obj) bool {
	_, ok := key.(Int)
	return ok
}

type ErrIndexOutOfBounds struct {
	Index int
	Size  int
}

func (err ErrIndexOutOfBounds) Error() string {
	return fmt.Sprintf("index out of bounds: index %d, size %d", err.Index, err.Size)
}

func (array *Array) Get(key types.Obj) (types.Obj, error) {
	index := key.(Int)
	if index < 0 || int(index) >= len(array.Slice) {
		return nil, ErrIndexOutOfBounds{int(index), len(array.Slice)}
	}

	return array.Slice[index], nil
}

func (array *Array) Set(key, value types.Obj) error {
	index := key.(Int)
	if index < 0 || int(index) >= len(array.Slice) {
		return ErrIndexOutOfBounds{int(index), len(array.Slice)}
	}

	array.Slice[index] = value
	return nil
}

func (array *Array) Add(other types.Obj) (types.Obj, bool) {
	otherArray, ok := other.(*Array)
	if !ok {
		return nil, false
	}

	newArray := &Array{
		Slice: make([]types.Obj, len(array.Slice)+len(otherArray.Slice)),
	}

	copy(newArray.Slice, array.Slice)
	copy(newArray.Slice[len(array.Slice):], otherArray.Slice)

	return newArray, true
}

func (array *Array) Mul(other types.Obj) (types.Obj, bool) {
	count, ok := other.(Int)
	if !ok {
		return nil, false
	}

	if count < 0 {
		return nil, false // TODO: Return invalid argument error
	}

	newArray := &Array{
		Slice: make([]types.Obj, len(array.Slice)*int(count)),
	}

	for i := 0; i < int(count); i++ {
		copy(newArray.Slice[i*len(array.Slice):], array.Slice)
	}

	return newArray, true
}

func (array *Array) Receive(other types.Obj) (types.Obj, bool) {
	array.Slice = append(array.Slice, other)

	return array, true
}

func (array *Array) Give() types.Obj {
	if len(array.Slice) == 0 {
		return Undefined{}
	}

	element := array.Slice[0]
	array.Slice = array.Slice[1:]

	return element
}
