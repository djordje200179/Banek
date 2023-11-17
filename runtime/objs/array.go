package objs

import "strings"

type Array struct {
	Slice []Obj
}

func (arr *Array) String() string {
	var sb strings.Builder

	sb.WriteByte('[')
	for i, obj := range arr.Slice {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(obj.String())
	}
	sb.WriteByte(']')

	return sb.String()
}
