package unaryops

import "banek/runtime/objs"

func negateInt(o objs.Obj) (objs.Obj, bool) {
	return objs.MakeInt(-o.AsInt()), true
}

func invertBool(o objs.Obj) (objs.Obj, bool) {
	return objs.MakeBool(!o.AsBool()), true
}
