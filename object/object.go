package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "int"
	BOOLEAN_OBJ = "bool"
	NULL_OBJ    = "null"
	ERROR_OBJ   = "error"
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

type TypeError struct {
	Message string
}

func (t *TypeError) Type() ObjectType {
	return ERROR_OBJ
}

func (t *TypeError) TypeIs(objectType ObjectType) bool {
	return t.Type() == objectType
}

func (t *TypeError) TypeNotIs(objectType ObjectType) bool {
	return t.Type() != objectType
}

func (t *TypeError) String() string {
	return fmt.Sprintf("TypeError: %s", t.Message)
}
