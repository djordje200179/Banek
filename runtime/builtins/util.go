package builtins

import (
	"banek/runtime/objs"
	"banek/runtime/types"
	"os"
)

func builtinExit(_ []types.Obj) (types.Obj, error) {
	os.Exit(0)
	return objs.Undefined{}, nil
}
