package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"weilang/ast"
)

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

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) TypeIs(objectType ObjectType) bool {
	return i.Type() == objectType
}

func (i *Integer) TypeNotIs(objectType ObjectType) bool {
	return i.Type() != objectType
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *Boolean) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct {
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) TypeIs(objectType ObjectType) bool {
	return n.Type() == objectType
}

func (n *Null) TypeNotIs(objectType ObjectType) bool {
	return n.Type() != objectType
}

func (n *Null) String() string {
	return "null"
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) TypeIs(objectType ObjectType) bool {
	return e.Type() == objectType
}

func (e *Error) TypeNotIs(objectType ObjectType) bool {
	return e.Type() != objectType
}

func (e *Error) String() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
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

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func NewFunction(fl *ast.FunctionLiteral, env *Environment) *Function {
	return &Function{
		Parameters: fl.Parameters,
		Body:       fl.Body,
		Env:        env,
	}
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) TypeIs(objectType ObjectType) bool {
	return f.Type() == objectType
}

func (f *Function) TypeNotIs(objectType ObjectType) bool {
	return f.Type() != objectType
}

func (f *Function) String() string {
	return fmt.Sprintf("<function at %p>", f)
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) TypeIs(objectType ObjectType) bool {
	return s.Type() == objectType
}

func (s *String) TypeNotIs(objectType ObjectType) bool {
	return s.Type() != objectType
}

func (s *String) String() string {
	return s.Value
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
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
	return "<builtin function>"
}

type List struct {
	Elements []Object
}

func (l *List) Type() ObjectType {
	return LIST_OBJ
}

func (l *List) TypeIs(objectType ObjectType) bool {
	return l.Type() == objectType
}

func (l *List) TypeNotIs(objectType ObjectType) bool {
	return l.Type() != objectType
}

func (l *List) String() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range l.Elements {
		elements = append(elements, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Dict struct {
	Pairs map[HashKey]HashPair
}

func (d *Dict) Type() ObjectType {
	return DICT_OBJ
}

func (d *Dict) TypeIs(objectType ObjectType) bool {
	return d.Type() == objectType
}

func (d *Dict) TypeNotIs(objectType ObjectType) bool {
	return d.Type() != objectType
}

func (d *Dict) String() string {
	var out bytes.Buffer

	var elements []string
	for _, pair := range d.Pairs {
		elements = append(elements, fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
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
