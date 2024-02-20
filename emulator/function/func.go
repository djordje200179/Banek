package function

import (
	"banek/runtime/objs"
)

type Obj struct {
	Index int

	Captures []*objs.Obj
}
