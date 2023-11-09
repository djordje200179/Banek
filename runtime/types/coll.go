package types

type Coll interface {
	Size() int

	CanIndex(key Obj) bool
	Get(key Obj) (Obj, error)
	Set(key, value Obj) error
}
