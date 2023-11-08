package objects

type String string

func (str String) Type() Type     { return TypeString }
func (str String) Clone() Object  { return str }
func (str String) String() string { return string(str) }

func (str String) Equals(other Object) bool {
	otherString, ok := other.(String)
	if !ok {
		return false
	}

	return str == otherString
}
