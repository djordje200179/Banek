package symbols

import (
	"banek/runtime/builtins"
	"slices"
)

type Builtin int

func (b Builtin) String() string { return builtins.Funcs[b].Name }
func (b Builtin) SymbolNode()    {}

func FindBuiltin(name string) (Builtin, bool) {
	i := slices.IndexFunc(builtins.Funcs[:], func(b builtins.Builtin) bool {
		return b.Name == name
	})

	if i == -1 {
		return -1, false
	}

	return Builtin(i), true
}
