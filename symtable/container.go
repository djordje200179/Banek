package symtable

import "banek/symtable/symbols"

type Container struct {
	outer *Container
	level int

	vars []Symbol
}

func (c *Container) insertVar(name string, mutable bool) Symbol {
	sym := symbols.Var{
		Name: name, Mutable: mutable,
		Index: len(c.vars), Level: c.level,
	}

	c.vars = append(c.vars, sym)

	return sym
}
