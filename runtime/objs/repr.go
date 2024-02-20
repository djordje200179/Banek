package objs

import (
	"strconv"
	"strings"
)

func (o Obj) String() string {
	switch o.Type {
	case Int:
		return strconv.Itoa(o.Int)
	case Bool:
		return strconv.FormatBool(o.Int != 0)
	case String:
		return o.AsString()
	case Array:
		var sb strings.Builder
		sb.WriteByte('[')
		for i, v := range o.AsArray() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(v.String())
		}
		sb.WriteByte(']')

		return sb.String()
	case Func:
		return "func"
	case Builtin:
		return "builtin"
	default:
		return "undefined"
	}
}
