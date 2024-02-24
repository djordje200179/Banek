package objs

import "fmt"

type NotIndexableError struct {
	Coll, Key Obj
}

func (err NotIndexableError) Error() string {
	return fmt.Sprintf("not indexable: %s, for key: %s", err.Coll.String(), err.Key.String())
}

func (o Obj) Get(key Obj) (Obj, error) {
	err := NotIndexableError{o, key}

	var elem Obj
	switch o.Type() {
	case Array:
		arr := o.AsArray()

		switch key.Type() {
		case Int:
			key := key.AsInt()

			if key < 0 || key >= len(arr) {
				return Obj{}, err
			}

			elem = arr[key]
		default:
			return Obj{}, err
		}
	default:
		return Obj{}, err
	}

	return elem, nil
}

func (o Obj) Set(key Obj, value Obj) error {
	err := NotIndexableError{o, key}

	switch o.Type() {
	case Array:
		arr := o.AsArray()
		switch key.Type() {
		case Int:
			key := key.AsInt()

			if key < 0 || key >= len(arr) {
				return err
			}

			arr[key] = value
		default:
			return err
		}
	default:
		return err
	}

	return nil
}
