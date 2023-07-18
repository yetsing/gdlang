package object

const weiName = "wei"

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	p := make(map[string]bool)
	return &Environment{store: s, propertys: p, outer: nil}
}

type Environment struct {
	store map[string]Object
	// propertys 保存声明是否是常量
	propertys map[string]bool
	outer     *Environment
}

func (e *Environment) Add(name string, val Object, constant bool) Object {
	if _, ok := e.store[name]; ok {
		return NewError("variable name '%s' redeclared in this block", name)
	}
	e.store[name] = val
	e.propertys[name] = constant
	return val
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) isConstant(name string) bool {
	if constant, ok := e.propertys[name]; ok {
		return constant
	}
	return false
}

func (e *Environment) Set(name string, val Object) Object {
	_, ok := e.store[name]
	if !ok {
		// 找不到，去上层环境设置
		if e.outer != nil {
			return e.outer.Set(name, val)
		}
		return NewError("undefined: '%s'", name)
	}
	if e.isConstant(name) {
		return NewError("cannot assign to constant: '%s'", name)
	}
	e.store[name] = val
	return val
}

// Pass 用于函数调用传值、for-in 传值
func (e *Environment) Pass(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) AddWei(val *wei) {
	e.Add(weiName, val, true)
}

func (e *Environment) GetFromWei(name string) Object {
	val, ok := e.Get(weiName)
	if ok {
		weiObj := val.(*wei)
		return weiObj.GetAttribute(name)
	}
	panic("unreachable: not found wei from environment")
}
