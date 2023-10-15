package objects

import "strconv"

type Integer int

func (integer Integer) Type() ObjectType { return IntegerType }

func (integer Integer) String() string {
	return strconv.Itoa(int(integer))
}
