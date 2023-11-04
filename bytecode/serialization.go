package bytecode

import (
	"banek/exec/objects"
	"encoding/gob"
)

func init() {
	gob.Register(objects.Array{})
	gob.Register(objects.Boolean(false))
	gob.Register(objects.BuiltinFunction{})
	gob.Register(objects.Integer(0))
	gob.Register(objects.String(""))
	gob.Register(objects.Undefined{})
	gob.Register(objects.Unknown)
	gob.Register(&Function{})
}
