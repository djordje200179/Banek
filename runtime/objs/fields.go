package objs

import "fmt"

type NotIndexableError struct {
	Coll, Key Obj
}

func (err NotIndexableError) Error() string {
	return fmt.Sprintf("not indexable: %s, for key: %s", err.Coll.String(), err.Key.String())
}

type IndexOutOfBoundsError struct {
	Arr Obj

	Index int
}

func (err IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("index out of bounds: %s, for index: %d", err.Arr.String(), err.Index)
}

func (o Obj) Get(key Obj) Obj {
	switch o.Type() {
	case Array:
		arr := o.AsArray()

		switch key.Type() {
		case Int:
			key := key.AsInt()

			if key < 0 || key >= len(arr) {
				panic(IndexOutOfBoundsError{o, key})
			}

			return arr[key]
		default:
			panic(NotIndexableError{o, key})
		}
	default:
		panic(NotIndexableError{o, key})
	}
}

func (o Obj) Set(key Obj, value Obj) {
	switch o.Type() {
	case Array:
		arr := o.AsArray()
		switch key.Type() {
		case Int:
			key := key.AsInt()

			if key < 0 || key >= len(arr) {
				panic(IndexOutOfBoundsError{o, key})
			}

			arr[key] = value
		default:
			panic(NotIndexableError{o, key})
		}
	default:
		panic(NotIndexableError{o, key})
	}
}
