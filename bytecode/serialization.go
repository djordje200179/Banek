package bytecode

import (
	"banek/runtime/builtins"
	"banek/runtime/objs"
	"encoding/gob"
)

func init() {
	gob.Register(objs.Array{})
	gob.Register(objs.Bool(false))
	gob.Register(builtins.BuiltinFunc{})
	gob.Register(objs.Int(0))
	gob.Register(objs.Str(""))
	gob.Register(objs.Undefined{})
	gob.Register(&Func{})
}
