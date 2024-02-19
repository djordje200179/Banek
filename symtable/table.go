package symtable

import (
	"banek/symtable/symbols"
	"fmt"
)

type Symbol interface {
	fmt.Stringer
	SymbolNode()
}

type Table struct {
	topScope     *scope
	topContainer *Container
}

func New() *Table {
	t := &Table{}

	t.topContainer = &Container{}
	t.topScope = &scope{nil, t.topContainer, make(map[string]Symbol)}

	return t
}

func (t *Table) Insert(name string, mutable bool) (Symbol, bool) {
	return t.topScope.insertVar(name, mutable)
}

func (t *Table) Lookup(name string) (Symbol, bool) {
	for s := t.topScope; s != nil; s = s.outer {
		if ident, ok := s.lookup(name); ok {
			return ident, true
		}
	}

	return symbols.FindBuiltin(name)
}

func (t *Table) OpenScope() {
	t.topScope = &scope{t.topScope, t.topContainer, make(map[string]Symbol)}
}

func (t *Table) OpenContainer() *Container {
	t.topContainer = &Container{t.topContainer, t.topContainer.level + 1, nil}
	t.topScope = &scope{t.topScope, t.topContainer, make(map[string]Symbol)}

	return t.topContainer
}

func (t *Table) CloseScope() {
	t.topScope = t.topScope.outer
}

func (t *Table) CloseContainer() {
	t.topContainer = t.topContainer.outer
	t.topScope = t.topScope.outer
}
