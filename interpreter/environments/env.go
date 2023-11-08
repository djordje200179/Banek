package environments

import (
	"banek/exec/errors"
	"banek/exec/objects"
	"slices"
	"sync"
)

type variable struct {
	objects.Object
	Mutable bool
}

type Env struct {
	keys   []string
	values []variable

	outer *Env
}

var envPool = sync.Pool{
	New: func() any {
		return &Env{}
	},
}

func New(outer *Env, capacity int) *Env {
	//env := envPool.Get().(*Env)
	env := new(Env)

	env.keys = make([]string, 0, capacity)
	env.values = make([]variable, 0, capacity)
	env.outer = outer

	return env
}

func (env *Env) Get(name string) (objects.Object, error) {
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

func (env *Env) Define(name string, value objects.Object, mutable bool) error {
	if slices.Index(env.keys, name) != -1 {
		return errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	env.keys = append(env.keys, name)
	env.values = append(env.values, variable{value, mutable})

	return nil
}

func (env *Env) Set(name string, value objects.Object) error {
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

func (env *Env) Delete(name string) error {
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

func (env *Env) Clear() {
	env.keys = nil
	env.values = nil
}
