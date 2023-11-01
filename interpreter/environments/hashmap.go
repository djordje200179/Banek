package environments

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

type hashmapEnvironment struct {
	values map[string]variable

	outer Environment
}

func NewHashmapEnvironment(outer Environment, capacity int) Environment {
	return &hashmapEnvironment{
		values: make(map[string]variable, capacity),
		outer:  outer,
	}
}

func (env *hashmapEnvironment) Get(name string) (objects.Object, error) {
	obj, ok := env.values[name]
	if ok {
		return obj.Object, nil
	} else if env.outer != nil {
		return env.outer.Get(name)
	} else {
		return nil, errors.ErrIdentifierNotDefined{Identifier: name}
	}
}

func (env *hashmapEnvironment) Define(name string, value objects.Object, mutable bool) error {
	if _, ok := env.values[name]; ok {
		return errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	env.values[name] = variable{value, mutable}

	return nil
}

func (env *hashmapEnvironment) Set(name string, value objects.Object) error {
	if varEntry, ok := env.values[name]; ok {
		if !varEntry.Mutable {
			return errors.ErrIdentifierNotMutable{Identifier: name}
		}

		env.values[name] = variable{value, true}

		return nil
	} else if env.outer != nil {
		return env.outer.Set(name, value)
	} else {
		return errors.ErrIdentifierNotDefined{Identifier: name}
	}
}

func (env *hashmapEnvironment) Delete(name string) error {
	if _, ok := env.values[name]; !ok {
		if env.outer != nil {
			return env.outer.Delete(name)
		} else {
			return errors.ErrIdentifierNotDefined{Identifier: name}
		}
	}

	delete(env.values, name)

	return nil
}

func (env *hashmapEnvironment) Clear() {
	clear(env.values)
}
