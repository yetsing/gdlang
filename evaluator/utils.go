package evaluator

import "weilang/object"

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	case object.NULL:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Integer:
			return obj.Value != 0
		case *object.String:
			return len(obj.Value) != 0
		case *object.List:
			return len(obj.Elements) != 0
		case *object.Dict:
			return len(obj.Pairs) != 0
		}
		return true
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return object.NULL
}

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.TypeIs(object.ERROR_OBJ)
	}
	return false
}
