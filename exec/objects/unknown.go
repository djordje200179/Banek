package objects

type Unknown struct{}

func (unknown Unknown) Type() Type     { return TypeUnknown }
func (unknown Unknown) Clone() Object  { return unknown }
func (unknown Unknown) String() string { return "<unknown>" }
