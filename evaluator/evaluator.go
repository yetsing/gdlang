package evaluator

import (
	"fmt"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"weilang/ast"
	"weilang/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// 语句
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatements(node, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		ret := env.Add(node.Name.Value, val, false)
		if IsError(ret) {
			return ret
		}

	case *ast.ConStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		ret := env.Add(node.Name.Value, val, true)
		if IsError(ret) {
			return ret
		}

	case *ast.AssignStatement:
		return evalAssignStatement(node, env)

	case *ast.IfStatement:
		return evalIfStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ContinueStatement:
		return object.CONTINUE_VALUE

	case *ast.BreakStatement:
		return object.BREAK_VALUE

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// 表达式
	case *ast.UnaryExpression:
		operand := Eval(node.Operand, env)
		if IsError(operand) {
			return operand
		}
		return evalUnaryExpression(node.Operator, operand)

	case *ast.BinaryOpExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalBinaryOpExpression(node.Operator, left, right)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.SubscriptionExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}
		return evalSubscriptionExpression(left, index)

	case *ast.AttributeExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		return evalAttributeExpression(left, node.Attribute.Value)

	case *ast.DictLiteral:
		return evalDictLiteral(node, env)

	case *ast.ListLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return object.NewList(elements)

	case *ast.FunctionLiteral:
		return object.NewFunction(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.StringLiteral:
		return object.NewString(node.Value)

	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)

	case *ast.Boolean:
		return object.NativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return object.NULL
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	blockEnv := object.NewEnclosedEnvironment(env)
	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		in := object.TypeIn(
			result,
			object.RETURN_VALUE_OBJ,
			object.ERROR_OBJ,
			object.CONTINUE_VALUE_OBJ,
			object.BREAK_VALUE_OBJ,
		)
		if in {
			return result
		}
	}

	return result
}

func evalAssignStatement(
	assign *ast.AssignStatement, env *object.Environment,
) object.Object {
	val := Eval(assign.Value, env)
	if IsError(val) {
		return val
	}
	switch obj := assign.Left.(type) {
	case *ast.Identifier:
		ret := env.Set(obj.Value, val)
		if IsError(ret) {
			return ret
		}
		return nil
	case *ast.SubscriptionExpression:
		left := Eval(obj.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(obj.Index, env)
		switch {
		case left.TypeIs(object.LIST_OBJ):
			listObj := left.(*object.List)
			return listObj.SetItem(index, val)
		case left.TypeIs(object.DICT_OBJ):
			dictObj := left.(*object.Dict)
			return dictObj.SetItem(index, val)
		default:
			return object.NewError("'%s' object is not subscriptable", left.Type())
		}
	default:
		return object.Unreachable("assign")
	}

}

func evalIfStatement(is *ast.IfStatement, env *object.Environment) object.Object {
	for _, branch := range is.IfBranches {
		condition := Eval(branch.Condition, env)
		if IsError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(branch.Body, env)
		}
	}

	if is.ElseBody != nil {
		return Eval(is.ElseBody, env)
	} else {
		return object.NULL
	}
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	for {
		condition := Eval(ws.Condition, env)
		if isTruthy(condition) {
			encolosedEnv := object.NewEnclosedEnvironment(env)
			val := Eval(ws.Body, encolosedEnv)
			switch val.(type) {
			case *object.ReturnValue, *object.Error:
				return val
			case *object.BreakValue:
				return object.NULL
			}
		} else {
			break
		}
	}
	return object.NULL
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return object.NewError("function expected %d arguments but got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		if IsError(evaluated) {
			return evaluated
		}
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.BoundBuiltinMethod:
		return fn.Fn(fn.This, args...)
	default:
		return object.NewError("not a function: '%s'", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, parameter := range fn.Parameters {
		env.Pass(parameter.Value, args[paramIdx])
	}

	return env
}

func evalSubscriptionExpression(left, index object.Object) object.Object {
	switch {
	case left.TypeIs(object.LIST_OBJ):
		listObj := left.(*object.List)
		return listObj.GetItem(index)
	case left.TypeIs(object.DICT_OBJ):
		dictObj := left.(*object.Dict)
		return dictObj.GetItem(index)
	case left.TypeIs(object.STRING_OBJ):
		return evalStringSubscriptionExpression(left, index)
	default:
		return object.NewError("'%s' object is not subscriptable", left.Type())
	}
}

func evalStringSubscriptionExpression(s, index object.Object) object.Object {
	if index.TypeNotIs(object.INTEGER_OBJ) {
		return object.NewError("string index must be integer")
	}

	idx := int(index.(*object.Integer).Value)
	strObj := s.(*object.String)
	length := strObj.Length
	if idx < 0 {
		idx += length
	}

	if idx < 0 || idx > length-1 {
		return object.NewError("string index out of range")
	}

	si := exutf8.RuneSubString(strObj.Value, idx, 1)
	return object.NewString(si)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalDictLiteral(node *ast.DictLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if IsError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError("unhashable type: '%s'", key.Type())
		}

		value := Eval(valueNode, env)
		if IsError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}
	return object.NewDict(pairs)
}

func evalUnaryExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "not":
		return evalNotOperatorExpression(right)
	case "-":
		return evalMinusUnaryOperatorExpression(right)
	case "+":
		return evalPlusUnaryOperatorExpression(right)
	case "~":
		return evalBitwiseNotOperatorExpression(right)
	default:
		return object.NewError("unsupported operand type for %s: '%s'", operator, right.Type())
	}
}

func evalNotOperatorExpression(operand object.Object) object.Object {
	return object.NativeBoolToBooleanObject(!isTruthy(operand))
}

func evalMinusUnaryOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("unsupported operand type for -: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(-value)
}

func evalPlusUnaryOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("unsupported operand type for +: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(value)
}

func evalBitwiseNotOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("bad operand type for unary +: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(^value)
}

func evalBinaryOpExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.TypeIs(object.INTEGER_OBJ) && right.TypeIs(object.INTEGER_OBJ):
		return evalIntegerBinaryOpExpression(operator, left, right)
	case left.TypeIs(object.STRING_OBJ) && right.TypeIs(object.STRING_OBJ):
		return evalStringBinaryOpExpression(operator, left, right)

	case operator == "==":
		return object.NativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return object.NativeBoolToBooleanObject(left != right)
	case operator == "and":
		return object.NativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "or":
		return object.NativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	default:
		msg := fmt.Sprintf("unsupported operand type for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
		return &object.Error{Message: msg}
	}
}

func evalIntegerBinaryOpExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return object.NewInteger(leftVal + rightVal)
	case "-":
		return object.NewInteger(leftVal - rightVal)
	case "*":
		return object.NewInteger(leftVal * rightVal)
	case "/":
		return object.NewInteger(leftVal / rightVal)
	case "%":
		return object.NewInteger(leftVal % rightVal)
	case ">>":
		return object.NewInteger(leftVal >> rightVal)
	case "<<":
		return object.NewInteger(leftVal << rightVal)
	case "&":
		return object.NewInteger(leftVal & rightVal)
	case "^":
		return object.NewInteger(leftVal ^ rightVal)
	case "|":
		return object.NewInteger(leftVal | rightVal)
	case "<":
		return object.NativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return object.NativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return object.NativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return object.NativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return object.NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return object.NativeBoolToBooleanObject(leftVal != rightVal)
	case "and":
		return object.NativeBoolToBooleanObject(leftVal != 0 && rightVal != 0)
	case "or":
		return object.NativeBoolToBooleanObject(leftVal != 0 || rightVal != 0)
	default:
		msg := fmt.Sprintf("unsupported operand type for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
		return &object.Error{Message: msg}
	}
}

func evalStringBinaryOpExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return object.NewString(leftVal + rightVal)

	case "==":
		return object.NativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return object.NativeBoolToBooleanObject(leftVal != rightVal)
	case "and":
		return object.NativeBoolToBooleanObject(len(leftVal) > 0 && len(rightVal) > 0)
	case "or":
		return object.NativeBoolToBooleanObject(len(leftVal) > 0 || len(rightVal) > 0)

	default:
		return object.NewError("unsupported operand type for %s: 'str' and 'str'", operator)
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("undefined: '%s'", node.Value)
}

func evalAttributeExpression(left object.Object, name string) object.Object {
	leftAttr, ok := left.(object.Attributable)
	if !ok {
		return object.NewError("'%s' object has not attribute '%s'", left.Type(), name)
	}

	return leftAttr.GetAttribute(name)
}
