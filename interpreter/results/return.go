package results

import "banek/runtime/types"

type Return struct {
	Value types.Obj
}

func (ret Return) String() string { return ret.Value.String() }
