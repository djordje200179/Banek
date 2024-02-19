package symtable

type scope struct {
	outer     *scope
	container *Container

	table map[string]Symbol
}

func (s *scope) insertVar(name string, mutable bool) (Symbol, bool) {
	if _, ok := s.table[name]; ok {
		return nil, false
	}

	sym := s.container.insertVar(name, mutable)
	s.table[name] = sym

	return sym, true
}

func (s *scope) lookup(name string) (Symbol, bool) {
	sym, ok := s.table[name]
	return sym, ok
}
