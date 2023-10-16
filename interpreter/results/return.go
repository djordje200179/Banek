package results

import (
	"banek/exec/objects"
)

type Return struct {
	Value objects.Object
}

func (ret Return) String() string { return ret.Value.String() }
