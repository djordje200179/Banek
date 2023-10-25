package objects

type unknown struct{}

func (unknown unknown) Type() string   { return "unknown" }
func (unknown unknown) String() string { return "unknown" }

var Unknown Object = unknown{}
