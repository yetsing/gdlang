package evaluator

import "weilang/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Name: "len",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.List:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return object.NewError("object of type '%s' has no len()", args[0].Type())
			}
		},
	},
}
