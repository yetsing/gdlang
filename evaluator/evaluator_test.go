package evaluator

import (
	"testing"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

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

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	return Eval(program)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
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
