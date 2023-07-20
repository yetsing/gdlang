package evaluator

import (
	"context"
	"fmt"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"weilang/ast"
	"weilang/object"
)

func Eval(
	ctx context.Context,
	node ast.Node,
	env *object.Environment,
) object.Object {
	switch node := node.(type) {

	// 语句
	case *ast.Program:
		return evalProgram(ctx, node, env)

	case *ast.BlockStatement:
		return evalBlockStatements(ctx, node, env)

	case *ast.VarStatement:
		val := Eval(ctx, node.Value, env)
		if IsError(val) {
			return val
		}
		ret := env.Add(node.Name.Value, val, false)
		if IsError(ret) {
			return ret
		}

	case *ast.ConStatement:
		val := Eval(ctx, node.Value, env)
		if IsError(val) {
			return val
		}
		ret := env.Add(node.Name.Value, val, true)
		if IsError(ret) {
			return ret
		}

	case *ast.AssignStatement:
		return evalAssignStatement(ctx, node, env)

	case *ast.IfStatement:
		return evalIfStatement(ctx, node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(ctx, node, env)

	case *ast.ContinueStatement:
		return object.CONTINUE_VALUE

	case *ast.BreakStatement:
		return object.BREAK_VALUE

	case *ast.ExpressionStatement:
		return Eval(ctx, node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(ctx, node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ForInStatement:
		return evalForInStatement(ctx, node, env)

	case *ast.WeiExportStatement:
		var names []string
		for _, ident := range node.Names {
			name := ident.Value
			if _, ok := env.Get(name); !ok {
				return object.NewError("undefined '%s'", name)
			}
			names = append(names, name)
		}
		return evalExport(ctx, env, node.Names)

	// 表达式
	case *ast.WeiImportExpression:
		filename := Eval(ctx, node.Filename, env)
		return evalImport(ctx, filename)

	case *ast.UnaryExpression:
		operand := Eval(ctx, node.Operand, env)
		if IsError(operand) {
			return operand
		}
		return evalUnaryExpression(ctx, node.Operator, operand)

	case *ast.WeiAttributeExpression:
		left, ok := env.Get("wei")
		if !ok {
			return object.Unreachable("undefined 'wei'")
		}
		return evalAttributeExpression(ctx, left, node.Attribute.Value)

	case *ast.BinaryOpExpression:
		left := Eval(ctx, node.Left, env)
		if IsError(left) {
			return left
		}
		right := Eval(ctx, node.Right, env)
		if IsError(right) {
			return right
		}
		return evalBinaryOpExpression(ctx, node.Operator, left, right)

	case *ast.CallExpression:
		function := Eval(ctx, node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(ctx, node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return evalFunction(ctx, function, args)

	case *ast.SubscriptionExpression:
		left := Eval(ctx, node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(ctx, node.Index, env)
		if IsError(index) {
			return index
		}
		return evalSubscriptionExpression(ctx, left, index)

	case *ast.AttributeExpression:
		left := Eval(ctx, node.Left, env)
		if IsError(left) {
			return left
		}
		return evalAttributeExpression(ctx, left, node.Attribute.Value)

	case *ast.DictLiteral:
		return evalDictLiteral(ctx, node, env)

	case *ast.ListLiteral:
		elements := evalExpressions(ctx, node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return object.NewList(elements)

	case *ast.FunctionLiteral:
		return object.NewFunction(node, env)

	case *ast.Identifier:
		return evalIdentifier(ctx, node, env)

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

func evalProgram(
	ctx context.Context,
	program *ast.Program,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(ctx, statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatements(
	ctx context.Context,
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	blockEnv := object.NewEnclosedEnvironment(env)
	for _, statement := range block.Statements {
		result = Eval(ctx, statement, blockEnv)

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
	ctx context.Context,
	assign *ast.AssignStatement, env *object.Environment,
) object.Object {
	val := Eval(ctx, assign.Value, env)
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
		left := Eval(ctx, obj.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(ctx, obj.Index, env)
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
	case *ast.AttributeExpression:
		left := Eval(ctx, obj.Left, env)
		if IsError(left) {
			return left
		}
		name := obj.Attribute.Value
		if attr, ok := left.(object.Attributable); ok {
			return attr.SetAttribute(name, val)
		}
		return object.NewError("'%s' object can not set attribute", left.Type())
	default:
		return object.Unreachable("assign")
	}

}

func evalIfStatement(
	ctx context.Context,
	is *ast.IfStatement,
	env *object.Environment,
) object.Object {
	for _, branch := range is.IfBranches {
		condition := Eval(ctx, branch.Condition, env)
		if IsError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(ctx, branch.Body, env)
		}
	}

	if is.ElseBody != nil {
		return Eval(ctx, is.ElseBody, env)
	} else {
		return object.NULL
	}
}

func evalWhileStatement(
	ctx context.Context,
	ws *ast.WhileStatement,
	env *object.Environment,
) object.Object {
	for {
		condition := Eval(ctx, ws.Condition, env)
		if isTruthy(condition) {
			encolosedEnv := object.NewEnclosedEnvironment(env)
			val := Eval(ctx, ws.Body, encolosedEnv)
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

func evalForInStatement(
	ctx context.Context,
	forInStmt *ast.ForInStatement,
	env *object.Environment,
) object.Object {
	obj := Eval(ctx, forInStmt.Expr, env)
	iterable, ok := obj.(object.Iterable)
	if !ok {
		return object.NewError("'%s' object is not iterable", obj.Type())
	}
	iterator := iterable.Iter()
	for {
		enclosedEnv := object.NewEnclosedEnvironment(env)
		nextVal := iterator.Next()
		if nextVal == object.StopIteration {
			break
		}
		if IsError(nextVal) {
			return nextVal
		}
		var values []object.Object
		switch nextVal := nextVal.(type) {
		case *object.Tuple:
			values = append(values, nextVal.Elements...)
		default:
			values = append(values, nextVal)
		}
		if len(forInStmt.Targets) != len(values) {
			return object.WrongNumberUnpack(len(values), len(forInStmt.Targets))
		}
		for i, target := range forInStmt.Targets {
			enclosedEnv.Add(target.Value, values[i], forInStmt.Con)
		}
		ret := Eval(ctx, forInStmt.Body, enclosedEnv)
		if IsError(ret) {
			return ret
		}
	}
	return nil
}

func evalFunction(
	ctx context.Context,
	fn object.Object,
	args []object.Object,
) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return object.NewError("function expected %d arguments but got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(ctx, fn.Body, extendedEnv)
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

func evalSubscriptionExpression(
	ctx context.Context,
	left, index object.Object,
) object.Object {
	switch {
	case left.TypeIs(object.LIST_OBJ):
		listObj := left.(*object.List)
		return listObj.GetItem(index)
	case left.TypeIs(object.DICT_OBJ):
		dictObj := left.(*object.Dict)
		return dictObj.GetItem(index)
	case left.TypeIs(object.STRING_OBJ):
		return evalStringSubscriptionExpression(ctx, left, index)
	default:
		return object.NewError("'%s' object is not subscriptable", left.Type())
	}
}

//goland:noinspection GoUnusedParameter
func evalStringSubscriptionExpression(
	ctx context.Context,
	s, index object.Object,
) object.Object {
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

func evalExpressions(
	ctx context.Context,
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(ctx, exp, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalDictLiteral(
	ctx context.Context,
	node *ast.DictLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(ctx, keyNode, env)
		if IsError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return object.NewError("unhashable type: '%s'", key.Type())
		}

		value := Eval(ctx, valueNode, env)
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

func evalUnaryExpression(
	ctx context.Context,
	operator string,
	right object.Object,
) object.Object {
	switch operator {
	case "not":
		return evalNotOperatorExpression(ctx, right)
	case "-":
		return evalMinusUnaryOperatorExpression(ctx, right)
	case "+":
		return evalPlusUnaryOperatorExpression(ctx, right)
	case "~":
		return evalBitwiseNotOperatorExpression(ctx, right)
	default:
		return object.NewError("unsupported operand type for %s: '%s'", operator, right.Type())
	}
}

//goland:noinspection GoUnusedParameter
func evalNotOperatorExpression(
	ctx context.Context,
	operand object.Object,
) object.Object {
	return object.NativeBoolToBooleanObject(!isTruthy(operand))
}

//goland:noinspection GoUnusedParameter
func evalMinusUnaryOperatorExpression(
	ctx context.Context,
	right object.Object,
) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("unsupported operand type for -: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(-value)
}

//goland:noinspection GoUnusedParameter
func evalPlusUnaryOperatorExpression(
	ctx context.Context,
	right object.Object,
) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("unsupported operand type for +: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(value)
}

//goland:noinspection GoUnusedParameter
func evalBitwiseNotOperatorExpression(
	ctx context.Context,
	right object.Object,
) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("bad operand type for unary +: '%s'", right.Type())
		return &object.Error{Message: message}
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(^value)
}

func evalBinaryOpExpression(
	ctx context.Context,
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.TypeIs(object.INTEGER_OBJ) && right.TypeIs(object.INTEGER_OBJ):
		return evalIntegerBinaryOpExpression(ctx, operator, left, right)
	case left.TypeIs(object.STRING_OBJ) && right.TypeIs(object.STRING_OBJ):
		return evalStringBinaryOpExpression(ctx, operator, left, right)

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

//goland:noinspection GoUnusedParameter
func evalIntegerBinaryOpExpression(
	ctx context.Context,
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

//goland:noinspection GoUnusedParameter
func evalStringBinaryOpExpression(
	ctx context.Context,
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

//goland:noinspection GoUnusedParameter
func evalIdentifier(
	ctx context.Context,
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return object.NewError("undefined: '%s'", node.Value)
}

//goland:noinspection GoUnusedParameter
func evalAttributeExpression(
	ctx context.Context,
	left object.Object,
	name string,
) object.Object {
	leftAttr, ok := left.(object.Attributable)
	if !ok {
		return object.NewError("'%s' object has not attribute '%s'", left.Type(), name)
	}

	return leftAttr.GetAttribute(name)
}
