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

func (array Array) CanAdd(other types.Obj) bool {
	_, ok := other.(Array)
	return ok
}

func (array Array) Add(other types.Obj) types.Obj {
	combinedArray := make(Array, len(array)+len(other.(Array)))

	copy(combinedArray, array)
	copy(combinedArray[len(array):], other.(Array))

	return combinedArray
}

func (array Array) CanMul(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (array Array) Mul(other types.Obj) types.Obj {
	count := int(other.(Int))
	if count < 0 {
		return Array(nil)
	}

	combinedArray := make(Array, len(array)*count)

	for i := 0; i < count; i++ {
		copy(combinedArray[i*len(array):], array)
	}

	return combinedArray
}
