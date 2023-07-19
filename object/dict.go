package object

type HashPair struct {
	Key   Object
	Value Object
}

type Dict struct {
	*attributeStore
	Pairs map[HashKey]HashPair
}

func NewDict(pairs map[HashKey]HashPair) *Dict {
	return &Dict{attributeStore: dictAttr, Pairs: pairs}
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
	visited := make(map[Object]bool)
	return objectString(d, visited)
}

func (d *Dict) GetItem(key Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("unhashable type: '%s'", key.Type())
	}

	pair, ok := d.Pairs[hashKey.HashKey()]
	if !ok {
		return NewError("key '%s' does not exist", key.String())
	}
	return pair.Value
}

func (d *Dict) SetItem(key, value Object) Object {
	hashKey, ok := key.(Hashable)
	if !ok {
		return NewError("unhashable type: '%s'", key.Type())
	}
	d.Pairs[hashKey.HashKey()] = HashPair{
		Key:   key,
		Value: value,
	}
	return nil
}

func (d *Dict) Iter() Iterator {
	return NewDictIterator(d)
}

func (d *Dict) GetAttribute(name string) Object {
	ret := d.attributeStore.get(d, name)
	if ret != nil {
		return ret
	}
	key := NewString(name)
	pair, ok := d.Pairs[key.HashKey()]
	if ok {
		return pair.Value
	}
	return attributeError(string(d.Type()), name)
}

func (d *Dict) SetAttribute(name string, value Object) Object {
	key := NewString(name)
	d.Pairs[key.HashKey()] = HashPair{
		Key:   key,
		Value: value,
	}
	return nil
}

// ================================
// dict 对象的内置属性和方法
// ================================
var dictAttr = &attributeStore{
	attribute: map[string]Object{
		"get": &BuiltinMethod{
			ctype: DICT_OBJ,
			name:  "get",
			Fn: func(obj Object, args ...Object) Object {
				argc := len(args)
				if argc == 0 || argc > 2 {
					return WrongNumberArgument2(argc, 1, 2)
				}
				key, ok := args[0].(Hashable)
				if !ok {
					return NewError("unhashable type: '%s'", args[0].Type())
				}
				var defaultValue Object
				if argc == 1 {
					defaultValue = NULL
				} else {
					defaultValue = args[1]
				}
				this := obj.(*Dict)
				pair, ok := this.Pairs[key.HashKey()]
				if ok {
					return pair.Value
				}
				return defaultValue
			},
		},
		"has": &BuiltinMethod{
			ctype: LIST_OBJ,
			name:  "has",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) != 1 {
					return WrongNumberArgument(len(args), 1)
				}
				this := obj.(*Dict)
				key, ok := args[0].(Hashable)
				if !ok {
					return NewError("unhashable type: '%s'", args[0].Type())
				}
				_, ok = this.Pairs[key.HashKey()]
				return NativeBoolToBooleanObject(ok)
			},
		},
		"pop": &BuiltinMethod{
			ctype: DICT_OBJ,
			name:  "pop",
			Fn: func(obj Object, args ...Object) Object {
				argc := len(args)
				if argc == 0 || argc > 2 {
					return WrongNumberArgument2(argc, 1, 2)
				}
				key, ok := args[0].(Hashable)
				if !ok {
					return NewError("unhashable type: '%s'", args[0].Type())
				}
				var defaultValue Object
				if argc == 1 {
					defaultValue = NULL
				} else {
					defaultValue = args[1]
				}
				this := obj.(*Dict)
				hk := key.HashKey()
				pair, ok := this.Pairs[hk]
				if ok {
					delete(this.Pairs, hk)
					return pair.Value
				}
				return defaultValue
			},
		},
		"setdefault": &BuiltinMethod{
			ctype: DICT_OBJ,
			name:  "setdefault",
			Fn: func(obj Object, args ...Object) Object {
				argc := len(args)
				if argc == 0 || argc > 2 {
					return WrongNumberArgument2(argc, 1, 2)
				}
				key, ok := args[0].(Hashable)
				if !ok {
					return NewError("unhashable type: '%s'", args[0].Type())
				}
				var defaultValue Object
				if argc == 1 {
					defaultValue = NULL
				} else {
					defaultValue = args[1]
				}
				this := obj.(*Dict)
				hk := key.HashKey()
				pair, ok := this.Pairs[hk]
				if ok {
					return pair.Value
				}
				this.Pairs[hk] = HashPair{
					Key:   args[0],
					Value: defaultValue,
				}
				return defaultValue
			},
		},
		"update": &BuiltinMethod{
			ctype: DICT_OBJ,
			name:  "update",
			Fn: func(obj Object, args ...Object) Object {
				if len(args) != 1 {
					return WrongNumberArgument(len(args), 1)
				}
				this := obj.(*Dict)
				other, ok := args[0].(*Dict)
				if !ok {
					return WrongArgumentTypeAt(args[0].Type(), 1)
				}
				for key, pair := range other.Pairs {
					this.Pairs[key] = pair
				}
				return this
			},
		},
	},
}
