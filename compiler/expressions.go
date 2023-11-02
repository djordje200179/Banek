package compiler

import (
	"banek/ast"
	"banek/ast/expressions"
	"banek/bytecode"
	"banek/bytecode/instruction"
	"banek/compiler/scopes"
	"banek/exec/errors"
	"banek/exec/objects"
)

func (compiler *compiler) compileExpression(expression ast.Expression) error {
	scope := compiler.topScope()

	switch expression := expression.(type) {
	case expressions.IntegerLiteral:
		integer := objects.Integer(expression)
		scope.EmitInstr(instruction.PushConst, compiler.addConstant(integer))
		return nil
	case expressions.BooleanLiteral:
		boolean := objects.Boolean(expression)
		scope.EmitInstr(instruction.PushConst, compiler.addConstant(boolean))
		return nil
	case expressions.StringLiteral:
		str := objects.String(expression)
		scope.EmitInstr(instruction.PushConst, compiler.addConstant(str))
		return nil
	case expressions.InfixOperation:
		return compiler.compileInfixOperation(expression)
	case expressions.PrefixOperation:
		return compiler.compilePrefixOperation(expression)
	case expressions.If:
		return compiler.compileIfExpression(expression)
	case expressions.ArrayLiteral:
		for _, element := range expression {
			err := compiler.compileExpression(element)
			if err != nil {
				return err
			}
		}

		scope.EmitInstr(instruction.NewArray, len(expression))

		return nil
	case expressions.CollectionAccess:
		err := compiler.compileExpression(expression.Collection)
		if err != nil {
			return err
		}

		err = compiler.compileExpression(expression.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.PushCollectionElement)

		return nil
	case expressions.Assignment:
		return compiler.compileAssigmentExpression(expression)
	case expressions.FunctionCall:
		for _, argument := range expression.Arguments {
			err := compiler.compileExpression(argument)
			if err != nil {
				return err
			}
		}

		err := compiler.compileExpression(expression.Function)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.Call, len(expression.Arguments))

		return nil
	case expressions.FunctionLiteral:
		return compiler.compileFunctionLiteral(expression)
	case expressions.Identifier:
		return compiler.compileIdentifier(expression)
	default:
		return ast.ErrUnknownExpression{Expression: expression}
	}
}

func (compiler *compiler) compileIfExpression(expression expressions.If) error {
	err := compiler.compileExpression(expression.Condition)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	firstPatchAddress := scope.CurrAddr()
	scope.EmitInstr(instruction.BranchIfFalse, 0)

	err = compiler.compileExpression(expression.Consequence)
	if err != nil {
		return err
	}

	elseAddress := scope.CurrAddr()

	branchSize := instruction.Branch.Info().Size()

	secondPatchAddress := elseAddress
	scope.EmitInstr(instruction.Branch, 0)
	elseAddress += branchSize

	err = compiler.compileExpression(expression.Alternative)
	if err != nil {
		return err
	}

	outAddress := scope.CurrAddr()
	scope.PatchInstrOperand(secondPatchAddress, 0, outAddress-secondPatchAddress-branchSize)

	scope.PatchInstrOperand(firstPatchAddress, 0, elseAddress-firstPatchAddress-instruction.BranchIfFalse.Info().Size())

	return nil
}

func (compiler *compiler) compileAssigmentExpression(expression expressions.Assignment) error {
	err := compiler.compileExpression(expression.Value)
	if err != nil {
		return err
	}

	scope := compiler.topScope()

	scope.EmitInstr(instruction.PushDuplicate)

	switch variable := expression.Variable.(type) {
	case expressions.Identifier:
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
			scope.EmitInstr(instruction.PopGlobal, varIndex)
			return nil
		} else if varScope == scope {
			scope.EmitInstr(instruction.PopLocal, varIndex)
			return nil
		}

		capturedVarLevel := len(compiler.scopes) - 2 - varScopeIndex

		functionScope, ok := varScope.(*scopes.Function)
		if !ok {
			block := varScope.(*scopes.Block)
			for {
				nextScope := block.Parent

				functionScope, ok = nextScope.(*scopes.Function)
				if ok {
					break
				}

				block = nextScope.(*scopes.Block)
			}
		}

		capturedVariableIndex := functionScope.AddCapturedVar(capturedVarLevel, varIndex)
		scope.EmitInstr(instruction.PopCaptured, capturedVariableIndex)

		return nil
	case expressions.CollectionAccess:
		err := compiler.compileExpression(variable.Collection)
		if err != nil {
			return err
		}

		err = compiler.compileExpression(variable.Key)
		if err != nil {
			return err
		}

		scope.EmitInstr(instruction.PopCollectionElement)

		return nil
	default:
		return errors.ErrInvalidOperand{Operation: "=", LeftOperand: objects.Unknown} // TODO: fix right operand
	}
}

func (compiler *compiler) compileFunctionLiteral(expression expressions.FunctionLiteral) error {
	funcScope := new(scopes.Function)

	paramNames := make([]string, len(expression.Parameters))
	for i, param := range expression.Parameters {
		paramNames[i] = param.String()
	}

	err := funcScope.AddParams(paramNames)
	if err != nil {
		return err
	}

	compiler.pushScope(funcScope)
	err = compiler.compileExpression(expression.Body)
	if err != nil {
		return err
	}
	funcScope.EmitInstr(instruction.Return)
	compiler.popScope()

	functionTemplate := funcScope.MakeFunction()

	functionIndex := compiler.addFunction(functionTemplate)

	scope := compiler.topScope()

	if functionTemplate.IsClosure() {
		scope.EmitInstr(instruction.NewFunction, functionIndex)
	} else {
		functionObject := &bytecode.Function{
			TemplateIndex: functionIndex,
		}

		scope.EmitInstr(instruction.PushConst, compiler.addConstant(functionObject))
	}

	return nil
}

func (compiler *compiler) compileIdentifier(expression expressions.Identifier) error {
	varName := expression.String()

	scope := compiler.topScope()

	if index := objects.BuiltinFindIndex(varName); index != -1 {
		scope.EmitInstr(instruction.PushBuiltin, index)
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
		scope.EmitInstr(instruction.PushGlobal, varIndex)
		return nil
	} else if varScope == scope {
		scope.EmitInstr(instruction.PushLocal, varIndex)
		return nil
	}

	capturedVariableLevel := len(compiler.scopes) - 2 - varScopeIndex

	functionScope, ok := varScope.(*scopes.Function)
	if !ok {
		block := varScope.(*scopes.Block)
		for {
			nextScope := block.Parent

			functionScope, ok = nextScope.(*scopes.Function)
			if ok {
				break
			}

			block = nextScope.(*scopes.Block)
		}
	}

	capturedVariableIndex := functionScope.AddCapturedVar(capturedVariableLevel, varIndex)
	scope.EmitInstr(instruction.PushCaptured, capturedVariableIndex)

	return nil
}
