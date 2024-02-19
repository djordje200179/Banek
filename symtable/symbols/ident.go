package symbols

type Ident string

func (i Ident) String() string { return string(i) }
func (i Ident) SymbolNode()    {}
