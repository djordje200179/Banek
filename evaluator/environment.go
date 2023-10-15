package evaluator

import (
	"fmt"
)

type Object interface {
	Type() string

	fmt.Stringer
}

type Variable struct {
	Object
	Mutable bool
}

type environment struct {
	values map[string]Variable

	outer *environment
}

func newEnvironment(outer *environment) *environment {
	return &environment{
		values: map[string]Variable{},
		outer:  outer,
	}
}

func (env *environment) Get(name string) (Object, error) {
	obj, ok := env.values[name]
	if ok {
		return obj.Object, nil
	} else if env.outer != nil {
		return env.outer.Get(name)
	} else {
		return nil, IdentifierNotDefinedError{name}
	}
}

func (env *environment) Define(name string, value Object, mutable bool) error {
	if _, ok := env.values[name]; ok {
		return IdentifierAlreadyDefinedError{name}
	}

	env.values[name] = Variable{value, mutable}

	return nil
}

func (env *environment) Set(name string, value Object) error {
	if variable, ok := env.values[name]; ok {
		if !variable.Mutable {
			return IdentifierNotMutableError{name}
		}

		env.values[name] = Variable{value, true}

		return nil
	} else if env.outer != nil {
		return env.outer.Set(name, value)
	} else {
		return IdentifierNotDefinedError{name}
	}
}

func (env *environment) Delete(name string) {
	delete(env.values, name)
}
