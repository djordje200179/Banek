package compiler

import (
	"banek/ast"
	"banek/ast/exprs"
	"banek/bytecode"
	"banek/bytecode/instrs"
	"banek/compiler/scopes"
	"banek/runtime/builtins"
	"banek/runtime/errors"
	"banek/runtime/objs"
)

func (compiler *compiler) compileExpr(expr ast.Expr) error {
	scope := compiler.topScope()

	switch expr := expr.(type) {
	case exprs.ConstLiteral:
		switch expr.Value.Tag {
		case objs.TypeUndefined:
			scope.EmitInstr(instrs.OpPushUndefined)
		case objs.TypeInt:
			integer := expr.Value.AsInt()
			switch integer {
			case 0:
				scope.EmitInstr(instrs.OpPush0)
			case 1:
				scope.EmitInstr(instrs.OpPush1)
			case 2:
				scope.EmitInstr(instrs.OpPush2)
			default:
				scope.EmitInstr(instrs.OpPushConst, compiler.addConst(expr.Value))
			}
		default:
			scope.EmitInstr(instrs.OpPushConst, compiler.addConst(expr.Value))
		}

		return nil
	case exprs.BinaryOp:
		return compiler.compileBinaryOp(expr)
	case exprs.UnaryOp:
		return compiler.compileUnaryOp(expr)
	case exprs.If:
		return compiler.compileIfExpr(expr)
	case exprs.ArrayLiteral:
		for _, elem := range expr {
			err := compiler.compileExpr(elem)
			if err != nil {
				return err
			}
		}

		scope.EmitInstr(instrs.OpNewArray, len(expr))

		return nil
	case exprs.CollIndex:
		err := compiler.compileExpr(expr.Coll)
		if err != nil {
			return err
		}

		err = compiler.compileExpr(expr.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instrs.OpPushCollElem)

		return nil
	case exprs.Assignment:
		return compiler.compileAssigment(expr)
	case exprs.FuncCall:
		for _, arg := range expr.Args {
			err := compiler.compileExpr(arg)
			if err != nil {
				return err
			}
		}

		err := compiler.compileExpr(expr.Func)
		if err != nil {
			return err
		}

		scope.EmitInstr(instrs.OpCall, len(expr.Args))

		return nil
	case exprs.FuncLiteral:
		return compiler.compileFuncLiteral(expr)
	case exprs.Identifier:
		return compiler.compileIdentifier(expr)
	default:
		return ast.ErrUnknownExpr{Expr: expr}
	}
}

func (compiler *compiler) compileIfExpr(expr exprs.If) error {
	err := compiler.compileExpr(expr.Cond)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instrs.OpBranchIfFalse, 0)

	err = compiler.compileExpr(expr.Consequence)
	if err != nil {
		return err
	}

	branchSize := instrs.OpBranch.Info().Size()
	branchIfFalseSize := instrs.OpBranchIfFalse.Info().Size()

	secondPatchAddr := scope.CurrAddr()
	scope.EmitInstr(instrs.OpBranch, 0)
	elseAddr := secondPatchAddr + branchSize

	err = compiler.compileExpr(expr.Alternative)
	if err != nil {
		return err
	}

	outAddr := scope.CurrAddr()

	scope.PatchInstrOperand(secondPatchAddr, 0, outAddr-secondPatchAddr-branchSize)
	scope.PatchInstrOperand(firstPatchAddr, 0, elseAddr-firstPatchAddr-branchIfFalseSize)

	return nil
}

