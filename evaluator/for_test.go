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
for (var i, e in m) {
n = n + e
}
n`, "'int' object is not iterable",
			true},
		{`
var m = [1, 2, 3]
var n = 0
for (con e in m) {
n = n + i + e
}
n`, "unpack got=2, want=1",
			true},
		{`
var m = [1, 2, 3]
var n = 0
for (con i, e in m) {
n = n + i + e
}
n`, 9,
			false},
		{`
		var m = {0:1, 1:2, 2:3}
		var n = 0
		for (con i, e in m) {
		n = n + i + e
		}
		n`, 9,
			false},
		{`
		var m = "abcd"
		var c = ""
		for (con i, e in m) {
		c = e
		}
		c`, "d",
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
