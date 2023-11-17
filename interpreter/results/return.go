package results

import (
	"banek/runtime/objs"
)

type Return struct {
	Value objs.Obj
}

func (ret Return) String() string { return ret.Value.String() }
