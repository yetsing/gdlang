package evaluator

import (
	"context"
	"weilang/ast"
	"weilang/object"
)

func evalClassDefine(
	ctx context.Context,
	state *WeiState,
	env *object.Environment,
	node *ast.ClassDefineStatement,
) object.Object {
	cls := object.NewClass(node.Name)
	err := env.Add(node.Name, cls, true)
	if IsError(err) {
		state.HandleError(err)
		return err
	}

	for _, statement := range node.Body.Statements {
		ret := evalClassStatement(ctx, state, env, cls, statement)
		if IsError(ret) {
			return ret
		}
	}
	return cls
}

func evalClassStatement(
	ctx context.Context,
	state *WeiState,
	env *object.Environment,
	class *object.Class,
	node ast.Statement,
) object.Object {
	state.UpdateLocation(node)

	switch node := node.(type) {
	case *ast.ClassVariableDeclarationStatement:
		var val object.Object
		if node.Expr != nil {
			val = Eval(ctx, state, node.Expr, env)
			if IsError(val) {
				return val
			}
		}
		if node.Class {
			ret := class.AddClassMember(node.Name.Value, val, node.Con)
			if IsError(ret) {
				state.HandleError(ret)
				return ret
			}
		} else {
			ret := class.AddMember(node.Name.Value, val, node.Con)
			if IsError(ret) {
				state.HandleError(ret)
				return ret
			}
		}
		return nil
	case *ast.ClassMethodDefineStatement:
		function := object.NewFunction(node.Function, env)
		if node.Class {
			ret := class.AddClassMethod(function.Name, function)
			if IsError(ret) {
				state.HandleError(ret)
				return ret
			}
		} else {
			ret := class.AddMethod(function.Name, function)
			if IsError(ret) {
				state.HandleError(ret)
				return ret
			}
		}
		return nil
	default:
		return state.Unreachable("unknown statement in class block")
	}
}

func evalClassCall(
	ctx context.Context,
	state *WeiState,
	class *object.Class,
	args []object.Object,
) object.Object {
	ins := object.NewInstance(class)
	initMethod := ins.GetMethod("__init__")
	if initMethod != nil {
		ret := evalFunction(ctx, state, initMethod, args)
		if IsError(ret) {
			return ret
		}
	} else {
		if len(args) != 0 {
			return state.WrongNumberArgument("__init__", len(args), 0)
		}
	}
	ret := ins.Ready()
	if IsError(ret) {
		state.HandleError(ret)
	}
	return ret
}
