package object

import "fmt"

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

type BuiltinMethodFunction func(obj Object, args ...Object) Object

type BuiltinMethod struct {
	ctype ObjectType
	name  string
	Fn    BuiltinMethodFunction
}

func (b *BuiltinMethod) Type() ObjectType {
	return BUILTIN_METHOD_OBJ
}

func (b *BuiltinMethod) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BuiltinMethod) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BuiltinMethod) String() string {
	return fmt.Sprintf("<builtin method '%s' of '%s' object>", b.name, b.ctype)
}

// BoundBuiltinMethod 绑定了实例变量的内置方法
type BoundBuiltinMethod struct {
	*BuiltinMethod
	This Object
}

func (b *BoundBuiltinMethod) Type() ObjectType {
	return BOUND_BUILTIN_METHOD_OBJ
}

func (b *BoundBuiltinMethod) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BoundBuiltinMethod) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BoundBuiltinMethod) String() string {
	return fmt.Sprintf("<bound builtin method '%s' of '%s' object>", b.name, b.ctype)
}

func (b *BoundBuiltinMethod) GetAttribute(name string) Object {
	if name == "__name__" {
		return NewString(b.name)
	}
	return attributeError(string(b.ctype), name)
}
