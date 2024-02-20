package objs

import "fmt"

type NotIndexableError struct {
	Coll, Key Obj
}

func (err NotIndexableError) Error() string {
	return fmt.Sprintf("not indexable: %s, for key: %s", err.Coll.String(), err.Key.String())
}

func (o Obj) GetField(key Obj) (Obj, error) {
	err := NotIndexableError{o, key}

	var elem Obj
	switch o.Type {
	case Array:
		switch key.Type {
		case Int:
			arr := o.AsArray()

			if key.Int < 0 || key.Int >= len(arr) {
				return Obj{}, err
			}

			elem = arr[key.Int]
		default:
			return Obj{}, err
		}
	default:
		return Obj{}, err
	}

	return elem, nil
}

func (o Obj) SetField(key Obj, value Obj) error {
	err := NotIndexableError{o, key}

	switch o.Type {
	case Array:
		switch key.Type {
		case Int:
			arr := o.AsArray()

			if key.Int < 0 || key.Int >= len(arr) {
				return err
			}

			arr[key.Int] = value
		default:
			return err
		}
	default:
		return err
	}

	return nil
}
