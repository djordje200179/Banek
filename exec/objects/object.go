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
