package interpreter

import (
	errors2 "banek/exec/errors"
	"banek/exec/objects"
)

type variable struct {
	objects.Object
	Mutable bool
}

type environment struct {
	values map[string]variable

	outer *environment
}

func newEnvironment(outer *environment) *environment {
	return &environment{
		values: map[string]variable{},
		outer:  outer,
	}
}

func (env *environment) Get(name string) (objects.Object, error) {
	obj, ok := env.values[name]
	if ok {
		return obj.Object, nil
	} else if env.outer != nil {
		return env.outer.Get(name)
	} else {
		return nil, errors2.IdentifierNotDefinedError{Identifier: name}
	}
}

func (env *environment) Define(name string, value objects.Object, mutable bool) error {
	if _, ok := env.values[name]; ok {
		return errors2.IdentifierAlreadyDefinedError{Identifier: name}
	}

	env.values[name] = variable{value, mutable}

	return nil
}

func (env *environment) Set(name string, value objects.Object) error {
	if varEntry, ok := env.values[name]; ok {
		if !varEntry.Mutable {
			return errors2.IdentifierNotMutableError{Identifier: name}
		}

		env.values[name] = variable{value, true}

		return nil
	} else if env.outer != nil {
		return env.outer.Set(name, value)
	} else {
		return errors2.IdentifierNotDefinedError{Identifier: name}
	}
}

func (env *environment) Delete(name string) {
	delete(env.values, name)
}
