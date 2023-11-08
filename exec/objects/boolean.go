package objects

import "strconv"

type Boolean bool

func (boolean Boolean) Type() Type     { return TypeBoolean }
func (boolean Boolean) Clone() Object  { return boolean }
func (boolean Boolean) String() string { return strconv.FormatBool(bool(boolean)) }

func (boolean Boolean) Equals(other Object) bool {
	otherBoolean, ok := other.(Boolean)
	if !ok {
		return false
	}

	return boolean == otherBoolean
}
