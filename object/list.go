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
		// 如果 list 里面的元素有自身，会导致无限递归
		// 所以需要判断一下是不是自己
		var es string
		if e == l {
			es = "[...]"
		} else {
			es = e.String()
		}
		elements = append(elements, es)
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

var (
	outOfRange           = NewError("list index out of range")
	assignmentOutOfRange = NewError("list assignment index out of range")
)

func (l *List) GetItem(index Object) Object {
	if index.TypeNotIs(INTEGER_OBJ) {
		return NewError("list index expect 'int', got '%s'", index.Type())
	}
	idx := int(index.(*Integer).Value)
	length := len(l.Elements)
	if idx < 0 {
		idx += length
	}
	if idx < 0 || idx >= length {
		return outOfRange
	}
	return l.Elements[idx]
}

func (l *List) SetItem(index, value Object) Object {
	if index.TypeNotIs(INTEGER_OBJ) {
		return NewError("list index expect 'int', got '%s'", index.Type())
	}
	idx := int(index.(*Integer).Value)
	length := len(l.Elements)
	if idx < 0 {
		idx += length
	}
	if idx < 0 || idx >= length {
		return assignmentOutOfRange
	}
	l.Elements[idx] = value
	return nil
}