func (compiler *compiler) compileAssigment(expr exprs.Assignment) error {
	err := compiler.compileExpr(expr.Value)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	scope.EmitInstr(instrs.OpPushDup)

	switch variable := expr.Var.(type) {
	case exprs.Identifier:
		varName := variable.String()

		var varScope scopes.Scope
		var varIndex, varScopeIndex int
		for i := len(compiler.scopes) - 1; i >= 0; i-- {
			variable, index := compiler.scopes[i].GetVar(varName)
			if index == -1 {
				continue
			}

			if !variable.Mutable {
				return errors.ErrIdentifierNotMutable{Identifier: varName}
			}

			varScope = compiler.scopes[i]
			varScopeIndex = i
			varIndex = index

			break
		}

		if varScope == nil {
			return errors.ErrIdentifierNotDefined{Identifier: varName}
		}

		if varScope.IsGlobal() {
			scope.EmitInstr(instrs.OpPopGlobal, varIndex)
			return nil
		} else if varScope == scope {
			switch varIndex {
			case 0:
				scope.EmitInstr(instrs.OpPopLocal0)
			case 1:
				scope.EmitInstr(instrs.OpPopLocal1)
			default:
				scope.EmitInstr(instrs.OpPopLocal, varIndex)
			}

			return nil
		}

		varScope.MarkCaptured()

		capturedVarLevel := len(compiler.scopes) - 2 - varScopeIndex

		capturedVarIndex := scope.GetFunc().AddCapturedVar(capturedVarLevel, varIndex)
		scope.EmitInstr(instrs.OpPopCaptured, capturedVarIndex)

		return nil
	case exprs.CollIndex:
		err := compiler.compileExpr(variable.Coll)
		if err != nil {
			return err
		}

		err = compiler.compileExpr(variable.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instrs.OpPopCollElem)

		return nil
	default:
		return ast.ErrInvalidAssignment{Variable: expr.Var}
	}
}

func (compiler *compiler) compileFuncLiteral(expr exprs.FuncLiteral) error {
	funcScope := new(scopes.Func)

	paramNames := make([]string, len(expr.Params))
	for i, param := range expr.Params {
		paramNames[i] = param.String()
	}

	err := funcScope.AddParams(paramNames)
	if err != nil {
		return err
	}

	compiler.pushScope(funcScope)
	err = compiler.compileExpr(expr.Body)
	if err != nil {
		return err
	}
	funcScope.EmitInstr(instrs.OpReturn)
	compiler.popScope()

	funcTemplate := funcScope.MakeFunction()

	funcIndex := compiler.addFunc(funcTemplate)

	scope := compiler.topScope()

	if funcTemplate.IsClosure() {
		scope.EmitInstr(instrs.OpNewFunc, funcIndex)
	} else {
		funcObj := &bytecode.Func{
			TemplateIndex: funcIndex,
		}

		scope.EmitInstr(instrs.OpPushConst, compiler.addConst(funcObj.MakeObj()))
	}

	return nil
}

func (compiler *compiler) compileIdentifier(expr exprs.Identifier) error {
	varName := expr.String()

	scope := compiler.topScope()

	if index := builtins.Find(varName); index != -1 {
		scope.EmitInstr(instrs.OpPushBuiltin, index)
		return nil
	}

	var varScope scopes.Scope
	var varIndex, varScopeIndex int
	for i := len(compiler.scopes) - 1; i >= 0; i-- {
		_, index := compiler.scopes[i].GetVar(varName)
		if index == -1 {
			continue
		}

		varScope = compiler.scopes[i]
		varScopeIndex = i
		varIndex = index

		break
	}

	if varScope == nil {
		return errors.ErrIdentifierNotDefined{Identifier: varName}
	}

	if varScope.IsGlobal() {
		scope.EmitInstr(instrs.OpPushGlobal, varIndex)
		return nil
	} else if varScope == scope {
		switch varIndex {
		case 0:
			scope.EmitInstr(instrs.OpPushLocal0)
		case 1:
			scope.EmitInstr(instrs.OpPushLocal1)
		default:
			scope.EmitInstr(instrs.OpPushLocal, varIndex)
		}

		return nil
	}

	varScope.MarkCaptured()

	capturedVarLevel := len(compiler.scopes) - 2 - varScopeIndex

	capturedVarIndex := scope.GetFunc().AddCapturedVar(capturedVarLevel, varIndex)
	scope.EmitInstr(instrs.OpPushCaptured, capturedVarIndex)

	return nil
}

func (compiler *compiler) compileBinaryOp(expr exprs.BinaryOp) error {
	err := compiler.compileExpr(expr.Left)
	if err != nil {
		return err
	}

	err = compiler.compileExpr(expr.Right)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instrs.OpBinaryOp, int(expr.Operator))

	return nil
}

func (compiler *compiler) compileUnaryOp(expr exprs.UnaryOp) error {
	err := compiler.compileExpr(expr.Operand)
	if err != nil {
		return err
	}

	container := compiler.topScope()
	container.EmitInstr(instrs.OpUnaryOp, int(expr.Operator))

	return nil
}
