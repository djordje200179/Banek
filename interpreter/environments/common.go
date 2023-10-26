package environments

import "banek/exec/objects"

type variable struct {
	objects.Object
	Mutable bool
}

type Environment interface {
	Get(name string) (objects.Object, error)
	Set(name string, value objects.Object) error
	Define(name string, value objects.Object, mutable bool) error
	Delete(name string) error
}

type EnvironmentFactory func(outer Environment, capacity int) Environment
