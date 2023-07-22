package evaluator

import (
	"testing"
	"weilang/object"
)

func TestDictBuiltinAttributeReference(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`{}.ddd`, "'dict' object has not attribute 'ddd'", true},
		//{`{'a': 1}.a`, 1, false},
		//
		//{`{}.get()`, "wrong number of arguments. got=0, want=1-2", true},
		//{`{}.get([])`, "unhashable type: 'list'", true},
		//{`{}.get(1, 2, 3)`, "wrong number of arguments. got=3, want=1-2", true},
		//{`{}.get(1)`, nil, false},
		//{`{}.get(1, "abc")`, "abc", false},
		//{`{1: 2}.get(1)`, 2, false},
		//
		//{`{}.has()`, "wrong number of arguments. got=0, want=1", true},
		//{`{}.has([])`, "unhashable type: 'list'", true},
		//{`{}.has(1, 2, 3)`, "wrong number of arguments. got=3, want=1", true},
		//{`{}.has(1)`, false, false},
		//{`{1: 2}.has(1)`, true, false},
		//
		//{`{}.pop()`, "wrong number of arguments. got=0, want=1-2", true},
		//{`{}.pop([])`, "unhashable type: 'list'", true},
		//{`{}.pop(1, 2, 3)`, "wrong number of arguments. got=3, want=1-2", true},
		//{`{}.pop(1)`, nil, false},
		//{`{1: 2}.pop(1)`, 2, false},
		//{`{1: 2}.pop(2, 2023)`, 2023, false},
		//
		//{`{}.setdefault()`, "wrong number of arguments. got=0, want=1-2", true},
		//{`{}.setdefault([])`, "unhashable type: 'list'", true},
		//{`{}.setdefault(1, 2, 3)`, "wrong number of arguments. got=3, want=1-2", true},
		//{`{}.setdefault(1)`, nil, false},
		//{`{1: 2}.setdefault(1)`, 2, false},
		//{`{1: 2}.setdefault(2, 2023)`, 2023, false},
		//{`var d = {1: 2}; d.setdefault(2, 2023); len(d)`, 2, false},
		//
		//{`{}.update()`, "wrong number of arguments. got=0, want=1", true},
		//{`{}.update(1)`, "wrong argument type: 'int' at 1", true},
		//{`var d = {}; d.update({}); len(d)`, 0, false},
		//{`var d = {}; d.update({1: 2}); d[1]`, 2, false},
		//{`var d = {1: 100}; d.update({1: 2}); d[1]`, 2, false},
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
			testNullObject(t, evaluated)
		}
	}
}
