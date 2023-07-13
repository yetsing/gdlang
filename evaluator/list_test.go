package evaluator

import (
	"testing"
	"weilang/object"
)

func TestListBuiltinAttributeReference(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`[].ddd(1)`, "'list' object has not attribute 'ddd'", true},
		{`[].w`, "'list' object has not attribute 'w'", true},

		{`var a = []; a.append(); a[-1]`, "want at least 1 arguments", true},
		{`var a = []; a.append(0); a[-1]`, 0, false},
		{`var a = []; a.append(0, 1); a[-1]`, 1, false},
		{`var a = []; a.append(0, 1, 2); a[-1]`, 2, false},
		{`var a = []; a.append(0, 1, 2, 'abc'); a[-1]`, "abc", false},
		{`var a = ["abc"]; a.append(0, 1, 2, 'ddd'); a.append('x'); a[0]`, "abc", false},

		{`var a = []; a.extend(); a[-1]`, "wrong number of arguments. got=0, want=1", true},
		{`var a = []; a.extend(1, 2); a[-1]`, "wrong number of arguments. got=2, want=1", true},
		{`var a = []; a.extend(1); a[-1]`, "wrong argument type: 'int' at 1", true},
		{`var a = []; a.extend([0]); a[-1]`, 0, false},
		{`var a = []; a.extend([0, 'abc']); a[-1]`, "abc", false},

		{`var a = []; a.insert();`, "wrong number of arguments. got=0, want=2", true},
		{`var a = []; a.insert(1);`, "wrong number of arguments. got=1, want=2", true},
		{`var a = []; a.insert(1, 2, 3);`, "wrong number of arguments. got=3, want=2", true},
		{`var a = []; a.insert('a', 2);`, "wrong argument type: 'str' at 1", true},
		{`var a = []; a.insert(0, 1); a[0]`, 1, false},
		{`var a = []; a.insert(10, 1); a[0]`, 1, false},
		{`var a = [1, 2, 3]; a.insert(1, 1); a[1]`, 1, false},
		{`var a = [1, 2, 3]; a.insert(0, 'abc'); a[0]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(-3, 'abc'); a[0]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(1, 'abc'); a[1]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(-2, 'abc'); a[1]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(2, 'abc'); a[2]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(-1, 'abc'); a[2]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(3, 'abc'); a[3]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(4, 'abc'); a[3]`, "abc", false},
		{`var a = [1, 2, 3]; a.insert(100, 1); a[-1]`, 1, false},
		{`var a = [1, 2, 3]; a.insert(-100, '1'); a[0]`, "1", false},

		{`var a = [1, 2, 3]; a.pop(1, 2);`, "wrong number of arguments. got=2, want=0-1", true},
		{`var a = [1, 2, 3]; a.pop('a');`, "wrong argument type: 'str' at 1", true},
		{`var a = [1, 2, 3]; a.pop();`, 3, false},
		{`var a = [1, 2, 3]; a.pop(); len(a)`, 2, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop()`, 2, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); len(a)`, 1, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); a[1]`, "list index out of range", true},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); a[0]`, 1, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); a.pop()`, 1, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); a.pop(); len(a)`, 0, false},
		{`var a = [1, 2, 3]; a.pop(); a.pop(); a.pop(); a.pop()`, "pop from empty list", true},
		{`var a = [1, 2, 3]; a.pop(0)`, 1, false},
		{`var a = [1, 2, 3]; a.pop(-3)`, 1, false},
		{`var a = [1, 2, 3]; a.pop(-4)`, 1, false},
		{`var a = [1, 2, 3]; a.pop(0); a[0]`, 2, false},
		{`var a = [1, 2, 3]; a.pop(0); len(a)`, 2, false},
		{`var a = [1, 2, 3]; a.pop(1)`, 2, false},
		{`var a = [1, 2, 3]; a.pop(-2)`, 2, false},
		{`var a = [1, 2, 3]; a.pop(1); a[1]`, 3, false},
		{`var a = [1, 2, 3]; a.pop(1); len(a)`, 2, false},
		{`var a = [1, 2, 3]; a.pop(2)`, 3, false},
		{`var a = [1, 2, 'abc']; a.pop(3)`, "abc", false},
		{`var a = [1, 2, 3]; a.pop(-1)`, 3, false},
		{`var a = [1, 2, 3]; a.pop(2); a[1]`, 2, false},
		{`var a = [1, 2, 3]; a.pop(2); len(a)`, 2, false},

		{`var a = [1, '2', true]; a.reverse(1); a[0]`, "wrong number of arguments. got=1, want=0", true},
		{`var a = [1, '2', true]; a.reverse(1, 2); a[0]`, "wrong number of arguments. got=2, want=0", true},
		{`var a = [1, '2', true]; a.reverse(); a[0]`, true, false},
		{`var a = [1, '2', true]; a.reverse(); a[1]`, "2", false},
		{`var a = [1, '2', true]; a.reverse(); a[2]`, 1, false},
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
