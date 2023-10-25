package expressions

import "banek/ast"

type Assignment struct {
	Variable ast.Expression
	Value    ast.Expression
}

func (assignment Assignment) String() string {
	return assignment.Variable.String() + " = " + assignment.Value.String()
}

func (assignment Assignment) IsConstant() bool {
	return assignment.Value.IsConstant()
}
