package objects

type String string

func (str String) Type() Type     { return TypeString }
func (str String) Clone() Object  { return str }
func (str String) String() string { return string(str) }
