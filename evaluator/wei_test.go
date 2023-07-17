package evaluator

import (
	"testing"
	"weilang/object"
)

func TestWeiOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`wei.hello`, "undefined: 'wei.hello'", true},
		{`wei.filename`, "", false},
		{`wei.import('notfound')`, "Not found module filename: notfound", true},
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
		case []string:
			list, ok := evaluated.(*object.List)
			if !ok {
				t.Errorf("object is not List. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(list.Elements) != len(expected) {
				t.Errorf("list length not equal. got=%d, want=%d", len(list.Elements), len(expected))
				continue
			}
			for i, s := range expected {
				ele := list.Elements[i]
				testStringObject(t, ele, s)
			}
		default:
			t.Errorf("invalid case")
		}
	}
}
