package objects

type Undefined struct{}

func (undefined Undefined) Type() string   { return "Undefined" }
func (undefined Undefined) Clone() Object  { return undefined }
func (undefined Undefined) String() string { return "Undefined" }
