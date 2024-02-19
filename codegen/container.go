package codegen

type container struct {
	level, index int
	vars         int

	previous *container
}
