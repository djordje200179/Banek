package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime/errors"
	"slices"
)

type Func struct {
	params   []string
	vars     []Var
	captures []bytecode.Capture

	isCaptured bool

	code bytecode.Code

	blocksCounter
}

func (scope *Func) AddParams(params []string) error {
	for i, firstParam := range params {
		for j := i + 1; j < len(params); j++ {
			secondParam := params[j]

			if firstParam == secondParam {
				return errors.ErrIdentifierAlreadyDefined{Identifier: firstParam}
			}
		}
	}

	scope.params = slices.Clone(params)

	scope.vars = make([]Var, len(params))
	for i, param := range params {
		scope.vars[i] = Var{Name: param}
	}

	return nil
}

func (scope *Func) AddVar(name string, mutable bool) (int, error) {
	if slices.ContainsFunc(scope.vars, func(v Var) bool {
		return v.Name == name
	}) {
		return 0, errors.ErrIdentifierAlreadyDefined{Identifier: name}
	}

	scope.vars = append(scope.vars, Var{
		Name:    name,
		Mutable: mutable,
	})

	return len(scope.vars) - 1, nil
}

func (scope *Func) GetVar(name string) (Var, int) {
	index := slices.IndexFunc(scope.vars, func(v Var) bool {
		return v.Name == name
	})
	if index == -1 {
		return Var{}, -1
	}

	return scope.vars[index], index
}

func (scope *Func) AddCapturedVar(level, index int) int {
	captureInfo := bytecode.Capture{Index: index, Level: level}

	if captureIndex := slices.Index(scope.captures, captureInfo); captureIndex != -1 {
		return captureIndex
	}

	scope.captures = append(scope.captures, captureInfo)

	return len(scope.captures) - 1
}

func (scope *Func) EmitInstr(opcode instrs.Opcode, operands ...int) {
	instr := instrs.MakeInstr(opcode, operands...)

	newCode := make(bytecode.Code, len(scope.code)+len(instr))
	copy(newCode, scope.code)
	copy(newCode[len(scope.code):], instr)

	scope.code = newCode
}

func (scope *Func) PatchInstrOperand(addr int, operandIndex int, newValue int) {
	op := instrs.Opcode(scope.code[addr])
	opInfo := op.Info()

	instCode := scope.code[addr : addr+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	copy(instCode[operandOffset:], instrs.MakeOperandValue(newValue, operandWidth))

}

func (scope *Func) CurrAddr() int {
	return len(scope.code)
}

func (scope *Func) IsGlobal() bool {
	return false
}

func (scope *Func) MarkCaptured() {
	scope.isCaptured = true
}

func (scope *Func) GetFunc() *Func {
	return scope
}

func (scope *Func) MakeFunction() bytecode.FuncTemplate {
	return bytecode.FuncTemplate{
		Code: scope.code,

		Params:    scope.params,
		NumLocals: len(scope.vars),
		Captures:  scope.captures,

		IsCaptured: scope.isCaptured,
	}
}
