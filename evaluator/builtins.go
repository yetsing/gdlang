package evaluator

import (
	"bytes"
	"fmt"
	"strconv"
	"weilang/object"
)

// fastabs 快速 abs 实现，代码来自 http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go.html
// https://github.com/cavaliercoder/go-abs
func fastabs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func abs(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewInteger(fastabs(arg.Value))
	default:
		return object.NewError("wrong argument type for abs(): '%s'", arg.Type())
	}
}

func bin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		s := strconv.FormatInt(arg.Value, 2)
		if s[0] == '-' {
			s = "-0b" + s[1:]
		} else {
			s = "0b" + s
		}
		return object.NewString(s)
	default:
		return object.NewError("wrong argument type for bin(): '%s'", arg.Type())
	}
}

func hex(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		s := strconv.FormatInt(arg.Value, 16)
		if s[0] == '-' {
			s = "-0x" + s[1:]
		} else {
			s = "0x" + s
		}
		return object.NewString(s)
	default:
		return object.NewError("wrong argument type for hex(): '%s'", arg.Type())
	}
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(arg.Length)}
	case *object.List:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Dict:
		return object.NewInteger(int64(len(arg.Pairs)))
	default:
		return object.NewError("wrong argument type for len(): '%s'", arg.Type())
	}
}

func oct(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Integer:
		s := strconv.FormatInt(arg.Value, 8)
		if s[0] == '-' {
			s = "-0o" + s[1:]
		} else {
			s = "0o" + s
		}
		return object.NewString(s)
	default:
		return object.NewError("wrong argument type for oct(): '%s'", arg.Type())
	}
}

func _print(args ...object.Object) object.Object {
	var out bytes.Buffer
	count := len(args)
	for i, arg := range args {
		out.WriteString(arg.String())
		// 如果不是最后一个元素，在后面加一个空格
		if i != count-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString("\n")
	fmt.Printf(out.String())
	return object.NULL
}

func _type(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	arg := args[0]
	switch arg := arg.(type) {
	case *object.Instance:
		return object.NewString(arg.ClassName())
	default:
		return object.NewString(string(arg.Type()))
	}
}

var builtins = map[string]*object.Builtin{
	"abs": {
		Name: "abs",
		Fn:   abs,
	},
	"bin": {
		Name: "bin",
		Fn:   bin,
	},
	// bool(object) -> bool
	// 将对象转化为 bool 值
	"bool": {
		Name: "bool",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.WrongNumberArgument(len(args), 1)
			}
			return object.NativeBoolToBooleanObject(isTruthy(args[0]))
		},
	},
	// ensure(condition, msg)
	// condition 为假时报错，错误信息为传入的 msg
	"ensure": {
		Name: "ensure",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.WrongNumberArgument(len(args), 2)
			}
			if !isTruthy(args[0]) {
				return object.NewError(args[1].String())
			}
			return nil
		},
	},
	"hex": {
		Name: "hex",
		Fn:   hex,
	},
	// int(object) -> int
	// 将对象转化为整数，支持传入字符串、数字
	"int": {
		Name: "int",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.WrongNumberArgument(len(args), 1)
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				v, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return object.NewError(err.Error())
				}
				return object.NewInteger(v)
			default:
				return object.WrongArgumentTypeAt(args[0].Type(), 0)
			}
		},
	},
	"len": {
		Name: "len",
		Fn:   _len,
	},
	"oct": {
		Name: "oct",
		Fn:   oct,
	},
	"print": {
		Name: "print",
		Fn:   _print,
	},
	"type": {
		Name: "type",
		Fn:   _type,
	},
}
