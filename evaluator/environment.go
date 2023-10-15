package evaluator

import (
	"fmt"
)

type Object interface {
	Type() string

	fmt.Stringer
}

type environment struct {
	values map[string]Object

	outer *environment
}

func newEnvironment(outer *environment) *environment {
	return &environment{
		values: map[string]Object{},
		outer:  outer,
	}
}

func (env *environment) Get(name string) (Object, bool) {
	obj, ok := env.values[name]
	if !ok && env.outer != nil {
		return env.outer.Get(name)
	}

	return obj, ok
}

func (env *environment) Set(name string, value Object) {
	env.values[name] = value
}

func (env *environment) IsDefined(name string) bool {
	_, ok := env.values[name]

	return ok
}

func (env *environment) Delete(name string) {
	delete(env.values, name)
}
