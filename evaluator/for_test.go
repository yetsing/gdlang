package evaluator

import "testing"

func TestForInStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		isError  bool
	}{
		{`
var m = 1
var n = 0
for (i, e in m) {
n = n + e
}
n`, "'int' object is not iterable",
			true},
		{`
var m = [1, 2, 3]
var n = 0
for (i, e in m) {
n = n + i + e
}
n`, 9,
			false},
		//		{`
		//var m = {0:1, 1:2, 2:3}
		//var n = 0
		//for (i, e in m) {
		//n = n + i + e
		//}
		//n`, 9,
		//			false},
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
