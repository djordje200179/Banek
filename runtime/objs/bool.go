package objs

import (
	"banek/runtime/types"
	"strconv"
)

type Bool bool

func (boolean Bool) Type() types.Type { return types.TypeBool }
func (boolean Bool) Clone() types.Obj { return boolean }
func (boolean Bool) String() string   { return strconv.FormatBool(bool(boolean)) }

func (boolean Bool) Equals(other types.Obj) bool {
	otherBool, ok := other.(Bool)
	if !ok {
		return false
	}

	return boolean == otherBool
}

func (boolean Bool) Not() types.Obj {
	return !boolean
}
