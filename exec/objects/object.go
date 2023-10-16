package objects

import (
	"fmt"
)

type Object interface {
	Type() string

	fmt.Stringer
}
