package symbols

type Var struct {
	Name    string
	Mutable bool

	Level, Index int
}

func (v Var) String() string { return v.Name }
func (v Var) SymbolNode()    {}
