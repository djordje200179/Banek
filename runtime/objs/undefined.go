package objs

import "banek/runtime/types"

type Undefined struct{}

func (undefined Undefined) Type() types.Type { return types.TypeUndefined }
func (undefined Undefined) Clone() types.Obj { return undefined }
func (undefined Undefined) String() string   { return "<undefined>" }

func (undefined Undefined) Equals(other types.Obj) bool {
	_, ok := other.(Undefined)
	return ok
}
