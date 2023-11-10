package types

type Adder interface {
	Add(other Obj) (Obj, bool)
}

type Subber interface {
	Sub(other Obj) (Obj, bool)
}

type Multer interface {
	Mul(other Obj) (Obj, bool)
}

type Diver interface {
	Div(other Obj) (Obj, bool)
}

type Modder interface {
	Mod(other Obj) (Obj, bool)
}

type Powwer interface {
	Pow(other Obj) (Obj, bool)
}

type Negater interface {
	Negate() Obj
}

type Notter interface {
	Not() Obj
}

type Lesser interface {
	Less(other Obj) (less, ok bool)
}
