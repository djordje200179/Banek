package scopes

import (
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/runtime"
)

type Stack struct {
	global Scope
	active Scope

	last *Scope
}

func NewStack(program *bytecode.Executable) *Stack {
	globals := make([]runtime.Obj, program.FuncPool[0].NumLocals)

	stack := &Stack{
		global: Scope{
			vars: globals,
			code: program.FuncPool[0].Code,
		},
	}
	stack.active = stack.global

	return stack
}

func (s *Stack) GetGlobal(index int) runtime.Obj        { return s.global.vars[index] }
func (s *Stack) SetGlobal(index int, value runtime.Obj) { s.global.vars[index] = value }

func (s *Stack) GetLocal(index int) runtime.Obj        { return s.active.vars[index] }
func (s *Stack) SetLocal(index int, value runtime.Obj) { s.active.vars[index] = value }

func (s *Stack) GetCaptured(index int) runtime.Obj        { return *s.active.function.Captures[index] }
func (s *Stack) SetCaptured(index int, value runtime.Obj) { *s.active.function.Captures[index] = value }

func (s *Stack) GetCapture(captureInfo bytecode.Capture) *runtime.Obj {
	varScope := s.last
	for j := 1; j < captureInfo.Level; j++ {
		varScope = varScope.parent
	}

	return &varScope.vars[captureInfo.Index]
}

func (s *Stack) NewScope(function *bytecode.Func, template *bytecode.FuncTemplate) []runtime.Obj {
	funcScope := scopePool.Get().(*Scope)
	*funcScope = s.active
	s.last = funcScope

	locals := newScopeVars(template.NumLocals)

	s.active = Scope{
		code:     template.Code,
		vars:     locals,
		function: function,
		template: template,
		parent:   s.last,
	}

	return locals
}

func (s *Stack) RestoreScope() {
	if !s.active.template.IsCaptured {
		freeScopeVars(s.active.vars)
	}

	restoredScope := s.last
	s.last = restoredScope.parent
	s.active = *restoredScope

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
