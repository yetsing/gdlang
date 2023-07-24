package evaluator

import (
	"context"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"weilang/ast"
	"weilang/object"
)

func Eval(
	ctx context.Context,
	state *WeiState,
	node ast.Node,
	env *object.Environment,
) object.Object {
	state.UpdateLocation(node)
	switch node := node.(type) {

	// 语句
	case *ast.Program:
		return evalProgram(ctx, state, node, env)

	case *ast.BlockStatement:
		return evalBlockStatements(ctx, state, node, env)

	case *ast.VarStatement:
		val := Eval(ctx, state, node.Value, env)
		if IsError(val) {
			return val
		}
		// 更新到赋值操作所在的行号
		state.UpdateLocation(node)
		ret := env.Add(node.Name.Value, val, false)
		if IsError(ret) {
			state.HandleError(ret)
			return ret
		}

	case *ast.ConStatement:
		val := Eval(ctx, state, node.Value, env)
		if IsError(val) {
			return val
		}
		// 更新到赋值操作所在的行号
		state.UpdateLocation(node)
		ret := env.Add(node.Name.Value, val, true)
		if IsError(ret) {
			state.HandleError(ret)
			return ret
		}

	case *ast.AssignStatement:
		return evalAssignStatement(ctx, state, node, env)

	case *ast.IfStatement:
		return evalIfStatement(ctx, state, node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(ctx, state, node, env)

	case *ast.ContinueStatement:
		return object.CONTINUE_VALUE

	case *ast.BreakStatement:
		return object.BREAK_VALUE

	case *ast.ExpressionStatement:
		return Eval(ctx, state, node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(ctx, state, node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ForInStatement:
		return evalForInStatement(ctx, state, node, env)

	case *ast.FunctionDefineStatement:
		function := object.NewFunction(node.Function, env)
		ret := env.Add(node.Function.Name, function, true)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

	case *ast.ClassDefineStatement:
		return evalClassDefine(ctx, state, env, node)

	case *ast.WeiExportStatement:
		return evalExport(ctx, state, env, node.Names)

	// 表达式
	case *ast.WeiImportExpression:
		return evalImport(ctx, state, node.Filename)

	case *ast.UnaryExpression:
		operand := Eval(ctx, state, node.Operand, env)
		if IsError(operand) {
			return operand
		}
		state.UpdateLocation(node)
		ret := evalUnaryExpression(ctx, node.Operator, operand)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

	case *ast.WeiAttributeExpression:
		left, ok := env.Get("wei")
		if !ok {
			return state.Unreachable("undefined 'wei'")
		}
		ret := evalAttributeExpression(ctx, left, node.Attribute.Value)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

	case *ast.BinaryOpExpression:
		switch node.Operator {
		case "and", "or":
			return evalLogicalOperation(ctx, state, env, node)
		default:
			left := Eval(ctx, state, node.Left, env)
			if IsError(left) {
				return left
			}
			right := Eval(ctx, state, node.Right, env)
			if IsError(right) {
				return right
			}
			state.UpdateLocation(node)
			return evalBinaryOpExpression(ctx, state, node.Operator, left, right)
		}

	case *ast.CallExpression:
		function := Eval(ctx, state, node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(ctx, state, node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		state.UpdateLocation(node)
		return evalFunction(ctx, state, function, args)

	case *ast.SubscriptionExpression:
		left := Eval(ctx, state, node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(ctx, state, node.Index, env)
		if IsError(index) {
			return index
		}
		state.UpdateLocation(node)
		ret := evalSubscriptionExpression(ctx, left, index)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

	case *ast.AttributeExpression:
		left := Eval(ctx, state, node.Left, env)
		if IsError(left) {
			return left
		}
		state.UpdateLocation(node)
		ret := evalAttributeExpression(ctx, left, node.Attribute.Value)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

	case *ast.DictLiteral:
		return evalDictLiteral(ctx, state, node, env)

	case *ast.ListLiteral:
		elements := evalExpressions(ctx, state, node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return object.NewList(elements)

	case *ast.FunctionLiteral:
		return object.NewFunction(node, env)

	case *ast.Identifier:
		ret := evalIdentifier(ctx, node, env)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret

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
	state *WeiState,
	program *ast.Program,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(ctx, state, statement, env)

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
	state *WeiState,
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	blockEnv := object.NewEnclosedEnvironment(env)
	for _, statement := range block.Statements {
		result = Eval(ctx, state, statement, blockEnv)

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
	state *WeiState,
	assign *ast.AssignStatement,
	env *object.Environment,
) object.Object {
	val := Eval(ctx, state, assign.Value, env)
	if IsError(val) {
		return val
	}
	switch obj := assign.Left.(type) {
	case *ast.Identifier:
		// 设置回赋值所在的行号
		state.UpdateLocation(assign)
		ret := env.Set(obj.Value, val)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret
	case *ast.SubscriptionExpression:
		left := Eval(ctx, state, obj.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(ctx, state, obj.Index, env)
		if IsError(index) {
			return index
		}
		// 设置回赋值所在的行号
		state.UpdateLocation(assign)
		switch left.Type() {
		case object.LIST_OBJ:
			listObj := left.(*object.List)
			return listObj.SetItem(index, val)
		case object.DICT_OBJ:
			dictObj := left.(*object.Dict)
			return dictObj.SetItem(index, val)
		default:
			return state.NewError("'%s' object does not support item assignment", left.Type())
		}
	case *ast.AttributeExpression:
		left := Eval(ctx, state, obj.Left, env)
		if IsError(left) {
			return left
		}
		// 设置回赋值所在的行号
		state.UpdateLocation(assign)
		name := obj.Attribute.Value
		if attr, ok := left.(object.Attributable); ok {
			ret := attr.SetAttribute(name, val)
			if IsError(ret) {
				state.HandleError(ret)
			}
			return ret
		}
		return state.NewError("'%s' object does not support set attribute", left.Type())
	default:
		// 设置回赋值所在的行号
		state.UpdateLocation(assign)
		return state.Unreachable("assign")
	}

}

func evalIfStatement(
	ctx context.Context,
	state *WeiState,
	is *ast.IfStatement,
	env *object.Environment,
) object.Object {
	for _, branch := range is.IfBranches {
		condition := Eval(ctx, state, branch.Condition, env)
		if IsError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(ctx, state, branch.Body, env)
		}
	}

	if is.ElseBody != nil {
		return Eval(ctx, state, is.ElseBody, env)
	} else {
		return object.NULL
	}
}

func evalWhileStatement(
	ctx context.Context,
	state *WeiState,
	ws *ast.WhileStatement,
	env *object.Environment,
) object.Object {
	for {
		condition := Eval(ctx, state, ws.Condition, env)
		if IsError(condition) {
			return condition
		}
		if isTruthy(condition) {
			encolosedEnv := object.NewEnclosedEnvironment(env)
			val := Eval(ctx, state, ws.Body, encolosedEnv)
			switch val.(type) {
			case *object.ReturnValue, *object.Error:
				return val
			case *object.BreakValue:
				return nil
			}
		} else {
			break
		}
	}
	return nil
}

func evalForInStatement(
	ctx context.Context,
	state *WeiState,
	forInStmt *ast.ForInStatement,
	env *object.Environment,
) object.Object {
	obj := Eval(ctx, state, forInStmt.Expr, env)
	// 设置 in 后表达式的行号
	state.UpdateLocation(forInStmt.Expr)
	iterable, ok := obj.(object.Iterable)
	if !ok {
		return state.NewError("'%s' object is not iterable", obj.Type())
	}
	iterator := iterable.Iter()
	for {
		// 设置 in 后表达式的行号
		state.UpdateLocation(forInStmt.Expr)
		enclosedEnv := object.NewEnclosedEnvironment(env)
		nextVal := iterator.Next()
		if nextVal == object.StopIteration {
			break
		}
		if IsError(nextVal) {
			state.HandleError(nextVal)
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
			state.UpdateLocation(forInStmt)
			ret := object.WrongNumberUnpack(len(values), len(forInStmt.Targets))
			state.HandleError(ret)
			return ret
		}
		for i, target := range forInStmt.Targets {
			enclosedEnv.Add(target.Value, values[i], forInStmt.Con)
		}
		ret := Eval(ctx, state, forInStmt.Body, enclosedEnv)
		if IsError(ret) {
			return ret
		}
	}
	return nil
}

func evalFunction(
	ctx context.Context,
	state *WeiState,
	fn object.Object,
	args []object.Object,
) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return state.NewError("function expected %d arguments but got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		location := fn.Body.GetFileLocation()
		state.CreateFrame(location.Filename, fn.Name)
		evaluated := Eval(ctx, state, fn.Body, extendedEnv)
		state.DestroyFrame()
		if IsError(evaluated) {
			return evaluated
		}
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		ret := fn.Fn(args...)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret
	case *object.BoundBuiltinMethod:
		ret := fn.Fn(fn.This, args...)
		if IsError(ret) {
			state.HandleError(ret)
		}
		return ret
	case *object.Class:
		return evalClassCall(ctx, state, fn, args)
	case *object.BoundMethod:
		function := fn.Function()
		if len(args) != len(function.Parameters) {
			return state.WrongNumberArgument(function.Name, len(args), len(function.Parameters))
		}
		extendedEnv := extendFunctionEnv(function, args)
		extendedEnv.Pass("this", fn.This(), true)
		extendedEnv.Pass("cls", fn.Class(), true)
		extendedEnv.Pass("super", fn.Super(), true)
		location := function.Body.GetFileLocation()
		state.CreateFrame(location.Filename, function.Name)
		evaluated := Eval(ctx, state, function.Body, extendedEnv)
		state.DestroyFrame()
		if IsError(evaluated) {
			return evaluated
		}
		return unwrapReturnValue(evaluated)
	case *object.BoundClassMethod:
		function := fn.Function()
		if len(args) != len(function.Parameters) {
			return state.WrongNumberArgument(function.Name, len(args), len(function.Parameters))
		}
		extendedEnv := extendFunctionEnv(function, args)
		extendedEnv.Pass("cls", fn.Class(), true)
		extendedEnv.Pass("super", fn.Super(), true)
		location := function.Body.GetFileLocation()
		state.CreateFrame(location.Filename, function.Name)
		evaluated := Eval(ctx, state, function.Body, extendedEnv)
		state.DestroyFrame()
		if IsError(evaluated) {
			return evaluated
		}
		return unwrapReturnValue(evaluated)
	default:
		return state.NewError("not a function: '%s'", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, parameter := range fn.Parameters {
		env.Pass(parameter.Value, args[paramIdx], false)
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
	state *WeiState,
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(ctx, state, exp, env)
		if IsError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalDictLiteral(
	ctx context.Context,
	state *WeiState,
	node *ast.DictLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(ctx, state, keyNode, env)
		if IsError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			state.UpdateLocation(keyNode)
			return state.NewError("unhashable type: '%s'", key.Type())
		}

		value := Eval(ctx, state, valueNode, env)
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
		return object.NewError("unsupported operand type for -: '%s'", right.Type())
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
		return object.NewError("unsupported operand type for +: '%s'", right.Type())
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
		return object.NewError("bad operand type for unary +: '%s'", right.Type())
	}

	value := right.(*object.Integer).Value
	return object.NewInteger(^value)
}

// 计算 and or 逻辑运算
func evalLogicalOperation(
	ctx context.Context,
	state *WeiState,
	env *object.Environment,
	node *ast.BinaryOpExpression,
) object.Object {
	if node.Operator == "and" {
		left := Eval(ctx, state, node.Left, env)
		if IsError(left) {
			return left
		}
		if !isTruthy(left) {
			return left
		}
		right := Eval(ctx, state, node.Right, env)
		if IsError(right) {
			return right
		}
		return right
	} else {
		left := Eval(ctx, state, node.Left, env)
		if IsError(left) {
			return left
		}
		if isTruthy(left) {
			return left
		}
		right := Eval(ctx, state, node.Right, env)
		if IsError(right) {
			return right
		}
		return right
	}
}

func evalBinaryOpExpression(
	ctx context.Context,
	state *WeiState,
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
	default:
		return state.NewError("unsupported operand type for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
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
		return object.NewError("unsupported operand type for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
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
