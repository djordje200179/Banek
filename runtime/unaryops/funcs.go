package unaryops

import "banek/runtime/objs"

func negateInt(o objs.Obj) objs.Obj {
	return objs.MakeInt(-o.AsInt())
}

func invertBool(o objs.Obj) objs.Obj {
	return objs.MakeBool(!o.AsBool())
}
