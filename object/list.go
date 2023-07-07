package object

import (
	"bytes"
	"strings"
)

type List struct {
	Elements []Object
}

func NewList(elements []Object) *List {
	return &List{Elements: elements}
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
