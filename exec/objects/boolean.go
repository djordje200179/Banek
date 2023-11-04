package objects

import "strconv"

type Boolean bool

func (boolean Boolean) Type() string   { return "boolean" }
func (boolean Boolean) Clone() Object  { return boolean }
func (boolean Boolean) String() string { return strconv.FormatBool(bool(boolean)) }
