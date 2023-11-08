package objects

type Undefined struct{}

func (undefined Undefined) Type() Type     { return TypeString }
func (undefined Undefined) Clone() Object  { return undefined }
func (undefined Undefined) String() string { return "<undefined>" }
