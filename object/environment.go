package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, constants: c, outer: nil}
}

type Environment struct {
	store map[string]Object
	// constants 保存声明为常量的名字
	constants map[string]bool
	outer     *Environment
}

func (e *Environment) Add(name string, val Object, isConstant bool) Object {
	if _, ok := e.store[name]; ok {
		return NewError("variable name '%s' redeclared in this block", name)
	}
	e.store[name] = val
	if isConstant {
		e.constants[name] = true
	}
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
	if _, ok := e.constants[name]; ok {
		return true
	}
	if e.outer != nil {
		return e.outer.isConstant(name)
	}
	return false
}

func (e *Environment) Set(name string, val Object) Object {
	if e.isConstant(name) {
		return NewError("cannot assign to constant: '%s'", name)
	}
	e.store[name] = val
	return val
}

// Pass 用于函数调用传值
func (e *Environment) Pass(name string, val Object) Object {
	e.store[name] = val
	return val
}
