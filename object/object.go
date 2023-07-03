package object

import (
	"fmt"
	"weilang/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "int"
	BOOLEAN_OBJ      = "bool"
	NULL_OBJ         = "null"
	ERROR_OBJ        = "error"
	RETURN_VALUE_OBJ = "return_value"
	FUNCTION_OBJ     = "function"
)

type Object interface {
	Type() ObjectType
	TypeIs(objectType ObjectType) bool
	TypeNotIs(objectType ObjectType) bool
	String() string
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
