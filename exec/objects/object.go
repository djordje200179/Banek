package objects

import (
	"fmt"
)

type Object interface {
	Type() Type
	Clone() Object
	Equals(other Object) bool
	fmt.Stringer
}

type Coll interface {
	Object

	Size() int

	AcceptsKey(key Object) bool
	Get(key Object) (Object, error)
	Set(key, value Object) error
}
