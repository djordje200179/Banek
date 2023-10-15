package evaluator

import (
	"banek/ast"
	"banek/ast/expressions"
	"strconv"
	"strings"
)

type Boolean bool

func (boolean Boolean) Type() string   { return "boolean" }
func (boolean Boolean) String() string { return strconv.FormatBool(bool(boolean)) }

type Integer int

func (integer Integer) Type() string   { return "integer" }
func (integer Integer) String() string { return strconv.Itoa(int(integer)) }

type Null struct{}

func (null Null) Type() string   { return "null" }
func (null Null) String() string { return "null" }

type Function struct {
	Parameters []expressions.Identifier
	Body       ast.Statement
	Env        *environment
}

func (function Function) Type() string { return "function" }

func (function Function) String() string {
	var sb strings.Builder

	sb.WriteString("fn(")
	for i, param := range function.Parameters {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(param.String())
	}
	sb.WriteString(") {\n")
	sb.WriteString(function.Body.String())
	sb.WriteString("\n}")

	return sb.String()
}

type String string

func (str String) Type() string   { return "string" }
func (str String) String() string { return string(str) }

type Array []Object

func (array Array) Type() string { return "array" }
func (array Array) String() string {
	var sb strings.Builder

	elements := make([]string, len(array))
	for i, element := range array {
		elements[i] = element.String()
	}

	sb.WriteByte('[')
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteByte(']')

	return sb.String()
}
