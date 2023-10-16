package objects

import (
	"fmt"
	"strings"
)

type Array []Object

func (array Array) Type() string { return "array" }

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

type IndexOutOfBoundsError struct {
	Index int
	Size  int
}

func (err IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("index out of bounds: index %d, size %d", err.Index, err.Size)
}
