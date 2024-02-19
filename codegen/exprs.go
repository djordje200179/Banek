package codegen

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/bytecode/instrs"
	"banek/symtable/symbols"
	"banek/tokens"
	"slices"
)

func (g *generator) compileExpr(expr ast.Expr) {
	switch expr := expr.(type) {
	case exprs.ArrayLiteral:
		g.compileArrayLiteral(expr)
	case exprs.BinaryOp:
		g.compileBinaryOp(expr)
	case exprs.BoolLiteral:
		g.compileBoolLiteral(expr)
	case exprs.CollIndex:
		g.compileCollIndex(expr)
	case exprs.FuncCall:
		g.compileFuncCall(expr)
	case exprs.FuncLiteral:
		g.compileFuncLiteral(expr)
	case exprs.Ident:
		g.compileIdent(expr)
	case exprs.If:
		g.compileIfExpr(expr)
	case exprs.IntLiteral:
		g.compileIntLiteral(expr)
	case exprs.StringLiteral:
		g.compileStringLiteral(expr)
	case exprs.UnaryOp:
		g.compileUnaryOp(expr)
	case exprs.UndefinedLiteral:
		g.compileUndefinedLiteral(expr)
	default:
		panic("unreachable")
	}
}

func (g *generator) compileArrayLiteral(expr exprs.ArrayLiteral) {
	for _, elem := range expr {
		g.compileExpr(elem)
	}

	g.emitInstr(instrs.OpMakeArray, len(expr))
}

func (g *generator) compileBinaryOp(expr exprs.BinaryOp) {
	g.compileExpr(expr.Left)
	g.compileExpr(expr.Right)

	switch expr.Operator {
	case tokens.Plus:
		g.emitInstr(instrs.OpBinaryAdd)
	case tokens.Minus:
		g.emitInstr(instrs.OpBinarySub)
	case tokens.Asterisk:
		g.emitInstr(instrs.OpBinaryMul)
	case tokens.Slash:
		g.emitInstr(instrs.OpBinaryDiv)
	case tokens.Percent:
		g.emitInstr(instrs.OpBinaryMod)
	case tokens.Equals:
		g.emitInstr(instrs.OpBinaryEq)
	case tokens.NotEquals:
		g.emitInstr(instrs.OpBinaryNe)
	case tokens.Greater:
		g.emitInstr(instrs.OpBinaryGt)
	case tokens.GreaterEquals:
		g.emitInstr(instrs.OpBinaryGe)
	case tokens.Less:
		g.emitInstr(instrs.OpBinaryLt)
	case tokens.LessEquals:
		g.emitInstr(instrs.OpBinaryLe)
	default:
		panic("unreachable")
	}

}

func (g *generator) compileCollIndex(expr exprs.CollIndex) {

}

func (g *generator) compileFuncCall(expr exprs.FuncCall) {
	for _, arg := range expr.Args {
		g.compileExpr(arg)
	}

	g.compileExpr(expr.Func)
	g.emitInstr(instrs.OpCall, len(expr.Args))
}

func (g *generator) compileIdent(expr exprs.Ident) {
	switch sym := expr.Symbol.(type) {
	case symbols.Builtin:
		g.emitInstr(instrs.OpPushBuiltin, int(sym))
	case symbols.Var:
		switch {
		case sym.Level == 0:
			g.emitInstr(instrs.OpPushGlobal, sym.Index)
		case sym.Level == g.level:
			switch sym.Index {
			case 0:
				g.emitInstr(instrs.OpPushLocal0)
			case 1:
				g.emitInstr(instrs.OpPushLocal1)
			case 2:
				g.emitInstr(instrs.OpPushLocal2)
			default:
				g.emitInstr(instrs.OpPushLocal, sym.Index)
			}
		default:
			g.emitInstr(instrs.OpPushCaptured, g.level-sym.Level, sym.Index)
		}
	default:
		panic("unreachable")
	}
}

func (g *generator) compileIfExpr(expr exprs.If) {
	g.compileExpr(expr.Cond)
	jumpPc := g.container.currAddr()
	g.emitInstr(instrs.OpBranchFalse, 0)

	g.compileExpr(expr.Cons)

	altJumpPc := g.container.currAddr()
	g.emitInstr(instrs.OpJump, 0)
	g.container.patchJumpOperand(jumpPc, 0)
	g.compileExpr(expr.Alt)
	g.container.patchJumpOperand(altJumpPc, 0)
}

func (g *generator) compileUnaryOp(expr exprs.UnaryOp) {
	g.compileExpr(expr.Operand)

	switch expr.Operator {
	case tokens.Bang:
		g.emitInstr(instrs.OpUnaryNot)
	case tokens.Minus:
		g.emitInstr(instrs.OpUnaryNeg)
	default:
		panic("unreachable")
	}
}

func (g *generator) compileBoolLiteral(expr exprs.BoolLiteral) {
	if expr {
		g.emitInstr(instrs.OpPushTrue)
	} else {
		g.emitInstr(instrs.OpPushFalse)
	}
}

func (g *generator) compileIntLiteral(expr exprs.IntLiteral) {
	switch int(expr) {
	case 0:
		g.emitInstr(instrs.OpPush0)
	case 1:
		g.emitInstr(instrs.OpPush1)
	case 2:
		g.emitInstr(instrs.OpPush2)
	case 3:
		g.emitInstr(instrs.OpPush3)
	case -1:
		g.emitInstr(instrs.OpPushN1)
	default:
		g.emitInstr(instrs.OpPushInt, int(expr))
	}
}

func (g *generator) compileStringLiteral(expr exprs.StringLiteral) {
	i := slices.Index(g.stringPool, string(expr))
	if i == -1 {
		i = len(g.stringPool)
		g.stringPool = append(g.stringPool, string(expr))
	}

	g.emitInstr(instrs.OpPushStr, i)
}

func (g *generator) compileUndefinedLiteral(expr exprs.UndefinedLiteral) {
	g.emitInstr(instrs.OpPushUndef)
}

func (g *generator) compileFuncLiteral(expr exprs.FuncLiteral) {
}

func (g *generator) compileStore(expr ast.Expr) {
	switch expr := expr.(type) {
	case exprs.Ident:
		var sym symbols.Var
		var ok bool
		if sym, ok = expr.Symbol.(symbols.Var); !ok {
			panic("unreachable")
		}

		switch {
		case sym.Level == 0:
			g.emitInstr(instrs.OpPopGlobal, sym.Index)
		case sym.Level == g.level:
			switch sym.Index {
			case 0:
				g.emitInstr(instrs.OpPopLocal0)
			case 1:
				g.emitInstr(instrs.OpPopLocal1)
			case 2:
				g.emitInstr(instrs.OpPopLocal2)
			default:
				g.emitInstr(instrs.OpPopLocal, sym.Index)
			}
		default:
			g.emitInstr(instrs.OpPopCaptured, g.level-sym.Level, sym.Index)
		}
	default:
		panic("unreachable")
	}
}
