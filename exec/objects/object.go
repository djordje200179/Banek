package objects

import (
	"fmt"
)

type Object interface {
	Type() string
	Clone() Object

	fmt.Stringer
}
