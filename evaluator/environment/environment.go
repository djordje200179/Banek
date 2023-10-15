package environment

import "banek/evaluator/objects"

type Environment struct {
	values map[string]objects.Object

	outer *Environment
}

func New(outer *Environment) *Environment {
	return &Environment{
		values: map[string]objects.Object{},
		outer:  outer,
	}
}

func (env *Environment) Get(name string) (objects.Object, bool) {
	obj, ok := env.values[name]
	if !ok && env.outer != nil {
		return env.outer.Get(name)
	}

	return obj, ok
}

func (env *Environment) Set(name string, value objects.Object) {
	env.values[name] = value
}

func (env *Environment) IsDefined(name string) bool {
	_, ok := env.values[name]

	return ok
}

func (env *Environment) Delete(name string) {
	delete(env.values, name)
}
