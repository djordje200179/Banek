package objects

type Undefined struct{}

func (undefined Undefined) Type() Type     { return TypeString }
func (undefined Undefined) Clone() Object  { return undefined }
func (undefined Undefined) String() string { return "<undefined>" }

func (undefined Undefined) Equals(other Object) bool {
	_, ok := other.(Undefined)
	return ok
}
