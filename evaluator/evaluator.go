package evaluator

import (
	"fmt"
	"weilang/ast"
	"weilang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// 语句
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// 表达式
	case *ast.UnaryExpression:
		operand := Eval(node.Operand)
		return evalUnaryExpression(node.Operator, operand)

	case *ast.BinaryOpExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalBinaryOpExpression(node.Operator, left, right)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
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
		return &object.TypeError{Message: fmt.Sprintf("unknown unary operator '%s'", operator)}
	}
}

func evalNotOperatorExpression(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		switch operand.(type) {
		case *object.Integer:
			obj := operand.(*object.Integer)
			return nativeBoolToBooleanObject(obj.Value == 0)
		}
		return FALSE
	}
}

func evalMinusUnaryOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("bad operand type for unary -: '%s'", right.Type())
		return &object.TypeError{Message: message}
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPlusUnaryOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("bad operand type for unary +: '%s'", right.Type())
		return &object.TypeError{Message: message}
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: value}
}

func evalBitwiseNotOperatorExpression(right object.Object) object.Object {
	if right.TypeNotIs(object.INTEGER_OBJ) {
		message := fmt.Sprintf("bad operand type for unary +: '%s'", right.Type())
		return &object.TypeError{Message: message}
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: ^value}
}

func evalBinaryOpExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerBinaryOpExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		msg := fmt.Sprintf("unsupported operand type(s) for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
		return &object.TypeError{Message: msg}
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
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}
	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		msg := fmt.Sprintf("unsupported operand type(s) for %s: '%s' and '%s'",
			operator, left.Type(), right.Type())
		return &object.TypeError{Message: msg}
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
