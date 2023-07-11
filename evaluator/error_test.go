package evaluator

import (
	"testing"
	"weilang/object"
)

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
