package object

import (
	"bytes"
	"strings"
)

type List struct {
	*attributeStore
	Elements []Object
}

func NewList(elements []Object) *List {
	return &List{
		attributeStore: listAttr,
		Elements:       elements,
	}
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

func (l *List) GetAttribute(name string) Object {
	ret := l.attributeStore.get(l, name)
	if ret != nil {
		return ret
	}
	return attributeError(string(l.Type()), name)
}

// ================================
// list 对象的内置属性和方法
// ================================
var listAttr = &attributeStore{
	attribute: map[string]Object{
		// list.append(*objs)
		"append": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "append",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) == 0 {
					return atLeastOneArgument
				}

				this := obj.(*List)
				this.Elements = append(this.Elements, args...)
				return this
			},
		},
		// list.extend(list2)
		"extend": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "extend",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) != 1 {
					return WrongNumberArgument(len(args), 1)
				}

				this := obj.(*List)
				arg, ok := args[0].(*List)
				if !ok {
					return wrongArgumentTypeAt(args[0].Type(), 1)
				}
				this.Elements = append(this.Elements, arg.Elements...)
				return this
			},
		},
		// list.insert(i, obj)
		"insert": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "insert",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) != 2 {
					return WrongNumberArgument(len(args), 2)
				}

				this := obj.(*List)
				index, ok := args[0].(*Integer)
				if !ok {
					return wrongArgumentTypeAt(args[0].Type(), 1)
				}
				val := args[1]
				length := len(this.Elements)
				idx := convertRange(int(index.Value), length)
				if idx >= length {
					this.Elements = append(this.Elements, val)
					return this
				}
				// 增加一个空间
				elements := append(this.Elements, NULL)
				copy(elements[idx+1:], elements[idx:])
				elements[idx] = val
				return this
			},
		},
		// list.pop() or list.pop(i)
		// list.pop() 弹出最后一个元素
		"pop": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "pop",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) > 1 {
					return WrongNumberArgument2(len(args), 0, 1)
				}
				this := obj.(*List)
				elements := this.Elements
				length := len(elements)
				if len(args) == 0 {
					if length == 0 {
						return NewError("pop from empty list")
					}
					ele := elements[length-1]
					this.Elements = elements[:length-1]
					return ele
				}
				arg, ok := args[0].(*Integer)
				if !ok {
					return wrongArgumentTypeAt(args[0].Type(), 1)
				}
				if length == 0 {
					return NewError("pop from empty list")
				}
				idx := convertRange(int(arg.Value), length)
				if idx >= length-1 {
					ele := elements[length-1]
					this.Elements = elements[:length-1]
					return ele
				}
				ele := elements[idx]
				this.Elements = append(elements[:idx], elements[idx+1:]...)
				return ele
			},
		},
	}}
