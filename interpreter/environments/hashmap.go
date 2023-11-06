package environments

import (
	"banek/exec/errors"
	"banek/exec/objects"
)

type hashmapEnv struct {
	values map[string]variable

	outer Env
}

func NewHashmapEnv(outer Env, capacity int) Env {
	return &hashmapEnv{
		values: make(map[string]variable, capacity),
		outer:  outer,
	}
}

func (env *hashmapEnv) Get(name string) (objects.Object, error) {
	obj, ok := env.values[name]
	if ok {
		return obj.Object, nil
	} else if env.outer != nil {
		return env.outer.Get(name)
	} else {
		return nil, errors.ErrIdentifierNotDefined{Identifier: name}
	}
}

func (env *hashmapEnv) Define(name string, value objects.Object, mutable bool) error {
	if _, ok := env.values[name]; ok {
		return errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	env.values[name] = variable{value, mutable}

	return nil
}

func (env *hashmapEnv) Set(name string, value objects.Object) error {
	if entry, ok := env.values[name]; ok {
		if !entry.Mutable {
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

func (env *hashmapEnv) Delete(name string) error {
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

func (env *hashmapEnv) Clear() {
	clear(env.values)
}
