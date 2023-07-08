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
	return object.NewString(string(arg.Type()))
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
	"len": {
		Name: "len",
		Fn:   _len,
	},
	"hex": {
		Name: "hex",
		Fn:   hex,
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
