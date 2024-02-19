package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime"
)

type Stack struct {
	global Scope
	active *Scope
}

func NewStack(entryFunc bytecode.FuncTemplate) *Stack {
	globals := make([]runtime.Obj, entryFunc.NumLocals)

	stack := &Stack{
		global: Scope{
			vars: globals,
			code: entryFunc.Code,
		},
	}
	stack.active = &stack.global

	return stack
}

func (s *Stack) GetGlobal(index int) runtime.Obj        { return s.global.vars[index] }
func (s *Stack) SetGlobal(index int, value runtime.Obj) { s.global.vars[index] = value }

func (s *Stack) GetLocal(index int) runtime.Obj        { return s.active.vars[index] }
func (s *Stack) SetLocal(index int, value runtime.Obj) { s.active.vars[index] = value }

func (s *Stack) GetCaptured(index int) runtime.Obj        { return *s.active.function.Captures[index] }
func (s *Stack) SetCaptured(index int, value runtime.Obj) { *s.active.function.Captures[index] = value }

func (s *Stack) GetCapture(captureInfo bytecode.Capture) *runtime.Obj {
	varScope := s.active
	for range captureInfo.Level {
		varScope = varScope.parent
	}

	return &varScope.vars[captureInfo.Index]
}

func (s *Stack) NewScope(function *bytecode.Func, template *bytecode.FuncTemplate) []runtime.Obj {
	locals := newScopeVars(template.NumLocals)

	funcScope := scopePool.Get().(*Scope)
	*funcScope = Scope{
		code:     template.Code,
		vars:     locals,
		function: function,
		template: template,
		parent:   s.active,
	}
	s.active = funcScope

	return locals
}

func (s *Stack) RestoreScope() {
	if !s.active.template.IsCaptured {
		freeScopeVars(s.active.vars)
	}

	restoredScope := s.active
	s.active = s.active.parent

	*restoredScope = Scope{}
	scopePool.Put(restoredScope)
}

func (s *Stack) ReadOpcode() instrs.Opcode {
	opcode := instrs.Opcode(s.active.code[s.active.pc])
	s.active.pc++

	return opcode
}

func (s *Stack) ReadOperand(op instrs.Opcode, index int) int {
	width := op.Info().Operands[index].Width

	operandSlice := s.active.code[s.active.pc : s.active.pc+width]
	operandValue := instrs.ReadOperandValue(operandSlice)
	s.active.pc += width

	return operandValue
}

func (s *Stack) MovePC(offset int) { s.active.pc += offset }
