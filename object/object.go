package object

type ObjectType string

const (
	INTEGER_OBJ              = "int"
	BOOLEAN_OBJ              = "bool"
	NULL_OBJ                 = "null"
	ERROR_OBJ                = "error"
	RETURN_VALUE_OBJ         = "return_value"
	FUNCTION_OBJ             = "function"
	STRING_OBJ               = "str"
	BUILTIN_OBJ              = "builtin"
	LIST_OBJ                 = "list"
	DICT_OBJ                 = "dict"
	CONTINUE_VALUE_OBJ       = "continue_value"
	BREAK_VALUE_OBJ          = "break_value"
	BUILTIN_METHOD_OBJ       = "builtin_method"
	BOUND_BUILTIN_METHOD_OBJ = "bound_builtin_method"
	MODULE_OBJ               = "module"
	WEI_OBJ                  = "wei"
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

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Attributable interface {
	GetAttribute(name string) Object
	SetAttribute(name string, value Object) Object
}

type attributeStore struct {
	attribute map[string]Object
}

func (a *attributeStore) get(object Object, name string) Object {
	val, ok := a.attribute[name]
	if ok {
		switch rval := val.(type) {
		case *BuiltinMethod:
			return &BoundBuiltinMethod{
				BuiltinMethod: rval,
				This:          object,
			}
		default:
			return rval
		}
	}
	return nil
}
