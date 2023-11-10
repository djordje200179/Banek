package objs

import (
	"banek/runtime/types"
	"fmt"
	"slices"
	"strings"
)

type Array []types.Obj

func (array Array) Type() types.Type { return types.TypeArray }
func (array Array) Clone() types.Obj { return slices.Clone(array) }

func (array Array) Equals(other types.Obj) bool {
	otherArray, ok := other.(Array)
	if !ok {
		return false
	}

	return slices.Equal(array, otherArray)
}

func (array Array) String() string {
	var sb strings.Builder

	elements := make([]string, len(array))
	for i, element := range array {
		elements[i] = element.String()
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteByte(']')

	return sb.String()
}

func (array Array) Size() int {
	return len(array)
}

func (array Array) CanIndex(key types.Obj) bool {
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

func (array Array) Get(key types.Obj) (types.Obj, error) {
	index := key.(Int)
	if index < 0 || int(index) >= len(array) {
		return nil, ErrIndexOutOfBounds{int(index), len(array)}
	}

	return array[index], nil
}

func (array Array) Set(key, value types.Obj) error {
	index := key.(Int)
	if index < 0 || int(index) >= len(array) {
		return ErrIndexOutOfBounds{int(index), len(array)}
	}

	array[index] = value
	return nil
}

func (array Array) Add(other types.Obj) (types.Obj, bool) {
	otherArray, ok := other.(Array)
	if !ok {
		return nil, false
	}

	combinedArray := make(Array, len(array)+len(otherArray))
	copy(combinedArray, array)
	copy(combinedArray[len(array):], otherArray)

	return combinedArray, true
}

func (array Array) Mul(other types.Obj) (types.Obj, bool) {
	count, ok := other.(Int)
	if !ok {
		return nil, false
	}

	if count < 0 {
		return Array(nil), true
	}

	combinedArray := make(Array, len(array)*int(count))

	for i := 0; i < int(count); i++ {
		copy(combinedArray[i*len(array):], array)
	}

	return combinedArray, true
}
