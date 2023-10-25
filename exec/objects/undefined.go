package objects

type undefined struct{}

func (undefined undefined) Type() string   { return "undefined" }
func (undefined undefined) String() string { return "undefined" }

var Undefined Object = undefined{}
