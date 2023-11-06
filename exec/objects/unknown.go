package objects

type Unknown struct{}

func (unknown Unknown) Type() string   { return "Unknown" }
func (unknown Unknown) Clone() Object  { return unknown }
func (unknown Unknown) String() string { return "Unknown" }
