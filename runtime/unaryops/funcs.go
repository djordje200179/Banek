package unaryops

import "banek/runtime/objs"

func negateInt(o objs.Obj) (objs.Obj, bool) {
	return objs.Obj{Int: -o.Int, Type: objs.Int}, true
}

func invertBool(o objs.Obj) (objs.Obj, bool) {
	return objs.Obj{Int: ^o.Int, Type: objs.Bool}, true
}
