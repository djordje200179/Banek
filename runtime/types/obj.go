package types

import "fmt"

type Obj interface {
	Type() Type
	Clone() Obj
	Equals(other Obj) bool
	fmt.Stringer
}
