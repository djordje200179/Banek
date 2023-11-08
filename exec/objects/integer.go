package objects

import "strconv"

type Integer int

func (integer Integer) Type() Type    { return TypeInteger }
func (integer Integer) Clone() Object { return integer }

func (integer Integer) Equals(other Object) bool {
	otherInteger, ok := other.(Integer)
	if !ok {
		return false
	}

	return integer == otherInteger
}

func (integer Integer) String() string { return strconv.Itoa(int(integer)) }
