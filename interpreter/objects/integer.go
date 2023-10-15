package objects

import "strconv"

type Integer int

func (integer Integer) Type() string   { return "integer" }
func (integer Integer) String() string { return strconv.Itoa(int(integer)) }
