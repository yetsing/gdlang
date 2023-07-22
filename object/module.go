package object

import "fmt"

type Module struct {
	filename string
	env      *Environment
	export   map[string]bool
}

func (m *Module) Type() ObjectType {
	return MODULE_OBJ
}

func (m *Module) TypeIs(objectType ObjectType) bool {
	return m.Type() == objectType
}

func (m *Module) TypeNotIs(objectType ObjectType) bool {
	return m.Type() != objectType
}

func (m *Module) String() string {
	return fmt.Sprintf("<module object at '%s'>", m.filename)
}

func (m *Module) Filename() string {
	return m.filename
}

func (m *Module) GetAttribute(name string) Object {
	if _, ok := m.export[name]; !ok {
		return attributeError(string(m.Type()), name)
	}
	if obj, ok := m.env.Get(name); ok {
		return obj
	}
	return attributeError(string(m.Type()), name)
}

func (m *Module) SetAttribute(name string, value Object) Object {
	if _, ok := m.export[name]; !ok {
		return attributeError(string(m.Type()), name)
	}
	return m.env.Set(name, value)
}

func (m *Module) GetEnv() *Environment {
	return m.env
}

func (m *Module) AddExport(name string) {
	m.export[name] = true
}

func NewModule(filename string) *Module {
	env := NewEnvironment()
	weiObj := newWei()
	weiObj.Add("filename", NewString(filename))
	env.AddWei(weiObj)
	export := make(map[string]bool)
	return &Module{
		filename: filename,
		env:      env,
		export:   export,
	}
}
