package objects

type Undefined struct{}

func (undefined Undefined) Type() string   { return "undefined" }
func (undefined Undefined) String() string { return "undefined" }
