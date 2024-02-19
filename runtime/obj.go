package runtime

import "fmt"

type Obj interface {
	fmt.Stringer
	Truthy() bool
	Clone() Obj
	Equals(other Obj) bool
}

type Coll interface {
	Obj
	Len() int
	Get(index Obj) (Obj, bool)
	Set(index Obj, value Obj) bool
}

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

type Comparer interface {
	Compare(other Obj) (int, bool)
}

type Negator interface {
	Neg() (Obj, bool)
}

type Notter interface {
	Not() (Obj, bool)
}
