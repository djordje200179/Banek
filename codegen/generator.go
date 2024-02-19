package codegen

import (
	"banek/ast"
	"banek/bytecode"
	"banek/bytecode/instrs"
	"slices"
)

func Generate(stmtChan <-chan ast.Stmt) bytecode.Executable {
	g := &generator{
		funcPool: make([]bytecode.FuncTemplate, 1),
	}
	g.container = &g.global

	for stmt := range stmtChan {
		g.compileStmt(stmt)
	}

	g.emitInstr(instrs.OpHalt)

	g.funcPool[0] = bytecode.FuncTemplate{
		Name:       "<entry>",
		NumLocals:  g.global.vars,
		IsCaptured: true,
	}

	return g.makeExecutable()
}

type generator struct {
	global container
	*container

	code instrs.Code

	stringPool []string
	funcPool   []bytecode.FuncTemplate
}

func (g *generator) makeExecutable() bytecode.Executable {
	return bytecode.Executable{
		Code: g.code,

		StringPool: g.stringPool,
		FuncPool:   g.funcPool,
	}
}

func (g *generator) emitInstr(opcode instrs.Opcode, operands ...int) {
	instr := instrs.MakeInstr(opcode, operands...)
	g.code = slices.Concat(g.code, instr)
}

func (g *generator) patchJumpOperand(addr int, operandIndex int) {
	op := instrs.Opcode(g.code[addr])
	opInfo := op.Info()

	instCode := g.code[addr : addr+opInfo.Size()]

	operandWidth := opInfo.Operands[operandIndex].Width
	operandOffset := opInfo.OperandOffset(operandIndex)

	offset := g.currAddr() - addr - opInfo.Size()
	copy(instCode[operandOffset:], instrs.MakeOperandValue(offset, operandWidth))
}

func (g *generator) currAddr() int {
	return len(g.code)
}
