package exprs

import (
	"banek/symtable"
)

type Ident struct {
	symtable.Symbol
}

func (i Ident) IsConst() bool { return false }
