package objs

import "bytes"

type TypeConfig struct {
	Stringer func(Obj) string

	Equaler func(Obj, Obj) bool

	Marshaller   func(Obj) ([]byte, error)
	Unmarshaller func(buf *bytes.Buffer) (Obj, error)
}

var Config = map[Tag]TypeConfig{}
