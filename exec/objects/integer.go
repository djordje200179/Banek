package objects

import "strconv"

type Integer int

func (integer Integer) Type() Type     { return TypeInteger }
func (integer Integer) Clone() Object  { return integer }
func (integer Integer) String() string { return strconv.Itoa(int(integer)) }
