package objects

import (
	"strconv"
)

type Boolean bool

func (boolean Boolean) Type() ObjectType { return BooleanType }

func (boolean Boolean) String() string {
	return strconv.FormatBool(bool(boolean))
}
