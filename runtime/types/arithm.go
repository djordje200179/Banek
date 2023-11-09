package types

type Adder interface {
	CanAdd(other Obj) bool
	Add(other Obj) Obj
}

type Subber interface {
	CanSub(other Obj) bool
	Sub(other Obj) Obj
}

type Multer interface {
	CanMul(other Obj) bool
	Mul(other Obj) Obj
}

type Diver interface {
	CanDiv(other Obj) bool
	Div(other Obj) Obj
}

type Modder interface {
	CanMod(other Obj) bool
	Mod(other Obj) Obj
}

type Powwer interface {
	CanPow(other Obj) bool
	Pow(other Obj) Obj
}

type Negater interface {
	Negate() Obj
}

type Notter interface {
	Not() Obj
}

type Lesser interface {
	CanLess(other Obj) bool
	Less(other Obj) bool
}
