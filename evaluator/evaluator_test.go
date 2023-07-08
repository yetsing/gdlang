package evaluator

import (
	"testing"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

func TestWhileStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
var a = 0;
var n = 5;
while (n > 0) {
a = a + n
n = n - 1
}
a`, 15},
		{`
var a = 0;
var n = 5;
while (n > 0) {
if (n % 2 == 0) {
 a = a + n
}
n = n - 1
}
a`, 6},
		{`
var a = 0;
var n = 5;
while (n > 0) {
if (n % 2 != 0) {
   n = n - 1
   continue
}
 a = a + n
n = n - 1
}
a`, 6},
		{`
		var a = 0;
		var n = 5;
		while (n > 0) {
			if (not(n % 2)) {
			 break
			}
			a = a + n
			n = n - 1
		}
		a`,
			5},
		{`
		var a = 0;
		var n = 5;
		while (1) {
			if (n < 0) {
				break
			}
			if (n % 2 == 0) {
				n = n - 1
				continue
				break
			}
			a = a + n
			n = n - 1
		}
		a`,
			9},
		{`
		var a = 0;
		var n = 5;
		while (1) {
			if (n < 0) {
				break
			}
			while (n % 2 == 0) {
				n = n - 1
				break
			}
			a = a + n
			n = n - 1
		}
		a`,
			8},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; fn() {var a = 0}; a", 25},
		{`
				var foo = fn(a) {
					if(1){
						var a = 2 
					}
					return a;
				}
				var b = foo(10)
				b`,
			10,
		},
		{`
var a = 10
if (a > 5) {
  a = 20
}
a`,
			20,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{\n(x + 2)\n}"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = fn(x) { return x; }; identity(5);", 5},
		{"var double = fn(x) { return x * 2; }; double(5);", 10},
		{"var add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"var add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { return x; }(5)", 5},
		{"var a = 10; fn(x) { var a = 4; }(10); a", 10},
		{"con a = 10; fn(a) { return a }(30)", 30},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}

	atests := []struct {
		input string
	}{
		{"var identity = fn(x) { x; }; identity(5);"},
	}

	for _, tt := range atests {
		got := testEval(t, tt.input)
		testNullObject(t, got)
	}
}

func TestClosures(t *testing.T) {
	input := `
var newAdder = fn(x) {
  return fn(y) { return x + y };
};

var addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(t, input), 4)
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"var a = 5; var b = a; b = a + b + 5; b;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestConStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"con a = 5; a;", 5},
		{"con a = 5 * 5; a;", 25},
		{"con a = 5; var b = a; b;", 5},
		{"con a = 5; con b = a; con c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},

		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5 % 2", 5 % 2},
		{"2 + 5 % 2", 2 + 5%2},
		{"3 / 2", 3 / 2},

		{"8 >> 1", 8 >> 1},
		{"8 << 1", 8 << 1},
		{"-8 >> 1", -8 >> 1},
		{"-8 << 1", -8 << 1},
		{"3 * 8 >> 1", 3 * (8 >> 1)},
		{"3 * 8 << 1", 3 * (8 << 1)},
		{"(3 * 8) >> 1", (3 * 8) >> 1},
		{"(3 * 8) << 1", (3 * 8) << 1},
		{"1 | 0", 1 | 0},
		{"1 & 0", 1 & 0},
		{"1 ^ 1", 1 ^ 1},
		{"2 + (1 | 0) + (1 & 0) * (1 ^ 1)", 2 + (1 | 0) + (1&0)*(1^1)},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 <= 2", true},
		{"1 > 2", false},
		{"1 >= 2", false},
		{"1 != 2", true},
		{"1 == 2", false},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},

		{"'' == ``", true},
		{`'' == ""`, true},
		{"'abcd' == `abcd`", true},
		{`'中文' == "中文"`, true},
		{`"中文" == "中文"`, true},
		{`"中文a" == "中文"`, false},
		{`"中文a" != "中文"`, true},
		{`"中文" != "中文"`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNotOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"not true", false},
		{"not false", true},
		{"not null", true},
		{"not 0", true},
		{"not 5", false},
		{"not true", false},
		{"not not false", false},
		{"not not 5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestNumberUnaryOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"-0", 0},
		{"5", 5},
		{"-0b0010", -0b0010},
		{"-0xabcd", -0xabcd},

		{"+5", 5},
		{"+0", 0},
		{"5", 5},
		{"+0b0010", 0b0010},
		{"+0xabcd", +0xabcd},

		{"~5", ^5},
		{"~0", ^0},
		{"5", 5},
		{"~0b0010", ^0b0010},
		{"~0xabcd", ^0xabcd},

		{"-~5", -^5},
		{"+~0", ^0},
		{"5", 5},
		{"-~0b0010", -^0b0010},
		{"-~0xabcd", -^0xabcd},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIfElseStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 < 2) { 10 } else if (2 < 3) { 20 } else {10}", 10},
		{"if (2 < 2) { 10 } else if (2 < 3) { 20 } else {30}", 20},
		{"if (2 < 2) { 10 } else if (3 < 3) { 20 } else {30}", 30},
		{`
var a = 5; 
if ("abc") {
  a = 30
}
a`,
			30,
		},
		{`
var a = 5; 
if ("") {
  a = 30
}
a`,
			5,
		},
		{`
var a = 5; 
if (123) {
  a = 30
}
a`,
			30,
		},
		{`
var a = 5; 
if (0) {
  a = 30
}
a`,
			5,
		},
		{`
var a = 5; 
if ([1]) {
  a = 30
}
a`,
			30,
		},
		{`
var a = 5; 
if ([]) {
  a = 30
}
a`,
			5,
		},
		{`
var a = 5; 
if ({1:1}) {
  a = 30
}
a`,
			30,
		},
		{`
var a = 5; 
if ({}) {
  a = 30
}
a`,
			5,
		},
		{`
var a = 5; 
if (not null) {
  a = 30
}
a`,
			30,
		},
		{`
var a = 5; 
if (null) {
  a = 30
}
a`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
			  if (10 > 1) {
				return 10
			  }
			
			  return 1
			}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"unsupported operand type for +: 'int' and 'bool'",
		},
		{
			"-true",
			"unsupported operand type for -: 'bool'",
		},
		{
			"true + false;",
			"unsupported operand type for +: 'bool' and 'bool'",
		},
		{
			"5; true + false; 5",
			"unsupported operand type for +: 'bool' and 'bool'",
		},
		{
			"if (10 > 1) { true + false; }",
			"unsupported operand type for +: 'bool' and 'bool'",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unsupported operand type for +: 'bool' and 'bool'",
		},
		{
			"return true + false;",
			"unsupported operand type for +: 'bool' and 'bool'",
		},
		{
			"if (10 > true) { true + false; }",
			"unsupported operand type for >: 'int' and 'bool'",
		},
		{
			"foobar",
			"undefined: 'foobar'",
		},
		{
			"var foo = fn(){}; foo(1);",
			"function expected 0 arguments but got 1",
		},
		{
			`"5" - "true";`,
			"unsupported operand type for -: 'str' and 'str'",
		},
		{
			"[1, 2, 3][3]",
			"list index out of range",
		},
		{
			"{[]: 1}",
			"unhashable type: 'list'",
		},
		{
			`{"foo": 5}["bar"]`,
			"key 'bar' does not exist",
		},
		{
			`{}["foo"]`,
			"key 'foo' does not exist",
		},
		{
			`{}[true]`,
			"key 'true' does not exist",
		},
		{
			`{}[fn(){}]`,
			"unhashable type: 'function'",
		},
		{
			`con a = 1; a = 2`,
			"cannot assign to constant: 'a'",
		},
		{
			`var a = 1; var a = 1`,
			"variable name 'a' redeclared in this block",
		},
		{
			`
if (1) {
  if (2) {
    if (3) {
       a = 1
    }
  }
}`,
			"undefined: 'a'",
		},
		{
			`while (1) {a}`,
			"undefined: 'a'",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(t, input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`abs(1)`, 1, false},
		{`abs(-1)`, 1, false},
		{`abs(-1, 2)`, "wrong number of arguments. got=2, want=1", true},
		{`abs('')`, "wrong argument type for abs(): 'str'", true},

		{`bin(3) == "0b11"`, true, false},
		{`bin(-10) == "-0b1010"`, true, false},
		{`bin(0)`, "0b0", false},
		{`bin(0, 0) == "0b0"`, "wrong number of arguments. got=2, want=1", true},
		{`bin("")`, "wrong argument type for bin(): 'str'", true},
		{`bin([])`, "wrong argument type for bin(): 'list'", true},

		{`len("")`, 0, false},
		{`len("中文")`, 2, false},
		{`len("a中文")`, 3, false},
		{`len("four")`, 4, false},
		{`len("hello world")`, 11, false},
		{`len(1)`, "wrong argument type for len(): 'int'", true},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1", true},
		{`len([])`, 0, false},
		{`len([1 + 1])`, 1, false},
		{`len(["one", "two"])`, 2, false},
		{`len(["one", "two",])`, 2, false},

		{`hex(255)`, "0xff", false},
		{`hex(-42)`, "-0x2a", false},
		{`hex(0x1234)`, "0x1234", false},
		{`hex(0x1234, 1)`, "wrong number of arguments. got=2, want=1", true},
		{`hex('')`, "wrong argument type for hex(): 'str'", true},

		{`oct(8)`, "0o10", false},
		{`oct(-56)`, "-0o70", false},
		{`oct(0o1234)`, "0o1234", false},
		{`oct(0o1234, 1)`, "wrong number of arguments. got=2, want=1", true},
		{`oct('')`, "wrong argument type for oct(): 'str'", true},

		{`print()`, object.NULL, false},
		{`print(1, -1, "abc", true, false, null, [], {})`, object.NULL, false},

		{`type(1) == "str"`, false, false},
		{`type('') == "str"`, true, false},
		{`type([])`, "list", false},
		{`type({})`, "dict", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			if tt.isError {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. got=%T (%+v)",
						evaluated, evaluated)
					continue
				}
				if errObj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q",
						expected, errObj.Message)
				}
			} else {
				strObj, ok := evaluated.(*object.String)
				if !ok {
					t.Errorf("object is not string. got=%T (%+v)",
						evaluated, evaluated)
					continue
				}
				if strObj.Value != expected {
					t.Errorf("wrong string value. expected=%q, got=%q",
						expected, strObj.Value)
				}
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		case *object.Null:
			testNullObject(t, evaluated)
		}
	}
}

func TestListLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.List)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestListSubscriptionExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{
			"[1, 2, 3][0]",
			1,
			false,
		},
		{
			"[1, 2, 3][1]",
			2,
			false,
		},
		{
			"[1, 2, 3][2]",
			3,
			false,
		},
		{
			"var i = 0; [1][i];",
			1,
			false,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
			false,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
			false,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
			false,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
			false,
		},

		{
			"[1, 2, 3][-1]",
			3,
			false,
		},

		{
			"[1, 2, 3][3]",
			"list index out of range",
			true,
		},

		{
			`"abc"[0]`,
			"a",
			false,
		},
		{
			`"中文"[0]`,
			"中",
			false,
		},
		{
			`"a中文"[1]`,
			"中",
			false,
		},
		{
			`"中文"[1]`,
			"文",
			false,
		},
		{
			`"中文"[2]`,
			"string index out of range",
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			if tt.isError {
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("expected a error. got=%T (%+v)", evaluated, evaluated)
					continue
				}
				if errObj.Message != expected {
					t.Errorf("expected error msg. expected=%q, got=%q", errObj.Message, expected)
				}
			} else {
				strObj, ok := evaluated.(*object.String)
				if !ok {
					t.Errorf("expected string. got=%T (%+v)", evaluated, evaluated)
					continue
				}
				if strObj.Value != expected {
					t.Errorf("expected string. expected=%q, got=%q", strObj.Value, expected)
				}
			}
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `var two = "two";
    {
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6 / 2,
        4: 4,
        true: 5,
        false: 6
    }`

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Dict)
	if !ok {
		t.Fatalf("Eval didn't return Dict. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		object.TRUE.HashKey():                      5,
		object.FALSE.HashKey():                     6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Dict has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

// evaluator/evaluator_test.go

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	t.Helper()
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	t.Helper()
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	t.Helper()
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
		return false
	}
	return true
}
