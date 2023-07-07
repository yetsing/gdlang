package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ        = "int"
	BOOLEAN_OBJ        = "bool"
	NULL_OBJ           = "null"
	ERROR_OBJ          = "error"
	RETURN_VALUE_OBJ   = "return_value"
	FUNCTION_OBJ       = "function"
	STRING_OBJ         = "str"
	BUILTIN_OBJ        = "builtin"
	LIST_OBJ           = "list"
	DICT_OBJ           = "dict"
	CONTINUE_VALUE_OBJ = "continue_value"
	BREAK_VALUE_OBJ    = "break_value"
)

type Object interface {
	Type() ObjectType
	TypeIs(objectType ObjectType) bool
	TypeNotIs(objectType ObjectType) bool
	String() string
}

func TypeIn(obj Object, a ...ObjectType) bool {
	if obj == nil {
		return false
	}
	for _, objectType := range a {
		if obj.TypeIs(objectType) {
			return true
		}
	}
	return false
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) TypeIs(objectType ObjectType) bool {
	return rv.Type() == objectType
}

func (rv *ReturnValue) TypeNotIs(objectType ObjectType) bool {
	return rv.Type() != objectType
}

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *Builtin) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *Builtin) String() string {
	return fmt.Sprintf("<builtin function %s>", b.Name)
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

//==========================
// 两个特殊值，用于处理 continue break 语句
//==========================

type ContinueValue struct {
}

func (c *ContinueValue) Type() ObjectType {
	return CONTINUE_VALUE_OBJ
}

func (c *ContinueValue) TypeIs(objectType ObjectType) bool {
	return c.Type() == objectType
}

func (c *ContinueValue) TypeNotIs(objectType ObjectType) bool {
	return c.Type() != objectType
}

func (c *ContinueValue) String() string {
	return "continue"
}

type BreakValue struct {
}

func (b *BreakValue) Type() ObjectType {
	return BREAK_VALUE_OBJ
}

func (b *BreakValue) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BreakValue) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BreakValue) String() string {
	return "break"
}
