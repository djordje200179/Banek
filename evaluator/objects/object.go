package objects

import "fmt"

type Object interface {
	Type() ObjectType
	fmt.Stringer
}
