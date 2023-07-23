package object

import (
	"fmt"
)

// Class 对应用户定义的类
type Class struct {
	Name string
	// parent 父类
	parent     *Class
	members    map[string]Object
	methods    map[string]*Function
	conMembers map[string]bool
	// 类属性和类方法
	classMembers    map[string]Object
	classMethods    map[string]*Function
	conClassMembers map[string]bool
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) TypeIs(objectType ObjectType) bool {
	return c.Type() == objectType
}

func (c *Class) TypeNotIs(objectType ObjectType) bool {
	return c.Type() != objectType
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %s>", c.Name)
}

func (c *Class) SetAttribute(name string, value Object) Object {
	if c.isConstantClassMember(name) {
		return NewError("cannot assign to constant attribute: '%s'", name)
	}
	if _, ok := c.classMembers[name]; ok {
		c.classMembers[name] = value
		return nil
	}
	if c.parent != nil {
		return c.parent.SetAttribute(name, value)
	}
	return attributeError(c.String(), name)
}

func (c *Class) GetAttribute(name string) Object {
	return c.getAttribute(c, name)
}

func (c *Class) getAttribute(cls *Class, name string) Object {
	if val, ok := c.classMembers[name]; ok {
		return val
	}

	if val, ok := c.classMethods[name]; ok {
		return &BoundClassMethod{
			define:   c,
			cls:      cls,
			function: val,
		}
	}
	if c.parent != nil {
		return c.parent.getAttribute(cls, name)
	}
	return attributeError(c.String(), name)
}

func (c *Class) getMethod(ins *Instance, name string) *BoundMethod {
	if val, ok := c.methods[name]; ok {
		return &BoundMethod{
			define:   c,
			this:     ins,
			function: val,
		}
	}
	if c.parent != nil {
		return c.parent.getMethod(ins, name)
	}
	return nil
}

func (c *Class) isConstantMember(name string) bool {
	if _, ok := c.conMembers[name]; ok {
		return true
	}
	return false
}

func (c *Class) isConstantClassMember(name string) bool {
	if _, ok := c.conClassMembers[name]; ok {
		return true
	}
	return false
}

func (c *Class) AddMember(name string, defaultValue Object, isConstant bool) Object {
	if _, ok := c.members[name]; ok {
		return NewError("'%s' redeclared in this block", name)
	}
	c.members[name] = defaultValue
	if isConstant {
		c.conMembers[name] = true
	}
	return nil
}

func (c *Class) AddMethod(name string, function *Function) Object {
	if _, ok := c.methods[name]; ok {
		return NewError("'%s' redeclared in this block", name)
	}
	c.methods[name] = function
	return nil
}

func (c *Class) AddClassMember(name string, defaultValue Object, isConstant bool) Object {
	if _, ok := c.classMembers[name]; ok {
		return NewError("'%s' redeclared in this block", name)
	}
	c.classMembers[name] = defaultValue
	if isConstant {
		c.conClassMembers[name] = true
	}
	return nil
}

func (c *Class) AddClassMethod(name string, function *Function) Object {
	if _, ok := c.classMethods[name]; ok {
		return NewError("'%s' redeclared in this block", name)
	}
	c.classMethods[name] = function
	return nil
}

func NewClass(name string, parent *Class) *Class {
	return &Class{
		Name:            name,
		parent:          parent,
		members:         make(map[string]Object),
		methods:         make(map[string]*Function),
		conMembers:      make(map[string]bool),
		classMembers:    make(map[string]Object),
		classMethods:    make(map[string]*Function),
		conClassMembers: make(map[string]bool),
	}
}

// Instance 对应用户创建的实例
type Instance struct {
	class   *Class
	members map[string]Object
	// 是否在初始化过程
	inInit bool
}

func (ins *Instance) Type() ObjectType {
	return INSTANCE_OBJ
}

func (ins *Instance) TypeIs(objectType ObjectType) bool {
	return ins.Type() == objectType
}

func (ins *Instance) TypeNotIs(objectType ObjectType) bool {
	return ins.Type() != objectType
}

func (ins *Instance) String() string {
	return fmt.Sprintf("<%s object at %p>", ins.class.Name, ins)
}

func (ins *Instance) SetAttribute(name string, value Object) Object {
	// 初始化过程中允许设置 con 声明的属性成员
	if !ins.inInit && ins.class.isConstantMember(name) {
		return NewError("cannot assign to constant attribute: '%s'", name)
	}

	if _, ok := ins.members[name]; ok {
		ins.members[name] = value
		return nil
	}
	return attributeError(ins.String(), name)
}

func (ins *Instance) GetAttribute(name string) Object {
	if val, ok := ins.members[name]; ok {
		return val
	}

	val := ins.class.getMethod(ins, name)
	if val != nil {
		return val
	}
	return attributeError(ins.String(), name)
}

