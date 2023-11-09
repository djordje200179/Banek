package objs

import (
	"banek/runtime/types"
	"strconv"
)

type Int int

func (integer Int) Type() types.Type { return types.TypeInt }
func (integer Int) Clone() types.Obj { return integer }

func (integer Int) Equals(other types.Obj) bool {
	otherInt, ok := other.(Int)
	if !ok {
		return false
	}

	return integer == otherInt
}

func (integer Int) String() string { return strconv.Itoa(int(integer)) }

func (integer Int) CanAdd(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Add(other types.Obj) types.Obj {
	otherInt := other.(Int)

	return integer + otherInt
}

func (integer Int) CanSub(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Sub(other types.Obj) types.Obj {
	otherInt := other.(Int)

	return integer - otherInt
}

func (integer Int) CanMul(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Mul(other types.Obj) types.Obj {
	otherInt := other.(Int)

	return integer * otherInt
}

func (integer Int) CanDiv(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Div(other types.Obj) types.Obj {
	otherInt := other.(Int)

	return integer / otherInt
}

func (integer Int) CanMod(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Mod(other types.Obj) types.Obj {
	otherInt := other.(Int)

	return integer % otherInt
}

func (integer Int) CanPow(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Pow(other types.Obj) types.Obj {
	otherInt := other.(Int)

	result := Int(1)
	if otherInt < 0 {
		for i := Int(0); i > otherInt; i-- {
			result /= integer
		}
	} else {
		for i := Int(0); i < otherInt; i++ {
			result *= integer
		}
	}

	return result
}

func (integer Int) Negate() types.Obj {
	return -integer
}

func (integer Int) CanLess(other types.Obj) bool {
	_, ok := other.(Int)
	return ok
}

func (integer Int) Less(other types.Obj) bool {
	otherInt := other.(Int)

	return integer < otherInt
}
