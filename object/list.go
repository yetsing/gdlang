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

// pop 弹出 idx 位置的元素，调用者需要保证 idx 在范围内
func (l *List) pop(idx int) Object {
	length := len(l.Elements)
	elements := l.Elements
	ele := elements[idx]
	if idx == length-1 {
		l.Elements = elements[:length-1]
		return ele
	}
	l.Elements = append(elements[:idx], elements[idx+1:]...)
	return ele
}

var (
	outOfRange           = NewError("list index out of range")
	popOutOfRange        = NewError("list pop index out of range")
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

func (l *List) SetAttribute(name string, _ Object) Object {
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
				this.Elements = elements
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
				idx := int(arg.Value)
				if idx < 0 {
					idx += length
				}
				if idx < 0 || idx >= length {
					return popOutOfRange
				}
				if idx == length-1 {
					ele := elements[length-1]
					this.Elements = elements[:length-1]
					return ele
				}
				ele := elements[idx]
				this.Elements = append(elements[:idx], elements[idx+1:]...)
				return ele
			},
		},
		// list.remove(obj)
		"remove": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "remove",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) != 1 {
					return WrongNumberArgument(len(args), 1)
				}
				this := obj.(*List)
				idx := -1
				for i, element := range this.Elements {
					if equal(element, args[0]) {
						idx = i
						break
					}
				}
				if idx != -1 {
					this.pop(idx)
					return this
				}
				return NewError("object not in list")
			},
		},
		// list.reverse()
		"reverse": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "reverse",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) > 0 {
					return WrongNumberArgument(len(args), 0)
				}
				this := obj.(*List)
				elements := this.Elements
				for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
					elements[i], elements[j] = elements[j], elements[i]
				}
				return this
			},
		},
	}}
