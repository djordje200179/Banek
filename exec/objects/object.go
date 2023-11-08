package objects

import (
	"fmt"
)

type Object interface {
	Type() Type
	Clone() Object

	fmt.Stringer
}
