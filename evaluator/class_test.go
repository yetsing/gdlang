package evaluator

import "testing"

func TestClassStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`
class Foo {
var a
con a}
Foo.a`,
			"'a' redeclared in this block",
			true},
		{`
class Foo {
var class.b = 1
con class.b = 2
}
Foo.a`,
			"'b' redeclared in this block",
			true},
		{`
class Foo {
fn d(){}
fn d(){}
}
Foo.a`,
			"'d' redeclared in this block",
			true},
		{`
class Foo {
fn class.ef(){}
fn class.ef(){}
}
Foo.a`,
			"'ef' redeclared in this block",
			true},
		{`
class Foo {}
Foo.a`,
			"'<class Foo>' object has not attribute 'a'",
			true},
		{`
class Foo {
	var abc
}
Foo.abc`,
			"'<class Foo>' object has not attribute 'abc'",
			true},
		{`
class Foo {
	var class.abc = 1
}
Foo.abc`,
			1,
			false},
		{`
class Foo {
	con class.abc = 1
}
Foo.abc = 2`,
			"cannot assign to constant attribute: 'abc'",
			true},
		{`
class Foo {
	var class.abc = 1
}
Foo.abc = 2
Foo.abc`,
			2,
			false},
		{`
class Foo {
	var abc = 1
}
var foo = Foo(1)
foo.abc`,
			"__init__ wrong number of arguments. got=1, want=0",
			true},
		{`
class Foo {
	var abc = 1
	fn __init__(a) {}
}
var foo = Foo()
foo.abc`,
			"__init__ wrong number of arguments. got=0, want=1",
			true},
		{`
class Foo {
	var abc
}
var foo = Foo()
`,
			"Foo object does not initialize attribute: 'abc'",
			true},
		{`
class Foo {
	var abc = 1
}
var foo = Foo()
foo.abc`,
			1,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn get() {return this.abc}
}
var foo = Foo(123)
foo.abc`,
			123,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn get() {return this.abc}
}
var foo = Foo(234)
foo.inc()
foo.get()`,
			235,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn class.get() {return 111}
}
var foo = Foo(234)
foo.inc()
Foo.get()`,
			111,
			false},
		{`
class Foo {
	var abc = 1
	var class.d = 2
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn class.get() {return cls.d + 100}
}
var foo = Foo(234)
Foo.get()`,
			102,
			false},
		{`
class Foo {
	var abc = 1
	var class.d = 1
	fn __init__(abc) {this.abc = abc}
	fn class.inc() {cls.d = cls.d + 1}
	fn class.get() {return cls.d + 100}
}
var foo = Foo(234)
Foo.inc()
Foo.get()`,
			102,
			false},
		{`
class Foo {
	var abc = 1
	con class.d = 1
	fn __init__(abc) {this.abc = abc}
	fn class.inc() {cls.d = cls.d + 1}
	fn class.get() {return cls.d + 100}
}
var foo = Foo(234)
Foo.inc()
Foo.get()`,
			"cannot assign to constant attribute: 'd'",
			true},

		{`
class Foo {
	var class.abc = 1
}
class Goo(Foo){}
Goo.abc`,
			1,
			false},
		{`
		class Foo {
			var class.abc = 1
		}
		class Goo(Foo){}
		Goo.abc = 2
		Goo.abc`,
			2,
			false},
		{`
		class Foo {
			var class.abc = 1
		}
		class Goo(Foo){}
		Goo.abc = 2
		Foo.abc`,
			2,
			false},
		{`
		class Foo {
			var class.abc = 1
		}
		class Goo(Foo){}
		class Hoo(Foo){}
		Goo.abc = 2
		Hoo.abc`,
			2,
			false},
		{`
class Foo {
	var abc = 1
}
class Goo(Foo) {}
var obj = Goo()
obj.abc`,
			1,
			false},
		{`
class Foo {
var abc = 1
}
class Goo(Foo) {
	fn __init__(n) {this.abc = n}
}
var obj = Goo(234)
obj.abc`,
			234,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn get() {return this.abc}
}
class Goo(Foo) {}
var obj = Goo(123)
obj.abc`,
			123,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn get() {return this.abc}
}
class Goo(Foo) {}
var obj = Goo(234)
obj.inc()
obj.get()`,
			235,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn get() {return this.abc}
}
class Goo(Foo) {
	fn inc() {this.abc = this.abc + 100}
}
var obj = Goo(234)
obj.inc()
obj.get()`,
			334,
			false},
		{`
class Foo {
	var class.abc = 1
}
class Goo(Foo) {
	fn class.get() {return super.abc}
}
Goo.get()`,
			1,
			false},
		{`
class Foo {
	var class.abc = 1
}
class Goo(Foo) {
	fn class.set() {super.abc = 2}
}
Goo.set()`,
			"super does not support set attribute",
			true},
		{`
class Foo {
	var class.abc = 1
	fn class.get() {return 0}
}
class Goo(Foo) {
	fn class.get() {return super.get() + 123}
}
Goo.get()`,
			123,
			false},
		{`
class Foo {
	var abc = 1
	fn __init__(abc) {this.abc = abc}
	fn inc() {this.abc = this.abc + 1}
	fn get() {return this.abc}
}
class Goo(Foo) {
	fn inc() {
	  super.inc()
	  this.abc = this.abc + 1
	}
}
var obj = Goo(234)
obj.inc()
obj.get()`,
			236,
			false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			if tt.isError {
				testErrorObject(t, evaluated, expected)
			} else {
				testStringObject(t, evaluated, expected)
			}
		default:
			t.Errorf("impossible type case")
		}
	}
}
