package objects

type String string

func (str String) Type() string   { return "string" }
func (str String) String() string { return string(str) }