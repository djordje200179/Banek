package objects

import (
	"fmt"
	"slices"
	"strings"
)

type Array []Object

func (array Array) Type() Type    { return TypeArray }
func (array Array) Clone() Object { return slices.Clone(array) }

func (array Array) Equals(other Object) bool {
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

func (array Array) AcceptsKey(key Object) bool {
	_, ok := key.(Integer)
	return ok
}

func (array Array) Get(key Object) (Object, error) {
	index := key.(Integer)
	if index < 0 || int(index) >= len(array) {
		return nil, ErrIndexOutOfBounds{int(index), len(array)}
	}

	return array[index], nil
}

func (array Array) Set(key, value Object) error {
	index := key.(Integer)
	if index < 0 || int(index) >= len(array) {
		return ErrIndexOutOfBounds{int(index), len(array)}
	}

	array[index] = value
	return nil
}

type ErrIndexOutOfBounds struct {
	Index int
	Size  int
}

func (err ErrIndexOutOfBounds) Error() string {
	return fmt.Sprintf("index out of bounds: index %d, size %d", err.Index, err.Size)
}
