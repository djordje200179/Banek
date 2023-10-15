package results

import "banek/interpreter/objects"

type Return struct {
	Value objects.Object
}

func (ret Return) String() string { return ret.Value.String() }