func (ins *Instance) SetMember(name string, value Object) {
	ins.members[name] = value
}

func (ins *Instance) GetMethod(name string) *BoundMethod {
	method := ins.class.getMethod(ins, name)
	return method
}

// Ready 检查实例初始化情况，并且做一些通用的初始化工作
func (ins *Instance) Ready() Object {
	for s, object := range ins.members {
		if object == nil {
			return NewError("%s object does not initialize attribute: '%s'", ins.class.Name, s)
		}
	}
	ins.inInit = false
	ins.SetMember("__class__", ins.class)
	return ins
}

func (ins *Instance) ClassName() string {
	return ins.class.Name
}

func NewInstance(class *Class) *Instance {
	var inheritList []*Class
	cls := class
	for cls != nil {
		inheritList = append(inheritList, cls)
		cls = cls.parent
	}

	m := make(map[string]Object, len(class.members))
	for i := len(inheritList) - 1; i >= 0; i-- {
		cls = inheritList[i]
		for s, object := range cls.members {
			m[s] = object
		}
	}
	ins := &Instance{
		class:   class,
		members: m,
	}
	return ins
}

// BoundClassMethod 绑定具体类的类方法
type BoundClassMethod struct {
	// define 方法定义所在的类
	define   *Class
	cls      *Class
	function *Function
}

func (b *BoundClassMethod) Type() ObjectType {
	return BOUND_CLASS_METHOD_OBJ
}

func (b *BoundClassMethod) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BoundClassMethod) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BoundClassMethod) String() string {
	return fmt.Sprintf("<class method '%s' of '%s'>", b.function.Name, b.cls.Name)
}

func (b *BoundClassMethod) Function() *Function {
	return b.function
}

func (b *BoundClassMethod) Class() *Class {
	return b.cls
}

func (b *BoundClassMethod) Super() *Super {
	return &Super{
		define: b.define,
		cls:    b.cls,
		this:   nil,
	}
}

// BoundMethod 绑定具体实例的方法
type BoundMethod struct {
	// define 方法定义所在的类
	define   *Class
	this     *Instance
	function *Function
}

func (b *BoundMethod) Type() ObjectType {
	return BOUND_METHOD_OBJ
}

func (b *BoundMethod) TypeIs(objectType ObjectType) bool {
	return b.Type() == objectType
}

func (b *BoundMethod) TypeNotIs(objectType ObjectType) bool {
	return b.Type() != objectType
}

func (b *BoundMethod) String() string {
	return fmt.Sprintf("<bound method '%s' of '%s'>", b.function.Name, b.this.String())
}

func (b *BoundMethod) Function() *Function {
	return b.function
}

func (b *BoundMethod) This() *Instance {
	return b.this
}

func (b *BoundMethod) Class() *Class {
	return b.this.class
}

func (b *BoundMethod) Super() *Super {
	return &Super{
		define: b.define,
		cls:    nil,
		this:   b.this,
	}
}

/*
super 指向代码块所在类的父类，只支持访问实例方法、类属性、类方法
例子如下

class A {
  fn get() {}
}
class B(A) {
  fn get() {super.get()}
}
class C(B) {
  fn get() {super.get()}
}

C.get 里面的 super 指着是 B
B.get 里面的 super 指着是 A

*/

// Super super 只支持访问实例方法、类属性、类方法，
// 因为实例属性是绑定在实例上的，他的属性值没有父类子类一说；而另外三个则是跟类有关系。
type Super struct {
	// super 代码所在的类
	define *Class
	// 当前调用的类和实例
	cls  *Class
	this *Instance
}

func (s *Super) Type() ObjectType {
	return SUPER_OBJ
}

func (s *Super) TypeIs(objectType ObjectType) bool {
	return s.Type() == objectType
}

func (s *Super) TypeNotIs(objectType ObjectType) bool {
	return s.Type() != objectType
}

func (s *Super) String() string {
	if s.cls != nil {
		return fmt.Sprintf("<%s super in %s)>", s.cls.String(), s.define.String())
	} else {
		return fmt.Sprintf("<%s super in %s)>", s.this.String(), s.define.String())
	}
}

// GetAttribute 只支持访问实例方法、类属性、类方法，
// 因为实例属性是绑定在实例上的，他的属性值没有父类子类一说；而另外三个则是跟类有关系。
func (s *Super) GetAttribute(name string) Object {
	// 这里的 parent 不可能是 nil ，我们会创建一个内置 object 作为继承链的终点
	parent := s.define.parent
	if s.cls != nil {
		// 类属性和类方法
		return parent.getAttribute(s.cls, name)
	} else {
		// 实例方法
		return parent.getMethod(s.this, name)
	}
}

//goland:noinspection GoUnusedParameter
func (s *Super) SetAttribute(name string, value Object) Object {
	return NewError("super does not support set attribute")
}
