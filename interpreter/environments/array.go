package environments

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"slices"
)

type arrayEnvironment struct {
	keys   []string
	values []variable

	outer Environment
}

func NewArrayEnvironment(outer Environment, capacity int) Environment {
	return &arrayEnvironment{
		keys:   make([]string, 0, capacity),
		values: make([]variable, 0, capacity),
		outer:  outer,
	}
}

func (env *arrayEnvironment) Get(name string) (objects.Object, error) {
	index := slices.Index(env.keys, name)
	if index == -1 {
		if env.outer != nil {
			return env.outer.Get(name)
		} else {
			return nil, errors.ErrIdentifierNotDefined{Identifier: name}
		}
	}

	return env.values[index].Object, nil
}

func (env *arrayEnvironment) Define(name string, value objects.Object, mutable bool) error {
	if slices.Index(env.keys, name) != -1 {
		return errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	env.keys = append(env.keys, name)
	env.values = append(env.values, variable{value, mutable})

	return nil
}

func (env *arrayEnvironment) Set(name string, value objects.Object) error {
	index := slices.Index(env.keys, name)
	if index == -1 {
		if env.outer != nil {
			return env.outer.Set(name, value)
		} else {
			return errors.ErrIdentifierNotDefined{Identifier: name}
		}
	}

	if !env.values[index].Mutable {
		return errors.ErrIdentifierNotMutable{Identifier: name}
	}

	env.values[index].Object = value

	return nil
}

func (env *arrayEnvironment) Delete(name string) error {
	index := slices.Index(env.keys, name)
	if index == -1 {
		if env.outer != nil {
			return env.outer.Delete(name)
		} else {
			return errors.ErrIdentifierNotDefined{Identifier: name}
		}
	}

	env.keys = slices.Delete(env.keys, index, index+1)
	env.values = slices.Delete(env.values, index, index+1)

	return nil
}

func (env *arrayEnvironment) Clear() {
	env.keys = nil
	env.values = nil
}
